package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"

	"github.com/tuupke/pixie/env"
)

var (
	pdfUnit        = env.StringFb("PDF_UNIT_SIZE", "mm")
	pdfSize        = env.StringFb("PDF_PAGE_SIZE", "A4")
	pdfFontDir     = env.String("PDF_FONT_DIR")
	pdfInLandscape = env.Bool("PDF_LANDSCAPE")

	pdfLeftMargin   = float64(env.IntFb("PDF_LEFT_MARGIN", 10))
	pdfTopMargin    = float64(env.IntFb("PDF_TOP_MARGIN", 15))
	pdfBottomMargin = float64(env.IntFb("PDF_BOTTOM_MARGIN", 6))
	pdfLineHeight   = float64(env.IntFb("PDF_LINE_HEIGHT", 8))

	imgDpi = float64(env.IntFb("IMAGE_PPI", 120))
)

func BannerPage(outWrite io.Writer, data map[string]string, keys ...string) error {
	orientation := "P"
	if pdfInLandscape {
		orientation = "L"
	}

	pdf := gofpdf.New(orientation, pdfUnit, pdfSize, pdfFontDir)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	yTop := pdfTopMargin
	for _, k := range keys {
		val, ok := data[k]
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
