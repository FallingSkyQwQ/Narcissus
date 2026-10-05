package ui

import (
	"sync"
	"time"
)

// Focusable is implemented by widgets that opt into keyboard focus and take
// part in TAB traversal. BaseWidget provides the opt-in; widgets enable it in
// their constructor. Widgets that do not implement it are skipped.
type Focusable interface {
	IsFocusable() bool
}

// focusRequester is implemented by native controls that can take toolkit
// focus. Backends implement it on controls that are focusable, so the focus
// manager can hand focus to the toolkit without every control growing the
// method.
type focusRequester interface {
	RequestFocus()
}

var (
	focusMu     sync.RWMutex
	focusTarget Widget
	focused     Widget
)

// FocusedWidget returns the widget that currently holds keyboard focus, or nil.
func FocusedWidget() Widget {
	focusMu.RLock()
	defer focusMu.RUnlock()
	return focused
}

// SetFocusRoot records the widget tree used for TAB traversal. The mount path
// calls it with the window content whenever a tree is (re)layed out.
func SetFocusRoot(root Widget) {
	focusMu.Lock()
	focusTarget = root
	focusMu.Unlock()
}

// CanFocus reports whether w can hold focus right now: it must opt in, be
// enabled and be visible.
func CanFocus(w Widget) bool {
	bw := widgetBase(w)
	if bw == nil {
		return false
	}
	return bw.IsFocusable() && bw.GetEnabled() && bw.GetVisible()
}

// RequestFocus moves keyboard focus to w. Passing nil clears focus. A focus
// change dispatches EventBlur to the previous widget and EventFocus to the new
// one, and asks the backend for native focus when the control supports it.
//
// It is a no-op when w is not focusable, so callers cannot park focus on a
// plain container or text label.
func RequestFocus(w Widget) {
	if w != nil && !CanFocus(w) {
		return
	}

	focusMu.Lock()
	prev := focused
	if prev == w {
		focusMu.Unlock()
		return
	}
	focused = w
	focusMu.Unlock()

	if prev != nil {
		dispatchFocusEvent(prev, EventBlur)
	}
	if w != nil {
		dispatchFocusEvent(w, EventFocus)
		requestNativeFocus(w)
	}
}

// notifyNativeFocus is called by a backend when a control gains toolkit focus
// (for example because the user clicked it). It updates the focused widget and
// dispatches the focus events without asking the toolkit for focus again.
func notifyNativeFocus(w Widget) {
	if w == nil {
		return
	}
	focusMu.Lock()
	prev := focused
	if prev == w {
		focusMu.Unlock()
		return
	}
	focused = w
	focusMu.Unlock()

	if prev != nil {
		dispatchFocusEvent(prev, EventBlur)
	}
	dispatchFocusEvent(w, EventFocus)
}

// notifyNativeBlur is called by a backend when a control loses toolkit focus.
// It only clears the focused widget when that widget still holds focus, so the
// enter/leave pairs a toolkit emits while moving focus do not double-dispatch.
func notifyNativeBlur(w Widget) {
	if w == nil {
		return
	}
	focusMu.Lock()
	if focused != w {
		focusMu.Unlock()
		return
	}
	focused = nil
	focusMu.Unlock()

	dispatchFocusEvent(w, EventBlur)
}

// FocusNext moves focus to the next focusable widget in document order,
// wrapping around at the end.
func FocusNext() { moveFocus(1) }

// FocusPrevious moves focus to the previous focusable widget, wrapping around.
func FocusPrevious() { moveFocus(-1) }

// FocusOrder returns the focusable widgets under root in document order.
func FocusOrder(root Widget) []Widget {
	var order []Widget
	var walk func(Widget)
	walk = func(w Widget) {
		if w == nil {
			return
		}
		if CanFocus(w) {
			order = append(order, w)
		}
		for _, child := range w.GetChildren() {
			walk(child)
		}
	}
	walk(root)
	return order
}

// DispatchKey routes a key event to the focused widget and then up its
// ancestors, stopping at the first handler that consumes it. TAB and
// Shift+TAB move focus when the focused widget does not handle them.
func DispatchKey(ev *KeyEvent) bool {
	if ev == nil {
		return false
	}

	if ev.Type == EventKeyDown && ev.Key == "Tab" {
		if w := FocusedWidget(); w != nil && w.HandleEvent(ev) {
			return true
		}
		if ev.Shift {
			FocusPrevious()
		} else {
			FocusNext()
		}
		return true
	}

	for w := FocusedWidget(); w != nil; w = w.GetParent() {
		if w.HandleEvent(ev) {
			return true
		}
		if ev.IsPropagationStopped() {
			return true
		}
	}
	return false
}

// NewKeyEvent builds a key event with the current timestamp and no target; the
// dispatch path fills in bubbling through the focused widget's ancestry.
func NewKeyEvent(kind EventType, key string, keyCode int, ctrl, shift, alt, meta bool) *KeyEvent {
	return &KeyEvent{
		BaseEvent: BaseEvent{Type: kind, Timestamp: time.Now().UnixMilli()},
		Key:       key,
		KeyCode:   keyCode,
		Ctrl:      ctrl,
		Shift:     shift,
		Alt:       alt,
		Meta:      meta,
	}
}

func dispatchFocusEvent(w Widget, kind EventType) {
	w.HandleEvent(&BaseEvent{Type: kind, Target: w, Timestamp: time.Now().UnixMilli()})
}

func requestNativeFocus(w Widget) {
	bw := widgetBase(w)
	if bw == nil {
		return
	}
	bw.mu.RLock()
	ctrl := bw.native
	bw.mu.RUnlock()
	if requester, ok := ctrl.(focusRequester); ok && requester != nil {
		requester.RequestFocus()
	}
}

func moveFocus(delta int) {
	focusMu.RLock()
	root := focusTarget
	current := focused
	focusMu.RUnlock()

	if root == nil {
		return
	}

	order := FocusOrder(root)
	if len(order) == 0 {
		return
	}

	index := -1
	for i, w := range order {
		if w == current {
			index = i
			break
		}
	}

	var next int
	switch {
	case index == -1 && delta > 0:
		next = 0
	case index == -1:
		next = len(order) - 1
	default:
		next = ((index+delta)%len(order) + len(order)) % len(order)
	}
	RequestFocus(order[next])
}
