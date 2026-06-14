package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// BorderStyle 边框样式
type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderSolid
	BorderDashed
	BorderDotted
)

// Border 定义边框
type Border struct {
	Width  float32
	Style  BorderStyle
	Color  Color
	Radius float32
}

// Insets 定义内边距或外边距
type Insets struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

// UniformInsets 创建统一的边距
func UniformInsets(value float32) Insets {
	return Insets{Top: value, Right: value, Bottom: value, Left: value}
}

// HorizontalInsets 创建水平边距
func HorizontalInsets(horizontal float32) Insets {
	return Insets{Top: 0, Right: horizontal, Bottom: 0, Left: horizontal}
}

// VerticalInsets 创建垂直边距
func VerticalInsets(vertical float32) Insets {
	return Insets{Top: vertical, Right: 0, Bottom: vertical, Left: 0}
}

// SymmetricInsets 创建对称边距
func SymmetricInsets(vertical, horizontal float32) Insets {
	return Insets{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

// FontWeight 字体粗细
type FontWeight int

const (
	FontWeightNormal FontWeight = iota
	FontWeightBold
	FontWeightLight
	FontWeightMedium
	FontWeightSemiBold
)

// FontStyle 字体样式
type FontStyle int

const (
	FontStyleNormal FontStyle = iota
	FontStyleItalic
)

// TextAlign 文本对齐方式
type TextAlign int

const (
	TextAlignLeft TextAlign = iota
	TextAlignCenter
	TextAlignRight
	TextAlignJustify
)

// Font 字体定义
type Font struct {
	Family     string
	Size       float32
	Weight     FontWeight
	Style      FontStyle
	LineHeight float32
}

// DefaultFont 返回默认字体
func DefaultFont() Font {
	return Font{
		Family:     "Segoe UI",
		Size:       14,
		Weight:     FontWeightNormal,
		Style:      FontStyleNormal,
		LineHeight: 1.5,
	}
}

// Shadow 阴影定义
type Shadow struct {
	Color      Color
	OffsetX    float32
	OffsetY    float32
	BlurRadius float32
	Spread     float32
}

// Style 组件样式定义
type Style struct {
	// 布局属性
	Margin  Insets
	Padding Insets

	// Flex 属性
	FlexGrow   float32
	FlexShrink float32
	FlexBasis  float32
	AlignSelf  flex.Align

	// 尺寸
	Width     float32
	Height    float32
	MinWidth  float32
	MinHeight float32
	MaxWidth  float32
	MaxHeight float32

	// 背景
	BackgroundColor Color
	BackgroundImage string

	// 边框
	Border       Border
	BorderTop    Border
	BorderRight  Border
	BorderBottom Border
	BorderLeft   Border

	// 文本
	Font      Font
	TextColor Color
	TextAlign TextAlign

	// 阴影
	Shadow Shadow

	// 透明度
	Opacity float32

	// 光标
	Cursor string

	// 其他
	Overflow string
}

// NewStyle 创建新样式
func NewStyle() *Style {
	return &Style{
		Margin:          Insets{},
		Padding:         Insets{},
		FlexGrow:        0,
		FlexShrink:      1,
		FlexBasis:       0,
		AlignSelf:       flex.AlignAuto,
		BackgroundColor: ColorTransparent,
		TextColor:       ColorBlack,
		Font:            DefaultFont(),
		TextAlign:       TextAlignLeft,
		Opacity:         1.0,
	}
}

// Clone 克隆样式
func (s *Style) Clone() *Style {
	if s == nil {
		return NewStyle()
	}
	cloned := *s
	return &cloned
}

// WithMargin 设置外边距（链式调用）
func (s *Style) WithMargin(margin Insets) *Style {
	s.Margin = margin
	return s
}

// WithPadding 设置内边距（链式调用）
func (s *Style) WithPadding(padding Insets) *Style {
	s.Padding = padding
	return s
}

// WithFlex 设置 Flex 属性（链式调用）
func (s *Style) WithFlex(grow, shrink, basis float32) *Style {
	s.FlexGrow = grow
	s.FlexShrink = shrink
	s.FlexBasis = basis
	return s
}

// WithFlexGrow 设置 FlexGrow（链式调用）
func (s *Style) WithFlexGrow(grow float32) *Style {
	s.FlexGrow = grow
	return s
}

// WithFlexShrink 设置 FlexShrink（链式调用）
func (s *Style) WithFlexShrink(shrink float32) *Style {
	s.FlexShrink = shrink
	return s
}

// WithBackgroundColor 设置背景色（链式调用）
func (s *Style) WithBackgroundColor(color Color) *Style {
	s.BackgroundColor = color
	return s
}

// WithBackgroundColorToken 使用主题颜色设置背景色（链式调用）
func (s *Style) WithBackgroundColorToken(token ColorToken) *Style {
	s.BackgroundColor = GetThemeColor(token)
	return s
}

// WithTextColor 设置文本颜色（链式调用）
func (s *Style) WithTextColor(color Color) *Style {
	s.TextColor = color
	return s
}

// WithTextColorToken 使用主题颜色设置文本颜色（链式调用）
func (s *Style) WithTextColorToken(token ColorToken) *Style {
	s.TextColor = GetThemeColor(token)
	return s
}

// WithBorder 设置边框（链式调用）
func (s *Style) WithBorder(border Border) *Style {
	s.Border = border
	return s
}

// WithBorderRadius 设置边框圆角（链式调用）
func (s *Style) WithBorderRadius(radius float32) *Style {
	s.Border.Radius = radius
	return s
}

// WithFont 设置字体（链式调用）
func (s *Style) WithFont(font Font) *Style {
	s.Font = font
	return s
}

// WithFontSize 设置字体大小（链式调用）
func (s *Style) WithFontSize(size float32) *Style {
	s.Font.Size = size
	return s
}

// WithOpacity 设置透明度（链式调用）
func (s *Style) WithOpacity(opacity float32) *Style {
	if opacity < 0 {
		opacity = 0
	} else if opacity > 1 {
		opacity = 1
	}
	s.Opacity = opacity
	return s
}

// Merge 合并另一个样式的属性（所有非默认值的属性都会被覆盖）
func (s *Style) Merge(other *Style) *Style {
	if other == nil {
		return s
	}

	// Only copy non-default values
	if other.Margin != (Insets{}) {
		s.Margin = other.Margin
	}
	if other.Padding != (Insets{}) {
		s.Padding = other.Padding
	}
	if other.FlexGrow != 0 {
		s.FlexGrow = other.FlexGrow
	}
	if other.FlexShrink != 1 {
		s.FlexShrink = other.FlexShrink
	}
	if other.FlexBasis != 0 {
		s.FlexBasis = other.FlexBasis
	}
	if other.AlignSelf != flex.AlignAuto {
		s.AlignSelf = other.AlignSelf
	}
	if other.Width != 0 {
		s.Width = other.Width
	}
	if other.Height != 0 {
		s.Height = other.Height
	}
	if other.MinWidth != 0 {
		s.MinWidth = other.MinWidth
	}
	if other.MinHeight != 0 {
		s.MinHeight = other.MinHeight
	}
	if other.MaxWidth != 0 {
		s.MaxWidth = other.MaxWidth
	}
	if other.MaxHeight != 0 {
		s.MaxHeight = other.MaxHeight
	}
	if other.BackgroundColor != (Color{}) {
		s.BackgroundColor = other.BackgroundColor
	}
	if other.BackgroundImage != "" {
		s.BackgroundImage = other.BackgroundImage
	}
	if other.Border != (Border{}) {
		s.Border = other.Border
	}
	if other.BorderTop != (Border{}) {
		s.BorderTop = other.BorderTop
	}
	if other.BorderRight != (Border{}) {
		s.BorderRight = other.BorderRight
	}
	if other.BorderBottom != (Border{}) {
		s.BorderBottom = other.BorderBottom
	}
	if other.BorderLeft != (Border{}) {
		s.BorderLeft = other.BorderLeft
	}
	if other.Font != (Font{}) {
		s.Font = other.Font
	}
	if other.TextColor != (Color{}) {
		s.TextColor = other.TextColor
	}
	if other.TextAlign != 0 {
		s.TextAlign = other.TextAlign
	}
	if other.Shadow != (Shadow{}) {
		s.Shadow = other.Shadow
	}
	if other.Opacity != 1.0 {
		// Clamp opacity to [0, 1] range
		clamped := other.Opacity
		if clamped < 0 {
			clamped = 0
		} else if clamped > 1 {
			clamped = 1
		}
		s.Opacity = clamped
	}
	if other.Cursor != "" {
		s.Cursor = other.Cursor
	}
	if other.Overflow != "" {
		s.Overflow = other.Overflow
	}

	return s
}

// ToFlexProperties 转换为 Flex 属性
func (s *Style) ToFlexProperties() flex.FlexProperties {
	return flex.FlexProperties{
		FlexGrow:   s.FlexGrow,
		FlexShrink: s.FlexShrink,
		FlexBasis:  s.FlexBasis,
		AlignSelf:  s.AlignSelf,
	}
}

// ToFlexItem 转换为 FlexItem
func (s *Style) ToFlexItem() flex.Item {
	return flex.Item{
		Width:         s.Width,
		Height:        s.Height,
		FlexGrow:      s.FlexGrow,
		FlexShrink:    s.FlexShrink,
		FlexBasis:     s.FlexBasis,
		AlignSelf:     s.AlignSelf,
		MarginTop:     s.Margin.Top,
		MarginRight:   s.Margin.Right,
		MarginBottom:  s.Margin.Bottom,
		MarginLeft:    s.Margin.Left,
		PaddingTop:    s.Padding.Top,
		PaddingRight:  s.Padding.Right,
		PaddingBottom: s.Padding.Bottom,
		PaddingLeft:   s.Padding.Left,
	}
}
