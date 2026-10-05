package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// Switch 是开关（切换）组件。
type Switch struct {
	*BaseWidget

	checked       bool
	checkedSignal *reactive.Signal[bool]
	onToggle      func(bool)
	label         string
}

// NewSwitch 创建开关组件
func NewSwitch() *Switch {
	s := &Switch{
		BaseWidget:    NewBaseWidget(),
		checkedSignal: reactive.NewSignal(false),
	}
	s.SetFocusable(true)
	s.style.TextColor = ColorBlack
	s.style.Font = DefaultFont()
	return s
}

// Checked 设置开关状态（链式调用）
func (s *Switch) Checked(checked bool) *Switch {
	old := s.checked
	if old != checked {
		s.checked = checked
		s.checkedSignal.Set(checked)
		syncWidget(s)
		if s.onToggle != nil {
			s.onToggle(checked)
		}
	}
	return s
}

// IsChecked 返回开关状态
func (s *Switch) IsChecked() bool { return s.checked }

// Toggle 切换开关状态
func (s *Switch) Toggle() *Switch { return s.Checked(!s.checked) }

// OnToggle 设置状态变更处理器（链式调用）
func (s *Switch) OnToggle(handler func(bool)) *Switch {
	s.onToggle = handler
	return s
}

// Label 设置可选标签文本（链式调用）
func (s *Switch) Label(label string) *Switch {
	s.label = label
	syncWidget(s)
	return s
}

// GetLabel 返回标签文本
func (s *Switch) GetLabel() string { return s.label }

// GetCheckedSignal 获取状态信号
func (s *Switch) GetCheckedSignal() *reactive.Signal[bool] { return s.checkedSignal }

// Width 设置宽度（链式调用）
func (s *Switch) Width(width float32) *Switch {
	s.style.Width = width
	return s
}

// Style 设置样式（链式调用）
func (s *Switch) Style(style *Style) *Switch {
	s.SetStyle(style)
	return s
}

// Measure 测量开关大小
func (s *Switch) Measure(constraints flex.Constraint) flex.Size {
	s.mu.Lock()
	defer s.mu.Unlock()

	width := float32(44) // track width
	height := float32(24)

	if s.label != "" {
		labelWidth, labelHeight := MeasureText(s.label, s.style.Font)
		width += 8 + labelWidth
		if labelHeight > height {
			height = labelHeight
		}
	}
	if s.style.Width > 0 {
		width = s.style.Width
	}
	if s.style.Height > 0 {
		height = s.style.Height
	}

	width = clampAxis(width, constraints.MinWidth, constraints.MaxWidth)
	height = clampAxis(height, constraints.MinHeight, constraints.MaxHeight)

	s.flexItem.SetMeasuredSize(width, height)
	s.measured = true
	return flex.Size{Width: width, Height: height}
}

// Layout 布局开关
func (s *Switch) Layout(x, y, width, height float32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rect.Position.X = x
	s.rect.Position.Y = y
	s.rect.Size.Width = width
	s.rect.Size.Height = height

	s.flexItem.Left = x
	s.flexItem.Top = y
	s.flexItem.Width = width
	s.flexItem.Height = height
}

// Render 渲染开关
func (s *Switch) Render() error { return renderWidget(s) }

// HandleEvent 处理事件：点击切换状态。
func (s *Switch) HandleEvent(event Event) bool {
	if !s.GetEnabled() || !s.GetVisible() {
		return false
	}
	if event.GetType() == EventClick {
		s.Toggle()
		return true
	}
	return s.BaseWidget.HandleEvent(event)
}
