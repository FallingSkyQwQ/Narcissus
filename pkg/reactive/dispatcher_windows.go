package reactive

import (
	"errors"
	"syscall"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
)

// WindowsDispatcher implements Dispatcher using Windows thread ID
type WindowsDispatcher struct {
	uiThreadID uint32
}

// NewWindowsDispatcher creates a new Windows dispatcher.
// It captures the current thread ID as the UI thread ID.
func NewWindowsDispatcher() (*WindowsDispatcher, error) {
	tid, _, err := procGetCurrentThreadId.Call()
	if err != nil && err != syscall.Errno(0) {
		return nil, err
	}
	return &WindowsDispatcher{
		uiThreadID: uint32(tid),
	}, nil
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
// If not on UI thread, returns an error (DispatcherQueue not yet implemented).
func (d *WindowsDispatcher) RunOnUI(fn func()) error {
	if d.IsUIThread() {
		fn()
		return nil
	}
	// TODO: Implement DispatcherQueue.TryEnqueue when WinRT bindings are ready
	return errors.New("cannot schedule to UI thread from non-UI thread: DispatcherQueue not implemented")
}
