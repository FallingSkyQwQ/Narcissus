package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// TabItem 是 Tabs 中的一个页签。
type TabItem struct {
	// Title is the label shown in the header.
	Title string
	// Content is the widget shown when the tab is active.
	Content Widget
}

// Tabs 是页签组件：顶部一排标题按钮，下方显示当前页内容。
//
// It is composed from existing widgets (a header Container of Buttons plus a
// content area) rather than a new native control, so it works on every backend
// without extra backend code.
type Tabs struct {
	*BaseWidget

	items   []TabItem
	current int

	header  *Container
	body    *Container
	buttons []*Button

	onChange func(index int)

	indexSignal *reactive.Signal[int]
}

// NewTabs 创建空的页签组件。
func NewTabs() *Tabs {
	t := &Tabs{
		BaseWidget:  NewBaseWidget(),
		header:      NewContainer(),
		body:        NewContainer(),
		indexSignal: reactive.NewSignal(-1),
	}
	t.header.Direction(flex.DirectionRow).Gap(4)
	t.body.Direction(flex.DirectionColumn)
	t.AddChild(t.header)
	t.AddChild(t.body)
	return t
}

// AddTab 追加一个页签（链式调用）。第一个页签会被自动选中。
func (t *Tabs) AddTab(title string, content Widget) *Tabs {
	index := len(t.items)
	t.items = append(t.items, TabItem{Title: title, Content: content})

	button := NewButton().Text(title)
	button.OnClick(func() { t.Select(index) })
	t.buttons = append(t.buttons, button)
	t.header.Add(button)

	if len(t.items) == 1 {
		t.Select(0)
	}
	return t
}

// SetTabs 用给定的页签替换全部内容（链式调用）
func (t *Tabs) SetTabs(items ...TabItem) *Tabs {
	t.items = nil
	t.buttons = nil
	for _, child := range t.GetChildren() {
		t.RemoveChild(child)
	}

	t.header = NewContainer()
	t.header.Direction(flex.DirectionRow).Gap(4)
	t.body = NewContainer()
	t.body.Direction(flex.DirectionColumn)
	t.AddChild(t.header)
	t.AddChild(t.body)
	t.current = -1
	t.indexSignal.Set(-1)

	for _, item := range items {
		t.AddTab(item.Title, item.Content)
	}
	return t
}

// Select 选中给定索引的页签（链式调用）
func (t *Tabs) Select(index int) *Tabs {
	if index < 0 || index >= len(t.items) {
		return t
	}
	if t.current == index && len(t.body.GetChildren()) > 0 {
		return t
	}

	// Swap the body content.
	for _, child := range t.body.GetChildren() {
		t.body.RemoveChild(child)
	}
	if content := t.items[index].Content; content != nil {
		t.body.AddChild(content)
	}

	t.current = index
	t.indexSignal.Set(index)
	syncWidget(t)

	if t.onChange != nil {
		t.onChange(index)
	}
	return t
}

// CurrentIndex 返回当前选中的页签索引，无页签时为 -1。
func (t *Tabs) CurrentIndex() int { return t.current }

// CurrentTitle 返回当前页签标题
func (t *Tabs) CurrentTitle() string {
	if t.current < 0 || t.current >= len(t.items) {
		return ""
	}
	return t.items[t.current].Title
}

// OnChange 设置页签切换处理器（链式调用）
func (t *Tabs) OnChange(handler func(index int)) *Tabs {
	t.onChange = handler
	return t
}

// GetIndexSignal 获取当前索引的信号
func (t *Tabs) GetIndexSignal() *reactive.Signal[int] { return t.indexSignal }

// Style 设置样式（链式调用）
func (t *Tabs) Style(style *Style) *Tabs {
	t.SetStyle(style)
	return t
}

// Measure 测量页签组件大小：标题行高度加内容区高度。
func (t *Tabs) Measure(constraints flex.Constraint) flex.Size {
	t.mu.Lock()
	defer t.mu.Unlock()

	headerHeight := t.header.Measure(constraints).Height
	inner := constraints
	if inner.MaxHeight > 0 {
		inner.MaxHeight -= headerHeight
		if inner.MaxHeight < 0 {
			inner.MaxHeight = 0
		}
	}
	bodySize := t.body.Measure(inner)

	size := flex.Size{Width: maxf(headerSize(t), bodySize.Width), Height: headerHeight + bodySize.Height}
	size.Width = clampAxis(size.Width, constraints.MinWidth, constraints.MaxWidth)
	size.Height = clampAxis(size.Height, constraints.MinHeight, constraints.MaxHeight)

	t.flexItem.SetMeasuredSize(size.Width, size.Height)
	t.measured = true
	return size
}

// Layout 布局页签组件：标题行在上，内容区占满剩余空间。
func (t *Tabs) Layout(x, y, width, height float32) {
	t.mu.Lock()
	t.rect.Position.X = x
	t.rect.Position.Y = y
	t.rect.Size.Width = width
	t.rect.Size.Height = height
	t.flexItem.Left = x
	t.flexItem.Top = y
	t.flexItem.Width = width
	t.flexItem.Height = height
	t.mu.Unlock()

	headerHeight := t.header.Measure(flex.Constraint{MaxWidth: width}).Height
	t.header.Layout(x, y, width, headerHeight)
	t.body.Layout(x, y+headerHeight, width, height-headerHeight)
}

// Render 渲染页签组件
func (t *Tabs) Render() error { return renderWidget(t) }

// headerSize returns the width the header would like for the current tabs.
func headerSize(t *Tabs) float32 {
	var width float32
	for i, button := range t.buttons {
		if i > 0 {
			width += t.header.columnGap
		}
		width += button.Measure(flex.Constraint{}).Width
	}
	return width
}

func maxf(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
