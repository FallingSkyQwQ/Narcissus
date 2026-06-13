package ui

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

var (
	globalWidgetID atomic.Uint64
)

// generateWidgetID 生成唯一的组件 ID
func generateWidgetID() string {
	id := globalWidgetID.Add(1)
	return fmt.Sprintf("widget_%d", id)
}

// BaseWidget 是 Widget 接口的基础实现
type BaseWidget struct {
	mu sync.RWMutex

	// 基础属性
	id string

	// 布局
	rect     flex.Rect
	flexItem flex.Item
	measured bool

	// 样式
	style *Style

	// 父子关系
	parent   Widget
	children []Widget

	// 事件
	events *EventRegistry

	// 状态
	visible bool
	enabled bool

	// 其他
	disposed bool
}

// NewBaseWidget 创建新的基础组件
func NewBaseWidget() *BaseWidget {
	return &BaseWidget{
		id:       generateWidgetID(),
		style:    NewStyle(),
		children: make([]Widget, 0),
		events:   NewEventRegistry(),
		visible:  true,
		enabled:  true,
		disposed: false,
	}
}

// GetID 获取组件 ID
func (bw *BaseWidget) GetID() string {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.id
}

// SetID 设置组件 ID
func (bw *BaseWidget) SetID(id string) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	bw.id = id
}

// GetRect 获取组件矩形区域
func (bw *BaseWidget) GetRect() flex.Rect {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.rect
}

// GetStyle 获取组件样式
func (bw *BaseWidget) GetStyle() *Style {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.style
}

// SetStyle 设置组件样式
func (bw *BaseWidget) SetStyle(style *Style) {
	bw.mu.Lock()
	bw.style = style
	bw.mu.Unlock()

	// 更新 flex 属性
	if style != nil {
		bw.flexItem.FlexGrow = style.FlexGrow
		bw.flexItem.FlexShrink = style.FlexShrink
		bw.flexItem.FlexBasis = style.FlexBasis
		bw.flexItem.AlignSelf = style.AlignSelf
		bw.flexItem.MarginTop = style.Margin.Top
		bw.flexItem.MarginRight = style.Margin.Right
		bw.flexItem.MarginBottom = style.Margin.Bottom
		bw.flexItem.MarginLeft = style.Margin.Left
		bw.flexItem.PaddingTop = style.Padding.Top
		bw.flexItem.PaddingRight = style.Padding.Right
		bw.flexItem.PaddingBottom = style.Padding.Bottom
		bw.flexItem.PaddingLeft = style.Padding.Left
	}
}

// GetParent 获取父组件
func (bw *BaseWidget) GetParent() Widget {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.parent
}

// SetParent 设置父组件
func (bw *BaseWidget) SetParent(parent Widget) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	bw.parent = parent
}

// GetChildren 获取子组件列表
func (bw *BaseWidget) GetChildren() []Widget {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	result := make([]Widget, len(bw.children))
	copy(result, bw.children)
	return result
}

// AddChild 添加子组件
func (bw *BaseWidget) AddChild(child Widget) {
	if child == nil {
		return
	}

	bw.mu.Lock()
	bw.children = append(bw.children, child)
	bw.mu.Unlock()

	child.SetParent(bw)
}

// RemoveChild 移除子组件
func (bw *BaseWidget) RemoveChild(child Widget) {
	if child == nil {
		return
	}

	bw.mu.Lock()
	for i, c := range bw.children {
		if c == child {
			bw.children = append(bw.children[:i], bw.children[i+1:]...)
			break
		}
	}
	bw.mu.Unlock()

	child.SetParent(nil)
}

// HandleEvent 处理事件
func (bw *BaseWidget) HandleEvent(event Event) bool {
	if !bw.GetEnabled() || !bw.GetVisible() {
		return false
	}
	return bw.events.Dispatch(event)
}

// AddEventHandler 添加事件处理器
func (bw *BaseWidget) AddEventHandler(eventType EventType, handler EventHandler) func() {
	return bw.events.AddHandler(eventType, handler)
}

// GetVisible 获取可见性
func (bw *BaseWidget) GetVisible() bool {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.visible
}

// SetVisible 设置可见性
func (bw *BaseWidget) SetVisible(visible bool) {
	bw.mu.Lock()
	bw.visible = visible
	bw.mu.Unlock()
}

// GetEnabled 获取启用状态
func (bw *BaseWidget) GetEnabled() bool {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.enabled
}

// SetEnabled 设置启用状态
func (bw *BaseWidget) SetEnabled(enabled bool) {
	bw.mu.Lock()
	bw.enabled = enabled
	bw.mu.Unlock()
}

// GetFlexProperties 获取 Flex 属性
func (bw *BaseWidget) GetFlexProperties() flex.FlexProperties {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.style.ToFlexProperties()
}

// Measure 测量组件大小（子类需要重写）
func (bw *BaseWidget) Measure(constraints flex.Constraint) flex.Size {
	bw.mu.Lock()
	defer bw.mu.Unlock()

	// 默认实现：使用样式中定义的宽高
	width := bw.style.Width
	height := bw.style.Height

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

	bw.flexItem.SetMeasuredSize(width, height)
	bw.measured = true

	return flex.Size{Width: width, Height: height}
}

// Layout 布局组件（子类需要重写）
func (bw *BaseWidget) Layout(x, y, width, height float32) {
	bw.mu.Lock()
	defer bw.mu.Unlock()

	bw.rect.Position.X = x
	bw.rect.Position.Y = y
	bw.rect.Size.Width = width
	bw.rect.Size.Height = height

	bw.flexItem.Left = x
	bw.flexItem.Top = y
	bw.flexItem.Width = width
	bw.flexItem.Height = height
}

// Render 渲染组件（子类需要重写）
func (bw *BaseWidget) Render() error {
	// 基类默认实现为空
	return nil
}

// Dispose 销毁组件
func (bw *BaseWidget) Dispose() {
	bw.mu.Lock()
	if bw.disposed {
		bw.mu.Unlock()
		return
	}
	bw.disposed = true
	bw.mu.Unlock()

	// 清理事件
	bw.events.Clear()

	// 递归销毁子组件
	bw.mu.RLock()
	children := make([]Widget, len(bw.children))
	copy(children, bw.children)
	bw.mu.RUnlock()

	for _, child := range children {
		if disposable, ok := child.(interface{ Dispose() }); ok {
			disposable.Dispose()
		}
	}

	// 从父组件中移除
	if bw.parent != nil {
		bw.parent.RemoveChild(bw)
	}
}

// IsDisposed 返回组件是否已销毁
func (bw *BaseWidget) IsDisposed() bool {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.disposed
}

// GetFlexItem 获取 FlexItem（供布局使用）
func (bw *BaseWidget) GetFlexItem() *flex.Item {
	return &bw.flexItem
}

// SetFlexItem 设置 FlexItem
func (bw *BaseWidget) SetFlexItem(item flex.Item) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	bw.flexItem = item
}
