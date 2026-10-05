package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// ScrollView 是带滚动条的容器，承载单个子组件。
//
// The native control is a scrolling surface (GtkScrolledWindow / ScrollViewer)
// and the child keeps the framework's computed rectangle inside it, so the
// toolkit only provides the scrollbars.
type ScrollView struct {
	*BaseWidget

	content Widget
}

// NewScrollView 创建空滚动容器。
func NewScrollView() *ScrollView {
	s := &ScrollView{BaseWidget: NewBaseWidget()}
	s.style.BackgroundColor = ColorTransparent
	return s
}

// SetContent 设置内容组件（链式调用）
func (s *ScrollView) SetContent(content Widget) *ScrollView {
	if s.content != nil {
		s.RemoveChild(s.content)
	}
	s.content = content
	if content != nil {
		s.AddChild(content)
	}
	return s
}

// Content 返回当前内容组件
func (s *ScrollView) Content() Widget { return s.content }

// Style 设置样式（链式调用）
func (s *ScrollView) Style(style *Style) *ScrollView {
	s.SetStyle(style)
	return s
}

// Size 设置尺寸（链式调用）
func (s *ScrollView) Size(width, height float32) *ScrollView {
	s.style.Width = width
	s.style.Height = height
	return s
}

// Measure 测量滚动容器大小
func (s *ScrollView) Measure(constraints flex.Constraint) flex.Size {
	s.mu.Lock()
	defer s.mu.Unlock()

	width := s.style.Width
	height := s.style.Height

	if s.content != nil {
		inner := constraints
		if s.style.Width > 0 {
			inner.MaxWidth = s.style.Width
			inner.MinWidth = s.style.Width
		}
		if s.style.Height > 0 {
			inner.MaxHeight = s.style.Height
			inner.MinHeight = s.style.Height
		}
		childSize := s.content.Measure(inner)
		if width <= 0 {
			width = childSize.Width
		}
		if height <= 0 {
			height = childSize.Height
		}
	}

	if width <= 0 {
		width = 200
	}
	if height <= 0 {
		height = 150
	}

	width = clampAxis(width, constraints.MinWidth, constraints.MaxWidth)
	height = clampAxis(height, constraints.MinHeight, constraints.MaxHeight)

	s.flexItem.SetMeasuredSize(width, height)
	s.measured = true
	return flex.Size{Width: width, Height: height}
}

// Layout 布局滚动容器；内容按其自然高度放置，超出部分由原生滚动条处理。
func (s *ScrollView) Layout(x, y, width, height float32) {
	s.mu.Lock()
	s.rect.Position.X = x
	s.rect.Position.Y = y
	s.rect.Size.Width = width
	s.rect.Size.Height = height
	s.flexItem.Left = x
	s.flexItem.Top = y
	s.flexItem.Width = width
	s.flexItem.Height = height
	s.mu.Unlock()

	if s.content == nil {
		return
	}
	contentHeight := s.content.Measure(flex.Constraint{MaxWidth: width}).Height
	s.content.Layout(x, y, width, contentHeight)
}

// Render 渲染滚动容器
func (s *ScrollView) Render() error { return renderWidget(s) }
