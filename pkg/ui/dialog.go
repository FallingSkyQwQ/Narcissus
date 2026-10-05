package ui

import (
	"sync"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// Dialog 是一个模态对话框组件。它不出现在布局树中：调用 Show 时通过后端的
// Presenter 能力呈现。
type Dialog struct {
	*BaseWidget

	mu       sync.Mutex
	title    string
	message  string
	buttons  []string
	onResult func(index int)
}

// NewDialog 创建对话框，title 为标题、message 为正文。
func NewDialog(title, message string) *Dialog {
	d := &Dialog{
		BaseWidget: NewBaseWidget(),
		title:      title,
		message:    message,
		buttons:    []string{"OK"},
	}
	return d
}

// Buttons 设置按钮文本（链式调用）
func (d *Dialog) Buttons(labels ...string) *Dialog {
	d.mu.Lock()
	d.buttons = append([]string(nil), labels...)
	d.mu.Unlock()
	return d
}

// OnResult 设置结果处理器（链式调用），参数为被点击按钮的索引。
func (d *Dialog) OnResult(handler func(index int)) *Dialog {
	d.mu.Lock()
	d.onResult = handler
	d.mu.Unlock()
	return d
}

// Title 返回对话框标题
func (d *Dialog) Title() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.title
}

// Message 返回对话框正文
func (d *Dialog) Message() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.message
}

// GetButtons 返回按钮文本副本
func (d *Dialog) GetButtons() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]string(nil), d.buttons...)
}

// Show 呈现对话框。没有后端或后端不支持 Presenter 时是空操作。
func (d *Dialog) Show() error {
	presenter := currentPresenter()
	if presenter == nil {
		return nil
	}

	d.mu.Lock()
	spec := DialogSpec{
		Title:    d.title,
		Message:  d.message,
		Buttons:  append([]string(nil), d.buttons...),
		OnResult: d.onResult,
	}
	d.mu.Unlock()

	return presenter.PresentDialog(spec)
}

// Style 设置样式（链式调用）
func (d *Dialog) Style(style *Style) *Dialog {
	d.SetStyle(style)
	return d
}

// Measure 让对话框在布局树中占据零尺寸。
func (d *Dialog) Measure(flex.Constraint) flex.Size {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.flexItem.SetMeasuredSize(0, 0)
	d.measured = true
	return flex.Size{}
}

// Layout 让对话框在布局树中保持零尺寸。
func (d *Dialog) Layout(x, y, width, height float32) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.rect.Position.X = x
	d.rect.Position.Y = y
	d.rect.Size.Width = 0
	d.rect.Size.Height = 0
}

// Render 渲染对话框：空操作，呈现由 Show 完成。
func (d *Dialog) Render() error { return nil }
