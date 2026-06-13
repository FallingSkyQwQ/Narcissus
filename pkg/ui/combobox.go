package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// ComboBox 是下拉选择组件
type ComboBox struct {
	*BaseWidget

	// 选项列表
	items []string

	// 当前选中索引
	selectedIndex int

	// 选中索引的状态信号
	selectedIndexSignal *reactive.Signal[int]

	// 变更事件处理器
	onChange func(int, string)

	// 占位符文本
	placeholder string

	// 是否展开
	isOpen bool
}

// NewComboBox 创建新的下拉选择组件
func NewComboBox() *ComboBox {
	cb := &ComboBox{
		BaseWidget:          NewBaseWidget(),
		items:               make([]string, 0),
		selectedIndex:       -1,
		selectedIndexSignal: reactive.NewSignal(-1),
	}

	// 设置默认样式
	cb.style.BackgroundColor = ColorWhite
	cb.style.TextColor = ColorBlack
	cb.style.Font = DefaultFont()
	cb.style.Padding = UniformInsets(8)
	cb.style.Border.Width = 1
	cb.style.Border.Color = ColorLightBorder
	cb.style.Border.Radius = 4

	return cb
}

// Items 设置选项列表（链式调用）
func (cb *ComboBox) Items(items ...string) *ComboBox {
	cb.items = make([]string, len(items))
	copy(cb.items, items)
	// 如果当前选中索引超出范围，重置为 -1
	if cb.selectedIndex >= len(cb.items) {
		cb.selectedIndex = -1
		cb.selectedIndexSignal.Set(-1)
	}
	return cb
}

// GetItems 获取选项列表
func (cb *ComboBox) GetItems() []string {
	result := make([]string, len(cb.items))
	copy(result, cb.items)
	return result
}

// SelectedIndex 设置选中索引（链式调用）
func (cb *ComboBox) SelectedIndex(index int) *ComboBox {
	if index < -1 || index >= len(cb.items) {
		index = -1
	}

	oldIndex := cb.selectedIndex
	cb.selectedIndex = index
	cb.selectedIndexSignal.Set(index)

	// 触发变更事件
	if oldIndex != index && cb.onChange != nil {
		var selectedValue string
		if index >= 0 && index < len(cb.items) {
			selectedValue = cb.items[index]
		}
		cb.onChange(index, selectedValue)
	}

	return cb
}

// GetSelectedIndex 获取选中索引
func (cb *ComboBox) GetSelectedIndex() int {
	return cb.selectedIndex
}

// GetSelectedValue 获取选中的值
func (cb *ComboBox) GetSelectedValue() string {
	if cb.selectedIndex >= 0 && cb.selectedIndex < len(cb.items) {
		return cb.items[cb.selectedIndex]
	}
	return ""
}

// GetSelectedIndexSignal 获取选中索引的状态信号
func (cb *ComboBox) GetSelectedIndexSignal() *reactive.Signal[int] {
	return cb.selectedIndexSignal
}

// Placeholder 设置占位符（链式调用）
func (cb *ComboBox) Placeholder(placeholder string) *ComboBox {
	cb.placeholder = placeholder
	return cb
}

// GetPlaceholder 获取占位符
func (cb *ComboBox) GetPlaceholder() string {
	return cb.placeholder
}

// OnChange 设置变更事件处理器（链式调用）
func (cb *ComboBox) OnChange(handler func(index int, value string)) *ComboBox {
	cb.onChange = handler
	return cb
}

// SelectItem 选择指定项
func (cb *ComboBox) SelectItem(index int) *ComboBox {
	cb.SelectedIndex(index)
	cb.isOpen = false
	return cb
}

// SelectValue 选择指定值
func (cb *ComboBox) SelectValue(value string) *ComboBox {
	for i, item := range cb.items {
		if item == value {
			cb.SelectItem(i)
			return cb
		}
	}
	return cb
}

