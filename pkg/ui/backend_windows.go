//go:build windows && amd64

package ui

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	syswinrt "github.com/deploymenttheory/go-bindings-win32/bindings/win32/system/winrt"
	"github.com/deploymenttheory/go-bindings-windowsappsdk/app"
	uidispatching "github.com/deploymenttheory/go-bindings-windowsappsdk/bindings/winui/ui/dispatching"
	uixaml "github.com/deploymenttheory/go-bindings-windowsappsdk/bindings/winui/ui/xaml"
	"github.com/deploymenttheory/go-bindings-winrt/bindings/runtime/winrt"
	wrtfoundation "github.com/deploymenttheory/go-bindings-winrt/bindings/winrt/foundation"
	wrtui "github.com/deploymenttheory/go-bindings-winrt/bindings/winrt/ui"
)

// This Windows backend renders the platform-independent widget tree with
// WinUI 3 (Microsoft.UI.Xaml) through the pure-Go bindings in
// github.com/deploymenttheory/go-bindings-windowsappsdk.
//
// The same rule as the GTK backend applies: the framework's flex engine owns
// layout, and WinUI only draws and dispatches events. A container maps to a
// Canvas and every control is placed at the coordinates the layout engine
// produced (Canvas.Left / Canvas.Top), so windows lay out identically to Linux.
//
// Startup order is owned by the library's app.Run: it locks the UI thread,
// enters a single-threaded apartment, bootstraps the Windows App SDK, calls
// Application.Start and hands back the first Window. Widget construction is
// therefore deferred to that callback, exactly like GTK defers to "activate".
//
// STATUS: cross-compile verified on Linux (GOOS=windows), NOT run on Windows.
// A built application also needs the Windows App SDK runtime, the bootstrapper
// DLL beside the executable, and a generated resources.pri -- see
// internal/build/wasdk.go and the README.

func newPlatformBackend() Backend {
	return &winBackend{}
}

type winBackend struct {
	mu         sync.Mutex
	uiThreadID uint64
	pending    []*winWindow
	dispatcher *uidispatching.IDispatcherQueue
}

func (b *winBackend) Name() string { return "winui3" }

// Init does nothing: the UI thread, apartment and bootstrap are all entered by
// app.Run during Run, which is the only order WinUI accepts.
func (b *winBackend) Init() error { return nil }

func (b *winBackend) CreateWindow(title string, width, height float32) (NativeWindow, error) {
	if width <= 0 {
		width = DefaultAppWidth
	}
	if height <= 0 {
		height = DefaultAppHeight
	}
	gw := &winWindow{
		backend: b,
		title:   title,
		width:   width,
		height:  height,
		root:    &winControl{kind: ControlContainer},
	}
	b.mu.Lock()
	b.pending = append(b.pending, gw)
	b.mu.Unlock()
	return gw, nil
}

func (b *winBackend) CreateControl(w Widget, kind ControlKind, props ControlProps) (NativeControl, error) {
	return newWinControl(w, kind, props)
}

func (b *winBackend) Run() error {
	b.uiThreadID = winGoroutineID()

	b.mu.Lock()
	windows := b.pending
	b.mu.Unlock()

	return app.Run(func(ready *app.Ready) error {
		if len(windows) == 0 {
			return errors.New("winui3: no window was created before Run")
		}
		// The first window is created by app.Run; any extra windows are not
		// supported by the underlying driver yet.
		if dispatcher, err := ready.Window.DispatcherQueue(); err == nil && dispatcher != nil {
			b.mu.Lock()
			b.dispatcher = dispatcher
			b.mu.Unlock()
		}
		return windows[0].build(ready)
	}, app.Options{})
}

func (b *winBackend) Quit() {
	// Closing the window ends the message loop; Application.Exit is wired to the
	// window's Closed event in winWindow.build.
}

