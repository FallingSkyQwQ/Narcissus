package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// Widget 是所有 UI 组件的基础接口
type Widget interface {
	// 基础属性
	GetID() string
	SetID(id string)

	// 布局相关 - 实现 flex.LayoutNode 接口
	Measure(constraints flex.Constraint) flex.Size
	Layout(x, y, width, height float32)
	GetRect() flex.Rect
	GetFlexProperties() flex.FlexProperties

	// 样式
	GetStyle() *Style
	SetStyle(style *Style)

	// 渲染
	Render() error

	// 父子关系
	GetParent() Widget
	SetParent(parent Widget)
	GetChildren() []Widget
	AddChild(child Widget)
	RemoveChild(child Widget)

	// 事件
	HandleEvent(event Event) bool
	AddEventHandler(eventType EventType, handler EventHandler) func()

	// 可见性
	GetVisible() bool
	SetVisible(visible bool)

	// 启用状态
	GetEnabled() bool
	SetEnabled(enabled bool)
}

// WidgetBase 是 Widget 接口中通用方法的辅助接口
type WidgetBase interface {
	// 初始化组件
	Init()
	// 销毁组件
	Dispose()
}
