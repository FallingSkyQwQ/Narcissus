package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// ProgressBar 是进度条组件，支持确定值和不确定（indeterminate）两种模式。
type ProgressBar struct {
	*BaseWidget

	min   float32
	max   float32
	value float32

	indeterminate bool
	showText      bool

	valueSignal *reactive.Signal[float32]
}

// NewProgressBar 创建进度条，默认范围 0-100、值 0。
func NewProgressBar() *ProgressBar {
	p := &ProgressBar{
		BaseWidget:  NewBaseWidget(),
		min:         0,
		max:         100,
		value:       0,
		valueSignal: reactive.NewSignal(float32(0)),
	}
	p.style.BackgroundColor = ColorLightDivider
	p.style.Border.Radius = 4
	return p
}

// Min 设置最小值（链式调用）
func (p *ProgressBar) Min(min float32) *ProgressBar {
	p.min = min
	if p.value < min {
		p.Value(min)
	} else {
		syncWidget(p)
	}
	return p
}

// Max 设置最大值（链式调用）
func (p *ProgressBar) Max(max float32) *ProgressBar {
	p.max = max
	if p.value > max {
		p.Value(max)
	} else {
		syncWidget(p)
	}
	return p
}

// Range 设置取值范围（链式调用）
func (p *ProgressBar) Range(min, max float32) *ProgressBar {
	p.min = min
	p.max = max
	p.Value(p.value)
	return p
}

// Value 设置当前值（链式调用）。值会被限制在 [min, max] 内。
func (p *ProgressBar) Value(value float32) *ProgressBar {
	if value < p.min {
		value = p.min
	} else if value > p.max {
		value = p.max
	}
	if p.value != value {
		p.value = value
		p.valueSignal.Set(value)
	}
	syncWidget(p)
	return p
}

// GetValue 获取当前值
func (p *ProgressBar) GetValue() float32 { return p.value }

// GetMin 获取最小值
func (p *ProgressBar) GetMin() float32 { return p.min }

// GetMax 获取最大值
func (p *ProgressBar) GetMax() float32 { return p.max }

// Percent 返回当前进度百分比（0-1）。
func (p *ProgressBar) Percent() float32 {
	if p.max <= p.min {
		return 0
	}
	return (p.value - p.min) / (p.max - p.min)
}

// Indeterminate 设置不确定模式（链式调用）
func (p *ProgressBar) Indeterminate(on bool) *ProgressBar {
	p.indeterminate = on
	syncWidget(p)
	return p
}

// IsIndeterminate 返回是否为不确定模式
func (p *ProgressBar) IsIndeterminate() bool { return p.indeterminate }

// ShowText 设置是否显示进度文本（链式调用）
func (p *ProgressBar) ShowText(on bool) *ProgressBar {
	p.showText = on
	syncWidget(p)
	return p
}

// IsShowingText 返回是否显示进度文本
func (p *ProgressBar) IsShowingText() bool { return p.showText }

// GetValueSignal 获取值的状态信号
func (p *ProgressBar) GetValueSignal() *reactive.Signal[float32] { return p.valueSignal }

// Width 设置宽度（链式调用）
func (p *ProgressBar) Width(width float32) *ProgressBar {
	p.style.Width = width
	return p
}

// Height 设置高度（链式调用）
func (p *ProgressBar) Height(height float32) *ProgressBar {
	p.style.Height = height
	return p
}

// Style 设置样式（链式调用）
func (p *ProgressBar) Style(style *Style) *ProgressBar {
	p.SetStyle(style)
	return p
}

// Measure 测量进度条大小
func (p *ProgressBar) Measure(constraints flex.Constraint) flex.Size {
	p.mu.Lock()
	defer p.mu.Unlock()

	width := p.style.Width
	if width <= 0 {
		width = 200
	}
	height := p.style.Height
	if height <= 0 {
		if p.showText {
			height = p.style.Font.Size * p.style.Font.LineHeight
		} else {
			height = 8
		}
	}

	width = clampAxis(width, constraints.MinWidth, constraints.MaxWidth)
	height = clampAxis(height, constraints.MinHeight, constraints.MaxHeight)

	p.flexItem.SetMeasuredSize(width, height)
	p.measured = true
	return flex.Size{Width: width, Height: height}
}

// Layout 布局进度条
func (p *ProgressBar) Layout(x, y, width, height float32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.rect.Position.X = x
	p.rect.Position.Y = y
	p.rect.Size.Width = width
	p.rect.Size.Height = height

	p.flexItem.Left = x
	p.flexItem.Top = y
	p.flexItem.Width = width
	p.flexItem.Height = height
}

// Render 渲染进度条
func (p *ProgressBar) Render() error { return renderWidget(p) }
