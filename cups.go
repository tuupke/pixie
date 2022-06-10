package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenPrinting/goipp"
	"github.com/jung-kurt/gofpdf"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/rs/zerolog/log"
)

type (
	readSeekWriter interface {
		io.ReadSeeker
		io.Writer
		io.Closer
	}
)

var num = new(uint32)

func simpleProxy(writer http.ResponseWriter, request *http.Request) {
	n := atomic.AddUint32(num, 1)
	requestedUrl := fmt.Sprintf("ipp://%s/%s", strings.TrimRight(request.Host, "/"), strings.Trim(request.URL.Path, "/"))
	from := btsReplace([]byte(requestedUrl))

	f, err := os.OpenFile(fmt.Sprintf("/tmp/pixie-test/req-%v-%v.txt", time.Now().Format(time.RFC3339), n), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Take the body and replace to the correct url
	bts, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(f, `Pre
requestUrl: %v
newReqUrl: %v
Headers: %v
Body pre:
`, request.Host, requestedUrl, request.Header)
	f.Write(bts)
	f.WriteString("\n--- END ---\nBody post:\n")
	bts = bytes.Replace(bts, from, to, -1)

	f.Write(bts)
	f.WriteString("\n--- END ---\n")

	// Construct new request
	hreq, err := http.NewRequest(request.Method, "http://"+printerTo, bytes.NewReader(bts))
	for k, values := range request.Header {
		for _, value := range values {
			hreq.Header.Add(k, strings.Replace(value, requestedUrl, printerTo, -1))
		}
	}

	hreq.Header.Set("Content-Length", strconv.Itoa(len(bts)))

	fmt.Fprintf(f, "Headers POST REPLACE: %v\n\n-- RESPONSE --", hreq.Header)

	resp, err := http.DefaultClient.Do(hreq)
	_ = hreq.Body.Close()
	if err != nil {
		panic(err)
	}

	bts, err = io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	fmt.Fprintf(f, `Pre
Headers: %v
Body pre:
`, resp.Header)
	f.Write(bts)
	f.WriteString("\n--- END ---\nBody post:\n")
	bts = bytes.Replace(bts, from, to, -1)

	f.Write(bts)
	f.WriteString("\n--- END ---\n")

	writer.WriteHeader(resp.StatusCode)
	for k, values := range resp.Header {
		for _, value := range values {
			if k == "Content-Length" {
				writer.Header().Set("Content-Length", strconv.Itoa(len(bts)))
			} else {
				writer.Header().Add(k, strings.Replace(value, printerTo, requestedUrl, -1))
			}
		}
	}
	fmt.Fprintf(f, "Headers POST REPLACE: %v\n\n-- RESPONSE --", writer.Header())

	if _, err := writer.Write(bts); err != nil {
		panic(err)
	}

}

// cupsHandler is the handler which proxies all cups requests and fixes
func cupsHandler(writer http.ResponseWriter, request *http.Request) {
	requestedUrl := fmt.Sprintf("ipp://%s/%s", strings.TrimRight(request.Host, "/"), strings.Trim(request.URL.Path, "/"))
	from := btsReplace([]byte(requestedUrl))

	bts, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	os.WriteFile(fmt.Sprintf("/tmp/pixie-test/req-%v", num), bts, 0755)

	// Replace the url
	bts = bytes.Replace(bts, from, to, -1)

	// Read the message
	var req goipp.Message
	if err := req.DecodeBytes(bts); err != nil {
		panic(err)
	}

	// In normal cases, simply proxy the request, only when a document is sent some magic is needed.
	var b io.Reader = bytes.NewBuffer(bts)
	// if goipp.Op(req.Code) == goipp.OpSendDocument && false {
	// 	log.Debug().Msg("handling print")
	//
	// 	// Search for the PDF header
	// 	sep := bytes.Index(bts, []byte("%PDF"))
	// 	if sep < 0 {
	// 		panic("Not found!")
	// 	}
	//
	// 	var rw *io.PipeWriter
	// 	b, rw = io.Pipe()
	//
	// 	go func() {
	// 		if _, err := rw.Write(bts[:sep]); err != nil {
	// 			panic(err)
	// 		}
	//
	// 		hash := sha1.New()
	// 		hash.Write([]byte(request.URL.Path))
	// 		pdfLocation := fmt.Sprintf("%v/%s.pdf", cacheFolder, base32.StdEncoding.EncodeToString(hash.Sum(nil)))
	//
	// 		stat, stErr := os.Stat(pdfLocation)
	// 		var skipBuild = stErr == nil && stat != nil && time.Now().Before(stat.ModTime().Add(timeoutPdf))
	//
	// 		// Test if the page exists and is recent enough
	// 		var oldPage readSeekWriter
	// 		oldPage, err = os.OpenFile(pdfLocation, os.O_RDWR|os.O_CREATE, 0755)
	// 		if err != nil {
	// 			skipBuild = false
	// 		}
	//
	// 		// If the cached page is not there, create a buffer which can be used to store the file in
	// 		if oldPage == nil {
	// 			oldPage, err = memFs.OpenFile(pdfLocation, os.O_RDWR|os.O_CREATE, 0755)
	// 			if err != nil {
	// 				// TODO log
	// 				panic(err)
	// 			}
	//
	// 			defer memFs.Remove(pdfLocation)
	// 		} else {
	// 			oldPage.Seek(0, 0)
	// 		}
	//
	// 		pdfParts := []io.ReadSeeker{oldPage, bytes.NewReader(bts[sep:])}
	// 		defer oldPage.Close()
	// 		if !skipBuild {
	// 			order, mp := requestToMap(request)
	// 			if err := renderPage(oldPage, order, mp); err != nil {
	// 				// TODO log
	// 				pdfParts = pdfParts[1:]
	// 			} else if _, err := oldPage.Seek(0, 0); err != nil {
	// 				// TODO log
	// 				pdfParts = pdfParts[1:]
	// 			}
	// 		}
	//
	// 		if err := pdfcpu.Merge(pdfParts, rw, nil); err != nil {
	// 			panic(err)
	// 		}
	//
	// 		_ = rw.Close()
	// 	}()
	// }

	hreq, err := http.NewRequest(request.Method, "http://"+printerTo, b)
	for k, values := range request.Header {
		for _, value := range values {
			hreq.Header.Add(k, strings.Replace(value, requestedUrl, printerTo, -1))
		}
	}

	resp, err := http.DefaultClient.Do(hreq)
	_ = hreq.Body.Close()
	if err != nil {
		panic(err)
	}

	bts, err = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	bts = bytes.Replace(bts, to, from, -1)

	writer.WriteHeader(resp.StatusCode)
	for k, values := range resp.Header {
		for _, value := range values {
			writer.Header().Add(k, strings.Replace(value, printerTo, requestedUrl, -1))
		}
	}

	if _, err := writer.Write(bts); err != nil {
		panic(err)
	}
}

// renderPage renders a single pdf, which will be prefixed to the actual page, i.e. a banner page
func renderPage(outWrite io.Writer, keys keyset, values values) error {
	orientation := "P"
	if pdfInLandscape {
		orientation = "L"
	}

	pdf := gofpdf.New(orientation, pdfUnit, pdfSize, pdfFontDir)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	yTop := pdfTopMargin
	for _, k := range keys {
		if strings.HasPrefix(k, "img") {
			// The value of this contains either a path, or a url. Load the image and add it
			image, err := loadImage(values[k])
			if err != nil {
				// TODO log
				continue
			}
			defer image.Close()
			mime, found := detectImageType(image)
			if !found {
				// TODO log
				continue
			}

			// Attempt to load the image
			iopts := gofpdf.ImageOptions{ReadDpi: true, ImageType: mime}
			opts := pdf.RegisterImageOptionsReader(k, iopts, image)
			opts.SetDpi(imgDpi)
			pdf.ImageOptions(k, pdfLeftMargin, yTop, 0, 0, true, iopts, 0, "")

			yTop += pdfBottomMargin + opts.Height()
		} else {
			pdf.Text(pdfLeftMargin, yTop+pdfLineHeight, fmt.Sprintf("%v: %v", k, values[k]))
			yTop += pdfBottomMargin + pdfLineHeight
		}
	}

	return pdf.Output(outWrite)
}

func detectImageType(f *os.File) (string, bool) {
	f.Seek(0, 0)
	defer f.Seek(0, 0)
	buff := make([]byte, 512)
	if _, err := f.Read(buff); err != nil {
		return "something went wrong", false
	}

	switch typ := http.DetectContentType(buff); typ {
	case "image/jpeg":
		fallthrough
	case "image/png":
		fallthrough
	case "image/gif":
		return strings.ToUpper(typ[6:]), true
	}

	return "unsupported", false
}

// imageLock keeps track of which images are currently being downloaded.
var imageLock sync.Map

func loadImage(loc string) (*os.File, error) {
	var file *os.File
	var err error

	if strings.HasPrefix(loc, "http") {
		hash := sha1.New()
		hash.Write([]byte(loc))
		var localPath string = fmt.Sprintf("%v/%v", cacheFolder, base32.StdEncoding.EncodeToString(hash.Sum(nil)))
		// Check if local image exists

		file, err = os.OpenFile(localPath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return file, err
		}

		// Check if the file exists
		_, err := os.Stat(localPath)
		fileExists := err != nil

		// Always download the file
		var wg sync.WaitGroup
		wg.Add(1)
		var ge error
		go func() {
			// Decrement wg, and update the map, signaling retrieval has finished
			defer wg.Done()
			defer imageLock.Store(loc, false)

			// Store that we are loading
			act, loaded := imageLock.LoadOrStore(loc, true)
			if b, _ := act.(bool); fileExists && loaded && b {
				// Only exit if:
				//   1. Something is loaded, signalling some other process might have started downloading, and
				//   2. The loaded value is true, and
				//   3. The file actually exists.

				return
			}

			// File needs to be downloaded
			resp, err := http.Get(loc)
			if err != nil {
				ge = err
				return
			}

			defer resp.Body.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				ge = err
				return
			}

			_, ge = file.Seek(0, io.SeekStart)
		}()

		// If the file does not exist, wait for the waitgroup
		if !fileExists {
			wg.Wait()
		}

		if ge != nil {
			return nil, ge
		}

	} else {

		file, err = os.Open(loc)

		if err != nil {
			return file, err
		}

		var f os.FileInfo
		if f, err = file.Stat(); err != nil && f.IsDir() {
			err = fmt.Errorf("file is a directory")
		}
	}

	return file, err
}

