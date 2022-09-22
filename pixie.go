package pixie

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"openticket.tech/lifecycle/v2"

	"github.com/spf13/afero"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"openticket.tech/db"
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

var orm *gorm.DB

func Orm() *gorm.DB {
	if orm == nil {
		var err error
		sql, err := db.Conn("DB")
		if err != nil {
			log.Fatal().Err(err).Msg("could not retrieve database")
		}
		lifecycle.EFinally(sql.Close)

		d := sqlite.Dialector{
			Conn: sql,
		}

		orm, err = gorm.Open(d, &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
		log.Err(err).Msg("loaded gorm")
		if err != nil {
			log.Fatal().Msg("gorm must boot")
		}
	}

	return orm
}
