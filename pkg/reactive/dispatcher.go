package reactive

import (
	"sync"
)

// Dispatcher handles UI thread marshaling
type Dispatcher interface {
	// RunOnUI schedules a function to run on the UI thread
	RunOnUI(fn func()) error
	// IsUIThread returns true if called from the UI thread
	IsUIThread() bool
}

var (
	globalDispatcher   Dispatcher
	globalDispatcherMu sync.RWMutex
)

// SetDispatcher sets the global dispatcher
func SetDispatcher(d Dispatcher) {
	globalDispatcherMu.Lock()
	defer globalDispatcherMu.Unlock()
	globalDispatcher = d
}

// RunOnUI schedules a function to run on the UI thread.
// If no dispatcher is set or already on UI thread, the function runs immediately.
// Returns an error if the dispatcher cannot schedule the function.
func RunOnUI(fn func()) error {
	globalDispatcherMu.RLock()
	d := globalDispatcher
	globalDispatcherMu.RUnlock()

	if d == nil || d.IsUIThread() {
		fn()
		return nil
	}
	return d.RunOnUI(fn)
}

// IsUIThread returns true if called from the UI thread.
// Returns true if no dispatcher is set.
func IsUIThread() bool {
	globalDispatcherMu.RLock()
	d := globalDispatcher
	globalDispatcherMu.RUnlock()

	if d == nil {
		return true
	}
	return d.IsUIThread()
}
