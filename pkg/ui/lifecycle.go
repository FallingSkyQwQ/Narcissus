package ui

import "sync"

var (
	readyMu        sync.Mutex
	readyCallbacks []func()
	readyFired     bool
)

// OnReady registers a function to run once the first window has been laid out
// and the toolkit is initialized. This is the earliest point at which
// toolkit-only information (such as the OS dark-mode preference) can be read.
//
// If the app is already ready, fn runs immediately. Callbacks run on the UI
// thread.
func OnReady(fn func()) {
	if fn == nil {
		return
	}
	readyMu.Lock()
	if readyFired {
		readyMu.Unlock()
		fn()
		return
	}
	readyCallbacks = append(readyCallbacks, fn)
	readyMu.Unlock()
}

// fireReady runs the pending ready callbacks exactly once; later calls are
// no-ops.
func fireReady() {
	readyMu.Lock()
	if readyFired {
		readyMu.Unlock()
		return
	}
	readyFired = true
	callbacks := readyCallbacks
	readyCallbacks = nil
	readyMu.Unlock()

	for _, fn := range callbacks {
		fn()
	}
}
