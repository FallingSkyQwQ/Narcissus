package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// Button 是按钮组件
type Button struct {
	*BaseWidget

	// 按钮文本
	text string

	// 点击事件处理器（简化版，直接存储回调）
	onClick func()
}

// NewButton 创建新的按钮组件
func NewButton() *Button {
	b := &Button{
		BaseWidget: NewBaseWidget(),
	}

	// 设置默认样式
	b.style.BackgroundColor = ColorPrimary
	b.style.TextColor = ColorWhite
	b.style.Font = DefaultFont()
	b.style.Padding = UniformInsets(12)
	b.style.Border.Radius = 4

	return b
}

// Text 设置按钮文本（链式调用）
func (b *Button) Text(text string) *Button {
	b.text = text
	return b
}

// GetText 获取按钮文本
func (b *Button) GetText() string {
	return b.text
}

// OnClick 设置点击事件处理器（链式调用）
func (b *Button) OnClick(handler func()) *Button {
	b.onClick = handler
	return b
}

// BackgroundColor 设置背景色（链式调用）
func (b *Button) BackgroundColor(color Color) *Button {
	b.style.BackgroundColor = color
	return b
}

// BackgroundColorToken 使用主题颜色设置背景色（链式调用）
func (b *Button) BackgroundColorToken(token ColorToken) *Button {
	b.style.BackgroundColor = GetThemeColor(token)
	return b
}

// TextColor 设置文本颜色（链式调用）
func (b *Button) TextColor(color Color) *Button {
	b.style.TextColor = color
	return b
}

// TextColorToken 使用主题颜色设置文本颜色（链式调用）
func (b *Button) TextColorToken(token ColorToken) *Button {
	b.style.TextColor = GetThemeColor(token)
	return b
}

// FontSize 设置字体大小（链式调用）
func (b *Button) FontSize(size float32) *Button {
	b.style.Font.Size = size
	return b
}

// Padding 设置内边距（链式调用）
func (b *Button) Padding(padding Insets) *Button {
	b.style.Padding = padding
	b.flexItem.PaddingTop = padding.Top
	b.flexItem.PaddingRight = padding.Right
	b.flexItem.PaddingBottom = padding.Bottom
	b.flexItem.PaddingLeft = padding.Left
	return b
}

// BorderRadius 设置边框圆角（链式调用）
func (b *Button) BorderRadius(radius float32) *Button {
	b.style.Border.Radius = radius
	return b
}

// FlexGrow 设置 FlexGrow（链式调用）
func (b *Button) FlexGrow(grow float32) *Button {
	b.style.FlexGrow = grow
	b.flexItem.FlexGrow = grow
	return b
}

// FlexShrink 设置 FlexShrink（链式调用）
func (b *Button) FlexShrink(shrink float32) *Button {
	b.style.FlexShrink = shrink
	b.flexItem.FlexShrink = shrink
	return b
}

// Margin 设置外边距（链式调用）
func (b *Button) Margin(margin Insets) *Button {
	b.style.Margin = margin
	b.flexItem.MarginTop = margin.Top
	b.flexItem.MarginRight = margin.Right
	b.flexItem.MarginBottom = margin.Bottom
	b.flexItem.MarginLeft = margin.Left
	return b
}

// Width 设置宽度（链式调用）
func (b *Button) Width(width float32) *Button {
	b.style.Width = width
	return b
}

// Height 设置高度（链式调用）
func (b *Button) Height(height float32) *Button {
	b.style.Height = height
	return b
}

// Size 设置尺寸（链式调用）
func (b *Button) Size(width, height float32) *Button {
	b.style.Width = width
	b.style.Height = height
	return b
}

// Style 设置样式（链式调用）
func (b *Button) Style(style *Style) *Button {
	b.SetStyle(style)
	return b
}

// measureText 测量文本尺寸
func (b *Button) measureText() (width, height float32) {
	if b.text == "" {
		return b.style.Padding.Left + b.style.Padding.Right,
			b.style.Font.Size*b.style.Font.LineHeight + b.style.Padding.Top + b.style.Padding.Bottom
	}

	// 简化计算：假设每个字符平均宽度为字体大小的 0.6 倍
	charWidth := b.style.Font.Size * 0.6
	textWidth := float32(len([]rune(b.text))) * charWidth
	textHeight := b.style.Font.Size * b.style.Font.LineHeight

	width = textWidth + b.style.Padding.Left + b.style.Padding.Right
	height = textHeight + b.style.Padding.Top + b.style.Padding.Bottom

	return width, height
}

// Measure 测量按钮大小
func (b *Button) Measure(constraints flex.Constraint) flex.Size {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 测量文本
	width, height := b.measureText()

	// 应用样式中的固定尺寸
	if b.style.Width > 0 {
		width = b.style.Width
	}
	if b.style.Height > 0 {
		height = b.style.Height
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

	b.flexItem.SetMeasuredSize(width, height)
	b.measured = true

	return flex.Size{Width: width, Height: height}
}

// Layout 布局按钮
func (b *Button) Layout(x, y, width, height float32) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.rect.Position.X = x
	b.rect.Position.Y = y
	b.rect.Size.Width = width
	b.rect.Size.Height = height

	b.flexItem.Left = x
	b.flexItem.Top = y
	b.flexItem.Width = width
	b.flexItem.Height = height
}

// Render 渲染按钮
func (b *Button) Render() error {
	if !b.GetVisible() {
		return nil
	}

	// 实际渲染应该调用底层 WinUI API 创建 Button
	// 这里只是占位实现
	return nil
}

// HandleEvent 处理事件
func (b *Button) HandleEvent(event Event) bool {
	if !b.GetEnabled() || !b.GetVisible() {
		return false
	}

	// 处理点击事件
	if event.GetType() == EventClick {
		// 调用注册的点击处理器
		if b.onClick != nil {
			b.onClick()
			return true
		}
	}

	// 委托给基类处理其他事件
	return b.BaseWidget.HandleEvent(event)
}

// Click 模拟点击按钮
func (b *Button) Click() {
	if b.onClick != nil {
		b.onClick()
	}
}
