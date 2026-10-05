package ui

// MenuItem is one entry in a Menu. A separator entry has Separator set and its
// other fields ignored.
type MenuItem struct {
	// Label is the entry text.
	Label string
	// Enabled reports whether the entry can be chosen.
	Enabled bool
	// Separator draws a divider instead of a labelled entry.
	Separator bool
}

// ToastSeverity selects the visual treatment of a Toast.
type ToastSeverity int

const (
	// ToastInfo is a neutral notification (the default).
	ToastInfo ToastSeverity = iota
	// ToastSuccess reports a completed action.
	ToastSuccess
	// ToastWarning reports a recoverable problem.
	ToastWarning
	// ToastError reports a failure.
	ToastError
)

// DialogSpec describes a modal dialog to present. It is produced by the Dialog
// widget and consumed by a Backend's Presenter.
type DialogSpec struct {
	// Title is the dialog heading.
	Title string
	// Message is the dialog body text.
	Message string
	// Buttons are the action labels, in order. The index chosen is reported
	// through OnResult.
	Buttons []string
	// OnResult is called with the chosen button index. It may be nil.
	OnResult func(index int)
}

// ToastSpec describes a transient notification to present.
type ToastSpec struct {
	// Message is the notification text.
	Message string
	// Severity selects the visual treatment.
	Severity ToastSeverity
	// DurationMS is how long the toast stays visible, in milliseconds. A
	// non-positive value lets the backend pick a default.
	DurationMS int
}

// Presenter is an optional capability a Backend may implement to show
// surfaces that are not part of the laid-out widget tree: modal dialogs and
// transient toasts. Menus do not need it because they are ordinary controls
// (see ControlMenu).
//
// When the active backend does not implement Presenter, the corresponding
// widget calls become no-ops rather than failing.
type Presenter interface {
	// PresentDialog shows a modal dialog.
	PresentDialog(spec DialogSpec) error
	// PresentToast shows a transient notification.
	PresentToast(spec ToastSpec) error
}

// currentPresenter returns the active backend's Presenter, or nil when the
// backend does not implement the interface.
func currentPresenter() Presenter {
	backend := CurrentBackend()
	if backend == nil {
		return nil
	}
	presenter, _ := backend.(Presenter)
	return presenter
}
