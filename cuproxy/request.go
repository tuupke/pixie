// Package web implements a web retrieval and caching mechanism
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/puzpuzpuz/xsync"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/tuupke/pixie/env"
	"github.com/tuupke/pixie/lifecycle"
)

type (
	etagPair struct {
		Key, Date string
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
	// etagCache is an etag cache. Stores the retrieved etag for some url
	etagCache  = xsync.NewMapOf[etagPair]()
	pixieNonce = func() string {
		nonce := env.String("WEBHOOK_REQUEST_NONCE")
		if nonce == "" {

			const pixieNonceAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
			src := rand.New(rand.NewSource(time.Now().UnixNano()))

			b := make([]byte, 32)
			for i := range b {
				b[i] = pixieNonceAlphabet[src.Intn(len(pixieNonceAlphabet))]
			}

			nonce = string(b)
		}

		zlog.Info().Str("nonce", nonce).Msg("X-Pixie-Nonce value")
		return nonce
	}()

	maxWebhookTime = env.DurationFb("WEBHOOK_MAX_DURATION", time.Second*30)
)

func Do(ctx context.Context, log zerolog.Logger, url, verb string, ip net.IP, requestBody io.Reader) (responseBody io.ReadCloser, responseType string, loaded bool, err error) {
	req, err := http.NewRequestWithContext(ctx, verb, url, requestBody)
	if err != nil {
		err = fmt.Errorf("cannot create request [%v] '%v'", verb, url)
		return
	}

	req.Header.Set("User-Agent", "Pixie/CupsProxy")
	req.Header.Set("Accept", "application/json, image/*")
	req.Header.Set("X-Pixie-Nonce", pixieNonce)
	if ip != nil {
		req.Header.Set("X-Forwarded-For", ip.String())
	}

	// Check for, and add, cached etag values
	etag, etagLoaded := etagCache.Load(url)
	if etagLoaded {
		req.Header.Add("If-None-Match", etag.Key)
		req.Header.Add("If-Modified-Since", etag.Date)
	}

	resp, err := http.DefaultClient.Do(req)
	var statusCode int
	if resp != nil {
		statusCode = resp.StatusCode
	}

	log.Err(err).Int("status", statusCode).Msg("called hook")
	if err != nil {
		err = fmt.Errorf("cannot create request [%v] '%v'", verb, url)
		return
	}

	if resp.StatusCode == http.StatusNotModified {
		// Nothing to do, not loaded, nor an error is thrown. Close body just in case and return
		err = resp.Body.Close()
		log.Err(err).Msg("status not changed, continuing")
		return
	}

	responseBody, responseType, loaded = resp.Body, resp.Header.Get("content-type"), true

	if key, date := resp.Header.Get("ETag"), resp.Header.Get("Date"); key != "" || date != "" {
		etagCache.Store(url, etagPair{Key: key, Date: date})
	}

	if resp.StatusCode/100 != 2 {
		// Only continue on success
		log.Warn().Msg("status not succesfull, not continuing to next hook")
		err = fmt.Errorf("non-successfull status code received (%v)", resp.StatusCode)
	}

	return
}

func (ep endpoint) executeRequest(ctx context.Context, log zerolog.Logger, data *Props) (io.ReadCloser, string, bool, error) {
	u, params := replaceParameters(ep.url, data, ep.name)
	log.Debug().Str("original", ep.url).Object("relevant-data", params).Str("result", u).Msg("replaced url")

	var reqBody io.Reader
	if ep.method != http.MethodGet {
		reqBody = data.json(map[string]string{
			"webhook_name":   ep.name,
			"webhook_method": ep.method,
			"webhook_url":    u,
		})
	}

	ctx, cancel := context.WithTimeout(ctx, maxWebhookTime)
	defer cancel()

	return Do(ctx, log, u, ep.method, data.ip, reqBody)
}

func (ep endpoint) handleResponse(log zerolog.Logger, respBody io.ReadCloser, respType string, data *Props) (extra endpoints) {
	extra = make(endpoints, 0, 3)
	respType = strings.Split(respType, ";")[0]
	log.Info().Str("response_type", respType).Msg("handling response")
	switch respType {
	case "image/jpeg", "image/png", "image/gif":
		// Store in file and
		fn := downloadTo + "/" + data.ip.String() + "." + respType[6:]
		f, err := os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
		log.Debug().Err(err).Str("filename", fn).Msg("storing image")
		if err != nil {
			log.Warn().Err(err).Str("filename", fn).Msg("could not open imagefile")
			break
		}

		numbts, err := io.Copy(f, respBody)
		log.Debug().Err(err).Int64("num-bytes", numbts).Str("filename", fn).Msg("written image")
		if err != nil {
			log.Warn().Str("filename", fn).Msg("image writing failed, not storing the path")
			break
		}

		data.Store(imageKey, fn)

		_ = f.Close()
	case "application/json":
		var jsonData map[string]interface{}
		err := json.NewDecoder(respBody).Decode(&jsonData)
		if err != nil {
			log.Err(err).Msg("could not decode json")
			break
		}

		var str string
		var ok bool
		for k, v := range jsonData {
			if str, ok = interfaceToString(v); !ok {
				continue
			}

			// TODO fix the imageKey handling, this is currently incorrect!
			if k == imageKey {
				// Check for valid url
				_, err := url.Parse(str)
				log.Debug().Err(err).Str("url", str).Str("key", k).Msg("url found to call, will call if valid")

				if err == nil {
					extra = append(extra, endpoint{
						method: http.MethodGet,
						name:   imageKey,
						url:    str,
					})
				}
			} else if keyTemplate != "" {
				data.Store("webhook_key", k)
				// Replace they key with the contents of the template
				orig := k
				var params mapWriter
				k, params = replaceParameters(keyTemplate, data, ep.name)
				log.Debug().Str("original", orig).Str("template", keyTemplate).Object("relevant-data", params).Str("result", k).Msg("filled template for key")
			}

			data.Store(k, str)
		}
	default:
		log.Warn().Str("content-type", respType).Msg("unsupported content type encountered, ignored")
	}

	return
}

func (ep endpoints) handle(c chan e, log zerolog.Logger, data *Props) func() {
	return func() {
		defer func(c chan e) { c <- empty }(c)
		eps := slices.Clone(ep)

		// Any endpoint called can have
		for len(eps) > 0 {
			newEps := make(endpoints, 0, len(ep))

			for _, end := range eps {
				log = log.With().Str("verb", end.method).Str("url", end.url).Bool("with-ip", data.ip != nil).Logger()

				respBody, respType, loaded, err := end.executeRequest(lifecycle.ApplicationContext(), log, data)
				log.Err(err).Bool("new data", loaded).Msg("request executed")
				if err != nil {
					return
				}

				// Nothing to do/update, continue to next in set
				if !loaded {
					log.Debug().Msg("no new data loaded, skipping response handling")
					continue
				}

				data.latestData = time.Now()
				newEps = append(newEps, end.handleResponse(log, respBody, respType, data)...)
			}

			eps = newEps
		}

	}
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
