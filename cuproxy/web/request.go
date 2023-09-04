// Package web implements a web retrieval and caching mechanism
package web

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/puzpuzpuz/xsync"
	"github.com/rs/zerolog"

	"github.com/tuupke/pixie/env"
)

type (
	etagPair struct {
		Key, Date string
	}
)

var (
	// etagCache is an etag cache. Stores the retrieved etag for some url
	etagCache  = xsync.NewMapOf[etagPair]()
	pixieNonce = func() string {
		nonce := env.String("PIXIE_HTTP_REQUEST_NONCE")
		if nonce != "" {
			return nonce
		}

		const pixieNonceAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
		src := rand.New(rand.NewSource(time.Now().UnixNano()))

		b := make([]byte, 32)
		for i := range b {
			b[i] = pixieNonceAlphabet[src.Intn(len(pixieNonceAlphabet))]
		}

		return string(b)
	}()
)

func Do(ctx context.Context, log zerolog.Logger, url, verb string, ip net.IP, requestBody io.Reader) (responseBody io.ReadCloser, responseType string, loaded bool, err error) {
	req, err := http.NewRequestWithContext(ctx, verb, url, requestBody)
	if err != nil {
		err = fmt.Errorf("cannot create request [%v] '%v'", verb, url)
		return
	}

	log = log.With().Str("verb", verb).Str("url", url).Bool("with_ip", ip != nil).Logger()

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
