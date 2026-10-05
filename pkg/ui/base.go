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
	visible   bool
	enabled   bool
	focusable bool

	// 无障碍语义（供辅助技术读取）
	accessibleName        string
	accessibleDescription string
	accessibleRole        AccessibleRole

	// 原生控件句柄（由平台后端创建，未挂载时为 nil）
	native NativeControl

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

// SetNativeControl 记录后端为该组件创建的原生控件句柄
func (bw *BaseWidget) SetNativeControl(control NativeControl) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	bw.native = control
}

// NativeControl 返回后端为该组件创建的原生控件句柄（未挂载时返回 nil）
func (bw *BaseWidget) NativeControl() NativeControl {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.native
}

// SetStyle 设置组件样式
func (bw *BaseWidget) SetStyle(style *Style) {
	if style == nil {
		return
	}

	bw.mu.Lock()
	defer bw.mu.Unlock()

	bw.style = style

	// 同步到原生控件（若已挂载）
	if bw.native != nil {
		bw.native.SetStyle(style)
	}

	// 更新 flex 属性
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

	removed := false
	bw.mu.Lock()
	for i, c := range bw.children {
		if c == child {
			bw.children = append(bw.children[:i], bw.children[i+1:]...)
			removed = true
			break
		}
	}
	bw.mu.Unlock()

	if removed {
		child.SetParent(nil)
	}
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
	native := bw.native
	bw.mu.Unlock()

	if native != nil {
		native.SetVisible(visible)
	}
}

// SetFocusable 设置组件是否可以获取键盘焦点
func (bw *BaseWidget) SetFocusable(focusable bool) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	bw.focusable = focusable
}

// IsFocusable 返回组件是否参与键盘焦点与 TAB 遍历
func (bw *BaseWidget) IsFocusable() bool {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.focusable
}

// SetAccessibility 一次性设置组件的无障碍角色、名称与描述。设置后会在下次
// 布局或状态同步时推送到原生控件，并被 AccessibilityTree 读取。
func (bw *BaseWidget) SetAccessibility(role AccessibleRole, name, description string) {
	bw.mu.Lock()
	bw.accessibleRole = role
	bw.accessibleName = name
	bw.accessibleDescription = description
	bw.mu.Unlock()
}

// SetAccessibleRole 设置无障碍角色（RoleNone 表示按组件类型推断）。
func (bw *BaseWidget) SetAccessibleRole(role AccessibleRole) {
	bw.mu.Lock()
	bw.accessibleRole = role
	bw.mu.Unlock()
}

// AccessibleRole 返回显式设置的无障碍角色，未设置时为 RoleNone。
func (bw *BaseWidget) AccessibleRole() AccessibleRole {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.accessibleRole
}

// SetAccessibleName 设置辅助技术朗读的名称。
func (bw *BaseWidget) SetAccessibleName(name string) {
	bw.mu.Lock()
	bw.accessibleName = name
	bw.mu.Unlock()
}

// AccessibleName 返回显式设置的无障碍名称。
func (bw *BaseWidget) AccessibleName() string {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.accessibleName
}

// SetAccessibleDescription 设置辅助技术朗读的补充描述。
func (bw *BaseWidget) SetAccessibleDescription(description string) {
	bw.mu.Lock()
	bw.accessibleDescription = description
	bw.mu.Unlock()
}

// AccessibleDescription 返回显式设置的无障碍描述。
func (bw *BaseWidget) AccessibleDescription() string {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return bw.accessibleDescription
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
	native := bw.native
	bw.mu.Unlock()

	if native != nil {
		native.SetEnabled(enabled)
	}
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
	parent := bw.parent
	bw.mu.Unlock()

	// 清理事件
	bw.events.Clear()

	// 释放原生控件
	if bw.native != nil {
		bw.native.Destroy()
		bw.native = nil
	}

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
	if parent != nil {
		parent.RemoveChild(bw)
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

// clampAxis 将数值限制在非零的 min/max 约束内。
func clampAxis(value, minValue, maxValue float32) float32 {
	if minValue > 0 && value < minValue {
		value = minValue
	}
	if maxValue > 0 && value > maxValue {
		value = maxValue
	}
	return value
}

// SetFlexItem 设置 FlexItem
func (bw *BaseWidget) SetFlexItem(item flex.Item) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	bw.flexItem = item
}
