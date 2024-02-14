package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/chebyrash/promise"
	"github.com/panjf2000/ants/v2"
	"github.com/puzpuzpuz/xsync"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasttemplate"

	"github.com/tuupke/pixie/env"
	"github.com/tuupke/pixie/lifecycle"
)

type (
	Props struct {
		ip net.IP

		*xsync.MapOf[string, string]

		latestData time.Time
	}

	// promiseInteraction is a promise used to interact with external data and the
	// banner pdf.
	promiseInteraction struct {
		callItIn   func()
		pdfPromise *promise.Promise[*os.File]
	}
)

var (
	cpuPool, ioPool  promise.Pool
	printKeys        = strings.Split(env.StringFb("PRINT_KEYS", "*"), ",")
	includeBasicAuth = env.Bool("BASIC_AUTH_IN_DATA")
	basicAuthUser    = env.StringFb("BASIC_AUTH_USERNAME", "ba_username")
	basicAuthPass    = env.StringFb("BASIC_AUTH_PASSWORD", "ba_password")
	alwaysFreshData  = env.Bool("BANNER_DATA_ALWAYS_FRESH")
)

func init() {
	cpuPoolA, err := ants.NewPool(runtime.NumCPU())
	if err != nil {
		panic(err)
	}

	ioPoolA, err := ants.NewPool(runtime.NumCPU() * 5)
	if err != nil {
		panic(err)
	}

	ioPool = promise.FromAntsPool(ioPoolA)
	cpuPool = promise.FromAntsPool(cpuPoolA)
}

type e struct{}

var empty = e{}

func loadValues(log zerolog.Logger, ctx *fasthttp.RequestCtx, jobId int32) promiseInteraction {
	// Load or create a Props instance
	data := LoadFromRequest(ctx)
	isInitial := data.latestData.IsZero() || alwaysFreshData

	log = log.With().IPAddr("for", data.ip).Int32("job-id", jobId).Logger()

	awaitCtx, cancel := context.WithCancel(lifecycle.ApplicationContext())
	waitFor := len(toCall)
	c := make(chan e, waitFor)
	log.Info().Int("num_hooks", waitFor).Msg("loading data")
	for _, set := range toCall {
		ioPool.Go(set.handle(c, log, data))
	}

	// fanin is a promise that awaits until all
	fanin := promise.New(func(resolve func(e), reject func(error)) {
		log.Info().Bool("will-wait", isInitial).Int("webhooks-to-finish", waitFor).Msg("awaiting finish")

		var ctxEnd bool
		for !ctxEnd && waitFor > 0 {
			if isInitial {
				<-c
				waitFor--
				log.Debug().Int("remaining", waitFor).Msg("waiting for more hooks")
			} else {
				select {
				case <-c:
					waitFor--
				case <-awaitCtx.Done():
					log.Warn().Msg("called in")
					ctxEnd = true
				}
			}
		}

		resolve(empty)
	})

	// computes result based on the fetched data, runs on cpuOptimizedPool
	pdfPromise := promise.ThenWithPool(fanin, lifecycle.ApplicationContext(), func(_ e) (*os.File, error) {
		// Load the stat on the pdf
		fn := pdfLocation + "/" + data.ip.String() + ".pdf"
		file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			// Something went really wrong here, unrecoverable
			return nil, fmt.Errorf("encountered error opening file '%v'; %w", fn, err)
		}

		if fi, err := file.Stat(); err == nil && fi.Size() > 0 && fi != nil && fi.ModTime().After(data.latestData) && !data.latestData.IsZero() {
			log.Info().Msg("reusing cached banner")
			return file, err
		}

		log.Err(file.Truncate(0)).Msg("creating new banner, truncated file")
		// Render the pdf, data is either up-to or out-of-date, we don't care!
		return file, BannerPage(log, file, data, printKeys...)
	}, cpuPool)

	return promiseInteraction{callItIn: cancel, pdfPromise: pdfPromise}
}

type mapWriter map[string]string

func (m mapWriter) MarshalZerologObject(e *zerolog.Event) {
	for k, v := range m {
		e.Str(k, v)
	}
}

func replaceParameters(template string, data *Props, webhookname string) (string, mapWriter) {
	// params stores the retrieved parameters, this trick works since a slice is a pointer type.
	d := make(mapWriter)
	return fasttemplate.New(template, "{{", "}}").ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		if tag == "webhook_name" {
			d[tag] = webhookname
			return w.Write([]byte(webhookname))
		}

		v, _ := data.Load(tag)
		d[tag] = v
		return w.Write([]byte(v))
	}), d
}

var (
	imageKey    = env.StringFb("IMAGE_KEY", "image")
	props       = xsync.NewMapOf[Props]()
	keyTemplate = env.String("WEBHOOK_KEY_TEMPLATE")
	downloadTo  = env.StringFb("WEBHOOK_TEMP_DIR", os.TempDir())

	toCallString = env.String("WEBHOOKS_TO_CALL")
	toCall       endpointsSet
)

