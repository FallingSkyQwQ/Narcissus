//go:build !linux && !windows

package ui

// newPlatformBackend returns a backend that reports the platform as
// unsupported. Add a new backend_<goos>.go file to support another toolkit.
func newPlatformBackend() Backend {
	return &unsupportedBackend{}
}

type unsupportedBackend struct{}

func (b *unsupportedBackend) Name() string { return "unsupported" }

func (b *unsupportedBackend) Init() error { return ErrBackendUnavailable }

func (b *unsupportedBackend) CreateWindow(string, float32, float32) (NativeWindow, error) {
	return nil, ErrBackendUnavailable
}

func (b *unsupportedBackend) CreateControl(Widget, ControlKind, ControlProps) (NativeControl, error) {
	return nil, ErrBackendUnavailable
}

func (b *unsupportedBackend) Run() error { return ErrBackendUnavailable }

func (b *unsupportedBackend) Quit() {}

func (b *unsupportedBackend) Post(fn func()) error {
	if fn != nil {
		fn()
	}
	return nil
}

func (b *unsupportedBackend) OnUIThread() bool { return true }