type (
	keyset []string
	values map[string]string
)

func lengthAwareReplace(body, search, replace []byte) []byte {
	// Find where the thing starts
	at := bytes.Index(body, search)
	if at < 0 {
		// Has not been found
		return body
	}

	var res = make([]byte, 0, len(body)-len(search)+len(replace))
	res = append(res, body[:at]...)
	res = append(res, replace...)
	res = append(res, body[at+len(search):]...)

	// Put the length
	binary.LittleEndian.PutUint16(res[at-2:at], uint16(len(replace)))
	return res
}

func convertPdfIfNeeded(rs readSeekWriter) (readSeekWriter, error) {
	i, err := pdfcpu.Info(rs, nil, nil)
	if err != nil {
		return rs, err
	}

	var version string
	var search = "PDF version"
	for _, v := range i {
		if strings.HasPrefix(v, search) {
			// Split off this part
			version = v[len(search)+2:]
			break
		}
	}

	if version == "" {
		// todo error
		return rs, errors.New("version cannot be found")
	}

	if version < "1.7" {
		// todo convert to new
		newWriter, err := memFs.OpenFile("", os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return newWriter, err
		}

		if err := pdfcpu.Trim(rs, newWriter, nil, nil); err != nil {
			return newWriter, err
		}

		rs = newWriter
	}
	return rs, nil
}

