package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"

	"github.com/tuupke/pixie/env"
)

var (
	pdfUnit        = env.StringFb("PDF_UNIT", "mm")
	pdfSize        = env.StringFb("PDF_PAGE_SIZE", "A4")
	pdfFontDir     = env.String("PDF_FONT_DIR")
	pdfInLandscape = env.Bool("PDF_LANDSCAPE")

	pdfLeftMargin   = float64(env.IntFb("PDF_LEFT_MARGIN", 10))
	pdfTopMargin    = float64(env.IntFb("PDF_TOP_MARGIN", 15))
	pdfBottomMargin = float64(env.IntFb("PDF_BOTTOM_MARGIN", 6))
	pdfLineHeight   = float64(env.IntFb("PDF_LINE_HEIGHT", 3))

	imgDpi = float64(env.IntFb("IMAGE_PPI", 120))
)

func BannerPage(log zerolog.Logger, outWrite io.Writer, data *Props, keys ...string) error {
	if len(keys) == 1 && keys[0] == "*" {
		keys = make([]string, 0, 100)
		data.Range(func(key, _ string) bool {
			keys = append(keys, key)
			return true
		})

		slices.Sort(keys)
	}

	orientation := "P"
	if pdfInLandscape {
		orientation = "L"
	}

	log.Info().Bool("landscape", pdfInLandscape).Int("num_keys", len(keys)).Msg("rendering new banner")

	pdf := gofpdf.New(orientation, pdfUnit, pdfSize, pdfFontDir)
	if bannerOnBack {
		// Add an empty page if the banner is supposed to be printed on the back. Assumes a duplexer is installed.
		pdf.AddPage()
	}
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	yTop := pdfTopMargin
	for _, k := range keys {
		val, ok := data.Load(k)
		if !ok {
			continue
		}

		// Note, the length check here to make sure we can safely extract the extension.
		// This does mean that the image name "/tmp/pixie/png" is invalid!
		if strings.HasPrefix(k, "img") && len(val) >= 4 {
			// The value of this contains either a path, or a url. Load the image and add it
			image, err := os.Open(val)
			if err != nil {
				fmt.Printf("Cannot open image ('%v') for key '%v'; %v", val, k, err)
				// TODO log
				continue
			}

			// They way the download works ensures the image is stored using the correct extension
			// Take the last 4 characters from the filename and strip the '.' period if needed.
			ext := val[len(val)-4:]
			if ext[0] == '.' {
				ext = ext[1:]
			}

			mime := "image/" + ext
			defer image.Close()

			// Attempt to load the image
			iopts := gofpdf.ImageOptions{ReadDpi: true, ImageType: mime}
			opts := pdf.RegisterImageOptionsReader(k, iopts, image)
			opts.SetDpi(imgDpi)
			pdf.ImageOptions(k, pdfLeftMargin, yTop, 0, 0, true, iopts, 0, "")

			yTop += pdfBottomMargin + opts.Height()
		} else {
			pdf.Text(pdfLeftMargin, yTop+pdfLineHeight, fmt.Sprintf("%v: %v", k, val))
			yTop += pdfBottomMargin + pdfLineHeight
		}
	}

	return pdf.Output(outWrite)
}
