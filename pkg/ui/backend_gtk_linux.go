//go:build linux

package ui

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// newPlatformBackend returns the default backend for Linux.
func newPlatformBackend() Backend {
	return &gtkBackend{}
}

// gtkBackend implements Backend on top of GTK4.
//
// Layout stays owned by the framework's flex engine: every control is placed at
// absolute coordinates inside a GtkFixed surface, so Linux renders identically
// to the other backends. GTK is only used to create and draw the native
// controls and to deliver events.
//
// GTK4 requires top-level windows to be created from the GtkApplication's
// "activate" handler, so CreateWindow builds the widget tree eagerly but the
// GtkApplicationWindow itself is materialized in onActivate.
type gtkBackend struct {
	app         *gtk.Application
	uiThreadID  uint64
	initialized bool
	activated   bool
	windows     []*gtkWindow
}

func (b *gtkBackend) Name() string { return "gtk4" }

func (b *gtkBackend) Init() error {
	if b.initialized {
		return nil
	}
	b.app = gtk.NewApplication("com.narcissus.app", gio.ApplicationDefaultFlags)
	if b.app == nil {
		return errors.New("gtk4: failed to create GtkApplication (is a display available?)")
	}
	b.app.ConnectActivate(b.onActivate)
	// Measure text with Pango so the flex engine lays out against the same
	// metrics GTK uses to draw labels.
	SetTextMeasurer(newGTKTextMeasurer())
	b.uiThreadID = currentGoroutineID()
	b.initialized = true
	return nil
}

func (b *gtkBackend) onActivate() {
	b.activated = true
	for _, w := range b.windows {
		w.materialize()
	}
}

func (b *gtkBackend) CreateWindow(title string, width, height float32) (NativeWindow, error) {
	if b.app == nil {
		return nil, errors.New("gtk4: backend not initialized")
	}
	if width <= 0 {
		width = DefaultAppWidth
	}
	if height <= 0 {
		height = DefaultAppHeight
	}

	// GTK widgets may only be created after gtk_init has run, which happens
	// inside GtkApplication.Run. The window therefore records its description
	// here and builds the widget tree on activate (see materialize).
	gw := &gtkWindow{
		backend: b,
		title:   title,
		width:   width,
		height:  height,
	}

	b.windows = append(b.windows, gw)
	if b.activated {
		gw.materialize()
	}

	return gw, nil
}

func (b *gtkBackend) CreateControl(w Widget, kind ControlKind, props ControlProps) (NativeControl, error) {
	return newGTKControl(b, w, kind, props)
}

func (b *gtkBackend) Run() error {
	if b.app == nil {
		return errors.New("gtk4: backend not initialized")
	}
	argv := os.Args
	if len(argv) == 0 {
		argv = []string{"narcissus"}
	}
	b.app.Run(argv)
	return nil
}

func (b *gtkBackend) Quit() {
	if b.app != nil {
		b.app.Quit()
	}
}

func (b *gtkBackend) Post(fn func()) error {
	if fn == nil {
		return nil
	}
	glib.IdleAdd(func() { fn() })
	return nil
}

func (b *gtkBackend) OnUIThread() bool {
	return currentGoroutineID() == b.uiThreadID
}

// gtkWindow implements NativeWindow.
type gtkWindow struct {
	backend *gtkBackend

	title   string
	width   float32
	height  float32
	content Widget
	visible bool
	built   bool

	win     *gtk.ApplicationWindow
	overlay *gtk.Overlay
	area    *gtk.DrawingArea
	root    *gtkControl
}

