package flex

import "testing"

func measuredItem(w, h float32) Item {
	item := Item{}
	item.SetMeasuredSize(w, h)
	return item
}

// withMeasured sets the measured size and returns the item for chaining.
func (i Item) withMeasured(w, h float32) Item {
	i.SetMeasuredSize(w, h)
	return i
}

func TestMeasureUsesMeasuredSizeWithoutNode(t *testing.T) {
	c := NewContainer()
	c.AddItem(measuredItem(30, 10))
	c.AddItem(measuredItem(40, 20))

	size := c.Measure(Constraint{MaxWidth: 1000, MaxHeight: 1000})
	if size.Width != 70 {
		t.Errorf("width = %g, want 70", size.Width)
	}
	if size.Height != 20 {
		t.Errorf("height = %g, want 20", size.Height)
	}
}

func TestMeasureIncludesMarginsAndGaps(t *testing.T) {
	c := NewContainer()
	c.ColumnGap = 5

	a := measuredItem(100, 50)
	a.MarginLeft, a.MarginRight = 5, 5
	a.MarginTop, a.MarginBottom = 2, 2
	b := measuredItem(100, 50)
	b.MarginLeft, b.MarginRight = 5, 5
	b.MarginTop, b.MarginBottom = 2, 2

	c.AddItem(a)
	c.AddItem(b)

	size := c.Measure(Constraint{MaxWidth: 1000, MaxHeight: 1000})
	if size.Width != 225 {
		t.Errorf("width = %g, want 225", size.Width)
	}
	if size.Height != 54 {
		t.Errorf("height = %g, want 54", size.Height)
	}
}

func TestLayoutAppliesMainMargins(t *testing.T) {
	c := NewContainer()
	a := measuredItem(50, 10)
	a.MarginLeft, a.MarginRight = 10, 10
	b := measuredItem(50, 10)
	b.MarginLeft, b.MarginRight = 10, 10
	c.AddItem(a)
	c.AddItem(b)

	c.Layout(0, 0, 500, 100)

	if c.Items[0].Left != 10 {
		t.Errorf("first item left = %g, want 10", c.Items[0].Left)
	}
	if c.Items[1].Left != 80 {
		t.Errorf("second item left = %g, want 80", c.Items[1].Left)
	}
}

func TestFlexBasisOverridesMeasuredSize(t *testing.T) {
	c := NewContainer()
	item := measuredItem(100, 20)
	item.FlexBasis = 40
	c.AddItem(item)

	c.Layout(0, 0, 500, 100)
	if c.Items[0].Width != 40 {
		t.Errorf("width = %g, want flex-basis 40", c.Items[0].Width)
	}
}

func TestFlexGrowDistributesRemainingSpace(t *testing.T) {
	c := NewContainer()
	a := measuredItem(50, 10)
	a.FlexGrow = 1
	b := measuredItem(50, 10)
	b.FlexGrow = 1
	c.AddItem(a)
	c.AddItem(b)

	c.Layout(0, 0, 200, 100)

	if c.Items[0].Width != 100 || c.Items[1].Width != 100 {
		t.Errorf("grown widths = %g/%g, want 100/100", c.Items[0].Width, c.Items[1].Width)
	}
	if c.Items[1].Left != 100 {
		t.Errorf("second item left = %g, want 100", c.Items[1].Left)
	}
}

func TestFlexShrink(t *testing.T) {
	c := NewContainer()
	c.AddItem(Item{FlexShrink: 1}.withMeasured(100, 10))
	c.AddItem(Item{FlexShrink: 1}.withMeasured(100, 10))

	c.Layout(0, 0, 100, 50)

	if c.Items[0].Width != 50 || c.Items[1].Width != 50 {
		t.Errorf("shrunk widths = %g/%g, want 50/50", c.Items[0].Width, c.Items[1].Width)
	}
}

func TestShrinkRespectsMinWidth(t *testing.T) {
	c := NewContainer()
	a := measuredItem(100, 10)
	a.MinWidth = 60
	a.FlexShrink = 1
	b := measuredItem(100, 10)
	b.FlexShrink = 1
	c.AddItem(a)
	c.AddItem(b)

	c.Layout(0, 0, 120, 50)

	if c.Items[0].Width != 60 {
		t.Errorf("first width = %g, want min 60", c.Items[0].Width)
	}
}

