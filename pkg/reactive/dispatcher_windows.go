//go:build windows

package reactive

import (
	"errors"
	"syscall"

	"github.com/FallingSkyQwQ/Narcissus/pkg/bridge/rt"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
)

// WindowsDispatcher implements Dispatcher using DispatcherQueue
type WindowsDispatcher struct {
	dispatcherQueue *rt.DispatcherQueue
	// hasDispatcherQueue indicates whether DispatcherQueue is available
	// If false, RunOnUI will return an error when not guaranteed to be on UI thread
	hasDispatcherQueue bool
	// uiThreadID stores the UI thread ID for thread checking
	uiThreadID uint32
}

// getCurrentThreadId returns the current thread ID
func getCurrentThreadId() uint32 {
	ret, _, _ := procGetCurrentThreadId.Call()
	return uint32(ret)
}

// NewWindowsDispatcher creates a new Windows dispatcher.
// It tries to get the DispatcherQueue for the current thread.
// If DispatcherQueue is not available (e.g., WinRT not initialized), it still creates
// a functional dispatcher but RunOnUI will be limited.
func NewWindowsDispatcher() (*WindowsDispatcher, error) {
	d := &WindowsDispatcher{
		hasDispatcherQueue: false,
		uiThreadID:         getCurrentThreadId(),
	}

	// Try to get the DispatcherQueue for the current thread
	// This may fail if WinRT is not initialized or not available
	dq, err := rt.GetForCurrentThread()
	if err == nil && dq != nil {
		d.dispatcherQueue = dq
		d.hasDispatcherQueue = true
	}

	return d, nil
}

// IsUIThread checks if the current thread is the UI thread by comparing thread IDs
func (d *WindowsDispatcher) IsUIThread() bool {
	return getCurrentThreadId() == d.uiThreadID
}

// RunOnUI schedules a function to run on the UI thread.
// If already on the UI thread, executes the function inline.
// Otherwise, uses DispatcherQueue to schedule on the UI thread.
// If DispatcherQueue is not available, returns an error.
func (d *WindowsDispatcher) RunOnUI(fn func()) error {
	// If already on UI thread, execute inline to avoid deadlock
	if d.IsUIThread() {
		fn()
		return nil
	}

	if !d.hasDispatcherQueue || d.dispatcherQueue == nil {
		return errors.New("cannot schedule to UI thread: DispatcherQueue not available")
	}

	// Use a channel to wait for the function to complete and capture any panic
	done := make(chan struct{})
	var runErr error

	enqueued, err := d.dispatcherQueue.TryEnqueue(func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				runErr = errors.New("panic in UI callback")
			}
		}()
		fn()
	})

	if err != nil {
		return err
	}

	if !enqueued {
		return errors.New("failed to enqueue callback to dispatcher queue")
	}

	// Wait for the callback to complete
	<-done

	return runErr
}
