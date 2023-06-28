package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fasthttp/router"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"github.com/tuupke/pixie/cuproxy/props"
	"github.com/tuupke/pixie/lifecycle"
)

var (
	printerTo  = "localhost:631/printers/Virtual_PDF_Printer"
	cupsListen = "ppp:6631"

	to = btsReplace([]byte("ipp://" + printerTo))
)

func main() {
	privateRtr := router.New()
	privateRtr.ANY("/{path:*}", cupsHandler)

	go func() {
		err := fasthttp.ListenAndServe(cupsListen, privateRtr.Handler)
		log.Err(err).Msg("started cups proxy")
		if err != nil {
			log.Fatal().Msg("bye")
		}
	}()

	log.Info().Msg("Booted")
	lifecycle.Finally(func() { log.Warn().Msg("Stopping") })
	lifecycle.StopListener()
}

func writeToFile(name string, contents []byte) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	f.Write(contents)
}

func cupsHandler(ctx *fasthttp.RequestCtx) {
	// write the current request
	body := ctx.Request.Body()

	path := bytes.Trim(ctx.Request.URI().Path(), "/")
	requestedUrl := fmt.Sprintf("ipp://%s/%s", cupsListen, []byte(path))

	from := btsReplace([]byte(requestedUrl))

	// Replace the url
	body = bytes.Replace(body, from, to, -1)

	var until = 100
	if len(body) < until {
		until = len(body)
	}

	// If we detect "Create-Job" start retrieving the info
	if bytes.Index(body[:until], []byte("Create-Job")) >= 0 {
		props.LoadFromRequest(ctx).Refresh()
	} else if bytes.Index(body[:until], []byte("Print-Job")) >= 0 {
		// If we detect "Print-Job", ensure the is data is retrieved
		props.LoadFromRequest(ctx).CallItIn(false)
	}

	// In normal cases, simply proxy the request, only when a document is sent some magic is needed.
	var b io.Reader = bytes.NewBuffer(body)
	hreq, err := http.NewRequest(string(ctx.Method()), "http://"+printerTo, b)

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		hreq.Header.Add(string(key), strings.Replace(string(value), requestedUrl, printerTo, -1))
	})

	resp, err := http.DefaultClient.Do(hreq)
	_ = hreq.Body.Close()
	if err != nil {
		panic(err)
	}

	body, err = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	// writeToFile(pref+"-resp-orig.bin", body)
	body = setReplace(body)
	// writeToFile(pref+"-resp-repl.bin", body)

	var rawHeaders = make(http.Header)
	ctx.SetStatusCode(resp.StatusCode)
	// writer.WriteHeader(resp.StatusCode)
	for k, values := range resp.Header {
		for _, value := range values {
			repl := strings.Replace(value, printerTo, requestedUrl, -1)
			rawHeaders.Add(k, repl)
			ctx.Response.Header.Add(k, repl)
		}
	}

	body = bytes.Replace(body, to, from, -1)

	// Write the body
	if _, err := ctx.Write(body); err != nil {
		panic(err)
	}
}

// btsReplace takes a byte string and converts it to the representation used in CUPS
func btsReplace(check []byte) []byte {
	checkLen := len(check)
	b := make([]byte, 2+checkLen)
	copy(b[2:], check)

	// Prepend big-endian length, it's encoded in two bytes
	binary.BigEndian.PutUint16(b, uint16(checkLen))
	return b
}

// setReplace replaces all relevant sets/collections in the response that are
// needed to convince cups that only PDF is supported. All properties starting
// with "document-format-" need to be replaced. To simplify even further, all
// properties that depict mime-types are replaced.
func setReplace(body []byte) []byte {
	// The format used by CUPS is a bit annoying to efficiently replace. For a single value:
	//   <type, 1B><length, 2B, big-endian><key-name, `length`B><value length, 2B,
	//   big-endian><value, `value length`B>
	// subsequent values are than 'appended' as follows:
	//   <type, 1B>\u0000\u0000<value length, 2B, big-endian><value, `value length`B>
	// The two nulls represent the 'empty key' and thus the value should be interpreted as being part of last key.

	// All keys storing mime-types should be set to 'application/pdf'. The type depicting a mime-type is "I"
	var result = make([]byte, len(body))
	// We know all our properties start with 0x49 0x00
	prefix := []byte("I\u0000")

	var matchedUntil, appendedUntil int
	pdfMime := "\u0000\u000Fapplication/pdf"

	var bddy []byte
	for len(body) > matchedUntil {
		next := bytes.Index(body[matchedUntil:], prefix)
		if next < 0 {
			break
		}

		// Copy everything up to the match and move to the match
		copy(result[appendedUntil:], body[matchedUntil:matchedUntil+next])
		appendedUntil += next
		matchedUntil += next

		// A property is found, the next four bytes must be the prefix and big-endian length. Then
		keyLength := binary.BigEndian.Uint16(body[matchedUntil+1 : matchedUntil+3])

		// copy the type (1 byte), the length (2 bytes), and the key (length bytes) to the result, and advance
		toCopy := 3 + int(keyLength)
		copy(result[appendedUntil:], body[matchedUntil:matchedUntil+toCopy])
		appendedUntil += toCopy
		matchedUntil += toCopy

		// A key is copied, for which the value always must be application/pdf.
		copy(result[appendedUntil:], pdfMime)
		appendedUntil += len(pdfMime)

		// The first value always exists, read the length and skip
		matchedUntil += 2 + int(binary.BigEndian.Uint16(body[matchedUntil:matchedUntil+2]))

		bddy = body[matchedUntil:]
		// Another value starts with the type ("I") and a continuation character ("\u0000\u0000"). Keep reading until the next few characters are not continuation.
		for len(body) >= matchedUntil+5 && bytes.Equal(body[matchedUntil:matchedUntil+3], []byte("I\u0000\u0000")) {
			matchedUntil += 5 + int(binary.BigEndian.Uint16(body[matchedUntil+3:matchedUntil+5]))
			bddy = body[matchedUntil:]
		}
	}

	bddy = bddy

	// Copy the remaining bytes
	copy(result[appendedUntil:], body[matchedUntil:])

	// Reslice, result is at most `len(body)` bytes long
	return result[:appendedUntil+len(body)-matchedUntil]
}
