package flex

import (
	"math"
	"testing"
)

func TestCalculateLinesNoWrap(t *testing.T) {
	container := NewContainer()
	container.Wrap = WrapNoWrap

	// Add items
	for i := 0; i < 5; i++ {
		item := Item{Node: &mockNode{size: Size{Width: 100, Height: 50}}}
		item.SetMeasuredSize(100, 50)
		container.Items = append(container.Items, item)
	}

	lines := calculateLines(container, 250)

	if len(lines) != 1 {
		t.Errorf("expected 1 line for nowrap, got %d", len(lines))
	}
	if len(lines[0].Items) != 5 {
		t.Errorf("expected 5 items in line, got %d", len(lines[0].Items))
	}
}

func TestCalculateLinesWithWrap(t *testing.T) {
	container := NewContainer()
	container.Wrap = WrapWrap

	// Add items that would exceed container width
	for i := 0; i < 5; i++ {
		item := Item{Node: &mockNode{size: Size{Width: 100, Height: 50}}}
		item.SetMeasuredSize(100, 50)
		container.Items = append(container.Items, item)
	}

	lines := calculateLines(container, 250)

	if len(lines) != 3 { // 2 items + 2 items + 1 item
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestCalculateJustifyOffset(t *testing.T) {
	tests := []struct {
		justify        Justify
		availableSpace float32
		itemCount      int
		expected       float32
	}{
		{JustifyFlexStart, 100, 3, 0},
		{JustifyFlexEnd, 100, 3, 100},
		{JustifyCenter, 100, 3, 50},
		{JustifySpaceBetween, 100, 3, 0},
		{JustifySpaceAround, 100, 3, 16.666667},
		{JustifySpaceEvenly, 100, 3, 25},
	}

	for _, tt := range tests {
		offset := calculateJustifyOffset(tt.justify, tt.availableSpace, tt.itemCount)
		if math.Abs(float64(offset-tt.expected)) > 0.01 {
			t.Errorf("justify %d: expected offset %f, got %f", tt.justify, tt.expected, offset)
		}
	}
}

func TestCalculateJustifyGap(t *testing.T) {
	tests := []struct {
		justify        Justify
		availableSpace float32
		itemCount      int
		expected       float32
	}{
		{JustifyFlexStart, 100, 3, 0},
		{JustifyFlexEnd, 100, 3, 0},
		{JustifyCenter, 100, 3, 0},
		{JustifySpaceBetween, 100, 3, 50},
		{JustifySpaceAround, 100, 3, 33.333333},
		{JustifySpaceEvenly, 100, 3, 25},
	}

	for _, tt := range tests {
		gap := calculateJustifyGap(tt.justify, tt.availableSpace, tt.itemCount)
		if math.Abs(float64(gap-tt.expected)) > 0.01 {
			t.Errorf("justify %d: expected gap %f, got %f", tt.justify, tt.expected, gap)
		}
	}
}

func TestCalculateAlignment(t *testing.T) {
	tests := []struct {
		align    Align
		itemSize float32
		lineSize float32
		expected float32
	}{
		{AlignFlexStart, 50, 100, 0},
		{AlignFlexEnd, 50, 100, 50},
		{AlignCenter, 50, 100, 25},
		{AlignStretch, 50, 100, 0},
	}

	for _, tt := range tests {
		offset := calculateAlignment(tt.align, tt.itemSize, tt.lineSize)
		if offset != tt.expected {
			t.Errorf("align %d: expected offset %f, got %f", tt.align, tt.expected, offset)
		}
	}
}

func TestDistributeExtraSpace(t *testing.T) {
	line := Line{
		Items: []*Item{
			{FlexGrow: 1, Width: 100, Height: 50},
			{FlexGrow: 2, Width: 100, Height: 50},
		},
	}

	distributeExtraSpace(&line, 90, true, DirectionRow)

	// Total flex grow = 3, available space = 90
	// Item 0 gets 1/3 * 90 = 30 extra
	// Item 1 gets 2/3 * 90 = 60 extra
	if line.Items[0].Width != 130 {
		t.Errorf("expected item 0 width 130, got %f", line.Items[0].Width)
	}
	if line.Items[1].Width != 160 {
		t.Errorf("expected item 1 width 160, got %f", line.Items[1].Width)
	}
}

func TestMax(t *testing.T) {
	if max(10, 20) != 20 {
		t.Error("max(10, 20) should be 20")
	}
	if max(20, 10) != 20 {
		t.Error("max(20, 10) should be 20")
	}
}
