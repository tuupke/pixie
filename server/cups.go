package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenPrinting/goipp"
	"github.com/jung-kurt/gofpdf"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/valyala/fasthttp"
)

type (
	readSeekWriter interface {
		io.ReadSeeker
		io.Writer
		io.Closer
	}
)

// btsReplace takes a byte string and converts it to the representation used in CUPS
func btsReplace(check []byte) []byte {
	checkLen := len(check)
	b := make([]byte, 2+checkLen)
	copy(b[2:], check)

	// Prepend bigendian length, it's encoded in two bytes
	binary.BigEndian.PutUint16(b, uint16(checkLen))
	return b
}

var (
	// num = new(uint32)

	printerTo = envStringFb("IPP_PRINTER_URL", "localhost:631/printers/Virtual_PDF_Printer")
	// mergeWebhook     = envBool("WEBHOOK_MERGE")
	ipAddressKeyname = envStringFb("WEBHOOK_IP_NAME", "team_ip_address")
	cupsListen       = envStringFb("CUPS_ADDR", ":631")

	pdfUnit        = envStringFb("PDF_UNIT_SIZE", "mm")
	pdfSize        = envStringFb("PDF_PAGE_SIZE", "A4")
	pdfFontDir     = envString("PDF_FONT_DIR")
	pdfInLandscape = envBool("PDF_LANDSCAPE")
	timeoutPdf     = envDurationFb("PDF_REFRESH_DURATION", time.Minute*5)
	cacheFolder    = envStringFb("IPP_CACHE_DIR", "/tmp/pixie/")

	ippDebugFolder = envString("IPP_DEBUG_FOLDER")
	ippRequestNum  = new(uint64)

	imgDpi = float64(envIntFb("IMAGE_PPI", 120))

	pdfLeftMargin   = float64(envIntFb("PDF_LEFT_MARGIN", 10))
	pdfTopMargin    = float64(envIntFb("PDF_TOP_MARGIN", 15))
	pdfBottomMargin = float64(envIntFb("PDF_BOTTOM_MARGIN", 6))
	pdfLineHeight   = float64(envIntFb("PDF_LINE_HEIGHT", 8))

	memFs afero.Fs

	to = btsReplace([]byte("ipp://" + printerTo))
)

func ensureFolder(folder *string, empty bool) {
	*folder = strings.TrimRight(*folder, "/")

	if empty {
		// remove and let rest of procedure recreate, don't really care about failing removal
		_ = os.RemoveAll(*folder)
	}

	stat, err := os.Stat(*folder)

	// If err is not nil, attempt creation
	if err != nil {
		// Attempt creation
		if err = os.MkdirAll(*folder, 0755); err != nil {
			panic(fmt.Errorf("'%v' does not exist and can't be created; %w", *folder, err))
		}

		stat, err = os.Stat(*folder)
	}

	if err != nil {
		panic(err)
	}

	if !stat.IsDir() {
		panic(fmt.Errorf("'%v' is not a directory", stat.Name()))
	}
}

func init() {
	ensureFolder(&cacheFolder, false)
	memFs = afero.NewMemMapFs()

	if ippDebugFolder != "" {
		ensureFolder(&ippDebugFolder, true)
	}
}

const createOpts = os.O_WRONLY | os.O_CREATE | os.O_TRUNC

