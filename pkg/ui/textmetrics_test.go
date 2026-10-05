package ui

import (
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

func TestMeasureTextEmpty(t *testing.T) {
	w, h := MeasureText("", DefaultFont())
	if w != 0 {
		t.Errorf("empty text width = %g, want 0", w)
	}
	if h <= 0 {
		t.Errorf("empty text height = %g, want > 0", h)
	}
}

func TestEstimateWideRunesAreWider(t *testing.T) {
	font := DefaultFont()
	latin, _ := MeasureText("aaaa", font)
	wide, _ := MeasureText("中中中中", font)
	if wide <= latin {
		t.Errorf("wide text (%g) should be wider than latin text (%g)", wide, latin)
	}
}

func TestMeasureTextRespectsNewlines(t *testing.T) {
	font := DefaultFont()
	_, one := MeasureText("hello", font)
	_, two := MeasureText("hello\nworld", font)
	if two <= one {
		t.Errorf("multi-line height %g should exceed single line %g", two, one)
	}
}

// fixedMeasurer reports a constant size per rune so tests can assert that
// widget measurement actually routes through the installed measurer.
type fixedMeasurer struct {
	perRune float32
	height  float32
}

func (m fixedMeasurer) MeasureText(text string, _ Font) (float32, float32) {
	return float32(len([]rune(text))) * m.perRune, m.height
}

func TestWidgetsUseInstalledMeasurer(t *testing.T) {
	SetTextMeasurer(fixedMeasurer{perRune: 10, height: 7})
	defer SetTextMeasurer(nil)

	txt := NewText("abc")
	got := txt.Measure(flex.Constraint{})
	if got.Width != 30 {
		t.Errorf("text width = %g, want 30", got.Width)
	}
	if got.Height != 7 {
		t.Errorf("text height = %g, want 7", got.Height)
	}

	btn := NewButton().Text("abc") // default padding 12 on each side
	got = btn.Measure(flex.Constraint{})
	if got.Width != 30+24 {
		t.Errorf("button width = %g, want 54", got.Width)
	}
}

// TestTextMeasureDoesNotCacheConstrainedSize guards against reusing a size that
// was clamped by MaxWidth the first time the widget was measured.
func TestTextMeasureDoesNotCacheConstrainedSize(t *testing.T) {
	SetTextMeasurer(fixedMeasurer{perRune: 10, height: 7})
	defer SetTextMeasurer(nil)

	txt := NewText("abcd")

	natural := txt.Measure(flex.Constraint{})
	if natural.Width != 40 {
		t.Fatalf("natural width = %g, want 40", natural.Width)
	}

	clamped := txt.Measure(flex.Constraint{MaxWidth: 15})
	if clamped.Width != 15 {
		t.Fatalf("clamped width = %g, want 15", clamped.Width)
	}

	again := txt.Measure(flex.Constraint{})
	if again.Width != 40 {
		t.Errorf("width after clamping = %g, want natural 40", again.Width)
	}
}
