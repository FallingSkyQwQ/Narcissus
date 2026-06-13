package ui

import (
	"sync"
	"sync/atomic"
)

// EventType 定义事件类型
type EventType int

const (
	// 鼠标事件
	EventClick EventType = iota
	EventDoubleClick
	EventMouseDown
	EventMouseUp
	EventMouseMove
	EventMouseEnter
	EventMouseLeave

	// 键盘事件
	EventKeyDown
	EventKeyUp

	// 输入事件
	EventChange
	EventInput
	EventFocus
	EventBlur

	// 触摸事件
	EventTouchStart
	EventTouchMove
	EventTouchEnd

	// 自定义事件
	EventCustom
)

// String 返回事件类型的字符串表示
func (et EventType) String() string {
	switch et {
	case EventClick:
		return "Click"
	case EventDoubleClick:
		return "DoubleClick"
	case EventMouseDown:
		return "MouseDown"
	case EventMouseUp:
		return "MouseUp"
	case EventMouseMove:
		return "MouseMove"
	case EventMouseEnter:
		return "MouseEnter"
	case EventMouseLeave:
		return "MouseLeave"
	case EventKeyDown:
		return "KeyDown"
	case EventKeyUp:
		return "KeyUp"
	case EventChange:
		return "Change"
	case EventInput:
		return "Input"
	case EventFocus:
		return "Focus"
	case EventBlur:
		return "Blur"
	case EventTouchStart:
		return "TouchStart"
	case EventTouchMove:
		return "TouchMove"
	case EventTouchEnd:
		return "TouchEnd"
	case EventCustom:
		return "Custom"
	default:
		return "Unknown"
	}
}

// Event 是所有事件的接口
type Event interface {
	GetType() EventType
	GetTarget() Widget
	GetTimestamp() int64
	PreventDefault()
	IsDefaultPrevented() bool
	StopPropagation()
	IsPropagationStopped() bool
}

// BaseEvent 是事件的基础实现
type BaseEvent struct {
	Type               EventType
	Target             Widget
	Timestamp          int64
	defaultPrevented   bool
	propagationStopped bool
}

// GetType 返回事件类型
func (e *BaseEvent) GetType() EventType {
	return e.Type
}

// GetTarget 返回事件目标
func (e *BaseEvent) GetTarget() Widget {
	return e.Target
}

// GetTimestamp 返回事件时间戳
func (e *BaseEvent) GetTimestamp() int64 {
	return e.Timestamp
}

// PreventDefault 阻止默认行为
func (e *BaseEvent) PreventDefault() {
	e.defaultPrevented = true
}

// IsDefaultPrevented 返回是否阻止了默认行为
func (e *BaseEvent) IsDefaultPrevented() bool {
	return e.defaultPrevented
}

// StopPropagation 停止事件传播
func (e *BaseEvent) StopPropagation() {
	e.propagationStopped = true
}

// IsPropagationStopped 返回是否停止了事件传播
func (e *BaseEvent) IsPropagationStopped() bool {
	return e.propagationStopped
}

// MouseEvent 鼠标事件
type MouseEvent struct {
	BaseEvent
	X      float32
	Y      float32
	Button int
}

// KeyEvent 键盘事件
type KeyEvent struct {
	BaseEvent
	KeyCode int
	Key     string
	Ctrl    bool
	Shift   bool
	Alt     bool
	Meta    bool
}

// ChangeEvent 变更事件
type ChangeEvent struct {
	BaseEvent
	Value    string
	OldValue string
}

// InputEvent 输入事件
type InputEvent struct {
	BaseEvent
	Value string
}

// EventHandler 事件处理函数类型
type EventHandler func(Event)

// EventRegistry 事件注册表
type EventRegistry struct {
	mu       sync.RWMutex
	handlers map[EventType]map[uint64]EventHandler
	nextID   atomic.Uint64
}

// NewEventRegistry 创建新的事件注册表
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		handlers: make(map[EventType]map[uint64]EventHandler),
	}
}

// AddHandler 添加事件处理器，返回取消订阅函数
func (er *EventRegistry) AddHandler(eventType EventType, handler EventHandler) func() {
	er.mu.Lock()
	defer er.mu.Unlock()

	if er.handlers[eventType] == nil {
		er.handlers[eventType] = make(map[uint64]EventHandler)
	}

	id := er.nextID.Add(1) - 1
	er.handlers[eventType][id] = handler

	return func() {
		er.mu.Lock()
		defer er.mu.Unlock()
		if er.handlers[eventType] != nil {
			delete(er.handlers[eventType], id)
		}
	}
}

// Dispatch 分发事件
func (er *EventRegistry) Dispatch(event Event) bool {
	er.mu.RLock()
	handlers := er.handlers[event.GetType()]
	if len(handlers) == 0 {
		er.mu.RUnlock()
		return false
	}

	// 复制处理器列表以避免在遍历期间修改
	handlerList := make([]EventHandler, 0, len(handlers))
	for _, handler := range handlers {
		handlerList = append(handlerList, handler)
	}
	er.mu.RUnlock()

	handled := false
	for _, handler := range handlerList {
		handler(event)
		handled = true
		if event.IsPropagationStopped() {
			break
		}
	}

	return handled
}

// HasHandlers 检查是否有指定类型的事件处理器
func (er *EventRegistry) HasHandlers(eventType EventType) bool {
	er.mu.RLock()
	defer er.mu.RUnlock()
	return len(er.handlers[eventType]) > 0
}

// Clear 清除所有事件处理器
func (er *EventRegistry) Clear() {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.handlers = make(map[EventType]map[uint64]EventHandler)
}
