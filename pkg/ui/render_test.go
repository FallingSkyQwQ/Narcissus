package ui

import (
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// fakeControl records the operations the renderer performs on a native control.
type fakeControl struct {
	kind      ControlKind
	props     ControlProps
	x, y      float32
	w, h      float32
	attached  bool
	destroyed bool
}

func (c *fakeControl) SetBounds(x, y, width, height float32) {
	c.x, c.y, c.w, c.h = x, y, width, height
}
func (c *fakeControl) SetVisible(bool)               {}
func (c *fakeControl) SetEnabled(bool)               {}
func (c *fakeControl) SetStyle(*Style)               {}
func (c *fakeControl) Update(props ControlProps)     { c.props = props }
func (c *fakeControl) AttachTo(parent NativeControl) { c.attached = parent != nil }
func (c *fakeControl) Destroy()                      { c.destroyed = true }

// fakeBackend is an in-memory Backend used to exercise the renderer without a
// display or toolkit.
type fakeBackend struct {
	created []*fakeControl
}

func (b *fakeBackend) Name() string { return "fake" }
func (b *fakeBackend) Init() error  { return nil }

func (b *fakeBackend) CreateWindow(string, float32, float32) (NativeWindow, error) {
	return nil, nil
}

func (b *fakeBackend) CreateControl(w Widget, kind ControlKind, props ControlProps) (NativeControl, error) {
	c := &fakeControl{kind: kind, props: props}
	b.created = append(b.created, c)
	return c, nil
}

func (b *fakeBackend) Run() error           { return nil }
func (b *fakeBackend) Quit()                {}
func (b *fakeBackend) Post(fn func()) error { fn(); return nil }
func (b *fakeBackend) OnUIThread() bool     { return true }

func TestLayoutAndMountCreatesAndPositionsControls(t *testing.T) {
	btn := NewButton().Text("Increment")
	txt := NewText("Count: 0")
	root := NewContainer().Direction(flex.DirectionColumn).Add(txt, btn)

	backend := &fakeBackend{}
	surface := &fakeControl{kind: ControlContainer}

	if err := layoutAndMount(backend, surface, root, 400, 300); err != nil {
		t.Fatalf("layoutAndMount: %v", err)
	}

	// One control per widget: container + text + button.
	if len(backend.created) != 3 {
		t.Fatalf("expected 3 controls, got %d", len(backend.created))
	}

	if btn.NativeControl() == nil {
		t.Fatal("button was not mounted")
	}
	if txt.NativeControl() == nil {
		t.Fatal("text was not mounted")
	}

	var buttonProps ControlProps
	found := false
	for _, c := range backend.created {
		if c.kind == ControlButton {
			buttonProps = c.props
			found = true
			if !c.attached {
				t.Error("button control was not attached to its parent surface")
			}
			if c.w <= 0 || c.h <= 0 {
				t.Errorf("button has non-positive size: %gx%g", c.w, c.h)
			}
		}
	}
	if !found {
		t.Fatal("no button control created")
	}
	if buttonProps.Text != "Increment" {
		t.Errorf("button props Text = %q, want %q", buttonProps.Text, "Increment")
	}
}

func TestRemountReusesExistingControls(t *testing.T) {
	root := NewContainer().Direction(flex.DirectionColumn).Add(NewText("hello"))
	backend := &fakeBackend{}

	if err := layoutAndMount(backend, &fakeControl{kind: ControlContainer}, root, 200, 100); err != nil {
		t.Fatalf("first mount: %v", err)
	}
	first := len(backend.created)
	if first != 2 {
		t.Fatalf("expected 2 controls after first mount, got %d", first)
	}

	if err := layoutAndMount(backend, &fakeControl{kind: ControlContainer}, root, 300, 200); err != nil {
		t.Fatalf("second mount: %v", err)
	}
	if len(backend.created) != first {
		t.Fatalf("relayout recreated controls: %d -> %d", first, len(backend.created))
	}
}

func TestSyncWidgetPushesStateToNativeControl(t *testing.T) {
	txt := NewText("before")
	backend := &fakeBackend{}
	if err := layoutAndMount(backend, &fakeControl{kind: ControlContainer}, txt, 200, 100); err != nil {
		t.Fatalf("mount: %v", err)
	}

	txt.Text("after")

	ctrl, ok := txt.NativeControl().(*fakeControl)
	if !ok {
		t.Fatal("unexpected native control type")
	}
	if ctrl.props.Text != "after" {
		t.Errorf("native props not synced: got %q, want %q", ctrl.props.Text, "after")
	}
}
