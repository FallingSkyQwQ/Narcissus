//go:build windows
// +build windows

package com

import (
	"fmt"
	"unsafe"
)

// GUID represents a Windows GUID/UUID
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// HRESULT represents COM result codes
type HRESULT int32

// HRESULT constants - using negative decimal values to avoid overflow
const (
	S_OK           HRESULT = 0
	S_FALSE        HRESULT = 1           // 0x00000001 - operation successful but returned false
	E_FAIL         HRESULT = -2147467259 // 0x80004005
	E_INVALIDARG   HRESULT = -2147024809 // 0x80070057
	E_NOINTERFACE  HRESULT = -2147467262 // 0x80004002
	E_OUTOFMEMORY  HRESULT = -2147024882 // 0x8007000E
	E_NOTIMPL      HRESULT = -2147467263 // 0x80004001
	E_ACCESSDENIED HRESULT = -2147024891 // 0x80070005
	E_UNEXPECTED   HRESULT = -2147418113 // 0x8000FFFF
)

// Error returns the error representation of the HRESULT
func (h HRESULT) Error() string {
	switch h {
	case S_OK:
		return "S_OK"
	case E_FAIL:
		return "E_FAIL"
	case E_INVALIDARG:
		return "E_INVALIDARG"
	case E_NOINTERFACE:
		return "E_NOINTERFACE"
	case E_OUTOFMEMORY:
		return "E_OUTOFMEMORY"
	case E_NOTIMPL:
		return "E_NOTIMPL"
	case E_ACCESSDENIED:
		return "E_ACCESSDENIED"
	case E_UNEXPECTED:
		return "E_UNEXPECTED"
	default:
		return fmt.Sprintf("HRESULT(0x%08X)", int32(h))
	}
}

// RefCount is used for COM reference counting
type RefCount int32

// IID_IInspectable is the interface ID for IInspectable
var IID_IInspectable = GUID{
	Data1: 0xAF86E2E0,
	Data2: 0xB12D,
	Data3: 0x4C6A,
	Data4: [8]byte{0x9C, 0x5A, 0xD7, 0xAA, 0x65, 0x10, 0x1E, 0x90},
}

// IID_IUnknown is the interface ID for IUnknown
var IID_IUnknown = GUID{
	Data1: 0x00000000,
	Data2: 0x0000,
	Data3: 0x0000,
	Data4: [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46},
}

// NewGUID creates a GUID from individual components
func NewGUID(data1 uint32, data2 uint16, data3 uint16, data4 [8]byte) GUID {
	return GUID{
		Data1: data1,
		Data2: data2,
		Data3: data3,
		Data4: data4,
	}
}

// ToBytes converts the GUID to a byte slice
func (g *GUID) ToBytes() []byte {
	b := make([]byte, 16)
	*(*uint32)(unsafe.Pointer(&b[0])) = g.Data1
	*(*uint16)(unsafe.Pointer(&b[4])) = g.Data2
	*(*uint16)(unsafe.Pointer(&b[6])) = g.Data3
	copy(b[8:], g.Data4[:])
	return b
}
