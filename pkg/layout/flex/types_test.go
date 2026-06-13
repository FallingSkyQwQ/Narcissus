package flex

import (
	"testing"
)

func TestConstraintDefiniteSize(t *testing.T) {
	c := Constraint{
		MinWidth:  100,
		MaxWidth:  100,
		MinHeight: 200,
		MaxHeight: 200,
	}
	w, h, ok := c.DefiniteSize()
	if !ok {
		t.Error("expected definite size")
	}
	if w != 100 || h != 200 {
		t.Errorf("expected (100, 200), got (%f, %f)", w, h)
	}
}

func TestConstraintIndefiniteSize(t *testing.T) {
	c := Constraint{
		MinWidth:  0,
		MaxWidth:  100,
		MinHeight: 0,
		MaxHeight: 200,
	}
	_, _, ok := c.DefiniteSize()
	if ok {
		t.Error("expected indefinite size")
	}
}

func TestSize(t *testing.T) {
	s := Size{Width: 100, Height: 200}
	if s.Width != 100 {
		t.Errorf("expected width 100, got %f", s.Width)
	}
	if s.Height != 200 {
		t.Errorf("expected height 200, got %f", s.Height)
	}
}

func TestPosition(t *testing.T) {
	p := Position{X: 10, Y: 20}
	if p.X != 10 {
		t.Errorf("expected X 10, got %f", p.X)
	}
	if p.Y != 20 {
		t.Errorf("expected Y 20, got %f", p.Y)
	}
}

func TestRect(t *testing.T) {
	r := Rect{
		Position: Position{X: 10, Y: 20},
		Size:     Size{Width: 100, Height: 200},
	}
	if r.Position.X != 10 || r.Position.Y != 20 {
		t.Error("position not set correctly")
	}
	if r.Size.Width != 100 || r.Size.Height != 200 {
		t.Error("size not set correctly")
	}
}