// materialize builds the GtkApplicationWindow and its widget tree. It must run
// inside the application's activate handler, after gtk_init.
func (gw *gtkWindow) materialize() {
	if gw.built {
		return
	}

	overlay := gtk.NewOverlay()
	// A drawing area as the base child gives us a reliable resize signal; it
	// does not paint and sits underneath the fixed surface.
	area := gtk.NewDrawingArea()
	overlay.SetChild(area)

	// The root surface (a GtkFixed) is added as an overlay child so it fills
	// the window while hosting absolutely positioned controls.
	root := newGTKSurface(gw.backend)
	overlay.AddOverlay(root.widget)

	win := gtk.NewApplicationWindow(gw.backend.app)
	win.SetTitle(gw.title)
	win.SetDefaultSize(int(gw.width), int(gw.height))
	win.SetChild(overlay)

	// Keyboard input is routed to the focused widget; TAB traversal and
	// bubbling live in the framework (see focus.go).
	key := gtk.NewEventControllerKey()
	key.ConnectKeyPressed(func(keyval, keycode uint, state gdk.ModifierType) bool {
		return dispatchGTKKey(EventKeyDown, keyval, keycode, state)
	})
	key.ConnectKeyReleased(func(keyval, keycode uint, state gdk.ModifierType) {
		dispatchGTKKey(EventKeyUp, keyval, keycode, state)
	})
	win.AddController(key)

	gw.overlay, gw.area, gw.root, gw.win = overlay, area, root, win
	gw.built = true

	area.ConnectResize(func(w, h int) {
		gw.relayout(float32(w), float32(h))
	})

	if gw.content != nil {
		gw.relayout(gw.width, gw.height)
	}
	if gw.visible {
		win.Present()
	}
}

func (gw *gtkWindow) Root() NativeControl { return gw.root }

func (gw *gtkWindow) SetTitle(title string) {
	gw.title = title
	if gw.win != nil {
		gw.win.SetTitle(title)
	}
}

func (gw *gtkWindow) SetSize(width, height float32) {
	if width <= 0 || height <= 0 {
		return
	}
	gw.width, gw.height = width, height
	if gw.win != nil {
		gw.win.SetDefaultSize(int(width), int(height))
	}
	gw.relayout(width, height)
}

func (gw *gtkWindow) SetContent(content Widget) error {
	gw.content = content
	if gw.built {
		gw.relayout(gw.width, gw.height)
	}
	return nil
}

func (gw *gtkWindow) relayout(width, height float32) {
	if gw.content == nil || gw.root == nil {
		return
	}
	if width <= 0 {
		width = gw.width
	}
	if height <= 0 {
		height = gw.height
	}
	if err := layoutAndMount(gw.backend, gw.root, gw.content, width, height); err != nil {
		fmt.Fprintf(os.Stderr, "narcissus: layout error: %v\n", err)
	}
}

func (gw *gtkWindow) Show() {
	gw.visible = true
	if gw.win != nil {
		gw.win.Present()
	}
}

func (gw *gtkWindow) Close() {
	if gw.win != nil {
		gw.win.Close()
	}
}

// gtkControl implements NativeControl around a single GTK widget.
type gtkControl struct {
	backend *gtkBackend
	kind    ControlKind

	widget gtk.Widgetter
	base   *gtk.Widget
	// surface is non-nil for containers: it is the GtkFixed that hosts children.
	surface *gtk.Fixed
	// parentFixed is the surface this control currently lives in.
	parentFixed *gtk.Fixed
	attached    bool

	className string
	provider  *gtk.CSSProvider
	applying  bool

	// Typed handles for property updates and event wiring.
	button   *gtk.Button
	label    *gtk.Label
	check    *gtk.CheckButton
	dropdown *gtk.DropDown
	scale    *gtk.Scale
	picture  *gtk.Picture
	entry    *gtk.Entry
	textView *gtk.TextView
	textBuf  *gtk.TextBuffer
}

// newGTKSurface wraps a fresh GtkFixed as the root surface of a window.
func newGTKSurface(b *gtkBackend) *gtkControl {
	fixed := gtk.NewFixed()
	base := gtk.BaseWidget(fixed)
	return &gtkControl{
		backend:   b,
		kind:      ControlContainer,
		widget:    fixed,
		base:      base,
		surface:   fixed,
		className: "narc-surface",
	}
}

