package main

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"time"

	githttp "github.com/AaronO/go-git-http"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/hashicorp/mdns"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"

	"github.com/tuupke/pixie/beanstalk"
	"github.com/tuupke/pixie/packets"
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
	printerTo        = envStringFb("IPP_PRINTER_URL", "localhost:631/printers/Virtual_PDF_Printer")
	webhook          = envStringFb("WEBHOOK_URL", "")
	webhookMethod    = envStringFb("WEBHOOK_METHOD", "GET")
	webhookTimeout   = envDurationFb("WEBHOOK_TIMEOUT", time.Millisecond*250)
	mergeWebhook     = envBool("WEBHOOK_MERGE")
	ipAddressKeyname = envStringFb("WEBHOOK_IP_NAME", "team_ip_address")

	pdfUnit        = envStringFb("PDF_UNIT_SIZE", "mm")
	pdfSize        = envStringFb("PDF_PAGE_SIZE", "A4")
	pdfFontDir     = envString("PDF_FONT_DIR")
	pdfInLandscape = envBool("PDF_LANDSCAPE")
	timeoutPdf     = envDurationFb("PDF_REFRESH_DURATION", time.Minute*5)
	cacheFolder    = envStringFb("IPP_CACHE_DIR", "/tmp/pixie/")

	imgDpi = float64(envIntFb("IMAGE_PPI", 120))

	pdfLeftMargin   = float64(envIntFb("PDF_LEFT_MARGIN", 10))
	pdfTopMargin    = float64(envIntFb("PDF_TOP_MARGIN", 15))
	pdfBottomMargin = float64(envIntFb("PDF_BOTTOM_MARGIN", 6))
	pdfLineHeight   = float64(envIntFb("PDF_LINE_HEIGHT", 8))

	keyLength = envIntFb("RSA_KEY_LENGTH", 1024)

	listen = envStringFb("LISTEN_ADDR", ":6632")

	ansibleRepoDir = envStringFb("REPO_DIR", "/tmp/repo")

	domjudgeUrl = envStringFb("DOMJUDGE_URL", "https://mart:**29~what~PAIR~never~16**@judge.swerc.eu")

	memFs afero.Fs

	beanstalkLocation = "127.0.0.1:11300"

	to = btsReplace([]byte("ipp://" + printerTo))
)

func ensureFolder(folder *string) {
	*folder = strings.TrimRight(*folder, "/")

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
	ensureFolder(&ansibleRepoDir)
	ensureFolder(&cacheFolder)

	memFs = afero.NewMemMapFs()
}

var cId = "practice"

func main() {
	// Get git handler to serve the repo
	git := githttp.New(ansibleRepoDir, "git")
	http.Handle("/git", git)

	// CUPS only uses a single endpoint, all requests to this handler will be rerouted to there.

	// TODO verify whether adding "/ipp" has worked
	http.HandleFunc("/", simpleProxy)

	// http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	// 	resp, _ := http.Get(domjudgeUrl + "/api/contests/practice/teams")
	// 	io.Copy(writer, resp.Body)
	// 	defer resp.Body.Close()
	// })

	_, err := netip.ParseAddrPort(beanstalkLocation)
	if err != nil {
		log.Fatal().Msg("cannot parse beanstalkd addr")
	}

	// Setup our service export
	host, _ := os.Hostname()
	info := []string{beanstalkLocation}
	service, err := mdns.NewMDNSService(host, "_beanstalk._pixie._tcp", "progcont.", "", 8000, nil, info)
	log.Err(err).Msg("created mDNS service")

	// Create the mDNS server, defer shutdown
	server, err := mdns.NewServer(&mdns.Config{Zone: service})

	defer server.Shutdown()
	log.Err(err).Msg("started mDNS advertisement")

	// publish()

	panic(http.ListenAndServe(listen, nil))
}

func publish() {
	// Connect to beanstalk
	conn := beanstalk.Connect("tcp", beanstalkLocation, "demeter")

	b := flatbuffers.NewBuilder(128)

	header, body := b.CreateString("header"), b.CreateString("body")
	packets.NotifyStart(b)
	packets.NotifyAddHeader(b, header)
	packets.NotifyAddBody(b, body)
	notif := packets.NotifyEnd(b)

	packets.CommandStart(b)
	packets.CommandAddCommandType(b, packets.CmdNotify)
	packets.CommandAddCommand(b, notif)
	b.Finish(packets.CommandEnd(b))

	id, err := conn.Put("demeter", b.FinishedBytes(), 0, time.Duration(0), 0)
	log.Err(err).Uint64("id", id).Msg("published")
}
