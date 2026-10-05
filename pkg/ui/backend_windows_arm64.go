//go:build windows && !amd64

package ui

// newPlatformBackend returns a backend that reports the platform as
// unsupported. The WinUI 3 bindings used by the amd64 backend are generated for
// windows/amd64 only, so other Windows architectures get a clear error instead
// of a build failure.
func newPlatformBackend() Backend {
	return &unsupportedWindowsBackend{}
}

type unsupportedWindowsBackend struct{}

func (b *unsupportedWindowsBackend) Name() string { return "unsupported (windows non-amd64)" }

func (b *unsupportedWindowsBackend) Init() error { return ErrBackendUnavailable }

func (b *unsupportedWindowsBackend) CreateWindow(string, float32, float32) (NativeWindow, error) {
	return nil, ErrBackendUnavailable
}

func (b *unsupportedWindowsBackend) CreateControl(Widget, ControlKind, ControlProps) (NativeControl, error) {
	return nil, ErrBackendUnavailable
}

func (b *unsupportedWindowsBackend) Run() error { return ErrBackendUnavailable }

func (b *unsupportedWindowsBackend) Quit() {}

func (b *unsupportedWindowsBackend) Post(fn func()) error {
	if fn != nil {
		fn()
	}
	return nil
}

func (b *unsupportedWindowsBackend) OnUIThread() bool { return true }
