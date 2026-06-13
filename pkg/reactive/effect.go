package reactive

import (
	"sync"
)

// Effect represents a reactive computation.
// This is a simplified implementation that runs once on creation.
type Effect struct {
	mu       sync.Mutex
	fn       func()
	disposed bool
	cleanup  func()
}

// NewEffect creates a new effect and runs it immediately.
// Note: This simplified implementation does not auto-track dependencies.
func NewEffect(fn func()) *Effect {
	e := &Effect{
		fn: fn,
	}
	e.run()
	return e
}

func (e *Effect) run() {
	e.mu.Lock()
	if e.disposed {
		e.mu.Unlock()
		return
	}
	fn := e.fn
	e.mu.Unlock()

	// Execute the effect function
	fn()
}

// Dispose stops the effect and calls the cleanup function if set
func (e *Effect) Dispose() {
	e.mu.Lock()
	e.disposed = true
	cleanup := e.cleanup
	e.cleanup = nil
	e.mu.Unlock()

	// Call cleanup if exists
	if cleanup != nil {
		cleanup()
	}
}

// SetCleanup sets a cleanup function to be called on Dispose.
// If the effect is already disposed, the cleanup is called immediately.
func (e *Effect) SetCleanup(cleanup func()) {
	e.mu.Lock()
	if e.disposed {
		// Already disposed, don't store the cleanup but call it immediately
		e.mu.Unlock()
		if cleanup != nil {
			cleanup()
		}
		return
	}
	e.cleanup = cleanup
	e.mu.Unlock()
}
