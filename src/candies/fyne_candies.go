package candies

import (
	"fyne.io/fyne"
)

// TextStyle is the candied fyne text style
type TextStyle fyne.TextStyle

// NewTextStyle creates..
func NewTextStyle() TextStyle {
	return TextStyle(fyne.TextStyle{
		Bold:      false,
		Italic:    false,
		Monospace: false,
	})
}

// SetBold return with bold set
func (ts TextStyle) SetBold() TextStyle {
	k := ts
	k.Bold = true
	return k
}

// SetItalic return a new textstyle
func (ts TextStyle) SetItalic() TextStyle {
	k := ts
	k.Italic = true
	return k
}

// SetMono return a new textstyle
func (ts TextStyle) SetMono() TextStyle {
	k := ts
	k.Monospace = true
	return k
}

// Fin convert it to fyne-acceptable state
func (ts TextStyle) Fin() fyne.TextStyle {
	return fyne.TextStyle(ts)
}