func (b *winBackend) Post(fn func()) error {
	if fn == nil {
		return nil
	}
	b.mu.Lock()
	dispatcher := b.dispatcher
	b.mu.Unlock()
	if dispatcher == nil {
		// Before the loop is up there is nowhere to post; run inline.
		fn()
		return nil
	}
	var handler *uidispatching.DispatcherQueueHandler
	handler, err := uidispatching.NewDispatcherQueueHandler(func() {
		if handler != nil {
			handler.Close()
		}
		fn()
	})
	if err != nil {
		return err
	}
	enqueued, err := dispatcher.TryEnqueue(handler)
	if err != nil {
		handler.Close()
		return err
	}
	if !enqueued {
		handler.Close()
		return errors.New("winui3: dispatcher refused the callback")
	}
	return nil
}

func (b *winBackend) OnUIThread() bool {
	return winGoroutineID() == b.uiThreadID
}

// winWindow implements NativeWindow over the Window app.Run hands back.
type winWindow struct {
	backend *winBackend
	title   string
	width   float32
	height  float32
	content Widget
	visible bool
	built   bool
	root    *winControl
	// Window embeds IWindow, so its methods are promoted onto the class.
	window *uixaml.Window
}

func (w *winWindow) Root() NativeControl { return w.root }

func (w *winWindow) SetTitle(title string) {
	w.title = title
	if w.window != nil {
		_ = w.window.SetTitle(title)
	}
}

func (w *winWindow) SetSize(width, height float32) {
	if width <= 0 || height <= 0 {
		return
	}
	w.width, w.height = width, height
	if w.built {
		w.relayout()
	}
}

func (w *winWindow) SetContent(content Widget) error {
	w.content = content
	if w.built {
		w.relayout()
	}
	return nil
}

func (w *winWindow) Show() {
	w.visible = true
	if w.window != nil {
		_ = w.window.Activate()
	}
}

func (w *winWindow) Close() {
	if w.window != nil {
		_ = w.window.Close()
	}
}

// build runs on the UI thread inside app.Run's callback. It creates every
// native control, mounts the widget tree and activates the window.
func (w *winWindow) build(ready *app.Ready) error {
	w.window = ready.Window
	if err := ready.Window.SetTitle(w.title); err != nil {
		return fmt.Errorf("winui3: set title: %w", err)
	}
	// Root surface: a Canvas that hosts absolutely positioned controls.
	if err := w.root.ensureCanvas(); err != nil {
		return err
	}

	w.built = true
	if w.content != nil {
		w.relayout()
	}

	// Closing the window must end the message loop, and only the UI thread may
	// do that, so the exit is wired to the window's own Closed event. The
	// handler stays registered for the window's lifetime.
	if _, err := app.On(ready.Window.AddClosed, uixaml.NewTypedEventHandlerOfObjectAndWindowEventArgs,
		func(_ *syswinrt.IInspectable, _ *uixaml.IWindowEventArgs) {
			_ = ready.Application.Exit()
		}); err != nil {
		return fmt.Errorf("winui3: wire window close: %w", err)
	}

	if w.root.asUIElement == nil {
		return errors.New("winui3: root surface was not created")
	}
	if err := app.With(w.root.asUIElement, ready.Window.SetContent); err != nil {
		return fmt.Errorf("winui3: set window content: %w", err)
	}
	if w.visible {
		return ready.Window.Activate()
	}
	return nil
}

func (w *winWindow) relayout() {
	if w.content == nil || w.root == nil {
		return
	}
	if err := layoutAndMount(w.backend, w.root, w.content, w.width, w.height); err != nil {
		fmt.Fprintf(os.Stderr, "narcissus: layout error: %v\n", err)
	}
}

