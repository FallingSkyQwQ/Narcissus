package ui

import (
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// focusControl is a fake control that records focus requests.
type focusControl struct {
	*fakeControl
	requested int
}

func (c *focusControl) RequestFocus() { c.requested++ }

// focusBackend hands out focusControls so tests can observe native focus.
type focusBackend struct {
	created []*focusControl
}

func (b *focusBackend) Name() string { return "focus-fake" }
func (b *focusBackend) Init() error  { return nil }

func (b *focusBackend) CreateWindow(string, float32, float32) (NativeWindow, error) {
	return nil, nil
}

func (b *focusBackend) CreateControl(w Widget, kind ControlKind, props ControlProps) (NativeControl, error) {
	c := &focusControl{fakeControl: &fakeControl{kind: kind, props: props}}
	b.created = append(b.created, c)
	return c, nil
}

func (b *focusBackend) Run() error           { return nil }
func (b *focusBackend) Quit()                {}
func (b *focusBackend) Post(fn func()) error { fn(); return nil }
func (b *focusBackend) OnUIThread() bool     { return true }

func resetFocusState() {
	RequestFocus(nil)
	SetFocusRoot(nil)
}

func mountForFocus(t *testing.T, root Widget) *focusBackend {
	t.Helper()
	backend := &focusBackend{}
	if err := layoutAndMount(backend, &fakeControl{kind: ControlContainer}, root, 400, 300); err != nil {
		t.Fatalf("layoutAndMount: %v", err)
	}
	return backend
}

func TestFocusOrderSkipsNonFocusable(t *testing.T) {
	defer resetFocusState()

	label := NewText("label")
	button := NewButton().Text("ok")
	check := NewCheckbox().Label("check")
	root := NewContainer().Direction(flex.DirectionColumn).Add(label, button, check)

	mountForFocus(t, root)

	order := FocusOrder(root)
	if len(order) != 2 {
		t.Fatalf("focus order length = %d, want 2", len(order))
	}
	if order[0] != Widget(button) || order[1] != Widget(check) {
		t.Errorf("unexpected focus order: %T/%T", order[0], order[1])
	}
}

func TestRequestFocusDispatchesEventsAndNativeFocus(t *testing.T) {
	defer resetFocusState()

	button := NewButton().Text("ok")
	other := NewButton().Text("other")
	root := NewContainer().Add(button, other)

	mountForFocus(t, root)

	var seen []EventType
	button.AddEventHandler(EventFocus, func(e Event) { seen = append(seen, e.GetType()) })
	button.AddEventHandler(EventBlur, func(e Event) { seen = append(seen, e.GetType()) })

	RequestFocus(button)
	if FocusedWidget() != Widget(button) {
		t.Fatal("button did not receive focus")
	}
	ctrl, ok := button.NativeControl().(*focusControl)
	if !ok {
		t.Fatal("unexpected native control type")
	}
	if ctrl.requested != 1 {
		t.Errorf("native focus requests = %d, want 1", ctrl.requested)
	}

	RequestFocus(other)
	if FocusedWidget() != Widget(other) {
		t.Fatal("focus did not move to other")
	}

	if len(seen) != 2 || seen[0] != EventFocus || seen[1] != EventBlur {
		t.Errorf("focus/blur events = %v, want [Focus Blur]", seen)
	}
}

func TestRequestFocusRejectsNonFocusable(t *testing.T) {
	defer resetFocusState()

	label := NewText("label")
	root := NewContainer().Add(label)
	mountForFocus(t, root)

	RequestFocus(label)
	if FocusedWidget() != nil {
		t.Errorf("non-focusable label took focus: %T", FocusedWidget())
	}
}

func TestFocusTraversalWrapsAround(t *testing.T) {
	defer resetFocusState()

	first := NewButton().Text("1")
	second := NewButton().Text("2")
	root := NewContainer().Add(first, second)
	mountForFocus(t, root)

	FocusNext()
	if FocusedWidget() != Widget(first) {
		t.Fatalf("first FocusNext did not focus the first widget")
	}
	FocusNext()
	if FocusedWidget() != Widget(second) {
		t.Fatalf("second FocusNext did not focus the second widget")
	}
	FocusNext()
	if FocusedWidget() != Widget(first) {
		t.Errorf("FocusNext did not wrap around, got %T", FocusedWidget())
	}
	FocusPrevious()
	if FocusedWidget() != Widget(second) {
		t.Errorf("FocusPrevious did not wrap back, got %T", FocusedWidget())
	}
}

func TestTabKeyMovesFocus(t *testing.T) {
	defer resetFocusState()

	first := NewButton().Text("1")
	second := NewButton().Text("2")
	root := NewContainer().Add(first, second)
	mountForFocus(t, root)

	DispatchKey(NewKeyEvent(EventKeyDown, "Tab", 0, false, false, false, false))
	if FocusedWidget() != Widget(first) {
		t.Fatalf("Tab did not focus the first widget")
	}
	DispatchKey(NewKeyEvent(EventKeyDown, "Tab", 0, false, false, false, false))
	if FocusedWidget() != Widget(second) {
		t.Fatalf("second Tab did not focus the second widget")
	}
	DispatchKey(NewKeyEvent(EventKeyDown, "Tab", 0, false, false, true, false))
	if FocusedWidget() != Widget(first) {
		t.Errorf("Shift+Tab did not move backwards, got %T", FocusedWidget())
	}
}

func TestKeyEventBubblesToAncestor(t *testing.T) {
	defer resetFocusState()

	button := NewButton().Text("ok")
	root := NewContainer().Add(button)
	mountForFocus(t, root)

	var received string
	root.AddEventHandler(EventKeyDown, func(e Event) {
		if ke, ok := e.(*KeyEvent); ok {
			received = ke.Key
		}
	})

	RequestFocus(button)
	DispatchKey(NewKeyEvent(EventKeyDown, "a", 0, false, false, false, false))

	if received != "a" {
		t.Errorf("ancestor received key %q, want %q", received, "a")
	}
}

func TestNotifyNativeFocusAndBlur(t *testing.T) {
	defer resetFocusState()

	button := NewButton().Text("ok")
	root := NewContainer().Add(button)
	mountForFocus(t, root)

	var blurred bool
	button.AddEventHandler(EventBlur, func(Event) { blurred = true })

	notifyNativeFocus(button)
	if FocusedWidget() != Widget(button) {
		t.Fatal("notifyNativeFocus did not set the focused widget")
	}
	notifyNativeBlur(button)
	if FocusedWidget() != nil {
		t.Errorf("notifyNativeBlur left focus on %T", FocusedWidget())
	}
	if !blurred {
		t.Error("notifyNativeBlur did not dispatch a blur event")
	}
}

func TestDisabledWidgetIsNotFocusable(t *testing.T) {
	defer resetFocusState()

	button := NewButton().Text("ok")
	button.SetEnabled(false)
	root := NewContainer().Add(button)
	mountForFocus(t, root)

	if CanFocus(button) {
		t.Error("disabled button should not be focusable")
	}
	if len(FocusOrder(root)) != 0 {
		t.Error("disabled button should be excluded from the focus order")
	}
}
