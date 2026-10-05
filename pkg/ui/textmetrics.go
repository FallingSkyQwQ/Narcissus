package ui

import (
	"strings"
	"sync"
)

// TextMeasurer reports the rendered size of a string for a given Font.
//
// It exists so that the framework's layout engine can measure text with the
// same font engine that actually draws it. Backends install their toolkit's
// measurer (for example Pango on GTK4); when none is installed, the package
// falls back to a conservative estimate so that headless builds and tests
// still produce a stable layout.
type TextMeasurer interface {
	// MeasureText returns the width and height the text occupies when drawn
	// with font, in pixels. Height covers every line of a multi-line string.
	MeasureText(text string, font Font) (width, height float32)
}

var (
	textMeasurerMu sync.RWMutex
	textMeasurer   TextMeasurer = estimateMeasurer{}
)

// SetTextMeasurer installs the measurer used by MeasureText. Passing nil
// restores the built-in estimator.
func SetTextMeasurer(m TextMeasurer) {
	textMeasurerMu.Lock()
	defer textMeasurerMu.Unlock()
	if m == nil {
		textMeasurer = estimateMeasurer{}
		return
	}
	textMeasurer = m
}

// TextMeasurerOrDefault returns the installed measurer, never nil.
func TextMeasurerOrDefault() TextMeasurer {
	textMeasurerMu.RLock()
	defer textMeasurerMu.RUnlock()
	if textMeasurer == nil {
		return estimateMeasurer{}
	}
	return textMeasurer
}

// MeasureText measures text with the installed measurer. It is the single
// entry point widgets use so that a backend can swap in real font metrics
// without every widget changing.
func MeasureText(text string, font Font) (width, height float32) {
	return TextMeasurerOrDefault().MeasureText(text, font)
}

// wrappedTextMeasurer is optionally implemented by a TextMeasurer that can
// account for line wrapping at a maximum width.
type wrappedTextMeasurer interface {
	MeasureTextWidth(text string, font Font, maxWidth float32) (width, height float32)
}

// MeasureTextWrapped measures text as it would occupy when wrapped to maxWidth.
// The installed measurer is used when it supports width-aware measurement;
// otherwise a greedy word-wrap over MeasureText is used. It returns maxWidth as
// the width when wrapping occurs, so callers allocate the full line width.
func MeasureTextWrapped(text string, font Font, maxWidth float32) (width, height float32) {
	if maxWidth <= 0 {
		return MeasureText(text, font)
	}
	if m, ok := TextMeasurerOrDefault().(wrappedTextMeasurer); ok {
		return m.MeasureTextWidth(text, font, maxWidth)
	}
	return estimateWrapped(text, font, maxWidth)
}

// estimateWrapped greedily wraps text by words for measurers without a
// width-aware implementation.
func estimateWrapped(text string, font Font, maxWidth float32) (width, height float32) {
	size := font.Size
	if size <= 0 {
		size = DefaultFont().Size
	}
	lineHeight := font.LineHeight
	if lineHeight <= 0 {
		lineHeight = 1.2
	}

	lines := 0
	for _, paragraph := range strings.Split(text, "\n") {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			lines++
			continue
		}
		current := ""
		for _, word := range words {
			candidate := word
			if current != "" {
				candidate = current + " " + word
			}
			candidateWidth, _ := MeasureText(candidate, font)
			if current != "" && candidateWidth > maxWidth {
				lines++
				current = word
			} else {
				current = candidate
			}
		}
		if current != "" {
			lines++
		}
	}
	if lines == 0 {
		lines = 1
	}
	return maxWidth, float32(lines) * size * lineHeight
}

// estimateMeasurer approximates text size from per-rune advance widths. It does
// not consult any font and is only used when a backend has not installed a real
// measurer, so layouts stay stable but not pixel-perfect.
type estimateMeasurer struct{}

func (estimateMeasurer) MeasureText(text string, font Font) (width, height float32) {
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

	var maxLine, current float32
	lines := 1
	for _, r := range text {
		if r == '\n' {
			if current > maxLine {
				maxLine = current
			}
			current = 0
			lines++
			continue
		}
		current += runeAdvance(r, size)
	}
	if current > maxLine {
		maxLine = current
	}

	return maxLine, float32(lines) * size * lineHeight
}

// runeAdvance estimates the advance width of a single rune at the given font
// size. The factors are rough averages for a typical sans-serif face.
func runeAdvance(r rune, size float32) float32 {
	switch {
	case r == '\t':
		return size * 1.6
	case r == ' ':
		return size * 0.28
	case isWideRune(r):
		// CJK ideographs, Hangul, full-width punctuation and emoji are
		// approximately square.
		return size
	case isNarrowRune(r):
		// Narrow Latin letters and punctuation.
		return size * 0.32
	case r < 0x20:
		// Control characters contribute no advance.
		return 0
	default:
		return size * 0.55
	}
}

// isNarrowRune reports whether r is one of the conventionally narrow glyphs.
func isNarrowRune(r rune) bool {
	switch r {
	case 'i', 'l', 'j', 'I', '!', '.', ',', ':', ';', '\'',
		'|', '`', 'f', 't', 'r', '(', ')', '[', ']', '{', '}':
		return true
	}
	return false
}

// isWideRune reports whether r belongs to an East Asian Wide/Fullwidth or
// emoji block, whose glyphs advance roughly one em.
func isWideRune(r rune) bool {
	switch {
	case 0x1100 <= r && r <= 0x115F, // Hangul Jamo
		0x2E80 <= r && r <= 0x303E,   // CJK radicals, Kangxi, CJK symbols
		0x3041 <= r && r <= 0x33FF,   // Hiragana, Katakana, CJK compat
		0x3400 <= r && r <= 0x4DBF,   // CJK Unified Ideographs Extension A
		0x4E00 <= r && r <= 0x9FFF,   // CJK Unified Ideographs
		0xA000 <= r && r <= 0xA4CF,   // Yi
		0xAC00 <= r && r <= 0xD7A3,   // Hangul syllables
		0xF900 <= r && r <= 0xFAFF,   // CJK compatibility ideographs
		0xFE30 <= r && r <= 0xFE4F,   // CJK compatibility forms
		0xFF00 <= r && r <= 0xFF60,   // Fullwidth forms
		0xFFE0 <= r && r <= 0xFFE6,   // Fullwidth signs
		0x1F300 <= r && r <= 0x1FAFF, // Emoji and symbols
		0x20000 <= r && r <= 0x3FFFD: // CJK Extension B+
		return true
	}
	return false
}