func newGTKControl(b *gtkBackend, w Widget, kind ControlKind, props ControlProps) (*gtkControl, error) {
	c := &gtkControl{
		backend:   b,
		kind:      kind,
		className: "narc-" + sanitizeClass(w.GetID()),
	}

	switch kind {
	case ControlContainer:
		fixed := gtk.NewFixed()
		c.widget = fixed
		c.surface = fixed
	case ControlButton:
		btn := gtk.NewButtonWithLabel(props.Text)
		c.button = btn
		c.widget = btn
		btn.ConnectClicked(func() {
			if c.applying {
				return
			}
			w.HandleEvent(newClickEvent(w))
		})
	case ControlText:
		label := gtk.NewLabel(props.Text)
		label.SetXAlign(0)
		label.SetWrap(true)
		c.label = label
		c.widget = label
	case ControlCheckbox:
		cb := gtk.NewCheckButtonWithLabel(props.Text)
		cb.SetActive(props.Checked)
		c.check = cb
		c.widget = cb
		cb.ConnectToggled(func() {
			if c.applying {
				return
			}
			if box, ok := w.(*Checkbox); ok {
				box.Checked(cb.Active())
			}
		})
	case ControlComboBox:
		dd := gtk.NewDropDownFromStrings(props.Items)
		if props.Selected >= 0 && props.Selected < len(props.Items) {
			dd.SetSelected(uint(props.Selected))
		}
		c.dropdown = dd
		c.widget = dd
		dd.NotifyProperty("selected", func() {
			if c.applying {
				return
			}
			if combo, ok := w.(*ComboBox); ok {
				combo.SelectedIndex(int(dd.Selected()))
			}
		})
	case ControlSlider:
		scale := gtk.NewScaleWithRange(gtk.OrientationHorizontal, props.Min, props.Max, props.Step)
		scale.SetValue(props.Value)
		c.scale = scale
		c.widget = scale
		scale.ConnectValueChanged(func() {
			if c.applying {
				return
			}
			if slider, ok := w.(*Slider); ok {
				slider.Value(float32(scale.Value()))
			}
		})
	case ControlImage:
		pic := gtk.NewPicture()
		pic.SetCanShrink(true)
		if props.Source != "" {
			pic.SetFilename(props.Source)
		}
		c.picture = pic
		c.widget = pic
	case ControlTextInput:
		if props.Multiline {
			tv := gtk.NewTextView()
			tv.SetWrapMode(gtk.WrapWordChar)
			buf := tv.Buffer()
			buf.SetText(props.Text)
			c.textView = tv
			c.textBuf = buf
			c.widget = tv
			buf.ConnectChanged(func() {
				if c.applying {
					return
				}
				if input, ok := w.(*TextInput); ok {
					input.SetText(bufferText(buf))
				}
			})
		} else {
			entry := gtk.NewEntry()
			entry.SetPlaceholderText(props.Placeholder)
			entry.SetText(props.Text)
			c.entry = entry
			c.widget = entry
			entry.ConnectChanged(func() {
				if c.applying {
					return
				}
				if input, ok := w.(*TextInput); ok {
					input.SetText(entry.Text())
				}
			})
		}
	default:
		return nil, fmt.Errorf("gtk4: unsupported control kind %s", kind)
	}

	c.base = gtk.BaseWidget(c.widget)
	c.applyProps(props)

	// Focusable controls join the toolkit's focus chain and report focus
	// changes back into the framework's focus manager.
	if bw := widgetBase(w); bw != nil && bw.IsFocusable() {
		c.base.SetCanFocus(true)
		c.base.SetFocusOnClick(true)
		focus := gtk.NewEventControllerFocus()
		focus.ConnectEnter(func() { notifyNativeFocus(w) })
		focus.ConnectLeave(func() { notifyNativeBlur(w) })
		c.base.AddController(focus)
	}

	return c, nil
}