// winControl implements NativeControl around a single WinUI control.
//
// Only the accessors relevant to a control's kind are set; each queries an
// interface on the underlying class and must be released, which app.With /
// app.Append / app.On already do.
type winControl struct {
	kind ControlKind

	asUIElement        func() (*uixaml.IUIElement, error)
	asFrameworkElement func() (*uixaml.IFrameworkElement, error)
	asControl          func() (*uixaml.IControl, error)
	asPanel            func() (*uixaml.IPanel, error)
	asContentControl   func() (*uixaml.IContentControl, error)
	asButtonBase       func() (*uixaml.IButtonBase, error)
	asToggleButton     func() (*uixaml.IToggleButton, error)
	asTextBlock        func() (*uixaml.ITextBlock, error)
	asComboBox         func() (*uixaml.IComboBox, error)
	asItemsControl     func() (*uixaml.IItemsControl, error)
	asSelector         func() (*uixaml.ISelector, error)
	asRangeBase        func() (*uixaml.IRangeBase, error)
	asTextBox          func() (*uixaml.ITextBox, error)
	asImage            func() (*uixaml.IImage, error)
	asProgressBar      func() (*uixaml.IProgressBar, error)
	asToggleSwitch     func() (*uixaml.IToggleSwitch, error)
	asRadioButton      func() (*uixaml.IRadioButton, error)

	attached  bool
	applying  bool
	items     *app.ItemsSource
	itemsKey  string
	sourceKey string

	// canvas is the concrete Canvas for container controls; it is kept alive for
	// the lifetime of the control tree.
	canvas *uixaml.Canvas
}

// ensureCanvas builds the root Canvas lazily (containers only).
func (c *winControl) ensureCanvas() error {
	if c.asPanel != nil {
		return nil
	}
	canvas, err := uixaml.NewCanvas()
	if err != nil {
		return fmt.Errorf("winui3: create Canvas: %w", err)
	}
	c.canvas = canvas
	c.asPanel = canvas.AsPanel
	c.asUIElement = canvas.AsUIElement
	c.asFrameworkElement = canvas.AsFrameworkElement
	return nil
}

