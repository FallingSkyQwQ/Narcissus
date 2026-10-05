package ui

import (
	"errors"
	"sync"
)

// ControlKind identifies the kind of native control a widget maps onto. A
// Backend uses it (together with a ControlProps snapshot) to create the
// matching native control for a platform.
type ControlKind int

const (
	// ControlContainer is an absolutely positioned surface that hosts children.
	ControlContainer ControlKind = iota
	// ControlButton maps to a push button.
	ControlButton
	// ControlText maps to a static text label.
	ControlText
	// ControlCheckbox maps to a check box.
	ControlCheckbox
	// ControlComboBox maps to a drop-down selector.
	ControlComboBox
	// ControlSlider maps to a range slider.
	ControlSlider
	// ControlImage maps to an image view.
	ControlImage
	// ControlTextInput maps to a single/multi-line text entry.
	ControlTextInput
)

// String returns a stable, human readable name for a ControlKind.
func (k ControlKind) String() string {
	switch k {
	case ControlContainer:
		return "container"
	case ControlButton:
		return "button"
	case ControlText:
		return "text"
	case ControlCheckbox:
		return "checkbox"
	case ControlComboBox:
		return "combobox"
	case ControlSlider:
		return "slider"
	case ControlImage:
		return "image"
	case ControlTextInput:
		return "textinput"
	default:
		return "unknown"
	}
}

// ControlProps is a toolkit-neutral snapshot of everything a Backend needs in
// order to create or refresh a native control. Widgets never talk to a toolkit
// directly; they expose their state through this struct.
type ControlProps struct {
	// Text is the primary text of the control (button label, text content,
	// check box label, place-holder is separate).
	Text string
	// Placeholder is the hint text for text inputs and combo boxes.
	Placeholder string
	// Checked reports the current check box state.
	Checked bool
	// Value, Min, Max and Step describe a slider.
	Value float64
	Min   float64
	Max   float64
	Step  float64
	// Items and Selected describe a combo box.
	Items    []string
	Selected int
	// Source is the image path or URL.
	Source string
	// Fit is the image fit mode (an ImageFit value).
	Fit int
	// Multiline reports whether a text input should grow to multiple lines.
	Multiline bool
	// ReadOnly reports whether a text input can be edited.
	ReadOnly bool
}

// NativeControl is a handle to a native control created by a Backend. All
// coordinates are in pixels relative to the control's parent surface.
type NativeControl interface {
	// SetBounds moves and resizes the control within its parent surface.
	SetBounds(x, y, width, height float32)
	// SetVisible shows or hides the control.
	SetVisible(visible bool)
	// SetEnabled enables or disables user interaction with the control.
	SetEnabled(enabled bool)
	// SetStyle applies the toolkit-independent style to the control.
	SetStyle(style *Style)
	// Update refreshes the control from a property snapshot.
	Update(props ControlProps)
	// AttachTo parents the control inside the given surface. It is a no-op
	// when the control is already attached to that parent.
	AttachTo(parent NativeControl)
	// Destroy releases the native resources of the control.
	Destroy()
}

// NativeWindow is a native top-level window produced by a Backend.
type NativeWindow interface {
	// Root returns the absolute-positioning surface that hosts the content
	// tree. Its origin is the top-left corner of the window's client area.
	Root() NativeControl
	// SetTitle updates the window title.
	SetTitle(title string)
	// SetSize requests a new client size in pixels.
	SetSize(width, height float32)
	// SetContent lays out and mounts the given widget tree as the window's
	// content. It is safe to call again to relayout an existing tree.
	SetContent(content Widget) error
	// Show makes the window visible.
	Show()
	// Close destroys the window.
	Close()
}

// Backend adapts the toolkit-independent widget model to a concrete platform
// UI toolkit (GTK4 on Linux, WinUI on Windows, ...). Implementations live in
// platform-tagged files inside this package so that application code only ever
// imports pkg/ui.
type Backend interface {
	// Name returns the backend identifier, e.g. "gtk4".
	Name() string
	// Init prepares the toolkit. It must be called on the UI thread.
	Init() error
	// CreateWindow creates a top-level window.
	CreateWindow(title string, width, height float32) (NativeWindow, error)
	// CreateControl creates the native control backing a widget. The backend
	// wires toolkit events back into w (for example a button click becomes an
	// EventClick dispatched through w.HandleEvent).
	CreateControl(w Widget, kind ControlKind, props ControlProps) (NativeControl, error)
	// Run enters the toolkit event loop and blocks until Quit is called.
	Run() error
	// Quit asks the event loop to stop.
	Quit()
	// Post schedules fn to run on the UI thread.
	Post(fn func()) error
	// OnUIThread reports whether the caller is on the UI thread.
	OnUIThread() bool
}

// ErrBackendUnavailable is returned when the current platform has no backend.
var ErrBackendUnavailable = errors.New("ui: no backend available for this platform")

var (
	backendMu     sync.RWMutex
	activeBackend Backend
)

// SetBackend installs the backend used by App and the widget mounting code.
// Applications rarely need this; it exists so tests and alternate front-ends
// can inject a custom backend.
func SetBackend(b Backend) {
	backendMu.Lock()
	defer backendMu.Unlock()
	activeBackend = b
}

// CurrentBackend returns the active backend, lazily creating the platform
// default on first use.
func CurrentBackend() Backend {
	backendMu.RLock()
	if activeBackend != nil {
		b := activeBackend
		backendMu.RUnlock()
		return b
	}
	backendMu.RUnlock()

	backendMu.Lock()
	defer backendMu.Unlock()
	if activeBackend == nil {
		activeBackend = newPlatformBackend()
	}
	return activeBackend
}
