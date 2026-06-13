package rt

import (
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/FallingSkyQwQ/Narcissus/pkg/bridge/com"
)

var (
	// DispatcherQueue class name
	dispatcherQueueClassName = "Windows.System.DispatcherQueue"

	// IID_IDispatcherQueue for IDispatcherQueue interface
	// {5FC1548B-0E41-5B91-B7D5-C1D4E9C0C010}
	IID_IDispatcherQueue = com.NewGUID(
		0x5FC1548B,
		0x0E41,
		0x5B91,
		[8]byte{0xB7, 0xD5, 0xC1, 0xD4, 0xE9, 0xC0, 0xC0, 0x10},
	)

	// IID_IDispatcherQueueStatics for static methods
	// {A8D3D4F0-B273-5CBF-B5E9-FB0DBCE0B000}
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
	globalDispatcherQueue   *DispatcherQueue
	globalDispatcherQueueMu sync.RWMutex
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

	// Create a COM callable wrapper for the callback
	delegate, err := createDispatcherQueueHandler(callback)
	if err != nil {
		return false, err
	}
	defer delegate.release()

	var result bool
	ret, _, _ := syscall.SyscallN(
		dq.vtable.TryEnqueueWithPriority,
		uintptr(unsafe.Pointer(dq.inspectable)),
		uintptr(priority),
		delegate.objectPtr(),
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

// Global vtable for dispatcherQueueHandler
// This is initialized once and shared across all handler instances
var (
	dispatcherQueueHandlerVTableInstance *dispatcherQueueHandlerVTable
	dispatcherQueueHandlerVTableOnce     sync.Once
)

// dispatcherQueueHandlerVTable defines the COM vtable for IDispatcherQueueHandler
// IDispatcherQueueHandler implements IUnknown + Invoke method
type dispatcherQueueHandlerVTable struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	Invoke         uintptr
}

// getDispatcherQueueHandlerVTable returns the singleton vtable for handlers
func getDispatcherQueueHandlerVTable() *dispatcherQueueHandlerVTable {
	dispatcherQueueHandlerVTableOnce.Do(func() {
		vtable := &dispatcherQueueHandlerVTable{}
		vtable.QueryInterface = syscall.NewCallback(dispatcherQueueHandlerQueryInterface)
		vtable.AddRef = syscall.NewCallback(dispatcherQueueHandlerAddRef)
		vtable.Release = syscall.NewCallback(dispatcherQueueHandlerRelease)
		vtable.Invoke = syscall.NewCallback(dispatcherQueueHandlerInvoke)
		dispatcherQueueHandlerVTableInstance = vtable
	})
	return dispatcherQueueHandlerVTableInstance
}

// dispatcherQueueHandler represents a COM callable wrapper for Go callbacks
type dispatcherQueueHandler struct {
	vtable   *dispatcherQueueHandlerVTable
	refCount int32
	callback func()
}

// createDispatcherQueueHandler creates a new COM callable wrapper for the callback
func createDispatcherQueueHandler(callback func()) (*dispatcherQueueHandler, error) {
	handler := &dispatcherQueueHandler{
		vtable:   getDispatcherQueueHandlerVTable(),
		refCount: 1,
		callback: callback,
	}
	return handler, nil
}

// objectPtr returns the pointer to use as IInspectable*
func (h *dispatcherQueueHandler) objectPtr() uintptr {
	if h == nil {
		return 0
	}
	return uintptr(unsafe.Pointer(h))
}

// release decrements the reference count
func (h *dispatcherQueueHandler) release() {
	if h == nil {
		return
	}
	if atomic.AddInt32(&h.refCount, -1) == 0 {
		// Object is being destroyed
		h.callback = nil
	}
}

// COM method implementations for dispatcherQueueHandler

// dispatcherQueueHandlerQueryInterface implements IUnknown::QueryInterface
func dispatcherQueueHandlerQueryInterface(this uintptr, riid unsafe.Pointer, ppvObject unsafe.Pointer) uintptr {
	// For now, we only support IUnknown
	// In a full implementation, we would check riid against IUnknown and IDispatcherQueueHandler
	if ppvObject == nil {
		return 0x80070057 // E_INVALIDARG
	}
	// Return this pointer as the interface
	*(*uintptr)(ppvObject) = this
	// Increment reference count
	handler := (*dispatcherQueueHandler)(unsafe.Pointer(this))
	atomic.AddInt32(&handler.refCount, 1)
	return 0 // S_OK
}

// dispatcherQueueHandlerAddRef implements IUnknown::AddRef
func dispatcherQueueHandlerAddRef(this uintptr) uintptr {
	handler := (*dispatcherQueueHandler)(unsafe.Pointer(this))
	newCount := atomic.AddInt32(&handler.refCount, 1)
	return uintptr(newCount)
}

// dispatcherQueueHandlerRelease implements IUnknown::Release
func dispatcherQueueHandlerRelease(this uintptr) uintptr {
	handler := (*dispatcherQueueHandler)(unsafe.Pointer(this))
	newCount := atomic.AddInt32(&handler.refCount, -1)
	if newCount == 0 {
		// Object is being destroyed
		handler.callback = nil
	}
	return uintptr(newCount)
}

// dispatcherQueueHandlerInvoke implements IDispatcherQueueHandler::Invoke
func dispatcherQueueHandlerInvoke(this uintptr) uintptr {
	handler := (*dispatcherQueueHandler)(unsafe.Pointer(this))
	if handler.callback != nil {
		handler.callback()
	}
	return 0 // S_OK
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
