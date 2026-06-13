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

// WindowsDispatcher implements Dispatcher using Windows thread ID and DispatcherQueue
type WindowsDispatcher struct {
	uiThreadID      uint32
	dispatcherQueue *rt.DispatcherQueue
	// hasDispatcherQueue indicates whether DispatcherQueue is available
	// If false, RunOnUI will return an error when called from non-UI thread
	hasDispatcherQueue bool
}

// NewWindowsDispatcher creates a new Windows dispatcher.
// It captures the current thread ID as the UI thread ID and tries to get the DispatcherQueue.
// If DispatcherQueue is not available (e.g., WinRT not initialized), it still creates
// a functional dispatcher that works for UI thread detection.
func NewWindowsDispatcher() (*WindowsDispatcher, error) {
	tid, _, err := procGetCurrentThreadId.Call()
	if err != nil && err != syscall.Errno(0) {
		return nil, err
	}

	d := &WindowsDispatcher{
		uiThreadID:         uint32(tid),
		hasDispatcherQueue: false,
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

// IsUIThread returns true if the current thread is the UI thread
func (d *WindowsDispatcher) IsUIThread() bool {
	tid, _, err := procGetCurrentThreadId.Call()
	if err != nil && err != syscall.Errno(0) {
		return false
	}
	return uint32(tid) == d.uiThreadID
}

// RunOnUI schedules a function to run on the UI thread.
// If already on UI thread, the function runs immediately.
// If not on UI thread and DispatcherQueue is available, uses TryEnqueue to schedule.
// If not on UI thread and DispatcherQueue is not available, returns an error.
func (d *WindowsDispatcher) RunOnUI(fn func()) error {
	if d.IsUIThread() {
		fn()
		return nil
	}

	if !d.hasDispatcherQueue || d.dispatcherQueue == nil {
		return errors.New("cannot schedule to UI thread from non-UI thread: DispatcherQueue not available")
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
