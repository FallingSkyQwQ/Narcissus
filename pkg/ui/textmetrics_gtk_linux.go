//go:build linux

package ui

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

// gtkTextMeasurer measures text with Pango, the same engine GTK4 uses to draw
// labels, so laid-out text matches rendered text. It lazily creates a scratch
// label: Pango layouts need a widget context, and GTK widgets may only be
// created after gtk_init, which happens inside the application's activate
// handler — by which point the first layout pass is about to run.
type gtkTextMeasurer struct {
	label *gtk.Label
}

// newGTKTextMeasurer returns a measurer for the GTK4 backend. The scratch label
// is deferred until the first measurement so construction is safe before
// GtkApplication.Run.
func newGTKTextMeasurer() TextMeasurer {
	return &gtkTextMeasurer{}
}

func (m *gtkTextMeasurer) scratch() *gtk.Label {
	if m.label == nil {
		m.label = gtk.NewLabel("")
	}
	return m.label
}

// MeasureText returns the natural (unwrapped) size of text.
func (m *gtkTextMeasurer) MeasureText(text string, font Font) (width, height float32) {
	return m.measure(text, font, 0)
}

// MeasureTextWidth returns the size text occupies when wrapped to maxWidth.
// It is picked up through the wrappedTextMeasurer interface.
func (m *gtkTextMeasurer) MeasureTextWidth(text string, font Font, maxWidth float32) (width, height float32) {
	return m.measure(text, font, maxWidth)
}

func (m *gtkTextMeasurer) measure(text string, font Font, maxWidth float32) (width, height float32) {
	size := font.Size
	if size <= 0 {
		size = DefaultFont().Size
	}
	lineHeight := font.LineHeight
	if lineHeight <= 0 {
		lineHeight = 1.2
	}

	if text == "" {
		return 0, size * lineHeight
	}

	layout := m.scratch().CreatePangoLayout(text)

	desc := pango.NewFontDescription()
	if font.Family != "" {
		desc.SetFamily(font.Family)
	}
	// Pango absolute sizes are in PANGO_SCALE units (1/1024 of a pixel), not
	// pixels, so the pixel size must be scaled up.
	desc.SetAbsoluteSize(float64(size) * float64(pango.SCALE))
	desc.SetWeight(pangoWeight(font.Weight))
	desc.SetStyle(pangoStyle(font.Style))
	layout.SetFontDescription(desc)

	if maxWidth > 0 {
		layout.SetWrap(pango.WrapWordChar)
		layout.SetWidth(int(float64(maxWidth) * float64(pango.SCALE)))
	}

	w, h := layout.PixelSize()
	lines := layout.LineCount()
	if lines < 1 {
		lines = 1
	}

	// Pango reports the font's natural height; bias it up toward the style's
	// line-height multiplier so it never underestimates what the CSS-rendered
	// label will occupy.
	height = float32(h)
	if min := float32(lines) * size * lineHeight; height < min {
		height = min
	}

	width = float32(w)
	if maxWidth > 0 && width > maxWidth {
		width = maxWidth
	}
	return width, height
}

func pangoWeight(w FontWeight) pango.Weight {
	switch w {
	case FontWeightLight:
		return pango.WeightLight
	case FontWeightMedium:
		return pango.WeightMedium
	case FontWeightSemiBold:
		return pango.WeightSemibold
	case FontWeightBold:
		return pango.WeightBold
	default:
		return pango.WeightNormal
	}
}

func pangoStyle(s FontStyle) pango.Style {
	if s == FontStyleItalic {
		return pango.StyleItalic
	}
	return pango.StyleNormal
}
