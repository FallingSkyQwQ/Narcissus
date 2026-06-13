package rt

import (
	"sync"
	"syscall"
	"unsafe"

	"github.com/FallingSkyQwQ/Narcissus/pkg/bridge/com"
)

var (
	// DispatcherQueue class name
	dispatcherQueueClassName = "Windows.System.DispatcherQueue"

	// IID_IDispatcherQueue for IDispatcherQueue interface
	IID_IDispatcherQueue = com.NewGUID(
		0x5FC1548B,
		0x0E41,
		0x5B91,
		[8]byte{0xB7, 0xD5, 0xC1, 0xD4, 0xE9, 0xC0, 0xC0, 0x10},
	)

	// IID_IDispatcherQueueStatics for static methods
	IID_IDispatcherQueueStatics = com.NewGUID(
		0xA8D3D4F0,
		0xB273,
		0x5CBF,
		[8]byte{0xB5, 0xE9, 0xFB, 0x0D, 0xBC, 0xE0, 0xB0, 0x00},
	)
)

// DispatcherQueuePriority represents the priority for dispatcher queue callbacks
type DispatcherQueuePriority int32

const (
	// DispatcherQueuePriority_Low is low priority
	DispatcherQueuePriority_Low DispatcherQueuePriority = -10
	// DispatcherQueuePriority_Normal is normal priority
	DispatcherQueuePriority_Normal DispatcherQueuePriority = 0
	// DispatcherQueuePriority_High is high priority
	DispatcherQueuePriority_High DispatcherQueuePriority = 10
)

// DispatcherQueue represents a Windows.System.DispatcherQueue instance
type DispatcherQueue struct {
	inspectable *com.IInspectable
	vtable      *dispatcherQueueVTable
}

type dispatcherQueueVTable struct {
	QueryInterface         uintptr
	AddRef                 uintptr
	Release                uintptr
	GetIids                uintptr
	GetRuntimeClassName    uintptr
	GetTrustLevel          uintptr
	TryEnqueue             uintptr
	TryEnqueueWithPriority uintptr
}

// DispatcherQueueController represents a DispatcherQueue controller
type DispatcherQueueController struct {
	inspectable *com.IInspectable
}

var (
	// Global dispatcher queue for the main thread
	globalDispatcherQueue     *DispatcherQueue
	globalDispatcherQueueMu   sync.RWMutex
	globalDispatcherQueueOnce sync.Once
)

// GetForCurrentThread gets the DispatcherQueue for the current thread
func GetForCurrentThread() (*DispatcherQueue, error) {
	// First try to get the cached global instance
	globalDispatcherQueueMu.RLock()
	cached := globalDispatcherQueue
	globalDispatcherQueueMu.RUnlock()
	if cached != nil {
		return cached, nil
	}

	// Get the statics interface
	statics, err := getDispatcherQueueStatics()
	if err != nil {
		return nil, err
	}
	defer statics.Release()

	// Call GetForCurrentThread
	var queue *com.IInspectable
	hr := statics.GetForCurrentThread(&queue)
	if hr != com.S_OK {
		return nil, hr
	}
	if queue == nil {
		return nil, syscall.EINVAL
	}

	return &DispatcherQueue{
		inspectable: queue,
		vtable:      (*dispatcherQueueVTable)(unsafe.Pointer(queue.Vtbl())),
	}, nil
}

// TryEnqueue schedules a callback on the dispatcher queue
func (dq *DispatcherQueue) TryEnqueue(callback func()) (bool, error) {
	return dq.TryEnqueueWithPriority(DispatcherQueuePriority_Normal, callback)
}

// TryEnqueueWithPriority schedules a callback with specified priority
func (dq *DispatcherQueue) TryEnqueueWithPriority(priority DispatcherQueuePriority, callback func()) (bool, error) {
	if dq.vtable == nil || dq.vtable.TryEnqueueWithPriority == 0 {
		return false, com.E_NOTIMPL
	}

	// Create a delegate to wrap the callback
	delegate, err := createDispatcherQueueHandler(callback)
	if err != nil {
		return false, err
	}
	defer delegate.Release()

	var result bool
	ret, _, _ := syscall.SyscallN(
		dq.vtable.TryEnqueueWithPriority,
		uintptr(unsafe.Pointer(dq.inspectable)),
		uintptr(priority),
		uintptr(unsafe.Pointer(delegate)),
		uintptr(unsafe.Pointer(&result)),
	)

	if com.HRESULT(ret) != com.S_OK {
		return false, com.HRESULT(ret)
	}

	return result, nil
}