func TestGrowRespectsMaxWidth(t *testing.T) {
	c := NewContainer()
	a := measuredItem(50, 10)
	a.FlexGrow = 1
	a.MaxWidth = 60
	b := measuredItem(50, 10)
	b.FlexGrow = 1
	c.AddItem(a)
	c.AddItem(b)

	c.Layout(0, 0, 300, 50)

	if c.Items[0].Width != 60 {
		t.Errorf("first width = %g, want max 60", c.Items[0].Width)
	}
	// b must absorb the space a could not take (50 + 190).
	if c.Items[1].Width != 240 {
		t.Errorf("second width = %g, want 240", c.Items[1].Width)
	}
}

func TestJustifySpaceBetween(t *testing.T) {
	c := NewContainer()
	c.Justify = JustifySpaceBetween
	c.AddItem(measuredItem(20, 10))
	c.AddItem(measuredItem(20, 10))

	c.Layout(0, 0, 100, 50)

	if c.Items[0].Left != 0 {
		t.Errorf("first left = %g, want 0", c.Items[0].Left)
	}
	if c.Items[1].Left != 80 {
		t.Errorf("second left = %g, want 80", c.Items[1].Left)
	}
}

func TestAlignStretchFillsCrossWhenIndefinite(t *testing.T) {
	c := NewContainer()
	c.AddItem(measuredItem(100, 0)) // no definite cross size

	c.Layout(0, 0, 200, 100)

	if c.Items[0].Height != 100 {
		t.Errorf("height = %g, want stretched 100", c.Items[0].Height)
	}
}

func TestAlignItemsCenter(t *testing.T) {
	c := NewContainer()
	c.Align = AlignCenter
	c.AddItem(measuredItem(100, 20))

	c.Layout(0, 0, 200, 100)

	if c.Items[0].Top != 40 {
		t.Errorf("top = %g, want centered 40", c.Items[0].Top)
	}
}

func TestAlignSelfOverridesContainerAlign(t *testing.T) {
	c := NewContainer()
	c.Align = AlignFlexStart
	item := measuredItem(100, 20)
	item.AlignSelf = AlignFlexEnd
	c.AddItem(item)

	c.Layout(0, 0, 200, 100)

	if c.Items[0].Top != 80 {
		t.Errorf("top = %g, want flex-end 80", c.Items[0].Top)
	}
}

func TestWrapBreaksLinesAndStacksCross(t *testing.T) {
	c := NewContainer()
	c.Wrap = WrapWrap
	// Keep lines at their natural cross size so their stacked positions are
	// not affected by align-content:stretch.
	c.AlignContent = AlignFlexStart
	for i := 0; i < 3; i++ {
		c.AddItem(measuredItem(100, 30))
	}

	// Two items fit on the first line (200 <= 250), the third wraps.
	c.Layout(0, 0, 250, 200)

	if c.Items[0].Top != 0 || c.Items[1].Top != 0 {
		t.Errorf("first line tops = %g/%g, want 0/0", c.Items[0].Top, c.Items[1].Top)
	}
	if c.Items[2].Top != 30 {
		t.Errorf("wrapped item top = %g, want 30", c.Items[2].Top)
	}
}

func TestColumnDirectionUsesVerticalMainAxis(t *testing.T) {
	c := NewContainer()
	c.Direction = DirectionColumn
	c.AddItem(measuredItem(50, 20))
	c.AddItem(measuredItem(60, 30))

	c.Layout(0, 0, 200, 200)

	if c.Items[0].Top != 0 || c.Items[1].Top != 20 {
		t.Errorf("column tops = %g/%g, want 0/20", c.Items[0].Top, c.Items[1].Top)
	}
}

// A content-sized item (auto cross size) stretches to the full line, even
// though its measured cross size is non-zero — matching CSS/Yoga
// align-items:stretch.
func TestAlignStretchAutoCrossFillsLine(t *testing.T) {
	c := NewContainer()
	c.Direction = DirectionColumn
	item := measuredItem(50, 20)
	item.WidthAuto = true
	c.AddItem(item)

	c.Layout(0, 0, 200, 100)

	if c.Items[0].Width != 200 {
		t.Errorf("width = %g, want stretched 200", c.Items[0].Width)
	}
}

// An item with a definite cross size keeps it under align-items:stretch.
func TestAlignStretchLeavesDefiniteCross(t *testing.T) {
	c := NewContainer()
	c.Direction = DirectionColumn
	c.AddItem(measuredItem(50, 20))

	c.Layout(0, 0, 200, 100)

	if c.Items[0].Width != 50 {
		t.Errorf("width = %g, want definite 50", c.Items[0].Width)
	}
}