func newWinControl(w Widget, kind ControlKind, props ControlProps) (*winControl, error) {
	c := &winControl{kind: kind}

	switch kind {
	case ControlContainer:
		if err := c.ensureCanvas(); err != nil {
			return nil, err
		}
	case ControlButton:
		btn, err := uixaml.NewButton()
		if err != nil {
			return nil, fmt.Errorf("winui3: create Button: %w", err)
		}
		c.asUIElement = btn.AsUIElement
		c.asFrameworkElement = btn.AsFrameworkElement
		c.asControl = btn.AsControl
		c.asContentControl = btn.AsContentControl
		c.asButtonBase = btn.AsButtonBase
	case ControlText:
		text, err := uixaml.NewTextBlock()
		if err != nil {
			return nil, fmt.Errorf("winui3: create TextBlock: %w", err)
		}
		c.asUIElement = text.AsUIElement
		c.asFrameworkElement = text.AsFrameworkElement
		c.asTextBlock = text.AsTextBlock
	case ControlCheckbox:
		check, err := uixaml.NewCheckBox()
		if err != nil {
			return nil, fmt.Errorf("winui3: create CheckBox: %w", err)
		}
		c.asUIElement = check.AsUIElement
		c.asFrameworkElement = check.AsFrameworkElement
		c.asControl = check.AsControl
		c.asContentControl = check.AsContentControl
		c.asToggleButton = check.AsToggleButton
	case ControlComboBox:
		combo, err := uixaml.NewComboBox()
		if err != nil {
			return nil, fmt.Errorf("winui3: create ComboBox: %w", err)
		}
		c.asUIElement = combo.AsUIElement
		c.asFrameworkElement = combo.AsFrameworkElement
		c.asControl = combo.AsControl
		c.asComboBox = combo.AsComboBox
		c.asItemsControl = combo.AsItemsControl
		c.asSelector = combo.AsSelector
	case ControlSlider:
		slider, err := uixaml.NewSlider()
		if err != nil {
			return nil, fmt.Errorf("winui3: create Slider: %w", err)
		}
		c.asUIElement = slider.AsUIElement
		c.asFrameworkElement = slider.AsFrameworkElement
		c.asControl = slider.AsControl
		c.asRangeBase = slider.AsRangeBase
	case ControlImage:
		img, err := uixaml.NewImage()
		if err != nil {
			return nil, fmt.Errorf("winui3: create Image: %w", err)
		}
		c.asUIElement = img.AsUIElement
		c.asFrameworkElement = img.AsFrameworkElement
		c.asImage = img.AsImage
	case ControlTextInput:
		box, err := uixaml.NewTextBox()
		if err != nil {
			return nil, fmt.Errorf("winui3: create TextBox: %w", err)
		}
		c.asUIElement = box.AsUIElement
		c.asFrameworkElement = box.AsFrameworkElement
		c.asControl = box.AsControl
		c.asTextBox = box.AsTextBox
	case ControlProgress:
		pb, err := uixaml.NewProgressBar()
		if err != nil {
			return nil, fmt.Errorf("winui3: create ProgressBar: %w", err)
		}
		c.asUIElement = pb.AsUIElement
		c.asFrameworkElement = pb.AsFrameworkElement
		c.asControl = pb.AsControl
		c.asRangeBase = pb.AsRangeBase
		c.asProgressBar = pb.AsProgressBar
	case ControlSwitch:
		sw, err := uixaml.NewToggleSwitch()
		if err != nil {
			return nil, fmt.Errorf("winui3: create ToggleSwitch: %w", err)
		}
		c.asUIElement = sw.AsUIElement
		c.asFrameworkElement = sw.AsFrameworkElement
		c.asControl = sw.AsControl
		c.asToggleSwitch = sw.AsToggleSwitch
	case ControlRadio:
		rb, err := uixaml.NewRadioButton()
		if err != nil {
			return nil, fmt.Errorf("winui3: create RadioButton: %w", err)
		}
		c.asUIElement = rb.AsUIElement
		c.asFrameworkElement = rb.AsFrameworkElement
		c.asControl = rb.AsControl
		c.asContentControl = rb.AsContentControl
		c.asToggleButton = rb.AsToggleButton
		c.asRadioButton = rb.AsRadioButton
	default:
		return nil, fmt.Errorf("winui3: unsupported control kind %s", kind)
	} // Event wiring. Each handler calls back into the widget model; the applying
	// flag keeps our own property writes from re-entering the model. The
	// interface a handler is registered through is released once the add call has
	// returned -- the runtime holds its own reference to the handler.
	switch kind {
	case ControlButton:
		base, err := c.asButtonBase()
		if err != nil {
			return nil, fmt.Errorf("winui3: query IButtonBase: %w", err)
		}
		_, err = app.On(base.AddClick, uixaml.NewRoutedEventHandler,
			func(_ *syswinrt.IInspectable, _ *uixaml.IRoutedEventArgs) {
				if c.applying {
					return
				}
				w.HandleEvent(newClickEvent(w))
			})
		base.Release()
		if err != nil {
			return nil, fmt.Errorf("winui3: wire Button.Click: %w", err)
		}
	case ControlCheckbox:
		toggle, err := c.asToggleButton()
		if err != nil {
			return nil, fmt.Errorf("winui3: query IToggleButton: %w", err)
		}
		_, errChecked := app.On(toggle.AddChecked, uixaml.NewRoutedEventHandler,
			func(_ *syswinrt.IInspectable, _ *uixaml.IRoutedEventArgs) {
				if c.applying {
					return
				}
				if box, ok := w.(*Checkbox); ok {
					box.Checked(true)
				}
			})
		_, errUnchecked := app.On(toggle.AddUnchecked, uixaml.NewRoutedEventHandler,
			func(_ *syswinrt.IInspectable, _ *uixaml.IRoutedEventArgs) {
				if c.applying {
					return
				}
				if box, ok := w.(*Checkbox); ok {
					box.Checked(false)
				}
			})
		toggle.Release()
		if err := app.All(errChecked, errUnchecked); err != nil {
			return nil, fmt.Errorf("winui3: wire CheckBox: %w", err)
		}
	case ControlComboBox:
		selector, err := c.asSelector()
		if err != nil {
			return nil, fmt.Errorf("winui3: query ISelector: %w", err)
		}
		_, err = app.On(selector.AddSelectionChanged, uixaml.NewSelectionChangedEventHandler,
			func(_ *syswinrt.IInspectable, _ *uixaml.ISelectionChangedEventArgs) {
				if c.applying {
					return
				}
				combo, ok := w.(*ComboBox)
				if !ok {
					return
				}
				index, err := withValue(c.asSelector, func(sel *uixaml.ISelector) (int32, error) {
					return sel.SelectedIndex()
				})
				if err == nil {
					combo.SelectedIndex(int(index))
				}
			})
		selector.Release()
		if err != nil {
			return nil, fmt.Errorf("winui3: wire ComboBox.SelectionChanged: %w", err)
		}
	case ControlSlider:
		rangeBase, err := c.asRangeBase()
		if err != nil {
			return nil, fmt.Errorf("winui3: query IRangeBase: %w", err)
		}
		_, err = app.On(rangeBase.AddValueChanged, uixaml.NewRangeBaseValueChangedEventHandler,
			func(_ *syswinrt.IInspectable, _ *uixaml.IRangeBaseValueChangedEventArgs) {
				if c.applying {
					return
				}
				slider, ok := w.(*Slider)
				if !ok {
					return
				}
				value, err := withValue(c.asRangeBase, func(rb *uixaml.IRangeBase) (float64, error) {
					return rb.Value()
				})
				if err == nil {
					slider.Value(float32(value))
				}
			})
		rangeBase.Release()
		if err != nil {
			return nil, fmt.Errorf("winui3: wire Slider.ValueChanged: %w", err)
		}
	case ControlTextInput:
		textBox, err := c.asTextBox()
		if err != nil {
			return nil, fmt.Errorf("winui3: query ITextBox: %w", err)
		}
		_, err = app.On(textBox.AddTextChanged, uixaml.NewTextChangedEventHandler,
			func(_ *syswinrt.IInspectable, _ *uixaml.ITextChangedEventArgs) {
				if c.applying {
					return
				}
				input, ok := w.(*TextInput)
				if !ok {
					return
				}
				value, err := withValue(c.asTextBox, func(box *uixaml.ITextBox) (string, error) {
					return box.Text()
				})
				if err == nil {
					input.SetText(value)
				}
			})
		textBox.Release()
		if err != nil {
			return nil, fmt.Errorf("winui3: wire TextBox.TextChanged: %w", err)
		}
	case ControlSwitch:
		toggle, err := c.asToggleSwitch()
		if err != nil {
			return nil, fmt.Errorf("winui3: query IToggleSwitch: %w", err)
		}
		_, err = app.On(toggle.AddToggled, uixaml.NewRoutedEventHandler,
			func(_ *syswinrt.IInspectable, _ *uixaml.IRoutedEventArgs) {
				if c.applying {
					return
				}
				switchWidget, ok := w.(*Switch)
				if !ok {
					return
				}
				on, readErr := withValue(c.asToggleSwitch, func(t *uixaml.IToggleSwitch) (bool, error) {
					return t.IsOn()
				})
				if readErr == nil {
					switchWidget.Checked(on)
				}
			})
		toggle.Release()
		if err != nil {
			return nil, fmt.Errorf("winui3: wire ToggleSwitch.Toggled: %w", err)
		}
	case ControlRadio:
		toggle, err := c.asToggleButton()
		if err != nil {
			return nil, fmt.Errorf("winui3: query IToggleButton: %w", err)
		}
		_, err = app.On(toggle.AddChecked, uixaml.NewRoutedEventHandler,
			func(_ *syswinrt.IInspectable, _ *uixaml.IRoutedEventArgs) {
				if c.applying {
					return
				}
				if radio, ok := w.(*RadioButton); ok {
					radio.Checked(true)
				}
			})
		toggle.Release()
		if err != nil {
			return nil, fmt.Errorf("winui3: wire RadioButton.Checked: %w", err)
		}
	}

	c.applyProps(props)
	return c, nil
}

