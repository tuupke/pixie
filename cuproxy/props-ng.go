package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/chebyrash/promise"
	"github.com/panjf2000/ants/v2"
	"github.com/puzpuzpuz/xsync"
	"github.com/rs/zerolog"
	log2 "github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasttemplate"

	"github.com/tuupke/pixie/cuproxy/web"
	"github.com/tuupke/pixie/env"
	"github.com/tuupke/pixie/lifecycle"
)

type Props struct {
	ip net.IP

	*xsync.MapOf[string, string]

	latestData time.Time
}

type promisePair struct {
	callItIn   func()
	pdfPromise *promise.Promise[*os.File]
}

var cpuPool, ioPool promise.Pool
var printKeys = strings.Split(env.StringFb("PRINT_KEYS", "*"), ",")

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

func loadValues(rctx *fasthttp.RequestCtx, jobId int32) promisePair {
	// Load or create a Props instance
	data := LoadFromRequest(rctx)
	isInitial := data.latestData.IsZero()

	log := log2.With().IPAddr("for", data.ip).Int32("job-id", jobId).Logger()

	awaitCtx, cancel := context.WithCancel(lifecycle.ApplicationContext())
	waitFor := toCall.size()
	c := make(chan e, waitFor)
	log.Info().Int("num_hooks", waitFor).Msg("loading data")
	for _, set := range toCall {
		ioPool.Go(func() {
			handleRequests(log, slices.Clone(set), data, lifecycle.ApplicationContext())(func(_ e) {
				c <- empty
			}, func(err error) {
				log.Err(err).Msg("received error calling webhook")
				c <- empty
			})
		})
	}

	fanin := promise.New(func(resolve func(e), reject func(error)) {
		log.Info().Bool("will_wait", isInitial).Int("webhooks-to-finish", waitFor).Msg("awaiting finish")

	outer:
		for {
			if isInitial {
				<-c
				waitFor--
				log.Debug().Int("remaining", waitFor).Msg("waiting for more hooks")
				if waitFor <= 0 {
					break outer
				}
			} else {
				select {
				case <-c:
					waitFor--
					if waitFor <= 0 {
						break outer
					}
				case <-awaitCtx.Done():
					log.Warn().Msg("called in")
					break outer
				}
			}
		}

		resolve(empty)
	})

	// computes result based on the fetched data, runs on cpuOptimizedPool
	pdfPromise := promise.ThenWithPool(fanin, lifecycle.ApplicationContext(), func(_ e) (*os.File, error) {
		// data.
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

	return promisePair{cancel, pdfPromise}
}

func handleRequests(log zerolog.Logger, set endpoints, data *Props, ctx context.Context) func(resolve func(e), reject func(error)) {
	return func(resolve func(e), reject func(error)) {
		for _, ep := range set {
			u, method := ep.url, ep.method

			t := fasttemplate.New(u, "{{", "}}")
			u = t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
				v, _ := data.Load(tag)
				return w.Write([]byte(v))
			})

			var reqBody io.Reader
			if method != http.MethodGet {
				reqBody = data.json()
			}

			ctx, cancel := context.WithTimeout(ctx, time.Second*2)
			respBody, respType, loaded, err := web.Do(ctx, log, u, method, data.ip, reqBody)
			cancel()

			if err != nil {
				reject(fmt.Errorf("executing request; %w", err))
			}

			if loaded {
				switch respType = strings.Split(respType, ";")[0]; respType {
				case "image/jpeg":
					fallthrough
				case "image/png":
					fallthrough
				case "image/gif":
					// Store in file and
					fn := downloadTo + "/" + data.ip.String() + "." + respType[6:]
					f, err := os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
					if err != nil {
						fmt.Printf("Cannot open file '%v' for storing uri '%v'\n", fn, u)
						break
					}

					if _, err := io.Copy(f, respBody); err != nil {
						fmt.Printf("Cannot write response from uri '%v' to file '%v'\n", u, fn)
					}

					data.Store(imageKey, fn)

					_ = f.Close()
				case "application/json":
					// In our case this is just json
					var jsonData map[string]interface{}
					if err := json.NewDecoder(respBody).Decode(&jsonData); err != nil {
						reject(fmt.Errorf("could not decode data; %w", err))
					}

					var str string
					var ok bool
					for k, v := range jsonData {
						if str, ok = interfaceToString(v); !ok {
							continue
						}

						if k == imageKey {
							set = append(set, endpoint{
								method: http.MethodGet,
								name:   imageKey,
								url:    str,
							})
						} else if prependPropsWithName {
							k = ep.name + "_" + k
						}

						data.Store(k, str)
					}
				default:
					fmt.Printf("unsupported content type '%v' found for url '%v'\n", respType, u)
				}
			}

			data.latestData = time.Now()
			resolve(empty)
		}
	}
}

type (
	endpoint struct {
		method string
		name   string
		url    string
	}

	endpoints    []endpoint
	endpointsSet []endpoints
)

var (
	imageKey             = env.StringFb("IMAGE_KEY", "image")
	props                = xsync.NewMapOf[Props]()
	prependPropsWithName = env.Bool("PREPEND_WEBHOOK_RESULTS_WITH_NAME")
	downloadTo           = env.StringFb("DOWNLOAD_DIR", "/tmp/pixie")

	toCallString = env.StringFb("WEBHOOKS_TO_CALL", "user;GET;https://jury:jury@www.domjudge.org/demoweb/api/user?strict=false|team;GET;https://jury:jury@www.domjudge.org/demoweb/api/teams/{{team_id}}?strict=false")
	toCall       endpointsSet
)

func (e endpointsSet) size() (size int) {
	for _, eps := range e {
		size += len(eps)
	}

	return
}

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
	_, user, pass := decodeBasicAuth(ctx.Request.Header.Peek("Authorization"))
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

func (p *Props) json() io.Reader {
	// TODO create a pool of buffers to use
	b := new(bytes.Buffer)
	b.WriteByte('{')

	var notFirst bool
	p.Range(func(key string, value string) bool {
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
		return true
	})

	b.WriteByte('}')

	return b
}

func interfaceToString(vi interface{}) (strVal string, ok bool) {
	ok = true
	switch v := vi.(type) {
	case string:
		strVal = v
	case int:
		strVal = strconv.Itoa(v)
	case uint8:
		strVal = strconv.Itoa(int(v))
	case int8:
		strVal = strconv.Itoa(int(v))
	case uint16:
		strVal = strconv.Itoa(int(v))
	case int16:
		strVal = strconv.Itoa(int(v))
	case uint32:
		strVal = strconv.Itoa(int(v))
	case int32:
		strVal = strconv.Itoa(int(v))
	case uint64:
		strVal = strconv.FormatUint(v, 10)
	case int64:
		strVal = strconv.FormatInt(v, 10)
	case bool:
		strVal = strconv.FormatBool(v)
	case float64:
		strVal = strconv.FormatFloat(v, 'f', 10, 64)
	case float32:
		strVal = strconv.FormatFloat(float64(v), 'f', 10, 32)
	case []string:
		strVal = strings.Join(v, ", ")

	default:
		ok = false
	}

	return
}
