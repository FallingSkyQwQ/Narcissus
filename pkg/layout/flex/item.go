package flex

// Item represents a single flex item
type Item struct {
	// Layout properties
	Width  float32
	Height float32
	Left   float32
	Top    float32

	// Flex properties
	FlexGrow   float32
	FlexShrink float32
	FlexBasis  float32

	// Alignment override
	AlignSelf Align

	// Margins
	MarginTop    float32
	MarginRight  float32
	MarginBottom float32
	MarginLeft   float32

	// Padding
	PaddingTop    float32
	PaddingRight  float32
	PaddingBottom float32
	PaddingLeft   float32

	// Measured size cache
	measuredWidth  float32
	measuredHeight float32

	// Reference to the actual widget/node
	Node LayoutNode
}

// LayoutNode is the interface that layoutable nodes must implement
type LayoutNode interface {
	// Measure returns the desired size given constraints
	Measure(constraints Constraint) Size
	// GetFlexProperties returns the flex properties for this node
	GetFlexProperties() FlexProperties
}

// FlexProperties holds the flex-related CSS properties
type FlexProperties struct {
	FlexGrow   float32
	FlexShrink float32
	FlexBasis  float32
	AlignSelf  Align
}

// SetMeasuredSize sets the measured size
func (i *Item) SetMeasuredSize(width, height float32) {
	i.measuredWidth = width
	i.measuredHeight = height
}

// GetMeasuredSize returns the measured size
func (i *Item) GetMeasuredSize() (float32, float32) {
	return i.measuredWidth, i.measuredHeight
}
