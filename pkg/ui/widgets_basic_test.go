package ui

import (
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

func TestProgressBarClampsValue(t *testing.T) {
	p := NewProgressBar().Range(0, 100).Value(150)
	if p.GetValue() != 100 {
		t.Errorf("value = %g, want clamped 100", p.GetValue())
	}
	p.Value(-10)
	if p.GetValue() != 0 {
		t.Errorf("value = %g, want clamped 0", p.GetValue())
	}
}

func TestProgressBarPercent(t *testing.T) {
	p := NewProgressBar().Range(0, 200).Value(50)
	if p.Percent() != 0.25 {
		t.Errorf("percent = %g, want 0.25", p.Percent())
	}
}

func TestProgressBarIndeterminateAndRender(t *testing.T) {
	p := NewProgressBar().Indeterminate(true).ShowText(true)
	if !p.IsIndeterminate() || !p.IsShowingText() {
		t.Fatal("indeterminate/showText flags not set")
	}
	if kindOf(p) != ControlProgress {
		t.Errorf("kindOf = %v, want progress", kindOf(p))
	}
	props := propsOf(p)
	if !props.Indeterminate || !props.ShowText {
		t.Errorf("props missing flags: %+v", props)
	}
}

func TestSwitchToggleAndCallback(t *testing.T) {
	s := NewSwitch()
	var seen []bool
	s.OnToggle(func(on bool) { seen = append(seen, on) })

	s.Toggle()
	if !s.IsChecked() {
		t.Fatal("switch should be on after toggle")
	}
	s.Toggle()
	if s.IsChecked() {
		t.Fatal("switch should be off after second toggle")
	}
	if len(seen) != 2 || seen[0] != true || seen[1] != false {
		t.Errorf("toggle callbacks = %v, want [true false]", seen)
	}
}

func TestSwitchClickToggles(t *testing.T) {
	s := NewSwitch()
	if !s.HandleEvent(newTestClickEvent(s)) {
		t.Fatal("click should be handled")
	}
	if !s.IsChecked() {
		t.Error("click did not toggle the switch")
	}
}

func TestRadioButtonGroupIsMutuallyExclusive(t *testing.T) {
	a := NewRadioButton().Group("test-group-mutex").Label("A")
	b := NewRadioButton().Group("test-group-mutex").Label("B")
	c := NewRadioButton().Group("test-group-mutex").Label("C")

	a.Checked(true)
	if !a.IsChecked() {
		t.Fatal("A should be checked")
	}

	c.Select()
	if !c.IsChecked() || a.IsChecked() || b.IsChecked() {
		t.Errorf("group state after selecting C: A=%v B=%v C=%v, want false/false/true",
			a.IsChecked(), b.IsChecked(), c.IsChecked())
	}
}

func TestRadioButtonOnSelectFiresOnce(t *testing.T) {
	r := NewRadioButton().Group("test-group-select")
	count := 0
	r.OnSelect(func() { count++ })
	r.Checked(true)
	r.Checked(true)
	if count != 1 {
		t.Errorf("onSelect fired %d times, want 1", count)
	}
}

func TestRadioButtonRenderProps(t *testing.T) {
	r := NewRadioButton().Group("g").Label("Option").Checked(true)
	if kindOf(r) != ControlRadio {
		t.Errorf("kindOf = %v, want radio", kindOf(r))
	}
	props := propsOf(r)
	if props.Group != "g" || props.Text != "Option" || !props.Checked {
		t.Errorf("unexpected radio props: %+v", props)
	}
}

func TestSwitchIsFocusable(t *testing.T) {
	s := NewSwitch()
	if !CanFocus(s) {
		t.Error("switch should be focusable")
	}
	if CanFocus(NewProgressBar()) {
		t.Error("progress bar should not be focusable")
	}
}

func TestProgressBarMinMaxClampValue(t *testing.T) {
	p := NewProgressBar().Range(0, 100).Value(50)

	p.Max(80) // value stays inside the range
	if p.GetValue() != 50 {
		t.Errorf("value after raising max = %g, want 50", p.GetValue())
	}

	p.Min(60) // value is below the new minimum
	if p.GetValue() != 60 {
		t.Errorf("value after raising min = %g, want clamped 60", p.GetValue())
	}

	q := NewProgressBar().Range(0, 100).Value(50)
	q.Max(40) // value is above the new maximum
	if q.GetValue() != 40 {
		t.Errorf("value after lowering max = %g, want clamped 40", q.GetValue())
	}
}

func TestProgressBarRangeReclamps(t *testing.T) {
	p := NewProgressBar().Range(0, 100).Value(90)
	p.Range(0, 50)
	if p.GetValue() != 50 {
		t.Errorf("value after narrowing range = %g, want clamped 50", p.GetValue())
	}
}

func TestProgressBarDegenerateRangePercent(t *testing.T) {
	p := NewProgressBar().Range(10, 10).Value(10)
	if p.Percent() != 0 {
		t.Errorf("percent for empty range = %g, want 0", p.Percent())
	}
}

func TestProgressBarMeasureDefaults(t *testing.T) {
	size := NewProgressBar().Measure(flex.Constraint{})
	if size.Width != 200 || size.Height != 8 {
		t.Errorf("default size = %gx%g, want 200x8", size.Width, size.Height)
	}
}

func TestProgressBarRenderPropsRange(t *testing.T) {
	p := NewProgressBar().Range(10, 50).Value(30)
	props := propsOf(p)
	if props.Min != 10 || props.Max != 50 || props.Value != 30 {
		t.Errorf("unexpected progress props: %+v", props)
	}
}

func TestSwitchCheckedSameValueIsSilent(t *testing.T) {
	s := NewSwitch()
	calls := 0
	s.OnToggle(func(bool) { calls++ })

	s.Checked(false) // already off
	if calls != 0 {
		t.Errorf("redundant Checked(false) fired %d callbacks", calls)
	}

	s.Checked(true)
	s.Checked(true) // already on
	if calls != 1 {
		t.Errorf("redundant Checked(true) fired %d extra callbacks, want 0", calls-1)
	}
}

func TestSwitchLabelAndProps(t *testing.T) {
	s := NewSwitch().Label("Enable")
	if s.GetLabel() != "Enable" {
		t.Errorf("label = %q, want Enable", s.GetLabel())
	}
	props := propsOf(s)
	if props.Text != "Enable" || props.Checked {
		t.Errorf("unexpected switch props: %+v", props)
	}
}

func TestSwitchMeasureDefaults(t *testing.T) {
	SetTextMeasurer(fixedMeasurer{perRune: 10, height: 7})
	defer SetTextMeasurer(nil)

	size := NewSwitch().Measure(flex.Constraint{})
	if size.Width != 44 || size.Height != 24 {
		t.Errorf("default size = %gx%g, want 44x24", size.Width, size.Height)
	}

	labelled := NewSwitch().Label("abc") // 3 runes * 10 = 30 wide
	size = labelled.Measure(flex.Constraint{})
	if size.Width != 44+8+30 {
		t.Errorf("labelled width = %g, want %g", size.Width, float32(44+8+30))
	}
}

func TestSwitchDisabledClickIgnored(t *testing.T) {
	s := NewSwitch()
	s.SetEnabled(false)
	if s.HandleEvent(newTestClickEvent(s)) {
		t.Error("disabled switch should not handle clicks")
	}
	if s.IsChecked() {
		t.Error("disabled switch should not toggle")
	}
}

func TestRadioButtonGroupReassignment(t *testing.T) {
	a := NewRadioButton().Group("rg-move").Checked(true)
	b := NewRadioButton().Group("rg-move")

	// a leaves the group, so selecting b must no longer clear it.
	a.Group("rg-moved")
	b.Select()

	if !a.IsChecked() {
		t.Error("a should stay checked after leaving the group")
	}
	if !b.IsChecked() {
		t.Error("b should be checked")
	}
}

func TestRadioButtonGroupMapTracksMembers(t *testing.T) {
	a := NewRadioButton().Group("rg-map")
	NewRadioButton().Group("rg-map")

	radioGroupsMu.Lock()
	members := len(radioGroups["rg-map"])
	radioGroupsMu.Unlock()
	if members != 2 {
		t.Fatalf("group members = %d, want 2", members)
	}

	a.Group("") // leaves the group entirely
	radioGroupsMu.Lock()
	members = len(radioGroups["rg-map"])
	radioGroupsMu.Unlock()
	if members != 1 {
		t.Errorf("group members after leaving = %d, want 1", members)
	}
}

func TestRadioButtonCheckedFalseClears(t *testing.T) {
	r := NewRadioButton().Group("rg-clear").Checked(true)
	r.Checked(false)
	if r.IsChecked() {
		t.Error("Checked(false) should clear the selection")
	}
}

func TestRadioButtonClickSelects(t *testing.T) {
	r := NewRadioButton().Group("rg-click")
	if !r.HandleEvent(newTestClickEvent(r)) {
		t.Fatal("click should be handled")
	}
	if !r.IsChecked() {
		t.Error("click did not select the radio button")
	}
}

func newTestClickEvent(target Widget) *BaseEvent {
	return &BaseEvent{Type: EventClick, Target: target}
}
