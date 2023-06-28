// Package props implements a property cache
package props

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/puzpuzpuz/xsync"
	"github.com/valyala/fasthttp"

	"github.com/valyala/fasttemplate"

	"github.com/tuupke/pixie/cuproxy/web"
	"github.com/tuupke/pixie/env"
	"github.com/tuupke/pixie/lifecycle"
)

type (
	Props struct {
		initialDone, loadingDone *sync.Once
		initial, loading         *sync.WaitGroup
		ip                       net.IP

		data *xsync.MapOf[string, string]
	}

	propPair struct {
		p    *Props
		urls endpoints
	}

	endpoint struct {
		method string
		name   string
		url    string
	}

	endpoints    []endpoint
	endpointsSet []endpoints
)

var (
	requestWorkers          = env.IntFb("REQUEST_WORKER_NUM", 15)
	imageKey                = env.StringFb("IMAGE_KEY", "image")
	workerChan, workerGroup = lifecycle.FinallyWorker[propPair]()
	props                   = xsync.NewMapOf[Props]()
	prependPropsWithName    = env.Bool("PREPEND_WEBHOOK_RESULTS_WITH_NAME")
	downloadTo              = env.StringFb("DOWNLOAD_DIR", "/tmp/pixie")

	toCallString = env.String("WEBHOOKS_TO_CALL")
	toCall       endpointsSet
)

// Variables used for checking validity
var (
	empty   = struct{}{}
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
	workerGroup.Add(requestWorkers)
	for i := 0; i < requestWorkers; i++ {
		go worker(workerChan)
	}

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

func worker(requests chan propPair) {
	defer workerGroup.Done()
	fmt.Println("started worker")
	for req := range requests {
		if len(req.urls) == 0 {
			req.p.loading.Done()
			continue
		}

		// Take the first url, issue request.
		u, method := req.urls[0].url, req.urls[0].method

		t := fasttemplate.New(u, "{{", "}}")
		u = t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
			v, _ := req.p.data.Load(tag)
			return w.Write([]byte(v))
		})

		var reqBody io.Reader
		if method != http.MethodGet {
			reqBody = req.p.json()
		}

		fmt.Println("Executing request request", u, method)
		ctx, cancel := context.WithTimeout(lifecycle.ApplicationContext(), time.Second*2)
		respBody, respType, loaded, err := web.Do(ctx, u, method, req.p.ip, reqBody)
		cancel()

		if err != nil {
			// TODO log
			fmt.Println("Error encountered", err)
			req.p.loading.Done()
			continue
		}

		if loaded {
			switch respType = strings.Split(respType, ";")[0]; respType {
			case "image/jpeg":
				fallthrough
			case "image/png":
				fallthrough
			case "image/gif":
				// Store in file and
				fn := downloadTo + "/" + req.p.ip.String() + "." + respType[6:]
				f, err := os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
				if err != nil {
					fmt.Printf("Cannot open file '%v' for storing uri '%v'\n", fn, u)
					break
				}

				if _, err := io.Copy(f, respBody); err != nil {
					fmt.Printf("Cannot write response from uri '%v' to file '%v'\n", u, fn)
				}

				req.p.data.Store(imageKey, fn)

				_ = f.Close()
			case "application/json":
				// In our case this is just json
				var data map[string]interface{}
				if err := json.NewDecoder(respBody).Decode(&data); err != nil {
					// TODO handle the error
					log.Printf("Found error decoding url '%v': '%v'", u, err)
					break
				}

				var str string
				var ok bool
				for k, v := range data {
					if str, ok = interfaceToString(v); !ok {
						continue
					}

					if k == imageKey {
						req.urls = append(req.urls, endpoint{
							method: http.MethodGet,
							name:   imageKey,
							url:    str,
						})
					} else if prependPropsWithName {
						k = req.urls[0].name + "_" + k
					}

					req.p.data.Store(k, str)
				}
			default:
				fmt.Printf("unsupported content type '%v' found for url '%v'\n", respType, u)
			}
		}

		// Handle the next url next urls
		req.urls = req.urls[1:]
		requests <- req

		if respBody != nil {
			_ = respBody.Close()
		}
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
	props, _ := (*xsync.MapOf[string, *Props])(props).LoadOrCompute(key, func() (p *Props) {
		p = &Props{
			initialDone: new(sync.Once),
			loadingDone: new(sync.Once),
			initial:     new(sync.WaitGroup),
			loading:     new(sync.WaitGroup),
			data:        xsync.NewMapOf[string](),
			ip:          ip,
		}

		for k, v := range baseData {
			p.data.Store(k, v)
		}

		p.initial.Add(1)
		return
	})

	return props
}

func (p *Props) Refresh() {
	p.refresh(toCall)
}

func (p *Props) refresh(urlSet endpointsSet) {
	p.loadingDone.Do(func() {
		defer func() {
			p.loadingDone = new(sync.Once)
		}()

		for _, set := range urlSet {
			p.loading.Add(1)
			workerChan <- propPair{
				p:    p,
				urls: set,
			}
		}
	})
}

func (p *Props) CallItIn(wait bool) {
	p.initialDone.Do(func() {
		p.loading.Wait()
		p.initial.Done()
	})
	p.initial.Wait()
	if wait {
		p.loading.Done()
	}
}

func (p *Props) Reduce(keys ...string) (mp map[string]string) {
	mp = make(map[string]string)
	for _, key := range keys {
		if data, ok := p.data.Load(key); ok {
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
	p.data.Range(func(key string, value string) bool {
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