// CupsHandler is the handler which proxies all cups requests and fixes
func CupsHandler(ctx *fasthttp.RequestCtx) {
	path := bytes.Trim(ctx.Request.URI().Path(), "/")
	requestedUrl := fmt.Sprintf("ipp://%s/%s", cupsListen, []byte(path))

	from := btsReplace([]byte(requestedUrl))

	btso := ctx.PostBody()

	var ippNum = atomic.AddUint64(ippRequestNum, 1)

	if ippDebugFolder != "" {
		log.Err(os.WriteFile(fmt.Sprintf("%v/%v-1.1-req-pre.header", ippDebugFolder, ippNum), ctx.Request.Header.Header(), 0x0755)).Msg("written header pre-replace ipp debug")
		log.Err(os.WriteFile(fmt.Sprintf("%v/%v-1.2-req-pre.bin", ippDebugFolder, ippNum), btso, 0x0755)).Msg("written request pre-replace ipp debug")
	}

	// Replace the url
	bts := bytes.Replace(btso, from, to, -1)

	// Read the message
	var req goipp.Message
	if err := req.DecodeBytes(bts); err != nil {
		panic(err)
	}

	// In normal cases, simply proxy the request, only when a document is sent some magic is needed.
	var b io.Reader = bytes.NewBuffer(bts)
	if goipp.Op(req.Code) == goipp.OpSendDocument {
		remoteIp := ctx.Conn().RemoteAddr().String()
		host, _, err := net.SplitHostPort(remoteIp)
		log.Err(err).Str("remote", remoteIp).Str("ip", host).Msg("received print from")
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		// Search for the PDF header
		sep := bytes.Index(bts, []byte("%PDF"))
		if sep < 0 {
			panic("Not found!")
		}

		var rw *io.PipeWriter
		b, rw = io.Pipe()

		go func() {
			if _, err := rw.Write(bts[:sep]); err != nil {
				panic(err)
			}

			hash := sha1.New()
			hash.Write([]byte(host))
			pdfLocation := fmt.Sprintf("%v/%s.pdf", cacheFolder, base32.StdEncoding.EncodeToString(hash.Sum(nil)))

			stat, stErr := os.Stat(pdfLocation)
			skipBuild := stErr == nil && stat != nil && time.Now().Before(stat.ModTime().Add(timeoutPdf))

			// Test if the page exists and is recent enough
			var oldPage readSeekWriter
			var err error
			oldPage, err = os.OpenFile(pdfLocation, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				skipBuild = false
			}

			// If the cached page is not there, create a buffer which can be used to store the file in
			if oldPage == nil {
				oldPage, err = memFs.OpenFile(pdfLocation, os.O_RDWR|os.O_CREATE, 0755)
				if err != nil {
					// TODO log
					panic(err)
				}

				defer memFs.Remove(pdfLocation)
			} else {
				oldPage.Seek(0, 0)
			}

			pdfParts := []io.ReadSeeker{oldPage, bytes.NewReader(bts[sep:])}
			defer oldPage.Close()
			if !skipBuild {
				order, mp := requestToMap(ctx.Request)

				var d ExternalData
				// Load from the orm
				err := orm.Model(Host{}).Select("external_data.*").
					Joins("join external_data  on external_data.host_id = hosts.guid").
					Where("hosts.primary_ip = ?", host).
					First(&d)

				if err == nil {
					order = append(order, "username", "teamname", "teamid", "location")
					mp["username"] = d.Username
					if d.Teamname != nil {
						mp["teamname"] = *d.Teamname
					}
					if d.TeamId != nil {
						mp["teamid"] = *d.TeamId
					}
					mp["location"] = fmt.Sprintf("(%v, %v, %v)", d.Location.X, d.Location.Y, d.Location.Rotation)
				}

				if err := renderPage(oldPage, order, mp); err != nil {
					// TODO log
					pdfParts = pdfParts[1:]
				} else if _, err := oldPage.Seek(0, 0); err != nil {
					// TODO log
					pdfParts = pdfParts[1:]
				}
			}

			if err := pdfcpu.Merge(pdfParts, rw, nil); err != nil {
				panic(err)
			}

			_ = rw.Close()
		}()
	}

	hreq, err := http.NewRequest(string(ctx.Method()), "http://"+printerTo, b)

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		hreq.Header.Add(string(key), strings.Replace(string(value), requestedUrl, printerTo, -1))
	})

	if ippDebugFolder != "" {
		log.Err(os.WriteFile(fmt.Sprintf("%v/%v-2.1-req-post.header", ippDebugFolder, ippNum), ctx.Request.Header.Header(), 0x0755)).Msg("written request header post-replace ipp debug")
		log.Err(os.WriteFile(fmt.Sprintf("%v/%v-2.2-req-post.bin", ippDebugFolder, ippNum), btso, 0x0755)).Msg("written request body post-replace ipp debug")
	}

	resp, err := http.DefaultClient.Do(hreq)
	_ = hreq.Body.Close()
	if err != nil {
		panic(err)
	}

	bts, err = io.ReadAll(resp.Body)

	_ = resp.Body.Close()

	if ippDebugFolder != "" {
		f, err := os.OpenFile(fmt.Sprintf("%v/%v-3.1-resp-pre.header", ippDebugFolder, ippNum), createOpts, 0x0755)
		if err == nil {
			defer f.Close()
			err = resp.Header.Clone().Write(f)
		}
		log.Err(err).Msg("written response header pre-replace ipp debug")
		log.Err(os.WriteFile(fmt.Sprintf("%v/%v-3.2-resp-pre.bin", ippDebugFolder, ippNum), bts, 0x0755)).Msg("written response body pre-replace ipp debug")
	}

	bts = regexp.MustCompile("document-format-supported[\\S\\s]*?D\\x00").ReplaceAll(bts, []byte("document-format-supported\x00\x0fapplication/pdfD\x00"))
	bts = bytes.Replace(bts, to, from, -1)

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

	if ippDebugFolder != "" {
		f, err := os.OpenFile(fmt.Sprintf("%v/%v-4.1-resp-post.header", ippDebugFolder, ippNum), createOpts, 0x755)
		if err == nil {
			defer f.Close()
			err = rawHeaders.Write(f)
		}
		log.Err(err).Msg("written response header post-replace ipp debug")
		log.Err(os.WriteFile(fmt.Sprintf("%v/%v-4.2-resp-post.bin", ippDebugFolder, ippNum), bts, 0x755)).Msg("written response body post-replace ipp debug")
	}

	if _, err := ctx.Write(bts); err != nil {
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
func requestToMap(request fasthttp.Request) (keyset, values) {
	var ms = make(values)
	var order = make(keyset, 0, 10)

	order = append(order, "imglogo")

	segments := strings.Split(strings.Trim(string(request.URI().Path()), "/"), "/")
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

	// Load from the orm
	// orm.Model(ExternalData{}).Where("ip_address", "=", request.)

	// if webhook != "" {
	// 	// Call the webhook
	// 	var buf = new(bytes.Buffer)
	// 	json.NewEncoder(buf)
	//
	// 	var err error
	// 	ms[ipAddressKeyname], _, err = net.SplitHostPort(request.RemoteAddr)
	// 	if err != nil {
	// 		log.Err(err).Str("remote_addr", request.RemoteAddr).Msg("could not decode remote ip address")
	// 		delete(ms, ipAddressKeyname)
	// 	}
	//
	// 	if err := json.NewEncoder(buf).Encode(ms); err != nil {
	// 		// TODO log
	// 		goto skip
	// 	}
	//
	// 	ctx, cancel := context.WithTimeout(ctx, webhookTimeout)
	// 	defer cancel()
	//
	// 	// Call the webhook
	// 	req, err := http.NewRequestWithContext(ctx, webhookMethod, webhook, buf)
	// 	if err != nil {
	// 		// TODO log
	// 		goto skip
	// 	}
	//
	// 	_ = req.Body.Close()
	// 	req.Header.Set("content-type", "application/json")
	// 	resp, err := http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		// TODO log
	// 		goto skip
	// 	}
	//
	// 	defer resp.Body.Close()
	// 	var mm map[string]string
	// 	if err := json.NewDecoder(resp.Body).Decode(&mm); err != nil {
	// 		// TODO log
	// 		goto skip
	// 	}
	//
	// 	if mergeWebhook {
	// 		for k, v := range mm {
	// 			ms[k] = v
	// 		}
	// 	} else {
	// 		ms = mm
	// 	}
	// }

	// skip:
	delete(ms, ipAddressKeyname)

	return order, ms
}
