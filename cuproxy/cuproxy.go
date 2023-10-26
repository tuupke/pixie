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
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"github.com/tuupke/pixie/env"
	"github.com/tuupke/pixie/lifecycle"
)

var (
	cupsListen = "ppp:6631"     // env.StringFb("LISTEN", ":631")
	printerTo  = "10.1.0.1:631" // env.String("PRINTER_TO") // "localhost:631/printers/Virtual_PDF_Printer"

	to = btsReplace([]byte("ipp://" + printerTo))

	dumpsPath        = env.StringFb("DUMP_IPP_CONTENTS", "")
	dumpReplacements = env.Bool("DUMP_REPLACEMENTS")
	dumpOriginal     = env.Bool("DUMP_ORIGINAL")
	appendBanner     = env.Bool("BANNER_APPEND")
	bannerOnBack     = env.Bool("BANNER_ON_BACK")
	seqId            = new(uint64)

	pdfLocation = strings.TrimRight(env.StringFb("PDF_LOCATION", os.TempDir()), "/")
)

func init() {
	// Load the logger
	l, err := zerolog.ParseLevel(env.StringFb("LOG_LEVEL", "info"))
	if err != nil {
		zlog.Fatal().Err(err).Msg("could not parse level")
	}

	zlog.Info().Stringer("level", l).Msg("using loglevel")

	zerolog.SetGlobalLevel(l)
}

func main() {
	if err := os.MkdirAll(pdfLocation, 0755); err != nil {
		zlog.Fatal().Err(err).Str("pdf-folder", pdfLocation).Msg("cannot create pdf-folder")
		panic(err)
	}

	privateRtr := router.New()
	privateRtr.PanicHandler = func(ctx *fasthttp.RequestCtx, i interface{}) {
		zlog.Error().Interface("error", i).Msg("received panic")
	}
	privateRtr.ANY("/{path:*}", cupsHandler)
	if dumpsPath != "" && (dumpReplacements || dumpOriginal) {
		dumpsPath += "/dumps"
		// Recreate the dumps folder by recreating it entirely.
		if err := os.RemoveAll(dumpsPath); err != nil && !os.IsNotExist(err) {
			zlog.Fatal().Err(err).Str("path", dumpsPath).Msg("could not delete dumps folder")
		}
		zlog.Info().Str("path", dumpsPath).Msg("deleted previous dumps folder")

		if err := os.MkdirAll(dumpsPath, 0755); err != nil {
			zlog.Fatal().Str("path", dumpsPath).Err(err).Msg("could not create dumps folder")
		}
		zlog.Info().Msg("recreated dumps folder")
	}

	go func() {
		err := fasthttp.ListenAndServe(cupsListen, privateRtr.Handler)
		zlog.Err(err).Msg("started cups proxy")
		if err != nil {
			zlog.Fatal().Msg("bye")
		}
	}()

	zlog.Info().Str("printer to", printerTo).Str("listen", cupsListen).Msg("Booted")
	lifecycle.Finally(func() { zlog.Warn().Msg("Stopping") })
	lifecycle.StopListener()
}

// Byter is an interface around the `Bytes` method returning a byte slice. Used
// for more consistent logging of request/response contents.
type (
	Byter interface {
		Bytes() []byte
	}
	byteSlice []byte
)

func (b byteSlice) Bytes() []byte {
	return b
}

func writeToFile[T Byter](seqId uint64, request, replaced bool, contents T) error {
	if dumpsPath == "" || (replaced && !dumpReplacements) || (!replaced && !dumpOriginal) {
		return nil
	}

	var req, typ = "res", "orig"
	if request {
		req = "req"
	}

	if replaced {
		typ = "repl"
	}

	name := fmt.Sprintf("%v/%v-%v-%v.bin", dumpsPath, seqId, req, typ)
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return fmt.Errorf("could not open dump-file '%v'; %w", name, err)
	}

	// Ignore the error since it is
	defer errCloserIgnore(f)
	if _, err := f.Write(contents.Bytes()); err != nil {
		return fmt.Errorf("could not write to dump-file '%v'; %w", name, err)
	}

	return nil
}

// errCloserIgnore is a simple wrapper around an io.Closer that
func errCloserIgnore(c io.Closer) {
	// Ignore these errors, consider logging them in the future. These resource leaks
	// should not be too problematic.
	_ = c.Close()
}

// requestPromises maps all cups job-id's to structs that aid with interacting
// with the data-retrieval promises.
var requestPromises = xsync.NewIntegerMapOf[int32, promiseInteraction]()

