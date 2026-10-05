package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// ListSelectionMode selects how many rows a List may have selected.
type ListSelectionMode int

const (
	// ListSingleSelection allows at most one selected row (the default).
	ListSingleSelection ListSelectionMode = iota
	// ListMultipleSelection allows several selected rows.
	ListMultipleSelection
)

// List 是可滚动的选项列表组件。
type List struct {
	*BaseWidget

	items    []string
	selected map[int]bool
	mode     ListSelectionMode

	onSelect func(index int, value string)

	selectionSignal *reactive.Signal[[]int]
}

// NewList 创建空列表。
func NewList() *List {
	l := &List{
		BaseWidget:      NewBaseWidget(),
		selected:        map[int]bool{},
		mode:            ListSingleSelection,
		selectionSignal: reactive.NewSignal([]int(nil)),
	}
	l.SetFocusable(true)
	l.style.BackgroundColor = ColorWhite
	l.style.Border.Width = 1
	l.style.Border.Color = ColorLightBorder
	l.style.Border.Radius = 4
	l.style.Padding = UniformInsets(4)
	return l
}

// Items 设置选项（链式调用）。会清除不再有效的选中项。
func (l *List) Items(items ...string) *List {
	l.items = append([]string(nil), items...)
	for index := range l.selected {
		if index >= len(l.items) {
			delete(l.selected, index)
		}
	}
	l.publishSelection()
	syncWidget(l)
	return l
}

// GetItems 返回选项副本
func (l *List) GetItems() []string {
	return append([]string(nil), l.items...)
}

// SelectionMode 设置选择模式（链式调用）
func (l *List) SelectionMode(mode ListSelectionMode) *List {
	l.mode = mode
	if mode == ListSingleSelection && len(l.selected) > 1 {
		keep := l.firstSelected()
		l.selected = map[int]bool{}
		if keep >= 0 {
			l.selected[keep] = true
		}
		l.publishSelection()
		syncWidget(l)
	}
	return l
}

// GetSelectionMode 返回选择模式
func (l *List) GetSelectionMode() ListSelectionMode { return l.mode }

// Select 选中给定索引（链式调用）
func (l *List) Select(index int) *List {
	if index < 0 || index >= len(l.items) {
		return l
	}
	if l.mode == ListSingleSelection {
		l.selected = map[int]bool{index: true}
	} else {
		l.selected[index] = true
	}
	l.publishSelection()
	syncWidget(l)
	if l.onSelect != nil {
		l.onSelect(index, l.items[index])
	}
	return l
}

// Deselect 取消选中给定索引（链式调用）
func (l *List) Deselect(index int) *List {
	if !l.selected[index] {
		return l
	}
	delete(l.selected, index)
	l.publishSelection()
	syncWidget(l)
	return l
}

// ClearSelection 清除所有选中项（链式调用）
func (l *List) ClearSelection() *List {
	if len(l.selected) == 0 {
		return l
	}
	l.selected = map[int]bool{}
	l.publishSelection()
	syncWidget(l)
	return l
}

// IsSelected 返回给定索引是否被选中
func (l *List) IsSelected(index int) bool { return l.selected[index] }

// SelectedIndices 返回按升序排列的选中索引
func (l *List) SelectedIndices() []int {
	indices := make([]int, 0, len(l.selected))
	for index := range l.selected {
		indices = append(indices, index)
	}
	sortInts(indices)
	return indices
}

// SelectedIndex 返回最小的选中索引，未选中时返回 -1。
func (l *List) SelectedIndex() int {
	indices := l.SelectedIndices()
	if len(indices) == 0 {
		return -1
	}
	return indices[0]
}

// SelectedValue 返回第一个选中项的值，未选中时返回空串。
func (l *List) SelectedValue() string {
	index := l.SelectedIndex()
	if index < 0 || index >= len(l.items) {
		return ""
	}
	return l.items[index]
}

// OnSelect 设置选中处理器（链式调用）
func (l *List) OnSelect(handler func(index int, value string)) *List {
	l.onSelect = handler
	return l
}

// GetSelectionSignal 获取选中索引集合的信号
func (l *List) GetSelectionSignal() *reactive.Signal[[]int] { return l.selectionSignal }

// Width 设置宽度（链式调用）
func (l *List) Width(width float32) *List {
	l.style.Width = width
	return l
}

// Height 设置高度（链式调用）
func (l *List) Height(height float32) *List {
	l.style.Height = height
	return l
}

// Size 设置尺寸（链式调用）
func (l *List) Size(width, height float32) *List {
	l.style.Width = width
	l.style.Height = height
	return l
}

// Style 设置样式（链式调用）
func (l *List) Style(style *Style) *List {
	l.SetStyle(style)
	return l
}

// rowHeight 返回单行的自然高度。
func (l *List) rowHeight() float32 {
	return l.style.Font.Size*l.style.Font.LineHeight + 8
}

// Measure 测量列表大小
func (l *List) Measure(constraints flex.Constraint) flex.Size {
	l.mu.Lock()
	defer l.mu.Unlock()

	width := l.style.Width
	if width <= 0 {
		width = 200
		for _, item := range l.items {
			itemWidth, _ := MeasureText(item, l.style.Font)
			if w := itemWidth + l.style.Padding.Left + l.style.Padding.Right + 8; w > width {
				width = w
			}
		}
	}

	height := l.style.Height
	if height <= 0 {
		rows := len(l.items)
		if rows == 0 {
			rows = 1
		}
		height = float32(rows)*l.rowHeight() + l.style.Padding.Top + l.style.Padding.Bottom
	}

	width = clampAxis(width, constraints.MinWidth, constraints.MaxWidth)
	height = clampAxis(height, constraints.MinHeight, constraints.MaxHeight)

	l.flexItem.SetMeasuredSize(width, height)
	l.measured = true
	return flex.Size{Width: width, Height: height}
}

// Layout 布局列表
func (l *List) Layout(x, y, width, height float32) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.rect.Position.X = x
	l.rect.Position.Y = y
	l.rect.Size.Width = width
	l.rect.Size.Height = height

	l.flexItem.Left = x
	l.flexItem.Top = y
	l.flexItem.Width = width
	l.flexItem.Height = height
}

// Render 渲染列表
func (l *List) Render() error { return renderWidget(l) }

// HandleEvent 处理事件：点击切换（多选）或选中（单选）当前项。
func (l *List) HandleEvent(event Event) bool {
	if !l.GetEnabled() || !l.GetVisible() {
		return false
	}
	if event.GetType() == EventClick {
		return true
	}
	return l.BaseWidget.HandleEvent(event)
}

func (l *List) publishSelection() {
	l.selectionSignal.Set(l.SelectedIndices())
}

func (l *List) firstSelected() int {
	index := -1
	for candidate := range l.selected {
		if index == -1 || candidate < index {
			index = candidate
		}
	}
	return index
}

// sortInts sorts a small slice of indices in ascending order.
func sortInts(values []int) {
	for i := 1; i < len(values); i++ {
		key := values[i]
		j := i - 1
		for j >= 0 && values[j] > key {
			values[j+1] = values[j]
			j--
		}
		values[j+1] = key
	}
}
