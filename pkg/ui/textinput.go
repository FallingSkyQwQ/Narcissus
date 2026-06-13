package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// TextInputType 定义输入框类型
type TextInputType int

const (
	// TextInputTypeText 普通文本
	TextInputTypeText TextInputType = iota
	// TextInputTypePassword 密码
	TextInputTypePassword
	// TextInputTypeEmail 邮箱
	TextInputTypeEmail
	// TextInputTypeNumber 数字
	TextInputTypeNumber
	// TextInputTypeURL URL
	TextInputTypeURL
	// TextInputTypeMultiline 多行文本
	TextInputTypeMultiline
)

// TextInput 是文本输入组件
type TextInput struct {
	*BaseWidget

	// 输入类型
	inputType TextInputType

	// 占位符文本
	placeholder string

	// 当前值
	value string

	// 值的状态信号
	valueSignal *reactive.Signal[string]

	// 变更事件处理器
	onChange func(string)

	// 是否只读
	readOnly bool

	// 最大长度
	maxLength int
}

// NewTextInput 创建新的文本输入组件
func NewTextInput() *TextInput {
	ti := &TextInput{
		BaseWidget:  NewBaseWidget(),
		inputType:   TextInputTypeText,
		valueSignal: reactive.NewSignal(""),
		maxLength:   0, // 0 表示无限制
	}

	// 设置默认样式
	ti.style.BackgroundColor = ColorWhite
	ti.style.TextColor = ColorBlack
	ti.style.Font = DefaultFont()
	ti.style.Padding = UniformInsets(8)
	ti.style.Border.Width = 1
	ti.style.Border.Color = ColorLightBorder
	ti.style.Border.Radius = 4

	return ti
}

// Value 设置值（链式调用）
func (ti *TextInput) Value(value string) *TextInput {
	oldValue := ti.value
	ti.value = value
	ti.valueSignal.Set(value)

	// 触发变更事件
	if oldValue != value && ti.onChange != nil {
		ti.onChange(value)
	}

	return ti
}

// GetValue 获取当前值
func (ti *TextInput) GetValue() string {
	return ti.value
}

// GetValueSignal 获取值的状态信号
func (ti *TextInput) GetValueSignal() *reactive.Signal[string] {
	return ti.valueSignal
}

// Placeholder 设置占位符（链式调用）
func (ti *TextInput) Placeholder(placeholder string) *TextInput {
	ti.placeholder = placeholder
	return ti
}

// GetPlaceholder 获取占位符
func (ti *TextInput) GetPlaceholder() string {
	return ti.placeholder
}

// InputType 设置输入类型（链式调用）
func (ti *TextInput) InputType(inputType TextInputType) *TextInput {
	ti.inputType = inputType
	return ti
}

// GetInputType 获取输入类型
func (ti *TextInput) GetInputType() TextInputType {
	return ti.inputType
}

// OnChange 设置变更事件处理器（链式调用）
func (ti *TextInput) OnChange(handler func(string)) *TextInput {
	ti.onChange = handler
	return ti
}

// ReadOnly 设置只读状态（链式调用）
func (ti *TextInput) ReadOnly(readOnly bool) *TextInput {
	ti.readOnly = readOnly
	return ti
}

// IsReadOnly 返回是否只读
func (ti *TextInput) IsReadOnly() bool {
	return ti.readOnly
}

// MaxLength 设置最大长度（链式调用）
func (ti *TextInput) MaxLength(maxLength int) *TextInput {
	ti.maxLength = maxLength
	return ti
}

// GetMaxLength 获取最大长度
func (ti *TextInput) GetMaxLength() int {
	return ti.maxLength
}

// Width 设置宽度（链式调用）
func (ti *TextInput) Width(width float32) *TextInput {
	ti.style.Width = width
	return ti
}

// Height 设置高度（链式调用）
func (ti *TextInput) Height(height float32) *TextInput {
	ti.style.Height = height
	return ti
}

// Size 设置尺寸（链式调用）
func (ti *TextInput) Size(width, height float32) *TextInput {
	ti.style.Width = width
	ti.style.Height = height
	return ti
}

// FlexGrow 设置 FlexGrow（链式调用）
func (ti *TextInput) FlexGrow(grow float32) *TextInput {
	ti.style.FlexGrow = grow
	ti.flexItem.FlexGrow = grow
	return ti
}

// FlexShrink 设置 FlexShrink（链式调用）
func (ti *TextInput) FlexShrink(shrink float32) *TextInput {
	ti.style.FlexShrink = shrink
	ti.flexItem.FlexShrink = shrink
	return ti
}