// requestToMap builds a values object from a request. It calls the webhook and
// merges the data in the request with the data returned by webhook.
func requestToMap(request *http.Request) (keyset, values) {
	var ms = make(values)
	var order = make(keyset, 0, 10)

	order = append(order, "imglogo")

	segments := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	for _, segment := range segments {
		// Split the segment, only add if the actually is something to add
		keyValues := strings.SplitN(segment, "=", 2)
		if len(keyValues) > 1 {
			order = append(order, keyValues[0])
			ms[keyValues[0]] = keyValues[1]
		} else {
			// TODO error
		}
	}

	if webhook != "" {
		// Call the webhook
		var buf = new(bytes.Buffer)
		json.NewEncoder(buf)

		var err error
		ms[ipAddressKeyname], _, err = net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			log.Err(err).Str("remote_addr", request.RemoteAddr).Msg("could not decode remote ip address")
			delete(ms, ipAddressKeyname)
		}

		if err := json.NewEncoder(buf).Encode(ms); err != nil {
			// TODO log
			goto skip
		}

		ctx, cancel := context.WithTimeout(request.Context(), webhookTimeout)
		defer cancel()

		// Call the webhook
		req, err := http.NewRequestWithContext(ctx, webhookMethod, webhook, buf)
		if err != nil {
			// TODO log
			goto skip
		}

		_ = req.Body.Close()
		req.Header.Set("content-type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// TODO log
			goto skip
		}

		defer resp.Body.Close()
		var mm map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&mm); err != nil {
			// TODO log
			goto skip
		}

		if mergeWebhook {
			for k, v := range mm {
				ms[k] = v
			}
		} else {
			ms = mm
		}
	}

skip:
	delete(ms, ipAddressKeyname)

	return order, ms
}
