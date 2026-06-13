package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// Slider 是滑块组件
type Slider struct {
	*BaseWidget

	// 最小值
	min float32

	// 最大值
	max float32

	// 当前值
	value float32

	// 步长
	step float32

	// 值的状态信号
	valueSignal *reactive.Signal[float32]

	// 变更事件处理器
	onChange func(float32)
}

// NewSlider 创建新的滑块组件
func NewSlider() *Slider {
	s := &Slider{
		BaseWidget:  NewBaseWidget(),
		min:         0,
		max:         100,
		value:       0,
		step:        1,
		valueSignal: reactive.NewSignal(float32(0)),
	}

	// 设置默认样式
	s.style.BackgroundColor = ColorTransparent

	return s
}

// Min 设置最小值（链式调用）
func (s *Slider) Min(min float32) *Slider {
	s.min = min
	// 确保 min <= max
	if s.min > s.max {
		s.max = s.min
	}
	// 确保当前值在有效范围内
	if s.value < s.min {
		s.Value(s.min)
	} else if s.value > s.max {
		s.Value(s.max)
	}
	return s
}

// GetMin 获取最小值
func (s *Slider) GetMin() float32 {
	return s.min
}

// Max 设置最大值（链式调用）
func (s *Slider) Max(max float32) *Slider {
	s.max = max
	// 确保 min <= max
	if s.max < s.min {
		s.min = s.max
	}
	// 确保当前值在有效范围内
	if s.value > s.max {
		s.Value(s.max)
	} else if s.value < s.min {
		s.Value(s.min)
	}
	return s
}

// GetMax 获取最大值
func (s *Slider) GetMax() float32 {
	return s.max
}

// Range 设置范围（链式调用）
func (s *Slider) Range(min, max float32) *Slider {
	s.min = min
	s.max = max
	// 确保当前值在有效范围内
	if s.value < min {
		s.Value(min)
	} else if s.value > max {
		s.Value(max)
	}
	return s
}

// Value 设置当前值（链式调用）
func (s *Slider) Value(value float32) *Slider {
	// 先限制在范围内
	if value < s.min {
		value = s.min
	} else if value > s.max {
		value = s.max
	}

	// 应用步长
	if s.step > 0 {
		steps := (value - s.min) / s.step
		value = s.min + float32(int(steps+0.5))*s.step
		// 再次限制在范围内，防止舍入误差
		if value < s.min {
			value = s.min
		} else if value > s.max {
			value = s.max
		}
	}

	oldValue := s.value
	s.value = value
	s.valueSignal.Set(value)

	// 触发变更事件
	if oldValue != value && s.onChange != nil {
		s.onChange(value)
	}

	return s
}

// GetValue 获取当前值
func (s *Slider) GetValue() float32 {
	return s.value
}

// GetValueSignal 获取值的状态信号
func (s *Slider) GetValueSignal() *reactive.Signal[float32] {
	return s.valueSignal
}

// Step 设置步长（链式调用）
func (s *Slider) Step(step float32) *Slider {
	s.step = step
	// 重新应用当前值以确保符合步长
	s.Value(s.value)
	return s
}

// GetStep 获取步长
func (s *Slider) GetStep() float32 {
	return s.step
}

// OnChange 设置变更事件处理器（链式调用）
func (s *Slider) OnChange(handler func(float32)) *Slider {
	s.onChange = handler
	return s
}

// FlexGrow 设置 FlexGrow（链式调用）
func (s *Slider) FlexGrow(grow float32) *Slider {
	s.style.FlexGrow = grow
	s.flexItem.FlexGrow = grow
	return s
}

// FlexShrink 设置 FlexShrink（链式调用）
func (s *Slider) FlexShrink(shrink float32) *Slider {
	s.style.FlexShrink = shrink
	s.flexItem.FlexShrink = shrink
	return s
}

// Margin 设置外边距（链式调用）
func (s *Slider) Margin(margin Insets) *Slider {
	s.style.Margin = margin
	s.flexItem.MarginTop = margin.Top
	s.flexItem.MarginRight = margin.Right
	s.flexItem.MarginBottom = margin.Bottom
	s.flexItem.MarginLeft = margin.Left
	return s
}

// Width 设置宽度（链式调用）
func (s *Slider) Width(width float32) *Slider {
	s.style.Width = width
	return s
}

// Style 设置样式（链式调用）
func (s *Slider) Style(style *Style) *Slider {
	s.SetStyle(style)
	return s
}

// SetValueByPosition 根据位置设置值（用于拖动）
func (s *Slider) SetValueByPosition(position float32) {
	// position 应该是 0-1 之间的值
	if position < 0 {
		position = 0
	} else if position > 1 {
		position = 1
	}

	rangeSize := s.max - s.min
	if rangeSize <= 0 {
		s.Value(s.min)
		return
	}

	value := s.min + position*rangeSize
	s.Value(value)
}

// GetValueAsPercentage 获取值的百分比位置（0-1）
func (s *Slider) GetValueAsPercentage() float32 {
	rangeSize := s.max - s.min
	if rangeSize <= 0 {
		return 0
	}
	return (s.value - s.min) / rangeSize
}

// Measure 测量滑块大小
func (s *Slider) Measure(constraints flex.Constraint) flex.Size {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 默认滑块尺寸
	width := s.style.Width
	if width == 0 {
		width = 200 // 默认宽度
	}
	height := float32(20) // 默认高度

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

	s.flexItem.SetMeasuredSize(width, height)
	s.measured = true

	return flex.Size{Width: width, Height: height}
}

// Layout 布局滑块
func (s *Slider) Layout(x, y, width, height float32) {
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

// Render 渲染滑块
func (s *Slider) Render() error {
	if !s.GetVisible() {
		return nil
	}

	// 实际渲染应该调用底层 WinUI API 创建 Slider
	// 这里只是占位实现
	return nil
}

// HandleEvent 处理事件
func (s *Slider) HandleEvent(event Event) bool {
	if !s.GetEnabled() || !s.GetVisible() {
		return false
	}

	// 处理鼠标/触摸事件来更新值
	switch e := event.(type) {
	case *MouseEvent:
		// 根据鼠标位置更新值
		if s.rect.Size.Width > 0 {
			relativeX := e.X - s.rect.Position.X
			position := relativeX / s.rect.Size.Width
			s.SetValueByPosition(position)
			return true
		}
	}

	// 委托给基类处理其他事件
	return s.BaseWidget.HandleEvent(event)
}
