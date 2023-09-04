package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/fasthttp/router"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/puzpuzpuz/xsync"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"github.com/tuupke/pixie/env"
	"github.com/tuupke/pixie/lifecycle"
)

var (
	printerTo  = "localhost:631/printers/Virtual_PDF_Printer"
	cupsListen = "127.0.0.1:6631"

	to    = btsReplace([]byte("ipp://" + printerTo))
	pref  = env.StringFb("DUMP_IPP", "")
	seqId = new(uint64)

	dumpReplacements = env.Bool("DUMP_REPLACEMENTS")
	dumpOriginal     = env.Bool("DUMP_ORIGINAL")

	pdfLocation = env.StringFb("PDF_LOCATION", os.TempDir())
)

func main() {
	if err := os.MkdirAll(pdfLocation, 0755); err != nil {
		panic(err)
	}

	privateRtr := router.New()
	privateRtr.PanicHandler = func(ctx *fasthttp.RequestCtx, i interface{}) {
		log.Error().Interface("error", i).Msg("received panic")
	}
	privateRtr.ANY("/{path:*}", cupsHandler)
	if pref != "" && (dumpReplacements || dumpOriginal) {
		pref += "/dumps"
		os.RemoveAll(pref)
		os.MkdirAll(pref, 0755)
	}

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

func writeToFile(seqId uint64, request, replaced bool, contents []byte) {
	if pref == "" || (replaced && !dumpReplacements) || (!replaced && !dumpOriginal) {
		return
	}

	var req, typ = "res", "orig"
	if request {
		req = "req"
	}

	if replaced {
		typ = "repl"
	}

	name := fmt.Sprintf("%v/%v-%v-%v.bin", pref, seqId, req, typ)
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	if _, err := f.Write(contents); err != nil {
		log.Err(err).Str("filename", name).Msg("could not write file")
	}
}

var requestPromise = xsync.NewIntegerMapOf[int32, promisePair]()

func cupsHandler(ctx *fasthttp.RequestCtx) {
	seqId := atomic.AddUint64(seqId, 1)

	// write the current request
	body := ctx.Request.Body()

	var operationId byte
	if len(body) >= 3 {
		operationId = body[3]
	}
	var jobId int32
	isCreate := operationId == 0x005
	isPrint := operationId == 0x02 || operationId == 0x06

	path := bytes.Trim(ctx.Request.URI().Path(), "/")
	requestedUrl := fmt.Sprintf("ipp://%s/%s", cupsListen, []byte(path))

	from := btsReplace([]byte(requestedUrl))

	// Replace the url
	writeToFile(seqId, true, false, body)
	body = bytes.Replace(body, from, to, -1)

	var b *bytes.Buffer
	if !isPrint {
		// In normal cases, simply proxy the request, only when a document is sent some magic is needed.
		b = bytes.NewBuffer(body)
	} else {
		// The not-so normal case
		var found bool
		jobId, found = extractInt("job-id", body)
		log := log.With().Int32("job-id", jobId).Bool("job-id-found", found).Logger()

		newB := bytes.NewBuffer(make([]byte, 0, 2048+len(body)))

		// PDF start with "%PDF" and end with "%%EOF"
		pdfStart := bytes.Index(body, []byte("%PDF"))
		pdfEnd := pdfStart + bytes.Index(body[pdfStart:], []byte("%%EOF")) + 5

		// The preamble must be kept, copy it over
		newB.Write(body[:pdfStart])

		// Store the pdf in a new reader, so the contents of `body` remain valid!
		pdfReader := bytes.NewReader(slices.Clone(body[pdfStart:pdfEnd]))
		// if body[pdfEnd] == 0x0A {
		// 	pdfEnd++
		// }

		var v promisePair
		log.Info().Msg("print triggered")
		if !found {
			// This is wierd
			log.Error().Msg("got print and cannot extract job-id, very wierd.")
			v = loadValues(ctx, jobId)
		} else {
			v, found = requestPromise.Load(jobId)
			if !found {
				log.Warn().Msg("did not find promise for job, i.e. job is unknown to proxy; trying to load new data")
				v = loadValues(ctx, jobId)
			}
		}

		v.callItIn()
		ifile, err := v.pdfPromise.Await(lifecycle.ApplicationContext())
		log.Err(err).Msg("got PDF to prepend")

		if ifile != nil {
			file := *ifile
			defer file.Close()

			tempReader := bytes.NewBuffer(make([]byte, 0, len(body)))

			err := pdfcpu.MergeRaw([]io.ReadSeeker{file, pdfReader}, tempReader, nil)
			log.Err(err).Msg("merged banner")
			if err == nil {
				io.Copy(newB, tempReader)
			} else {
				newB.Write(body[pdfStart:pdfEnd])
			}
		}

		newB.Write(body[pdfEnd:])
		b = newB
	}

	writeToFile(seqId, true, true, b.Bytes())

	hreq, _ := http.NewRequest(string(ctx.Method()), "http://"+printerTo, b)

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		hreq.Header.Add(string(key), strings.Replace(string(value), requestedUrl, printerTo, -1))
	})

	resp, err := http.DefaultClient.Do(hreq)
	_ = hreq.Body.Close()
	if err != nil {
		panic(err)
	}

	body, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	body = setReplace(body)

	var rawHeaders = make(http.Header)
	ctx.SetStatusCode(resp.StatusCode)
	for k, values := range resp.Header {
		for _, value := range values {
			repl := strings.Replace(value, printerTo, requestedUrl, -1)
			rawHeaders.Add(k, repl)
			ctx.Response.Header.Add(k, repl)
		}
	}

	writeToFile(seqId, false, false, body)
	body = bytes.Replace(body, to, from, -1)
	writeToFile(seqId, false, true, body)

	if isCreate {
		var found bool
		jobId, found = extractInt("job-id", body)
		if !found {
			// Weirdness
			log.Error().Msg("cannot deduce job-id where it should be present!")
		} else {
			var pp promisePair
			if isCreate {
				pp = loadValues(ctx, jobId)
			}

			// Store as job
			requestPromise.Store(jobId, pp)
		}
	}

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

		// Another value starts with the type ("I") and a continuation character ("\u0000\u0000"). Keep reading until the next few characters are not continuation.
		for len(body) >= matchedUntil+5 && bytes.Equal(body[matchedUntil:matchedUntil+3], []byte("I\u0000\u0000")) {
			matchedUntil += 5 + int(binary.BigEndian.Uint16(body[matchedUntil+3:matchedUntil+5]))
		}
	}

	// Copy the remaining bytes
	copy(result[appendedUntil:], body[matchedUntil:])

	// Reslice, result is at most `len(body)` bytes long
	return result[:appendedUntil+len(body)-matchedUntil]
}

func extractInt(name string, body []byte) (value int32, found bool) {
	bts := make([]byte, 3, len(name)+3)
	bts[0] = '!'
	binary.BigEndian.PutUint16(bts[1:], uint16(len(name)))
	bts = append(bts, []byte(name)...)

	index := bytes.Index(body, bts)
	if index < 0 {
		return
	}

	found = true

	// Read the next byte, this is the length
	start := index + len(bts)
	length := int(binary.BigEndian.Uint16(body[start:]))
	// Length should be smaller than 4 bytes
	if length > 4 || length <= 0 {
		// Give warning
		panic("oops")
	}

	value = int32(binary.BigEndian.Uint32(body[start+2:]))

	return
}
