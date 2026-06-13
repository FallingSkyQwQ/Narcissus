package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// Checkbox 是复选框组件
type Checkbox struct {
	*BaseWidget

	// 是否选中
	checked bool

	// 选中状态信号
	checkedSignal *reactive.Signal[bool]

	// 标签文本
	label string

	// 变更事件处理器
	onChange func(bool)
}

// NewCheckbox 创建新的复选框组件
func NewCheckbox() *Checkbox {
	cb := &Checkbox{
		BaseWidget:    NewBaseWidget(),
		checked:       false,
		checkedSignal: reactive.NewSignal(false),
	}

	// 设置默认样式
	cb.style.TextColor = ColorBlack
	cb.style.Font = DefaultFont()

	return cb
}

// Checked 设置选中状态（链式调用）
func (cb *Checkbox) Checked(checked bool) *Checkbox {
	oldChecked := cb.checked
	cb.checked = checked
	cb.checkedSignal.Set(checked)

	// 触发变更事件
	if oldChecked != checked && cb.onChange != nil {
		cb.onChange(checked)
	}

	return cb
}

// IsChecked 获取选中状态
func (cb *Checkbox) IsChecked() bool {
	return cb.checked
}

// GetCheckedSignal 获取选中状态信号
func (cb *Checkbox) GetCheckedSignal() *reactive.Signal[bool] {
	return cb.checkedSignal
}

// Label 设置标签文本（链式调用）
func (cb *Checkbox) Label(label string) *Checkbox {
	cb.label = label
	return cb
}

// GetLabel 获取标签文本
func (cb *Checkbox) GetLabel() string {
	return cb.label
}

// OnChange 设置变更事件处理器（链式调用）
func (cb *Checkbox) OnChange(handler func(bool)) *Checkbox {
	cb.onChange = handler
	return cb
}

// Toggle 切换选中状态
func (cb *Checkbox) Toggle() *Checkbox {
	cb.Checked(!cb.checked)
	return cb
}

// FlexGrow 设置 FlexGrow（链式调用）
func (cb *Checkbox) FlexGrow(grow float32) *Checkbox {
	cb.style.FlexGrow = grow
	cb.flexItem.FlexGrow = grow
	return cb
}

// FlexShrink 设置 FlexShrink（链式调用）
func (cb *Checkbox) FlexShrink(shrink float32) *Checkbox {
	cb.style.FlexShrink = shrink
	cb.flexItem.FlexShrink = shrink
	return cb
}

// Margin 设置外边距（链式调用）
func (cb *Checkbox) Margin(margin Insets) *Checkbox {
	cb.style.Margin = margin
	cb.flexItem.MarginTop = margin.Top
	cb.flexItem.MarginRight = margin.Right
	cb.flexItem.MarginBottom = margin.Bottom
	cb.flexItem.MarginLeft = margin.Left
	return cb
}

// TextColor 设置文本颜色（链式调用）
func (cb *Checkbox) TextColor(color Color) *Checkbox {
	cb.style.TextColor = color
	return cb
}

// FontSize 设置字体大小（链式调用）
func (cb *Checkbox) FontSize(size float32) *Checkbox {
	cb.style.Font.Size = size
	return cb
}

// Style 设置样式（链式调用）
func (cb *Checkbox) Style(style *Style) *Checkbox {
	cb.SetStyle(style)
	return cb
}

// Measure 测量复选框大小
func (cb *Checkbox) Measure(constraints flex.Constraint) flex.Size {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// 复选框尺寸：方框 + 标签
	boxSize := cb.style.Font.Size * 1.2
	width := boxSize + 8 // 方框和标签间距
	height := boxSize

	// 如果有标签，计算标签宽度
	if cb.label != "" {
		charWidth := cb.style.Font.Size * 0.6
		labelWidth := float32(len([]rune(cb.label))) * charWidth
		width += labelWidth
		height = max(height, cb.style.Font.Size*cb.style.Font.LineHeight)
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

	cb.flexItem.SetMeasuredSize(width, height)
	cb.measured = true

	return flex.Size{Width: width, Height: height}
}

// Layout 布局复选框
func (cb *Checkbox) Layout(x, y, width, height float32) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.rect.Position.X = x
	cb.rect.Position.Y = y
	cb.rect.Size.Width = width
	cb.rect.Size.Height = height

	cb.flexItem.Left = x
	cb.flexItem.Top = y
	cb.flexItem.Width = width
	cb.flexItem.Height = height
}

// Render 渲染复选框
func (cb *Checkbox) Render() error {
	if !cb.GetVisible() {
		return nil
	}

	// 实际渲染应该调用底层 WinUI API 创建 CheckBox
	// 这里只是占位实现
	return nil
}

// HandleEvent 处理事件
func (cb *Checkbox) HandleEvent(event Event) bool {
	if !cb.GetEnabled() || !cb.GetVisible() {
		return false
	}

	// 处理点击事件
	if event.GetType() == EventClick {
		cb.Toggle()
		return true
	}

	// 委托给基类处理其他事件
	return cb.BaseWidget.HandleEvent(event)
}
