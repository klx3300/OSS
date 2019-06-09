package osstheme

import (
	"evbus"
	"image/color"
	"logger"
	"os"
	"staticres"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// CNFont is the font theme for chinese chars
type CNFont struct{}

const fontName = "cnfont.ttc"

var theFont *fyne.StaticResource

var log = logger.Logger{
	LogLevel: 0,
	Name:     "CNFonter:",
}

// Initialization the themes
func Initialization(bus *evbus.EventBus) {
	log.Debugln("Initializing..")
	fsres, reserr := staticres.GetResource("cnfont", fontName)
	if reserr != nil {
		log.FailLn("Unable to read font from disk:" + reserr.Error())
		log.FailLn("Unrecoverable failure happened during initialization process.")
		log.FailLn("::QUIT.")
		os.Exit(-1)
	}
	theFont = fyne.NewStaticResource(fsres.Name, fsres.Content)
}

// BackgroundColor is required
func (CNFont) BackgroundColor() color.Color {
	return theme.DarkTheme().BackgroundColor()
}

// ButtonColor is required
func (CNFont) ButtonColor() color.Color {
	return theme.DarkTheme().ButtonColor()
}

// HyperlinkColor is required
func (CNFont) HyperlinkColor() color.Color {
	return theme.DarkTheme().HyperlinkColor()
}

// TextColor is required
func (CNFont) TextColor() color.Color {
	return theme.DarkTheme().TextColor()
}

// PlaceHolderColor is required
func (CNFont) PlaceHolderColor() color.Color {
	return theme.DarkTheme().PlaceHolderColor()
}

// PrimaryColor is required
func (CNFont) PrimaryColor() color.Color {
	return theme.DarkTheme().PrimaryColor()
}

// FocusColor is required
func (CNFont) FocusColor() color.Color {
	return theme.DarkTheme().FocusColor()
}

// ScrollBarColor is required
func (CNFont) ScrollBarColor() color.Color {
	return theme.DarkTheme().ScrollBarColor()
}

// TextSize is required
func (CNFont) TextSize() int {
	return theme.DarkTheme().TextSize()
}

// TextFont is required
func (CNFont) TextFont() fyne.Resource {
	return theFont
}

// TextBoldFont is required
func (CNFont) TextBoldFont() fyne.Resource {
	return theFont
}

// TextItalicFont is required
func (CNFont) TextItalicFont() fyne.Resource {
	return theFont
}

// TextBoldItalicFont is required
func (CNFont) TextBoldItalicFont() fyne.Resource {
	return theFont
}

// TextMonospaceFont is required
func (CNFont) TextMonospaceFont() fyne.Resource {
	return theFont
}

// Padding is required
func (CNFont) Padding() int {
	return theme.DarkTheme().Padding()
}

// IconInlineSize is required
func (CNFont) IconInlineSize() int {
	return theme.DarkTheme().IconInlineSize()
}

// ScrollBarSize is required
func (CNFont) ScrollBarSize() int {
	return theme.DarkTheme().ScrollBarSize()
}