func (c *winControl) AttachTo(parent NativeControl) {
	p, ok := parent.(*winControl)
	if !ok || p == nil || p.asPanel == nil || c.asUIElement == nil {
		return
	}
	if c.attached {
		return
	}
	if err := app.Append(p.asPanel, c.asUIElement); err != nil {
		return
	}
	c.attached = true
}

func (c *winControl) SetBounds(x, y, width, height float32) {
	if c.asUIElement != nil {
		if err := app.With(c.asUIElement, func(element *uixaml.IUIElement) error {
			statics, err := uixaml.CanvasStatics()
			if err != nil {
				return err
			}
			defer statics.Release()
			return app.All(
				statics.SetLeft(element, float64(x)),
				statics.SetTop(element, float64(y)),
			)
		}); err != nil {
			return
		}
	}
	if width > 0 && height > 0 && c.asFrameworkElement != nil {
		_ = app.With(c.asFrameworkElement, func(fe *uixaml.IFrameworkElement) error {
			return app.All(fe.SetWidth(float64(width)), fe.SetHeight(float64(height)))
		})
	}
}

func (c *winControl) SetVisible(visible bool) {
	if c.asUIElement == nil {
		return
	}
	_ = app.With(c.asUIElement, func(element *uixaml.IUIElement) error {
		visibility := uixaml.VisibilityCollapsed
		if visible {
			visibility = uixaml.VisibilityVisible
		}
		return element.SetVisibility(visibility)
	})
}

