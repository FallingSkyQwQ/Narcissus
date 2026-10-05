package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// Text 是文本显示组件
type Text struct {
	*BaseWidget

	// 文本内容
	text string

	// 文本测量缓存
	measuredWidth  float32
	measuredHeight float32
}

// NewText 创建新的文本组件
func NewText(text string) *Text {
	t := &Text{
		BaseWidget: NewBaseWidget(),
		text:       text,
	}

	// 设置默认样式
	t.style.TextColor = ColorBlack
	t.style.Font = DefaultFont()

	return t
}

// Text 设置文本内容（链式调用）
func (t *Text) Text(text string) *Text {
	t.text = text
	// 清除测量缓存，下次测量时重新计算
	t.measured = false
	syncWidget(t)
	return t
}

// GetText 获取文本内容
func (t *Text) GetText() string {
	return t.text
}

// TextColor 设置文本颜色（链式调用）
func (t *Text) TextColor(color Color) *Text {
	t.style.TextColor = color
	syncWidget(t)
	return t
}

// TextColorToken 使用主题颜色设置文本颜色（链式调用）
func (t *Text) TextColorToken(token ColorToken) *Text {
	t.style.TextColor = GetThemeColor(token)
	return t
}

// FontSize 设置字体大小（链式调用）
func (t *Text) FontSize(size float32) *Text {
	t.style.Font.Size = size
	t.measured = false
	syncWidget(t)
	return t
}

// FontWeight 设置字体粗细（链式调用）
func (t *Text) FontWeight(weight FontWeight) *Text {
	t.style.Font.Weight = weight
	t.measured = false
	syncWidget(t)
	return t
}

// FontFamily 设置字体族（链式调用）
func (t *Text) FontFamily(family string) *Text {
	t.style.Font.Family = family
	t.measured = false
	return t
}

// TextAlign 设置文本对齐方式（链式调用）
func (t *Text) TextAlign(align TextAlign) *Text {
	t.style.TextAlign = align
	return t
}

// LineHeight 设置行高（链式调用）
func (t *Text) LineHeight(height float32) *Text {
	t.style.Font.LineHeight = height
	t.measured = false
	return t
}

// FlexGrow 设置 FlexGrow（链式调用）
func (t *Text) FlexGrow(grow float32) *Text {
	t.style.FlexGrow = grow
	t.flexItem.FlexGrow = grow
	return t
}

// FlexShrink 设置 FlexShrink（链式调用）
func (t *Text) FlexShrink(shrink float32) *Text {
	t.style.FlexShrink = shrink
	t.flexItem.FlexShrink = shrink
	return t
}

// Margin 设置外边距（链式调用）
func (t *Text) Margin(margin Insets) *Text {
	t.style.Margin = margin
	t.flexItem.MarginTop = margin.Top
	t.flexItem.MarginRight = margin.Right
	t.flexItem.MarginBottom = margin.Bottom
	t.flexItem.MarginLeft = margin.Left
	return t
}

// Padding 设置内边距（链式调用）
func (t *Text) Padding(padding Insets) *Text {
	t.style.Padding = padding
	t.flexItem.PaddingTop = padding.Top
	t.flexItem.PaddingRight = padding.Right
	t.flexItem.PaddingBottom = padding.Bottom
	t.flexItem.PaddingLeft = padding.Left
	return t
}

// Style 设置样式（链式调用）
func (t *Text) Style(style *Style) *Text {
	t.SetStyle(style)
	return t
}

// measureText 返回当前字体下文本的自然尺寸。它委托给已安装的
// TextMeasurer，因此后端可以接入系统真实的字体度量。
func (t *Text) measureText() (width, height float32) {
	return MeasureText(t.text, t.style.Font)
}

// Measure 测量文本组件大小
func (t *Text) Measure(constraints flex.Constraint) flex.Size {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 缓存自然尺寸（未施加约束），约束在每次调用时重新应用，避免同一个
	// 文本在不同约束下复用被裁剪过的旧尺寸。
	if !t.measured {
		t.measuredWidth, t.measuredHeight = t.measureText()
		t.measured = true
	}
	width, height := t.measuredWidth, t.measuredHeight

	// Text wider than the available space wraps: re-measure against the
	// maximum width so the height covers every wrapped line. Without this the
	// native label would wrap while the framework reserved a single line,
	// making following widgets overlap.
	if constraints.MaxWidth > 0 && width > constraints.MaxWidth {
		width, height = MeasureTextWrapped(t.text, t.style.Font, constraints.MaxWidth)
	}

	// 应用约束
	if constraints.MaxWidth > 0 && width > constraints.MaxWidth {
		width = constraints.MaxWidth
	}
	if constraints.MinWidth > 0 && width < constraints.MinWidth {
		width = constraints.MinWidth
	}
	if constraints.MaxHeight > 0 && height > constraints.MaxHeight {
		height = constraints.MaxHeight
	}
	if constraints.MinHeight > 0 && height < constraints.MinHeight {
		height = constraints.MinHeight
	}

	t.flexItem.SetMeasuredSize(width, height)

	return flex.Size{Width: width, Height: height}
}

// Layout 布局文本组件
func (t *Text) Layout(x, y, width, height float32) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.rect.Position.X = x
	t.rect.Position.Y = y
	t.rect.Size.Width = width
	t.rect.Size.Height = height

	t.flexItem.Left = x
	t.flexItem.Top = y
	t.flexItem.Width = width
	t.flexItem.Height = height
}

// Render 渲染文本组件：将当前状态同步到后端已创建的原生控件
func (t *Text) Render() error {
	return renderWidget(t)
}
