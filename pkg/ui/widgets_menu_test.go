package ui

import (
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

func TestMenuDefaultsAndLabel(t *testing.T) {
	m := NewMenu("Actions")
	if m.GetLabel() != "Actions" {
		t.Errorf("label = %q, want Actions", m.GetLabel())
	}
	if len(m.GetItems()) != 0 {
		t.Errorf("new menu should have no items, got %d", len(m.GetItems()))
	}
	if kindOf(m) != ControlMenu {
		t.Errorf("kindOf = %v, want menu", kindOf(m))
	}
}

func TestMenuItemsReplacesAndCopies(t *testing.T) {
	m := NewMenu("Actions").Items(
		MenuItem{Label: "Open", Enabled: true},
		MenuItem{Separator: true},
		MenuItem{Label: "Save", Enabled: false},
	)

	items := m.GetItems()
	if len(items) != 3 {
		t.Fatalf("items = %d, want 3", len(items))
	}
	if items[0].Label != "Open" || !items[0].Enabled {
		t.Errorf("items[0] = %+v, want {Open enabled}", items[0])
	}
	if !items[1].Separator {
		t.Errorf("items[1] should be a separator: %+v", items[1])
	}
	if items[2].Enabled {
		t.Errorf("items[2] should be disabled: %+v", items[2])
	}

	// Mutating the returned slice must not affect the widget.
	items[0].Label = "mutated"
	if m.GetItems()[0].Label != "Open" {
		t.Error("GetItems did not return a copy")
	}

	// Calling Items again replaces, not appends.
	m.Items(MenuItem{Label: "Only", Enabled: true})
	if got := m.GetItems(); len(got) != 1 || got[0].Label != "Only" {
		t.Errorf("items after replace = %+v, want a single Only entry", got)
	}
}

func TestMenuAddItemAndSeparator(t *testing.T) {
	m := NewMenu("Actions").AddItem("Cut").AddSeparator().AddItem("Copy")

	items := m.GetItems()
	if len(items) != 3 {
		t.Fatalf("items = %d, want 3", len(items))
	}
	if items[0].Label != "Cut" || !items[0].Enabled {
		t.Errorf("AddItem should default to enabled: %+v", items[0])
	}
	if !items[1].Separator {
		t.Errorf("AddSeparator should append a separator: %+v", items[1])
	}
	if items[2].Label != "Copy" || !items[2].Enabled {
		t.Errorf("items[2] = %+v, want {Copy enabled}", items[2])
	}
}

func TestMenuSelectFiresOnSelect(t *testing.T) {
	m := NewMenu("Actions").
		AddItem("First").
		AddItem("Second")

	var gotIndex int
	var gotItem MenuItem
	calls := 0
	m.OnSelect(func(index int, item MenuItem) {
		calls++
		gotIndex, gotItem = index, item
	})

	m.Select(1)
	if calls != 1 || gotIndex != 1 || gotItem.Label != "Second" {
		t.Errorf("onSelect = (%d calls, %d, %+v), want (1, 1, Second)", calls, gotIndex, gotItem)
	}
}

func TestMenuSelectIgnoresSeparatorsAndOutOfRange(t *testing.T) {
	m := NewMenu("Actions").AddItem("Only").AddSeparator()
	calls := 0
	m.OnSelect(func(int, MenuItem) { calls++ })

	m.Select(1) // separator
	if calls != 0 {
		t.Errorf("selecting a separator fired onSelect %d times", calls)
	}
	m.Select(-1)
	m.Select(5)
	if calls != 0 {
		t.Errorf("out-of-range Select fired onSelect %d times", calls)
	}
}

func TestMenuSelectIgnoresDisabledItems(t *testing.T) {
	m := NewMenu("Actions").Items(MenuItem{Label: "Ghost", Enabled: false})
	calls := 0
	m.OnSelect(func(int, MenuItem) { calls++ })

	m.Select(0)
	if calls != 0 {
		t.Errorf("selecting a disabled item fired onSelect %d times", calls)
	}
}

func TestMenuRenderProps(t *testing.T) {
	m := NewMenu("Actions").AddItem("Open").AddSeparator().AddItem("Save")

	props := propsOf(m)
	if props.Text != "Actions" {
		t.Errorf("props.Text = %q, want Actions", props.Text)
	}
	items := props.MenuItems
	if len(items) != 3 || items[0].Label != "Open" || !items[1].Separator || items[2].Label != "Save" {
		t.Errorf("props.MenuItems = %+v, want the three menu entries", items)
	}
}

func TestMenuMeasureReservesArrowSpace(t *testing.T) {
	SetTextMeasurer(fixedMeasurer{perRune: 10, height: 7})
	defer SetTextMeasurer(nil)

	m := NewMenu("abc") // 3 runes * 10 = 30 wide, 7 tall
	size := m.Measure(flex.Constraint{})
	if size.Width != 30+24 {
		t.Errorf("menu width = %g, want %g (label + drop-down arrow)", size.Width, 30.0+24.0)
	}
	if size.Height != 7 {
		t.Errorf("menu height = %g, want 7", size.Height)
	}

	m.Style(NewStyle().WithFontSize(7)) // no effect on explicit size path
	m.style.Width = 100
	m.style.Height = 30
	size = m.Measure(flex.Constraint{})
	if size.Width != 100 || size.Height != 30 {
		t.Errorf("explicit size = %gx%g, want 100x30", size.Width, size.Height)
	}
}

func TestMenuMountedUpdatePushesItems(t *testing.T) {
	m := NewMenu("Actions").AddItem("Open")
	backend := &fakeBackend{}
	if err := layoutAndMount(backend, &fakeControl{kind: ControlContainer}, m, 200, 100); err != nil {
		t.Fatalf("layoutAndMount: %v", err)
	}

	ctrl, ok := m.NativeControl().(*fakeControl)
	if !ok {
		t.Fatal("unexpected native control type")
	}
	if len(ctrl.props.MenuItems) != 1 {
		t.Fatalf("mounted props items = %d, want 1", len(ctrl.props.MenuItems))
	}

	// A later AddItem must reach the already-mounted native control.
	m.AddItem("Save")
	if len(ctrl.props.MenuItems) != 2 || ctrl.props.MenuItems[1].Label != "Save" {
		t.Errorf("mounted props after AddItem = %+v, want two entries ending in Save", ctrl.props.MenuItems)
	}
	if ctrl.props.Text != "Actions" {
		t.Errorf("mounted props Text = %q, want Actions", ctrl.props.Text)
	}
}

func TestMenuLayoutRecordsRect(t *testing.T) {
	m := NewMenu("Actions")
	m.Layout(10, 20, 120, 40)

	rect := m.GetRect()
	if rect.Position.X != 10 || rect.Position.Y != 20 || rect.Size.Width != 120 || rect.Size.Height != 40 {
		t.Errorf("rect = %+v, want 10,20 120x40", rect)
	}
	if m.Render() != nil {
		t.Error("Render on an unmounted menu should be a successful no-op")
	}
}