func (c *winControl) SetEnabled(enabled bool) {
	if c.asControl == nil {
		return
	}
	_ = app.With(c.asControl, func(control *uixaml.IControl) error {
		return control.SetIsEnabled(enabled)
	})
}

func (c *winControl) SetStyle(style *Style) {
	if style == nil {
		return
	}
	if c.asControl != nil {
		_ = app.With(c.asControl, func(control *uixaml.IControl) error {
			errs := []error{}
			if style.Font.Size > 0 {
				errs = append(errs, control.SetFontSize(float64(style.Font.Size)))
			}
			if hasInsets(style.Padding) {
				errs = append(errs, control.SetPadding(thicknessFor(style.Padding)))
			}
			if style.TextColor.A > 0 {
				brush, err := solidColorBrush(style.TextColor)
				if err == nil {
					defer brush.Release()
					errs = append(errs, control.SetForeground(brush))
				}
			}
			if style.BackgroundColor.A > 0 {
				brush, err := solidColorBrush(style.BackgroundColor)
				if err == nil {
					defer brush.Release()
					errs = append(errs, control.SetBackground(brush))
				}
			}
			return app.All(errs...)
		})
		return
	}
	if c.asTextBlock != nil {
		_ = app.With(c.asTextBlock, func(text *uixaml.ITextBlock) error {
			errs := []error{}
			if style.Font.Size > 0 {
				errs = append(errs, text.SetFontSize(float64(style.Font.Size)))
			}
			if style.TextColor.A > 0 {
				brush, err := solidColorBrush(style.TextColor)
				if err == nil {
					defer brush.Release()
					errs = append(errs, text.SetForeground(brush))
				}
			}
			return app.All(errs...)
		})
		return
	}
	if c.asPanel != nil && style.BackgroundColor.A > 0 {
		_ = app.With(c.asPanel, func(panel *uixaml.IPanel) error {
			brush, err := solidColorBrush(style.BackgroundColor)
			if err != nil {
				return nil
			}
			defer brush.Release()
			return panel.SetBackground(brush)
		})
	}
}