// applyProps pushes a property snapshot onto the concrete GTK widgets.
func (c *gtkControl) applyProps(props ControlProps) {
	switch c.kind {
	case ControlButton:
		if c.button != nil {
			c.button.SetLabel(props.Text)
		}
	case ControlText:
		if c.label != nil {
			c.label.SetText(props.Text)
		}
	case ControlCheckbox:
		if c.check != nil {
			c.check.SetActive(props.Checked)
			c.check.SetLabel(props.Text)
		}
	case ControlComboBox:
		if c.dropdown != nil {
			c.dropdown.SetModel(gtk.NewStringList(props.Items))
			if props.Selected >= 0 && props.Selected < len(props.Items) {
				c.dropdown.SetSelected(uint(props.Selected))
			}
		}
	case ControlSlider:
		if c.scale != nil {
			c.scale.SetRange(props.Min, props.Max)
			c.scale.SetValue(props.Value)
		}
	case ControlImage:
		if c.picture != nil && props.Source != "" {
			c.picture.SetFilename(props.Source)
			c.picture.SetContentFit(contentFitFor(ImageFit(props.Fit)))
		}
	case ControlTextInput:
		if c.entry != nil {
			c.entry.SetPlaceholderText(props.Placeholder)
			c.entry.SetText(props.Text)
		}
		if c.textBuf != nil && bufferText(c.textBuf) != props.Text {
			c.textBuf.SetText(props.Text)
		}
	}
}

func (c *gtkControl) AttachTo(parent NativeControl) {
	p, ok := parent.(*gtkControl)
	if !ok || p == nil || p.surface == nil || c.widget == nil {
		return
	}
	if c.attached && c.parentFixed == p.surface {
		return
	}
	if c.attached && c.parentFixed != nil {
		c.parentFixed.Remove(c.widget)
		c.attached = false
	}
	p.surface.Put(c.widget, 0, 0)
	c.parentFixed = p.surface
	c.attached = true
}

func (c *gtkControl) SetBounds(x, y, width, height float32) {
	if c.base != nil && width > 0 && height > 0 {
		c.base.SetSizeRequest(int(width), int(height))
	}
	if c.parentFixed != nil {
		if !c.attached {
			c.parentFixed.Put(c.widget, float64(x), float64(y))
			c.attached = true
		} else {
			c.parentFixed.Move(c.widget, float64(x), float64(y))
		}
	}
}

func (c *gtkControl) SetVisible(visible bool) {
	if c.base != nil {
		c.base.SetVisible(visible)
	}
}

func (c *gtkControl) SetEnabled(enabled bool) {
	if c.base != nil {
		c.base.SetSensitive(enabled)
	}
}

func (c *gtkControl) Update(props ControlProps) {
	if c.base == nil {
		return
	}
	c.applying = true
	defer func() { c.applying = false }()
	c.applyProps(props)
}

func (c *gtkControl) SetStyle(style *Style) {
	if c.base == nil || style == nil {
		return
	}
	if c.kind == ControlContainer && !hasBoxStyle(style) {
		// Plain layout surfaces don't need a CSS rule.
		return
	}
	if c.provider == nil {
		c.provider = gtk.NewCSSProvider()
		display := c.base.Display()
		if display == nil {
			display = gdk.DisplayGetDefault()
		}
		if display != nil {
			gtk.StyleContextAddProviderForDisplay(display, c.provider, gtk.STYLE_PROVIDER_PRIORITY_USER)
		}
		c.base.AddCSSClass(c.className)
	}
	c.provider.LoadFromString(cssFor(c.className, style))
}

// RequestFocus gives toolkit focus to the control. It satisfies the
// focusRequester interface consumed by the focus manager.
func (c *gtkControl) RequestFocus() {
	if c.base != nil {
		c.base.GrabFocus()
	}
}

func (c *gtkControl) Destroy() {
	if c.attached && c.parentFixed != nil {
		c.parentFixed.Remove(c.widget)
	}
	c.attached = false
	c.parentFixed = nil
}