// Open 展开下拉框
func (cb *ComboBox) Open() *ComboBox {
	cb.isOpen = true
	return cb
}

// Close 收起下拉框
func (cb *ComboBox) Close() *ComboBox {
	cb.isOpen = false
	return cb
}

// IsOpen 返回下拉框是否展开
func (cb *ComboBox) IsOpen() bool {
	return cb.isOpen
}

// Toggle 切换下拉框展开状态
func (cb *ComboBox) Toggle() *ComboBox {
	cb.isOpen = !cb.isOpen
	return cb
}

// FlexGrow 设置 FlexGrow（链式调用）
func (cb *ComboBox) FlexGrow(grow float32) *ComboBox {
	cb.style.FlexGrow = grow
	cb.flexItem.FlexGrow = grow
	return cb
}

// FlexShrink 设置 FlexShrink（链式调用）
func (cb *ComboBox) FlexShrink(shrink float32) *ComboBox {
	cb.style.FlexShrink = shrink
	cb.flexItem.FlexShrink = shrink
	return cb
}

// Margin 设置外边距（链式调用）
func (cb *ComboBox) Margin(margin Insets) *ComboBox {
	cb.style.Margin = margin
	cb.flexItem.MarginTop = margin.Top
	cb.flexItem.MarginRight = margin.Right
	cb.flexItem.MarginBottom = margin.Bottom
	cb.flexItem.MarginLeft = margin.Left
	return cb
}

// Width 设置宽度（链式调用）
func (cb *ComboBox) Width(width float32) *ComboBox {
	cb.style.Width = width
	return cb
}

// BackgroundColor 设置背景色（链式调用）
func (cb *ComboBox) BackgroundColor(color Color) *ComboBox {
	cb.style.BackgroundColor = color
	return cb
}

// TextColor 设置文本颜色（链式调用）
func (cb *ComboBox) TextColor(color Color) *ComboBox {
	cb.style.TextColor = color
	return cb
}

// FontSize 设置字体大小（链式调用）
func (cb *ComboBox) FontSize(size float32) *ComboBox {
	cb.style.Font.Size = size
	return cb
}

// BorderRadius 设置边框圆角（链式调用）
func (cb *ComboBox) BorderRadius(radius float32) *ComboBox {
	cb.style.Border.Radius = radius
	return cb
}

// Style 设置样式（链式调用）
func (cb *ComboBox) Style(style *Style) *ComboBox {
	cb.SetStyle(style)
	return cb
}

// Measure 测量下拉框大小
func (cb *ComboBox) Measure(constraints flex.Constraint) flex.Size {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	var width, height float32

	// 使用样式中设置的宽度
	if cb.style.Width > 0 {
		width = cb.style.Width
	} else {
		// 计算最小宽度：最长选项的宽度
		maxTextWidth := float32(0)
		charWidth := cb.style.Font.Size * 0.6
		for _, item := range cb.items {
			textWidth := float32(len([]rune(item))) * charWidth
			if textWidth > maxTextWidth {
				maxTextWidth = textWidth
			}
		}
		// 如果没有选项，使用默认宽度
		if maxTextWidth == 0 {
			maxTextWidth = 100
		}
		width = maxTextWidth + cb.style.Padding.Left + cb.style.Padding.Right + 30 // 30 是下拉箭头空间
	}

	// 计算高度
	height = cb.style.Font.Size*cb.style.Font.LineHeight + cb.style.Padding.Top + cb.style.Padding.Bottom

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

// Layout 布局下拉框
func (cb *ComboBox) Layout(x, y, width, height float32) {
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

// Render 渲染下拉框
func (cb *ComboBox) Render() error {
	if !cb.GetVisible() {
		return nil
	}

	// 实际渲染应该调用底层 WinUI API 创建 ComboBox
	// 这里只是占位实现
	return nil
}

// HandleEvent 处理事件
func (cb *ComboBox) HandleEvent(event Event) bool {
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