func (c *winControl) Update(props ControlProps) {
	c.applying = true
	defer func() { c.applying = false }()
	c.applyProps(props)
}

func (c *winControl) Destroy() {
	if c.items != nil {
		c.items.Close()
		c.items = nil
	}
}

// applyProps pushes a property snapshot onto the native controls.
func (c *winControl) applyProps(props ControlProps) {
	switch c.kind {
	case ControlButton:
		if c.asContentControl != nil {
			_ = app.SetContent(c.asContentControl, props.Text)
		}
	case ControlText:
		if c.asTextBlock != nil {
			_ = app.With(c.asTextBlock, func(text *uixaml.ITextBlock) error {
				return text.SetText(props.Text)
			})
		}
	case ControlCheckbox:
		if c.asContentControl != nil {
			_ = app.SetContent(c.asContentControl, props.Text)
		}
		if c.asToggleButton != nil {
			checked, err := app.BoxAs[uixaml.IReferenceOfBool](props.Checked, &uixaml.IID_IReferenceOfBool)
			if err == nil {
				_ = app.With(c.asToggleButton, func(toggle *uixaml.IToggleButton) error {
					return toggle.SetIsChecked(checked)
				})
				checked.Release()
			}
		}
	case ControlComboBox:
		c.setComboItems(props.Items)
		if props.Placeholder != "" && c.asComboBox != nil {
			_ = app.With(c.asComboBox, func(combo *uixaml.IComboBox) error {
				return combo.SetPlaceholderText(props.Placeholder)
			})
		}
		if c.asSelector != nil && props.Selected >= 0 && props.Selected < len(props.Items) {
			_ = app.With(c.asSelector, func(selector *uixaml.ISelector) error {
				return selector.SetSelectedIndex(int32(props.Selected))
			})
		}
	case ControlSlider:
		if c.asRangeBase != nil {
			_ = app.With(c.asRangeBase, func(rb *uixaml.IRangeBase) error {
				return app.All(
					rb.SetMinimum(props.Min),
					rb.SetMaximum(props.Max),
					rb.SetValue(props.Value),
				)
			})
		}
	case ControlTextInput:
		if c.asTextBox != nil {
			_ = app.With(c.asTextBox, func(box *uixaml.ITextBox) error {
				return app.All(
					box.SetText(props.Text),
					box.SetPlaceholderText(props.Placeholder),
					box.SetIsReadOnly(props.ReadOnly),
				)
			})
		}
	case ControlImage:
		c.setImageSource(props.Source)
	case ControlProgress:
		if c.asRangeBase != nil {
			_ = app.With(c.asRangeBase, func(rb *uixaml.IRangeBase) error {
				return app.All(
					rb.SetMinimum(props.Min),
					rb.SetMaximum(props.Max),
					rb.SetValue(props.Value),
				)
			})
		}
		if c.asProgressBar != nil {
			_ = app.With(c.asProgressBar, func(pb *uixaml.IProgressBar) error {
				return pb.SetIsIndeterminate(props.Indeterminate)
			})
		}
	case ControlSwitch:
		if c.asToggleSwitch != nil {
			_ = app.With(c.asToggleSwitch, func(t *uixaml.IToggleSwitch) error {
				return t.SetIsOn(props.Checked)
			})
		}
	case ControlRadio:
		if c.asContentControl != nil {
			_ = app.SetContent(c.asContentControl, props.Text)
		}
		if c.asToggleButton != nil {
			checked, err := app.BoxAs[uixaml.IReferenceOfBool](props.Checked, &uixaml.IID_IReferenceOfBool)
			if err == nil {
				_ = app.With(c.asToggleButton, func(toggle *uixaml.IToggleButton) error {
					return toggle.SetIsChecked(checked)
				})
				checked.Release()
			}
		}
		if c.asRadioButton != nil && props.Group != "" {
			_ = app.With(c.asRadioButton, func(rb *uixaml.IRadioButton) error {
				return rb.SetGroupName(props.Group)
			})
		}
	}
}

