package flex

// Justify defines alignment along the main axis
type Justify int

const (
	JustifyFlexStart Justify = iota
	JustifyFlexEnd
	JustifyCenter
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// Align defines alignment along the cross axis
type Align int

const (
	AlignAuto Align = iota
	AlignFlexStart
	AlignFlexEnd
	AlignCenter
	AlignStretch
	AlignBaseline

	// The space-* values are only meaningful for align-content; they
	// distribute leftover cross-axis space between lines.
	AlignSpaceBetween
	AlignSpaceAround
	AlignSpaceEvenly
)
