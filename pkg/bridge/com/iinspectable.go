//go:build windows
// +build windows

package com

import (
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// IInspectableVtbl defines the virtual function table for IInspectable
type IInspectableVtbl struct {
	QueryInterface      uintptr
	AddRef              uintptr
	Release             uintptr
	GetIids             uintptr
	GetRuntimeClassName uintptr
	GetTrustLevel       uintptr
}

// IInspectable is the base interface for all WinRT objects
type IInspectable struct {
	vtbl *IInspectableVtbl
}

// QueryInterface queries for a specific interface
func (i *IInspectable) QueryInterface(riid *GUID) (*IInspectable, error) {
	if i == nil || i.vtbl == nil {
		return nil, E_INVALIDARG
	}
	if riid == nil {
		return nil, E_INVALIDARG
	}

	var result *IInspectable
	ret, _, _ := syscall.SyscallN(
		i.vtbl.QueryInterface,
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(&result)),
	)
	if HRESULT(ret) != S_OK {
		return nil, HRESULT(ret)
	}

	if result != nil {
		return result, nil
	}

	return nil, E_NOINTERFACE
}

// AddRef increments the reference count
func (i *IInspectable) AddRef() uint32 {
	if i == nil || i.vtbl == nil {
		return 0
	}
	ret, _, _ := syscall.SyscallN(
		i.vtbl.AddRef,
		uintptr(unsafe.Pointer(i)),
	)
	return uint32(ret)
}

// Release decrements the reference count
func (i *IInspectable) Release() uint32 {
	if i == nil || i.vtbl == nil {
		return 0
	}
	ret, _, _ := syscall.SyscallN(
		i.vtbl.Release,
		uintptr(unsafe.Pointer(i)),
	)
	return uint32(ret)
}

// GetIids gets the interface IDs implemented by this object
func (i *IInspectable) GetIids() ([]GUID, error) {
	if i == nil || i.vtbl == nil {
		return nil, E_INVALIDARG
	}
	var count uint32
	var iids *GUID
	ret, _, _ := syscall.SyscallN(
		i.vtbl.GetIids,
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&count)),
		uintptr(unsafe.Pointer(&iids)),
	)
	if HRESULT(ret) != S_OK {
		return nil, HRESULT(ret)
	}
	if count == 0 || iids == nil {
		return nil, nil
	}

	// Allocate a new Go slice and copy the GUIDs from COM-allocated memory
	result := make([]GUID, count)
	comSlice := unsafe.Slice(iids, count)
	copy(result, comSlice)

	// Free the COM-allocated memory
	windows.CoTaskMemFree(unsafe.Pointer(iids))

	return result, nil
}

// GetRuntimeClassName gets the runtime class name as HSTRING handle
// Note: The returned HSTRING must be converted using HStringToString from rt package
// and then deleted using DeleteHString
func (i *IInspectable) GetRuntimeClassName() (uintptr, error) {
	if i == nil || i.vtbl == nil {
		return 0, E_INVALIDARG
	}
	var hstr uintptr
	ret, _, _ := syscall.SyscallN(
		i.vtbl.GetRuntimeClassName,
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&hstr)),
	)
	if HRESULT(ret) != S_OK {
		return 0, HRESULT(ret)
	}
	return hstr, nil
}

// GetTrustLevel gets the trust level
func (i *IInspectable) GetTrustLevel() (uint32, error) {
	if i == nil || i.vtbl == nil {
		return 0, E_INVALIDARG
	}
	var level uint32
	ret, _, _ := syscall.SyscallN(
		i.vtbl.GetTrustLevel,
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&level)),
	)
	if HRESULT(ret) != S_OK {
		return 0, HRESULT(ret)
	}
	return level, nil
}

// NewIInspectable creates a new IInspectable without automatic finalization.
// IMPORTANT: Caller is responsible for calling Release() when done.
//
// Thread safety: COM Release calls must occur on the correct thread/apartment.
// The Go garbage collector may run finalizers on any goroutine, which can
// lead to cross-apartment COM calls and crashes. Therefore, this function
// does NOT set a finalizer. Callers must explicitly manage the lifecycle
// using Release() or use NewIInspectableWithFinalizer() if they understand
// the threading implications.
func NewIInspectable(ptr unsafe.Pointer) *IInspectable {
	if ptr == nil {
		return nil
	}
	return (*IInspectable)(ptr)
}

// NewIInspectableWithFinalizer creates a new IInspectable with automatic Release on GC.
// WARNING: This is unsafe for COM objects with apartment threading requirements.
// The finalizer may call Release() from an arbitrary goroutine/thread, which can
// cause crashes if the COM object requires a specific apartment. Only use this
// for COM objects that are thread-safe or free-threaded (marshaled).
//
// For most WinRT/COM objects, prefer NewIInspectable() and explicit Release().
func NewIInspectableWithFinalizer(ptr unsafe.Pointer) *IInspectable {
	obj := NewIInspectable(ptr)
	if obj != nil {
		// Note: This is potentially unsafe - see function documentation
		runtime.SetFinalizer(obj, (*IInspectable).Release)
	}
	return obj
}

// UnsafePtr returns the unsafe pointer to the IInspectable
func (i *IInspectable) UnsafePtr() unsafe.Pointer {
	return unsafe.Pointer(i)
}

// Vtbl returns the vtable pointer
func (i *IInspectable) Vtbl() *IInspectableVtbl {
	if i == nil {
		return nil
	}
	return i.vtbl
}
