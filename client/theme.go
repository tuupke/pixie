package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	purple = &color.NRGBA{R: 128, G: 0, B: 128, A: 255}
	orange = &color.NRGBA{R: 198, G: 123, B: 0, A: 255}
	grey   = &color.Gray{Y: 123}
)

// customTheme is a simple demonstration of a bespoke theme loaded by a Fyne app.
type customTheme struct {
}

func (ct) Color(c fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameBackground:
		return purple
	case theme.ColorNameButton, theme.ColorNameDisabled:
		return color.Black
	case theme.ColorNamePlaceHolder, theme.ColorNameScrollBar:
		return grey
	case theme.ColorNamePrimary, theme.ColorNameHover, theme.ColorNameFocus:
		return orange
	case theme.ColorNameShadow:
		return &color.RGBA{R: 0xcc, G: 0xcc, B: 0xcc, A: 0xcc}
	default:
		return color.White
	}
}

type ct struct{ fyne.Theme }

func newCustomTheme() fyne.Theme {
	return ct{fyne.CurrentApp().Settings().Theme()}
}
