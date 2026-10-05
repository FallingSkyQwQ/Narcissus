package ui

import (
	"sync"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// DefaultToastDurationMS is the fallback visible duration for a Toast.
const DefaultToastDurationMS = 3000

// Toast 是一个瞬时通知组件。它不出现在布局树中：调用 Show 时通过后端的
// Presenter 能力呈现。
type Toast struct {
	*BaseWidget

	mu         sync.Mutex
	message    string
	severity   ToastSeverity
	durationMS int
}

// NewToast 创建通知，message 为文本。
func NewToast(message string) *Toast {
	return &Toast{
		BaseWidget: NewBaseWidget(),
		message:    message,
		severity:   ToastInfo,
		durationMS: DefaultToastDurationMS,
	}
}

// Severity 设置严重级别（链式调用）
func (t *Toast) Severity(severity ToastSeverity) *Toast {
	t.mu.Lock()
	t.severity = severity
	t.mu.Unlock()
	return t
}

// DurationMS 设置可见时长（毫秒，链式调用）
func (t *Toast) DurationMS(duration int) *Toast {
	t.mu.Lock()
	t.durationMS = duration
	t.mu.Unlock()
	return t
}

// Message 返回通知文本
func (t *Toast) Message() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.message
}

// GetSeverity 返回严重级别
func (t *Toast) GetSeverity() ToastSeverity {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.severity
}

// GetDurationMS 返回可见时长
func (t *Toast) GetDurationMS() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.durationMS
}

// Show 呈现通知。没有后端或后端不支持 Presenter 时是空操作。
func (t *Toast) Show() error {
	presenter := currentPresenter()
	if presenter == nil {
		return nil
	}

	t.mu.Lock()
	spec := ToastSpec{
		Message:    t.message,
		Severity:   t.severity,
		DurationMS: t.durationMS,
	}
	t.mu.Unlock()

	return presenter.PresentToast(spec)
}

// Style 设置样式（链式调用）
func (t *Toast) Style(style *Style) *Toast {
	t.SetStyle(style)
	return t
}

// Measure 让通知在布局树中占据零尺寸。
func (t *Toast) Measure(flex.Constraint) flex.Size {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.flexItem.SetMeasuredSize(0, 0)
	t.measured = true
	return flex.Size{}
}

// Layout 让通知在布局树中保持零尺寸。
func (t *Toast) Layout(x, y, width, height float32) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.rect.Position.X = x
	t.rect.Position.Y = y
	t.rect.Size.Width = 0
	t.rect.Size.Height = 0
}

// Render 渲染通知：空操作，呈现由 Show 完成。
func (t *Toast) Render() error { return nil }