// Variables used for checking validity
var (
	methods = map[string]struct{}{
		http.MethodGet:     empty,
		http.MethodHead:    empty,
		http.MethodPost:    empty,
		http.MethodPut:     empty,
		http.MethodPatch:   empty,
		http.MethodDelete:  empty,
		http.MethodConnect: empty,
		http.MethodOptions: empty,
		http.MethodTrace:   empty,
	}
)

func parseToCallString(toCallString string) (e endpointsSet, gerr error) {
	// Parse the toCallString

	eps := strings.Split(toCallString, "&&")
	e = make([]endpoints, len(eps))
	for k, set := range eps {
		ep := strings.Split(set, "|")
		e[k] = make([]endpoint, len(ep))
		for kk, end := range ep {
			s := strings.SplitN(end, ";", 3)
			if len(s) < 3 {
				gerr = fmt.Errorf("invalid webhook spec found, expected 3 parts: '%v'", e)
				return
			}

			// Attempt to parse the url
			_, err := url.Parse(s[2])
			if err != nil {
				gerr = fmt.Errorf("could not parse url '%v'; %w", s[2], err)
				return
			}

			if _, ok := methods[s[1]]; !ok {
				gerr = fmt.Errorf("could not parse method '%v' for url '%v', expected values look like GET, POST, DELETE", s[1], s[2])
				return
			}

			e[k][kk] = endpoint{
				name:   s[0],
				method: s[1],
				url:    s[2],
			}
		}
	}

	return
}

func init() {
	if len(toCallString) > 0 {
		var err error
		toCall, err = parseToCallString(toCallString)
		if err != nil {
			panic(err)
		}
	}

	if downloadTo == "" {
		downloadTo = os.TempDir()
		fmt.Println("Empty download dir, using", downloadTo)
	}

	if err := os.MkdirAll(downloadTo, 0755); err != nil {
		panic(fmt.Errorf("could not create download folder '%v'; %w", downloadTo, err))
	}
}

func decodeBasicAuth(auth []byte) (ok bool, user, pass string) {
	i := bytes.IndexByte(auth, ' ')
	if i == -1 || !bytes.EqualFold(auth[:i], []byte("basic")) {
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(string(auth[i+1:]))
	if err != nil {
		return
	}

	credentials := bytes.Split(decoded, []byte(":"))
	if len(credentials) <= 1 {
		return
	}

	user, pass, ok = string(credentials[0]), string(credentials[1]), true

	return
}

func LoadFromRequest(ctx *fasthttp.RequestCtx) *Props {
	// Construct the key from the ip, basic-auth username, and basic-auth password.
	ok, user, pass := decodeBasicAuth(ctx.Request.Header.Peek("Authorization"))
	ip := ctx.RemoteIP()

	var baseData = make(map[string]string)
	segments := strings.Split(strings.Trim(string(ctx.Request.URI().Path()), "/"), "/")
	for _, segment := range segments {
		// Split the segment, only add if the actually is something to add
		keyValues := strings.SplitN(segment, "=", 2)
		if len(keyValues) > 1 {
			baseData[keyValues[0]] = keyValues[1]
		}
	}

	if ok && includeBasicAuth {
		baseData[basicAuthUser] = user
		baseData[basicAuthPass] = pass
	}

	baseData["requesting_ip"] = ip.String()
	return Load(ip, baseData, user, pass, ctx.Request.URI().String())
}

func Load(ip net.IP, baseData map[string]string, segments ...string) *Props {
	key := ip.String() + "::" + strings.Join(segments, "::")
	props, load := (*xsync.MapOf[string, *Props])(props).LoadOrStore(key, &Props{
		ip:    ip,
		MapOf: xsync.NewMapOf[string](),
	})

	if !load {
		for k, v := range baseData {
			props.Store(k, v)
		}
	}

	return props
}

func (p *Props) Reduce(keys ...string) (mp map[string]string) {
	mp = make(map[string]string)
	for _, key := range keys {
		if data, ok := p.Load(key); ok {
			mp[key] = data
		}
	}

	return
}

func (p *Props) json(extra map[string]string) io.Reader {
	// TODO create a pool of buffers to use
	b := new(bytes.Buffer)
	b.WriteByte('{')

	var notFirst bool
	write := func(key, value string) bool {
		if notFirst {
			b.WriteByte(',')
		}
		b.WriteByte('"')
		b.WriteString(strings.Replace(key, `"`, `\"`, -1))
		b.WriteByte('"')
		b.WriteString(":")

		b.WriteByte('"')
		b.WriteString(strings.Replace(value, `"`, `\"`, -1))
		b.WriteByte('"')

		notFirst = true
		return notFirst
	}

	p.Range(write)

	for k, v := range extra {
		_ = write(k, v)
	}

	b.WriteByte('}')

	return b
}
