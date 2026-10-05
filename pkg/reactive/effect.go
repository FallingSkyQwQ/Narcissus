package reactive

import (
	"sync"
)

// Effect represents a reactive side effect.
//
// The effect function runs immediately and then re-runs automatically whenever
// any signal or computed it read changes. Runs are batched: several updates in
// the same Batch (or a burst outside one) collapse into a single re-run.
type Effect struct {
	dep *dependent

	mu       sync.Mutex
	fn       func()
	cleanup  func()
	disposed bool
}

// NewEffect creates a new effect and runs it immediately.
func NewEffect(fn func()) *Effect {
	e := &Effect{dep: newDependent(), fn: fn}
	e.dep.flush = e.run
	e.run()
	return e
}

func (e *Effect) run() {
	e.mu.Lock()
	if e.disposed {
		e.mu.Unlock()
		return
	}
	cleanup := e.cleanup
	e.cleanup = nil
	fn := e.fn
	e.mu.Unlock()

	// A re-run first tears down the previous run's cleanup, mirroring how
	// effects manage subscriptions.
	if cleanup != nil {
		cleanup()
	}

	e.dep.clearInvalid()
	track(e.dep, fn)
}

// Dispose stops the effect, detaches it from its dependencies and calls the
// cleanup function if set. Calling it more than once is safe.
func (e *Effect) Dispose() {
	e.mu.Lock()
	if e.disposed {
		e.mu.Unlock()
		return
	}
	e.disposed = true
	cleanup := e.cleanup
	e.cleanup = nil
	e.mu.Unlock()

	e.dep.dispose()

	if cleanup != nil {
		cleanup()
	}
}

// SetCleanup sets a cleanup function to be called before the next run and on
// Dispose. If the effect is already disposed, the cleanup is called immediately.
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
