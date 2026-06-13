package flex

// Size represents a 2D size
type Size struct {
	Width  float32
	Height float32
}

// Position represents a 2D position
type Position struct {
	X float32
	Y float32
}

// Rect represents a rectangle with position and size
type Rect struct {
	Position Position
	Size     Size
}

// Constraint represents layout constraints passed during measure
type Constraint struct {
	MinWidth  float32
	MaxWidth  float32
	MinHeight float32
	MaxHeight float32
}

// DefiniteSize returns true if both dimensions have definite sizes
func (c Constraint) DefiniteSize() (float32, float32, bool) {
	if c.MinWidth == c.MaxWidth && c.MinHeight == c.MaxHeight {
		return c.MaxWidth, c.MaxHeight, true
	}
	return 0, 0, false
}
