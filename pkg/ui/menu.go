package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// Menu 是一个展开菜单的按钮组件。
type Menu struct {
	*BaseWidget

	label string
	items []MenuItem

	onSelect func(index int, item MenuItem)
}

// NewMenu 创建菜单组件，label 为按钮上的文本。
func NewMenu(label string) *Menu {
	m := &Menu{
		BaseWidget: NewBaseWidget(),
		label:      label,
	}
	m.SetFocusable(true)
	m.style.TextColor = ColorBlack
	m.style.Font = DefaultFont()
	return m
}

// Items 设置菜单项（链式调用）
func (m *Menu) Items(items ...MenuItem) *Menu {
	m.items = append([]MenuItem(nil), items...)
	syncWidget(m)
	return m
}

// GetItems 返回菜单项副本
func (m *Menu) GetItems() []MenuItem {
	return append([]MenuItem(nil), m.items...)
}

// GetLabel 返回按钮文本
func (m *Menu) GetLabel() string { return m.label }

// AddItem 追加一个可点击菜单项（链式调用）
func (m *Menu) AddItem(label string) *Menu {
	m.items = append(m.items, MenuItem{Label: label, Enabled: true})
	syncWidget(m)
	return m
}

// AddSeparator 追加一个分隔线（链式调用）
func (m *Menu) AddSeparator() *Menu {
	m.items = append(m.items, MenuItem{Separator: true})
	syncWidget(m)
	return m
}

// OnSelect 设置菜单项选择处理器（链式调用）
func (m *Menu) OnSelect(handler func(index int, item MenuItem)) *Menu {
	m.onSelect = handler
	return m
}

// Select 以编程方式触发某个菜单项，主要用于测试与无障碍调用。
func (m *Menu) Select(index int) {
	if index < 0 || index >= len(m.items) {
		return
	}
	item := m.items[index]
	if item.Separator || !item.Enabled {
		return
	}
	if m.onSelect != nil {
		m.onSelect(index, item)
	}
}

// Style 设置样式（链式调用）
func (m *Menu) Style(style *Style) *Menu {
	m.SetStyle(style)
	return m
}

// Measure 测量菜单按钮大小
func (m *Menu) Measure(constraints flex.Constraint) flex.Size {
	m.mu.Lock()
	defer m.mu.Unlock()

	labelWidth, labelHeight := MeasureText(m.label, m.style.Font)
	// Room for the drop-down arrow.
	width := labelWidth + 24 + m.style.Padding.Left + m.style.Padding.Right
	height := labelHeight + m.style.Padding.Top + m.style.Padding.Bottom

	if m.style.Width > 0 {
		width = m.style.Width
	}
	if m.style.Height > 0 {
		height = m.style.Height
	}

	width = clampAxis(width, constraints.MinWidth, constraints.MaxWidth)
	height = clampAxis(height, constraints.MinHeight, constraints.MaxHeight)

	m.flexItem.SetMeasuredSize(width, height)
	m.measured = true
	return flex.Size{Width: width, Height: height}
}

// Layout 布局菜单按钮
func (m *Menu) Layout(x, y, width, height float32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rect.Position.X = x
	m.rect.Position.Y = y
	m.rect.Size.Width = width
	m.rect.Size.Height = height

	m.flexItem.Left = x
	m.flexItem.Top = y
	m.flexItem.Width = width
	m.flexItem.Height = height
}

// Render 渲染菜单按钮
func (m *Menu) Render() error { return renderWidget(m) }
