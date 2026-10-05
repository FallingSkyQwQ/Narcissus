package ui

import "testing"

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

func newTestClickEvent(target Widget) *BaseEvent {
	return &BaseEvent{Type: EventClick, Target: target}
}