// Margin 设置外边距（链式调用）
func (ti *TextInput) Margin(margin Insets) *TextInput {
	ti.style.Margin = margin
	ti.flexItem.MarginTop = margin.Top
	ti.flexItem.MarginRight = margin.Right
	ti.flexItem.MarginBottom = margin.Bottom
	ti.flexItem.MarginLeft = margin.Left
	return ti
}

// Padding 设置内边距（链式调用）
func (ti *TextInput) Padding(padding Insets) *TextInput {
	ti.style.Padding = padding
	return ti
}

// BorderRadius 设置边框圆角（链式调用）
func (ti *TextInput) BorderRadius(radius float32) *TextInput {
	ti.style.Border.Radius = radius
	return ti
}

// BackgroundColor 设置背景色（链式调用）
func (ti *TextInput) BackgroundColor(color Color) *TextInput {
	ti.style.BackgroundColor = color
	return ti
}

// TextColor 设置文本颜色（链式调用）
func (ti *TextInput) TextColor(color Color) *TextInput {
	ti.style.TextColor = color
	return ti
}

// Style 设置样式（链式调用）
func (ti *TextInput) Style(style *Style) *TextInput {
	ti.SetStyle(style)
	return ti
}

// SetText 设置文本（触发变更事件）
func (ti *TextInput) SetText(text string) {
	// 检查最大长度限制
	if ti.maxLength > 0 && len([]rune(text)) > ti.maxLength {
		runes := []rune(text)
		text = string(runes[:ti.maxLength])
	}

	oldValue := ti.value
	ti.value = text
	ti.valueSignal.Set(text)

	// 触发变更事件
	if oldValue != text && ti.onChange != nil {
		ti.onChange(text)
	}
}

// Clear 清空输入
func (ti *TextInput) Clear() {
	ti.SetText("")
}

// Focus 聚焦输入框
func (ti *TextInput) Focus() {
	// 实际实现应该调用底层 WinUI API
}

// Blur 失焦输入框
func (ti *TextInput) Blur() {
	// 实际实现应该调用底层 WinUI API
}

// Measure 测量输入框大小
func (ti *TextInput) Measure(constraints flex.Constraint) flex.Size {
	ti.mu.Lock()
	defer ti.mu.Unlock()

	var width, height float32

	// 使用样式中设置的尺寸
	if ti.style.Width > 0 {
		width = ti.style.Width
	} else {
		// 默认宽度
		width = 200
	}

	if ti.style.Height > 0 {
		height = ti.style.Height
	} else {
		// 根据输入类型计算默认高度
		if ti.inputType == TextInputTypeMultiline {
			height = ti.style.Font.Size*ti.style.Font.LineHeight*3 + ti.style.Padding.Top + ti.style.Padding.Bottom
		} else {
			height = ti.style.Font.Size*ti.style.Font.LineHeight + ti.style.Padding.Top + ti.style.Padding.Bottom
		}
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

	ti.flexItem.SetMeasuredSize(width, height)
	ti.measured = true

	return flex.Size{Width: width, Height: height}
}

// Layout 布局输入框
func (ti *TextInput) Layout(x, y, width, height float32) {
	ti.mu.Lock()
	defer ti.mu.Unlock()

	ti.rect.Position.X = x
	ti.rect.Position.Y = y
	ti.rect.Size.Width = width
	ti.rect.Size.Height = height

	ti.flexItem.Left = x
	ti.flexItem.Top = y
	ti.flexItem.Width = width
	ti.flexItem.Height = height
}

// Render 渲染输入框
func (ti *TextInput) Render() error {
	if !ti.GetVisible() {
		return nil
	}

	// 实际渲染应该调用底层 WinUI API 创建 TextBox
	// 这里只是占位实现
	return nil
}

// HandleEvent 处理事件
func (ti *TextInput) HandleEvent(event Event) bool {
	if !ti.GetEnabled() || !ti.GetVisible() {
		return false
	}

	switch event.GetType() {
	case EventFocus:
		// 处理聚焦事件
		return true
	case EventBlur:
		// 处理失焦事件
		return true
	case EventInput:
		// 处理输入事件
		if inputEvent, ok := event.(*InputEvent); ok {
			ti.SetText(inputEvent.Value)
			return true
		}
	case EventChange:
		// 处理变更事件
		if changeEvent, ok := event.(*ChangeEvent); ok {
			ti.SetText(changeEvent.Value)
			return true
		}
	}

	// 委托给基类处理其他事件
	return ti.BaseWidget.HandleEvent(event)
}
