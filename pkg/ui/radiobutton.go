package ui

import (
	"sync"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// radioGroups tracks the members of each named radio group so selecting one
// button can clear the others. Access is guarded by radioGroupsMu.
var (
	radioGroupsMu sync.Mutex
	radioGroups   = map[string][]*RadioButton{}
)

// RadioButton 是单选按钮，同组内只能选中一个。
type RadioButton struct {
	*BaseWidget

	checked       bool
	group         string
	label         string
	checkedSignal *reactive.Signal[bool]
	onSelect      func()
}

// NewRadioButton 创建单选按钮
func NewRadioButton() *RadioButton {
	r := &RadioButton{
		BaseWidget:    NewBaseWidget(),
		checkedSignal: reactive.NewSignal(false),
	}
	r.SetFocusable(true)
	r.style.TextColor = ColorBlack
	r.style.Font = DefaultFont()
	return r
}

// Group 设置所属分组（链式调用）。同组内在 Checked(true) 时互斥。
func (r *RadioButton) Group(name string) *RadioButton {
	radioGroupsMu.Lock()
	if r.group != "" {
		r.removeFromGroupLocked()
	}
	r.group = name
	if name != "" {
		radioGroups[name] = append(radioGroups[name], r)
	}
	radioGroupsMu.Unlock()
	return r
}

// GetGroup 返回所属分组名
func (r *RadioButton) GetGroup() string { return r.group }

// Label 设置标签文本（链式调用）
func (r *RadioButton) Label(label string) *RadioButton {
	r.label = label
	syncWidget(r)
	return r
}

// GetLabel 返回标签文本
func (r *RadioButton) GetLabel() string { return r.label }

// Checked 设置选中状态（链式调用）。选中时会清除同组的其他按钮。
func (r *RadioButton) Checked(checked bool) *RadioButton {
	if !checked {
		if r.checked {
			r.checked = false
			r.checkedSignal.Set(false)
			syncWidget(r)
		}
		return r
	}

	radioGroupsMu.Lock()
	peers := append([]*RadioButton(nil), radioGroups[r.group]...)
	radioGroupsMu.Unlock()

	becameChecked := !r.checked
	if becameChecked {
		r.checked = true
		r.checkedSignal.Set(true)
	}
	syncWidget(r)

	for _, peer := range peers {
		if peer == r || !peer.checked {
			continue
		}
		peer.checked = false
		peer.checkedSignal.Set(false)
		syncWidget(peer)
	}

	if becameChecked && r.onSelect != nil {
		r.onSelect()
	}
	return r
}

// IsChecked 返回是否选中
func (r *RadioButton) IsChecked() bool { return r.checked }

// Select 选中该按钮
func (r *RadioButton) Select() *RadioButton { return r.Checked(true) }

// OnSelect 设置选中处理器（链式调用）
func (r *RadioButton) OnSelect(handler func()) *RadioButton {
	r.onSelect = handler
	return r
}

// GetCheckedSignal 获取选中状态信号
func (r *RadioButton) GetCheckedSignal() *reactive.Signal[bool] { return r.checkedSignal }

// removeFromGroupLocked removes r from its current group; caller holds the lock.
func (r *RadioButton) removeFromGroupLocked() {
	peers := radioGroups[r.group]
	for i, peer := range peers {
		if peer == r {
			radioGroups[r.group] = append(peers[:i], peers[i+1:]...)
			break
		}
	}
	if len(radioGroups[r.group]) == 0 {
		delete(radioGroups, r.group)
	}
}

// Width 设置宽度（链式调用）
func (r *RadioButton) Width(width float32) *RadioButton {
	r.style.Width = width
	return r
}

// Style 设置样式（链式调用）
func (r *RadioButton) Style(style *Style) *RadioButton {
	r.SetStyle(style)
	return r
}

// Measure 测量单选按钮大小
func (r *RadioButton) Measure(constraints flex.Constraint) flex.Size {
	r.mu.Lock()
	defer r.mu.Unlock()

	boxSize := r.style.Font.Size * 1.2
	width := boxSize + 8
	height := boxSize

	if r.label != "" {
		labelWidth, labelHeight := MeasureText(r.label, r.style.Font)
		width += labelWidth
		if labelHeight > height {
			height = labelHeight
		}
	}
	// Reserve a little vertical room for the toolkit's taller radio control.
	height += 6
	if r.style.Width > 0 {
		width = r.style.Width
	}
	if r.style.Height > 0 {
		height = r.style.Height
	}

	width = clampAxis(width, constraints.MinWidth, constraints.MaxWidth)
	height = clampAxis(height, constraints.MinHeight, constraints.MaxHeight)

	r.flexItem.SetMeasuredSize(width, height)
	r.measured = true
	return flex.Size{Width: width, Height: height}
}

// Layout 布局单选按钮
func (r *RadioButton) Layout(x, y, width, height float32) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rect.Position.X = x
	r.rect.Position.Y = y
	r.rect.Size.Width = width
	r.rect.Size.Height = height

	r.flexItem.Left = x
	r.flexItem.Top = y
	r.flexItem.Width = width
	r.flexItem.Height = height
}

// Render 渲染单选按钮
func (r *RadioButton) Render() error { return renderWidget(r) }

// HandleEvent 处理事件：点击选中。
func (r *RadioButton) HandleEvent(event Event) bool {
	if !r.GetEnabled() || !r.GetVisible() {
		return false
	}
	if event.GetType() == EventClick {
		r.Checked(true)
		return true
	}
	return r.BaseWidget.HandleEvent(event)
}
