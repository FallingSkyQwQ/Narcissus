package flex

import (
	"testing"
)

type mockNode struct {
	size Size
}

func (m *mockNode) Measure(constraints Constraint) Size {
	return m.size
}

func (m *mockNode) GetFlexProperties() FlexProperties {
	return FlexProperties{}
}

func TestNewContainer(t *testing.T) {
	c := NewContainer()
	if c == nil {
		t.Error("NewContainer returned nil")
	}
	if len(c.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(c.Items))
	}
	if c.Direction != DirectionRow {
		t.Error("expected default direction Row")
	}
	if c.Wrap != WrapNoWrap {
		t.Error("expected default wrap NoWrap")
	}
}

func TestContainerAddItem(t *testing.T) {
	c := NewContainer()
	c.AddItem(Item{Width: 100, Height: 50})
	c.AddItem(Item{Width: 80, Height: 60})

	if len(c.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(c.Items))
	}
}

func TestContainerMeasure(t *testing.T) {
	container := NewContainer()
	container.AddItem(Item{
		Node: &mockNode{size: Size{Width: 100, Height: 50}},
	})
	container.AddItem(Item{
		Node: &mockNode{size: Size{Width: 80, Height: 60}},
	})

	size := container.Measure(Constraint{
		MaxWidth:  1000,
		MaxHeight: 1000,
	})

	if size.Width != 180 {
		t.Errorf("expected width 180, got %f", size.Width)
	}
	if size.Height != 60 {
		t.Errorf("expected height 60, got %f", size.Height)
	}
}

func TestContainerLayout(t *testing.T) {
	container := NewContainer()
	container.AddItem(Item{
		Node: &mockNode{size: Size{Width: 100, Height: 50}},
	})
	container.AddItem(Item{
		Node: &mockNode{size: Size{Width: 80, Height: 60}},
	})

	// Measure first
	container.Measure(Constraint{
		MaxWidth:  1000,
		MaxHeight: 1000,
	})

	// Then layout
	container.Layout(0, 0, 1000, 1000)

	if container.Items[0].Left != 0 {
		t.Errorf("expected first item at x=0, got %f", container.Items[0].Left)
	}
	if container.Items[1].Left != 100 {
		t.Errorf("expected second item at x=100, got %f", container.Items[1].Left)
	}
}

func TestItemMeasuredSize(t *testing.T) {
	item := Item{}
	item.SetMeasuredSize(100, 200)

	w, h := item.GetMeasuredSize()
	if w != 100 || h != 200 {
		t.Errorf("expected measured size (100, 200), got (%f, %f)", w, h)
	}
}

func TestFlexProperties(t *testing.T) {
	props := FlexProperties{
		FlexGrow:   1,
		FlexShrink: 0.5,
		FlexBasis:  100,
		AlignSelf:  AlignCenter,
	}

	if props.FlexGrow != 1 {
		t.Errorf("expected FlexGrow 1, got %f", props.FlexGrow)
	}
	if props.FlexShrink != 0.5 {
		t.Errorf("expected FlexShrink 0.5, got %f", props.FlexShrink)
	}
	if props.FlexBasis != 100 {
		t.Errorf("expected FlexBasis 100, got %f", props.FlexBasis)
	}
	if props.AlignSelf != AlignCenter {
		t.Errorf("expected AlignSelf Center, got %d", props.AlignSelf)
	}
}