// hasBoxStyle reports whether a container has any visual style worth emitting.
func hasBoxStyle(s *Style) bool {
	return s.BackgroundColor.A > 0 || s.Border.Width > 0 || s.Border.Radius > 0
}

// cssFor builds the CSS rule applied to a control.
func cssFor(className string, s *Style) string {
	var b strings.Builder
	fmt.Fprintf(&b, ".%s {", className)
	if s.BackgroundColor.A > 0 {
		fmt.Fprintf(&b, "background-color: %s;", s.BackgroundColor.Hex())
	}
	if s.TextColor.A > 0 {
		fmt.Fprintf(&b, "color: %s;", s.TextColor.Hex())
	}
	if s.Padding.Top > 0 || s.Padding.Right > 0 || s.Padding.Bottom > 0 || s.Padding.Left > 0 {
		fmt.Fprintf(&b, "padding: %gpx %gpx %gpx %gpx;",
			s.Padding.Top, s.Padding.Right, s.Padding.Bottom, s.Padding.Left)
	}
	if s.Border.Radius > 0 {
		fmt.Fprintf(&b, "border-radius: %gpx;", s.Border.Radius)
	}
	if s.Border.Width > 0 {
		fmt.Fprintf(&b, "border: %gpx solid %s;", s.Border.Width, s.Border.Color.Hex())
	}
	if s.Font.Size > 0 {
		fmt.Fprintf(&b, "font-size: %gpx;", s.Font.Size)
	}
	if s.Font.Family != "" {
		fmt.Fprintf(&b, "font-family: \"%s\";", s.Font.Family)
	}
	if s.Font.Weight == FontWeightBold {
		b.WriteString("font-weight: bold;")
	}
	if s.Font.Style == FontStyleItalic {
		b.WriteString("font-style: italic;")
	}
	if s.Font.LineHeight > 0 {
		fmt.Fprintf(&b, "line-height: %g;", s.Font.LineHeight)
	}
	if s.Opacity > 0 && s.Opacity < 1 {
		fmt.Fprintf(&b, "opacity: %g;", s.Opacity)
	}
	b.WriteString("}")
	return b.String()
}

func contentFitFor(fit ImageFit) gtk.ContentFit {
	switch fit {
	case ImageFitCover:
		return gtk.ContentFitCover
	case ImageFitFill:
		return gtk.ContentFitFill
	case ImageFitNone:
		return gtk.ContentFitScaleDown
	case ImageFitScaleDown:
		return gtk.ContentFitScaleDown
	default:
		return gtk.ContentFitContain
	}
}

func bufferText(buf *gtk.TextBuffer) string {
	if buf == nil {
		return ""
	}
	return buf.Text(buf.StartIter(), buf.EndIter(), true)
}

// dispatchGTKKey translates a GTK key event into a framework KeyEvent and
// routes it through the focus manager.
func dispatchGTKKey(kind EventType, keyval, keycode uint, state gdk.ModifierType) bool {
	ev := NewKeyEvent(
		kind,
		gdk.KeyvalName(keyval),
		int(keycode),
		state&gdk.ControlMask != 0,
		state&gdk.ShiftMask != 0,
		state&gdk.AltMask != 0,
		state&gdk.MetaMask != 0,
	)
	return DispatchKey(ev)
}

func newClickEvent(target Widget) *BaseEvent {
	return &BaseEvent{
		Type:      EventClick,
		Target:    target,
		Timestamp: time.Now().UnixMilli(),
	}
}

// sanitizeClass converts a widget ID into a valid CSS class suffix.
func sanitizeClass(id string) string {
	var b strings.Builder
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "widget"
	}
	return b.String()
}

// currentGoroutineID extracts the current goroutine ID from the runtime stack.
// It is used to answer "am I on the toolkit UI thread?".
func currentGoroutineID() uint64 {
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