// setComboItems rebuilds the ItemsSource only when the list actually changed.
func (c *winControl) setComboItems(items []string) {
	if c.asItemsControl == nil {
		return
	}
	key := strings.Join(items, "\x00")
	if key == c.itemsKey {
		return
	}
	c.itemsKey = key

	source, err := app.NewStringItemsSource(items, winrt.CollectionIIDs{
		Iterable:   uixaml.IID_IIterableOfObject,
		Iterator:   uixaml.IID_IIteratorOfObject,
		VectorView: uixaml.IID_IVectorViewOfObject,
		Vector:     uixaml.IID_IVectorOfObject,
	})
	if err != nil {
		return
	}
	if old := c.items; old != nil {
		old.Close()
	}
	c.items = source
	_ = app.With(c.asItemsControl, func(control *uixaml.IItemsControl) error {
		return control.SetItemsSource(source.Inspectable())
	})
}

func (c *winControl) setImageSource(source string) {
	if c.asImage == nil || source == "" || source == c.sourceKey {
		return
	}
	c.sourceKey = source

	uri, err := wrtfoundation.CreateUri(source)
	if err != nil {
		return
	}
	defer uri.Release()

	bitmap, err := uixaml.NewBitmapImage()
	if err != nil {
		return
	}
	// Uri embeds IUriRuntimeClass, which is exactly what SetUriSource takes.
	if err := app.With(bitmap.AsBitmapImage, func(image *uixaml.IBitmapImage) error {
		return image.SetUriSource(&uri.IUriRuntimeClass)
	}); err != nil {
		return
	}
	imageSource, err := bitmap.AsImageSource()
	if err != nil {
		return
	}
	defer imageSource.Release()

	_ = app.With(c.asImage, func(image *uixaml.IImage) error {
		return image.SetSource(imageSource)
	})
}

// solidColorBrush creates a WinUI brush from a framework Color.
func solidColorBrush(color Color) (*uixaml.IBrush, error) {
	brush, err := uixaml.NewSolidColorBrush()
	if err != nil {
		return nil, err
	}
	if err := app.With(brush.AsSolidColorBrush, func(scb *uixaml.ISolidColorBrush) error {
		return scb.SetColor(wrtui.Color{A: color.A, R: color.R, G: color.G, B: color.B})
	}); err != nil {
		return nil, err
	}
	return brush.AsBrush()
}

// withValue queries an interface, runs fn against it, and releases it. It is the
// value-returning counterpart to app.With, which only returns an error.
func withValue[T app.Releaser, V any](get func() (T, error), fn func(T) (V, error)) (V, error) {
	var zero V
	value, err := get()
	if err != nil {
		return zero, err
	}
	defer value.Release()
	return fn(value)
}

// newClickEvent builds the click event dispatched through the widget model.
func newClickEvent(target Widget) *BaseEvent {
	return &BaseEvent{
		Type:      EventClick,
		Target:    target,
		Timestamp: time.Now().UnixMilli(),
	}
}

func thicknessFor(insets Insets) uixaml.Thickness {
	return uixaml.Thickness{
		Left:   float64(insets.Left),
		Top:    float64(insets.Top),
		Right:  float64(insets.Right),
		Bottom: float64(insets.Bottom),
	}
}

func hasInsets(insets Insets) bool {
	return insets.Top != 0 || insets.Right != 0 || insets.Bottom != 0 || insets.Left != 0
}

func winGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	const prefix = "goroutine "
	stack := string(buf[:n])
	if !strings.HasPrefix(stack, prefix) {
		return 0
	}
	rest := stack[len(prefix):]
	end := strings.IndexByte(rest, ' ')
	if end < 0 {
		end = len(rest)
	}
	id, _ := strconv.ParseUint(rest[:end], 10, 64)
	return id
}