func cupsHandler(ctx *fasthttp.RequestCtx) {
	seqId := atomic.AddUint64(seqId, 1)
	body := slices.Clone(ctx.Request.Body())

	// Extract the type of operation. Consider extracting operation into a custom type for improved readability.
	var operationId byte
	if len(body) > 3 {
		operationId = body[3]
	}
	var jobId int32
	isCreate := operationId == 0x05
	isPrint := operationId == 0x02 || operationId == 0x06

	// Construct a logger
	path := bytes.Trim(ctx.Request.URI().Path(), "/")
	requestedUrl := fmt.Sprintf("ipp://%s/%s", cupsListen, path)
	log := zlog.With().IPAddr("ip", ctx.RemoteIP()).Str("url", requestedUrl).Uint64("seq-id", seqId).Logger()

	// Replace the url in the body, the actual printer's url is constant so does not
	// need to be rebuilt every request.
	from := btsReplace([]byte(requestedUrl))
	log.Debug().Err(writeToFile[byteSlice](seqId, true, false, body)).Msg("written original request")
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

		var v promiseInteraction
		log.Info().Msg("print triggered")
		if !found {
			// This is wierd
			log.Error().Msg("got print and cannot extract job-id, very wierd.")
			v = loadValues(log, ctx, jobId)
		} else {
			v, found = requestPromises.Load(jobId)
			if !found {
				log.Warn().Msg("did not find promise for job, i.e. job is unknown to proxy; trying to load new data")
				v = loadValues(log, ctx, jobId)
			}
		}

		v.callItIn()
		filePointer, err := v.pdfPromise.Await(lifecycle.ApplicationContext())
		log.Err(err).Msg("got PDF to prepend")

		if filePointer != nil {
			file := *filePointer
			defer errCloserIgnore(file)

			tempReader := bytes.NewBuffer(make([]byte, 0, len(body)))
			pdf := []io.ReadSeeker{file, pdfReader}
			if appendBanner {
				// Flip the pdfs to stitch
				pdf[0], pdf[1] = pdf[1], pdf[0]
			}

			err := pdfcpu.MergeRaw([]io.ReadSeeker{file, pdfReader}, tempReader, nil)
			log.Err(err).Msg("merged banner with main print")
			if err == nil {
				num, err := io.Copy(newB, tempReader)
				log.Debug().Err(err).Int64("num", num).Msg("copied buffer")
			} else {
				newB.Write(body[pdfStart:pdfEnd])
			}
		}

		num, err := newB.Write(body[pdfEnd:])
		log.Err(err).Int("num", num).Msg("written rest of request to new body")
		b = newB
	}

	log.Debug().Err(writeToFile(seqId, true, true, b)).Msg("written replaced request")

	// Construct the proxy request, since `b` is always a *bytes.Buffer and these get closed automatically
	proxiedRequest, err := http.NewRequest(string(ctx.Method()), "http://"+printerTo, b)
	log.Err(err).Msg("created request to proxy")
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		proxiedRequest.Header.Add(string(key), strings.Replace(string(value), requestedUrl[5:], printerTo, -1))
	})

	resp, err := http.DefaultClient.Do(proxiedRequest)
	log.Debug().Err(err).Msg("proxied request")
	err = proxiedRequest.Body.Close()
	log.Debug().Err(err).Msg("closed request body")
	if err != nil {
		panic(err)
	}

	body, err = io.ReadAll(resp.Body)
	log.Debug().Err(err).Msg("read entire response body")
	err = resp.Body.Close()
	log.Debug().Err(err).Msg("closed response body")
	body = setReplace(body)

	ctx.SetStatusCode(resp.StatusCode)
	for k, values := range resp.Header {
		for _, value := range values {
			repl := strings.Replace(value, printerTo, requestedUrl[5:], -1)
			ctx.Response.Header.Add(k, repl)
		}
	}

	log.Debug().Err(writeToFile[byteSlice](seqId, false, false, body)).Msg("written original response")
	body = bytes.Replace(body, to, from, -1)
	log.Debug().Msg("replaced response body")
	log.Debug().Err(writeToFile[byteSlice](seqId, false, true, body)).Msg("written replaced response")

	if isCreate {
		var found bool
		jobId, found = extractInt("job-id", body)
		if !found {
			// Weirdness happens here
			log.Warn().Msg("cannot deduce job-id though it should be present!")
		} else {
			var pp promiseInteraction
			if isCreate {
				pp = loadValues(log, ctx, jobId)
			}

			// Store as job
			requestPromises.Store(jobId, pp)
		}
	}

	// Write the body
	if _, err := ctx.Write(body); err != nil {
		panic(err)
	}

	log.Debug().Msg("written proxied-body")
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

	// Re-slice the byte slice to the correct length.
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
