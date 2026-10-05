package ui

import (
	"sync"

	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// Default window size used when App is created without an explicit size.
const (
	DefaultAppWidth  float32 = 800
	DefaultAppHeight float32 = 600
)

// App boots the widget tree on the current platform backend. It owns the
// native window and drives the toolkit event loop.
type App struct {
	mu sync.Mutex

	title         string
	width         float32
	height        float32
	content       Widget
	backend       Backend
	window        NativeWindow
	quitOnce      sync.Once
	quitRequested bool
}

// AppOption customizes an App at construction time.
type AppOption func(*App)

// WithSize sets the initial window client size in pixels.
func WithSize(width, height float32) AppOption {
	return func(a *App) {
		if width > 0 {
			a.width = width
		}
		if height > 0 {
			a.height = height
		}
	}
}

// WithBackend injects an explicit backend, bypassing platform auto-detection.
func WithBackend(b Backend) AppOption {
	return func(a *App) {
		a.backend = b
	}
}

// NewApp creates a new application with the given window title.
func NewApp(title string, opts ...AppOption) *App {
	a := &App{
		title:  title,
		width:  DefaultAppWidth,
		height: DefaultAppHeight,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(a)
		}
	}
	return a
}

// SetContent sets the root widget of the window.
func (a *App) SetContent(content Widget) *App {
	a.mu.Lock()
	a.content = content
	window := a.window
	a.mu.Unlock()

	if window != nil && content != nil {
		_ = window.SetContent(content)
	}
	return a
}

// Content returns the current root widget.
func (a *App) Content() Widget {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.content
}

// Backend returns the backend in use, or nil before Run.
func (a *App) Backend() Backend {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.backend
}

// Run initializes the backend, mounts the content and blocks until the window
// is closed or Quit is called.
func (a *App) Run() error {
	a.mu.Lock()
	backend := a.backend
	title := a.title
	width := a.width
	height := a.height
	content := a.content
	a.mu.Unlock()

	if backend == nil {
		backend = CurrentBackend()
	}
	if backend == nil {
		return ErrBackendUnavailable
	}
	if err := backend.Init(); err != nil {
		return err
	}

	a.mu.Lock()
	a.backend = backend
	a.mu.Unlock()

	// Route reactive updates to the toolkit UI thread.
	reactive.SetDispatcher(&backendDispatcher{backend: backend})

	window, err := backend.CreateWindow(title, width, height)
	if err != nil {
		return err
	}

	a.mu.Lock()
	a.window = window
	a.mu.Unlock()

	if content != nil {
		if err := window.SetContent(content); err != nil {
			return err
		}
	}

	window.Show()
	return backend.Run()
}

// Quit stops the event loop and closes the window.
func (a *App) Quit() {
	a.quitOnce.Do(func() {
		a.mu.Lock()
		a.quitRequested = true
		window := a.window
		backend := a.backend
		a.mu.Unlock()

		if window != nil {
			window.Close()
		}
		if backend != nil {
			backend.Quit()
		}
	})
}

// backendDispatcher adapts a Backend to reactive.Dispatcher so that reactive
// updates are marshaled onto the toolkit's UI thread.
type backendDispatcher struct {
	backend Backend
}

func (d *backendDispatcher) RunOnUI(fn func()) error {
	if fn == nil {
		return nil
	}
	return d.backend.Post(fn)
}

func (d *backendDispatcher) IsUIThread() bool {
	return d.backend.OnUIThread()
}
