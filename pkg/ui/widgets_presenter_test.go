package ui

import (
	"errors"
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// presenterBackend is an in-memory Backend that also implements the
// optional Presenter capability. It records every dialog and toast it
// is asked to show so tests can assert on the specs.
type presenterBackend struct {
	fakeBackend

	dialogs []DialogSpec
	toasts  []ToastSpec
	err     error
}

func (b *presenterBackend) PresentDialog(spec DialogSpec) error {
	b.dialogs = append(b.dialogs, spec)
	return b.err
}

func (b *presenterBackend) PresentToast(spec ToastSpec) error {
	b.toasts = append(b.toasts, spec)
	return b.err
}

// useBackend installs b as the active backend for the duration of the
// test and restores the previous one afterwards.
func useBackend(t *testing.T, b Backend) {
	t.Helper()
	previous := CurrentBackend()
	SetBackend(b)
	t.Cleanup(func() { SetBackend(previous) })
}

func TestDialogDefaults(t *testing.T) {
	d := NewDialog("Confirm", "Delete file?")
	if d.Title() != "Confirm" {
		t.Errorf("title = %q, want Confirm", d.Title())
	}
	if d.Message() != "Delete file?" {
		t.Errorf("message = %q, want Delete file?", d.Message())
	}
	if got := d.GetButtons(); len(got) != 1 || got[0] != "OK" {
		t.Errorf("default buttons = %v, want [OK]", got)
	}
}

func TestDialogButtonsReplacesAndCopies(t *testing.T) {
	d := NewDialog("t", "m").Buttons("Yes", "No", "Cancel")

	buttons := d.GetButtons()
	if len(buttons) != 3 || buttons[0] != "Yes" || buttons[2] != "Cancel" {
		t.Fatalf("buttons = %v, want [Yes No Cancel]", buttons)
	}

	// Mutating the returned slice must not affect the widget.
	buttons[0] = "mutated"
	if d.GetButtons()[0] != "Yes" {
		t.Error("GetButtons did not return a copy")
	}

	// Calling Buttons again replaces, not appends.
	d.Buttons("OK")
	if got := d.GetButtons(); len(got) != 1 || got[0] != "OK" {
		t.Errorf("buttons after replace = %v, want [OK]", got)
	}
}

func TestDialogShowWithoutPresenterIsNoOp(t *testing.T) {
	// fakeBackend does not implement Presenter, so Show must succeed
	// without presenting anything.
	useBackend(t, &fakeBackend{})

	d := NewDialog("t", "m")
	if err := d.Show(); err != nil {
		t.Errorf("Show without presenter = %v, want nil", err)
	}
}

func TestDialogShowPresentsSpec(t *testing.T) {
	b := &presenterBackend{}
	useBackend(t, b)

	d := NewDialog("Confirm", "Delete?").Buttons("Delete", "Cancel")
	if err := d.Show(); err != nil {
		t.Fatalf("Show: %v", err)
	}
	if len(b.dialogs) != 1 {
		t.Fatalf("presented dialogs = %d, want 1", len(b.dialogs))
	}

	spec := b.dialogs[0]
	if spec.Title != "Confirm" {
		t.Errorf("spec.Title = %q, want Confirm", spec.Title)
	}
	if spec.Message != "Delete?" {
		t.Errorf("spec.Message = %q, want Delete?", spec.Message)
	}
	if len(spec.Buttons) != 2 || spec.Buttons[0] != "Delete" || spec.Buttons[1] != "Cancel" {
		t.Errorf("spec.Buttons = %v, want [Delete Cancel]", spec.Buttons)
	}
}

func TestDialogOnResultReceivesButtonIndex(t *testing.T) {
	b := &presenterBackend{}
	useBackend(t, b)

	d := NewDialog("t", "m").Buttons("One", "Two", "Three")
	got := -1
	d.OnResult(func(index int) { got = index })

	if err := d.Show(); err != nil {
		t.Fatalf("Show: %v", err)
	}
	spec := b.dialogs[0]
	if spec.OnResult == nil {
		t.Fatal("OnResult was not propagated to the spec")
	}

	// The backend reports the chosen button by index.
	spec.OnResult(2)
	if got != 2 {
		t.Errorf("onResult index = %d, want 2", got)
	}
}

func TestDialogShowReturnsPresenterError(t *testing.T) {
	wantErr := errors.New("present failed")
	b := &presenterBackend{err: wantErr}
	useBackend(t, b)

	if err := NewDialog("t", "m").Show(); !errors.Is(err, wantErr) {
		t.Errorf("Show error = %v, want %v", err, wantErr)
	}
}

func TestDialogStaysOutOfLayout(t *testing.T) {
	d := NewDialog("t", "m")

	size := d.Measure(flex.Constraint{MaxWidth: 400, MaxHeight: 300})
	if size.Width != 0 || size.Height != 0 {
		t.Errorf("measured size = %gx%g, want 0x0", size.Width, size.Height)
	}

	d.Layout(10, 20, 400, 300)
	rect := d.GetRect()
	if rect.Position.X != 10 || rect.Position.Y != 20 {
		t.Errorf("position = %v, want (10, 20)", rect.Position)
	}
	if rect.Size.Width != 0 || rect.Size.Height != 0 {
		t.Errorf("rect size = %gx%g, want 0x0", rect.Size.Width, rect.Size.Height)
	}

	if d.Render() != nil {
		t.Error("Render should be a successful no-op")
	}
}

func TestToastDefaults(t *testing.T) {
	toast := NewToast("Saved")
	if toast.Message() != "Saved" {
		t.Errorf("message = %q, want Saved", toast.Message())
	}
	if toast.GetSeverity() != ToastInfo {
		t.Errorf("severity = %v, want ToastInfo", toast.GetSeverity())
	}
	if toast.GetDurationMS() != DefaultToastDurationMS {
		t.Errorf("duration = %d, want %d", toast.GetDurationMS(), DefaultToastDurationMS)
	}
}

func TestToastSeverityAndDuration(t *testing.T) {
	toast := NewToast("x").Severity(ToastError).DurationMS(1500)
	if toast.GetSeverity() != ToastError {
		t.Errorf("severity = %v, want ToastError", toast.GetSeverity())
	}
	if toast.GetDurationMS() != 1500 {
		t.Errorf("duration = %d, want 1500", toast.GetDurationMS())
	}
}

func TestToastShowWithoutPresenterIsNoOp(t *testing.T) {
	useBackend(t, &fakeBackend{})

	if err := NewToast("x").Show(); err != nil {
		t.Errorf("Show without presenter = %v, want nil", err)
	}
}

func TestToastShowPresentsSpec(t *testing.T) {
	b := &presenterBackend{}
	useBackend(t, b)

	toast := NewToast("Done").Severity(ToastSuccess).DurationMS(1200)
	if err := toast.Show(); err != nil {
		t.Fatalf("Show: %v", err)
	}
	if len(b.toasts) != 1 {
		t.Fatalf("presented toasts = %d, want 1", len(b.toasts))
	}

	spec := b.toasts[0]
	if spec.Message != "Done" {
		t.Errorf("spec.Message = %q, want Done", spec.Message)
	}
	if spec.Severity != ToastSuccess {
		t.Errorf("spec.Severity = %v, want ToastSuccess", spec.Severity)
	}
	if spec.DurationMS != 1200 {
		t.Errorf("spec.DurationMS = %d, want 1200", spec.DurationMS)
	}
}

func TestToastShowReturnsPresenterError(t *testing.T) {
	wantErr := errors.New("toast failed")
	b := &presenterBackend{err: wantErr}
	useBackend(t, b)

	if err := NewToast("x").Show(); !errors.Is(err, wantErr) {
		t.Errorf("Show error = %v, want %v", err, wantErr)
	}
}

func TestToastStaysOutOfLayout(t *testing.T) {
	toast := NewToast("x")

	size := toast.Measure(flex.Constraint{MaxWidth: 400, MaxHeight: 300})
	if size.Width != 0 || size.Height != 0 {
		t.Errorf("measured size = %gx%g, want 0x0", size.Width, size.Height)
	}

	toast.Layout(5, 6, 400, 300)
	rect := toast.GetRect()
	if rect.Position.X != 5 || rect.Position.Y != 6 {
		t.Errorf("position = %v, want (5, 6)", rect.Position)
	}
	if rect.Size.Width != 0 || rect.Size.Height != 0 {
		t.Errorf("rect size = %gx%g, want 0x0", rect.Size.Width, rect.Size.Height)
	}

	if toast.Render() != nil {
		t.Error("Render should be a successful no-op")
	}
}
