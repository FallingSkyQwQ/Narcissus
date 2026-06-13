package flex

// Direction defines the main axis direction
type Direction int

const (
	DirectionRow Direction = iota
	DirectionRowReverse
	DirectionColumn
	DirectionColumnReverse
)

// Wrap defines how items wrap
type Wrap int

const (
	WrapNoWrap Wrap = iota
	WrapWrap
	WrapWrapReverse
)