// Release releases the dispatcher queue
func (dq *DispatcherQueue) Release() {
	if dq.inspectable != nil {
		dq.inspectable.Release()
		dq.inspectable = nil
	}
}

// dispatcherQueueStatics wraps the static interface
type dispatcherQueueStatics struct {
	inspectable *com.IInspectable
	vtable      *dispatcherQueueStaticsVTable
}

type dispatcherQueueStaticsVTable struct {
	QueryInterface          uintptr
	AddRef                  uintptr
	Release                 uintptr
	GetIids                 uintptr
	GetRuntimeClassName     uintptr
	GetTrustLevel           uintptr
	GetForCurrentThread     uintptr
	GetForCurrentThreadID   uintptr
	CreateOnDedicatedThread uintptr
}

func getDispatcherQueueStatics() (*dispatcherQueueStatics, error) {
	// Create HSTRING for class name
	classHStr, err := NewHStringFromString(dispatcherQueueClassName)
	if err != nil {
		return nil, err
	}
	defer DeleteHString(classHStr)

	// Get activation factory
	var factory *com.IInspectable
	ret, _, _ := procRoGetActivationFactory.Call(
		uintptr(classHStr),
		uintptr(unsafe.Pointer(&IID_IDispatcherQueueStatics)),
		uintptr(unsafe.Pointer(&factory)),
	)
	if com.HRESULT(ret) != com.S_OK {
		return nil, com.HRESULT(ret)
	}
	if factory == nil {
		return nil, com.E_FAIL
	}

	return &dispatcherQueueStatics{
		inspectable: factory,
		vtable:      (*dispatcherQueueStaticsVTable)(unsafe.Pointer(factory.Vtbl())),
	}, nil
}

func (s *dispatcherQueueStatics) GetForCurrentThread(queue **com.IInspectable) com.HRESULT {
	if s.vtable == nil || s.vtable.GetForCurrentThread == 0 {
		return com.E_NOTIMPL
	}
	ret, _, _ := syscall.SyscallN(
		s.vtable.GetForCurrentThread,
		uintptr(unsafe.Pointer(s.inspectable)),
		uintptr(unsafe.Pointer(queue)),
	)
	return com.HRESULT(ret)
}

func (s *dispatcherQueueStatics) Release() {
	if s.inspectable != nil {
		s.inspectable.Release()
	}
}

// dispatcherQueueHandler is a callback wrapper
type dispatcherQueueHandler struct {
	inspectable *com.IInspectable
	callback    func()
}

// IID_IDispatcherQueueHandler for the handler interface
var IID_IDispatcherQueueHandler = com.NewGUID(
	0x5C7BAA00,
	0x0A56,
	0x5B50,
	[8]byte{0xB8, 0xA3, 0xE5, 0xE5, 0xF5, 0xE3, 0xE3, 0xE3},
)

func createDispatcherQueueHandler(callback func()) (*dispatcherQueueHandler, error) {
	// For now, return a simple wrapper - in a full implementation,
	// this would create a COM callable wrapper
	return &dispatcherQueueHandler{
		callback: callback,
	}, nil
}

func (h *dispatcherQueueHandler) Release() {
	// Nothing to release for now
}

// SetGlobalDispatcherQueue sets the global dispatcher queue for the main thread
func SetGlobalDispatcherQueue(dq *DispatcherQueue) {
	globalDispatcherQueueMu.Lock()
	defer globalDispatcherQueueMu.Unlock()
	globalDispatcherQueue = dq
}

// GetGlobalDispatcherQueue gets the global dispatcher queue
func GetGlobalDispatcherQueue() *DispatcherQueue {
	globalDispatcherQueueMu.RLock()
	defer globalDispatcherQueueMu.RUnlock()
	return globalDispatcherQueue
}

var procRoGetActivationFactory = modcombase.NewProc("RoGetActivationFactory")
