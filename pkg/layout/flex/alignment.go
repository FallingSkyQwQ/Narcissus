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
)
