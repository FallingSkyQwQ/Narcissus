//go:build windows
// +build windows

package rt

import (
	"unicode/utf16"
	"unsafe"

	"github.com/FallingSkyQwQ/Narcissus/pkg/bridge/com"
)

var (
	procWindowsCreateString       = modcombase.NewProc("WindowsCreateString")
	procWindowsDeleteString       = modcombase.NewProc("WindowsDeleteString")
	procWindowsGetStringRawBuffer = modcombase.NewProc("WindowsGetStringRawBuffer")
)

// NewHStringFromString creates an HSTRING from a Go string
func NewHStringFromString(s string) (HString, error) {
	utf16Str := utf16.Encode([]rune(s))
	if len(utf16Str) == 0 {
		// Empty string - WindowsCreateString with null pointer and 0 length
		var hstr HString
		ret, _, _ := procWindowsCreateString.Call(
			uintptr(0),
			uintptr(0),
			uintptr(unsafe.Pointer(&hstr)),
		)
		if com.HRESULT(ret) != com.S_OK {
			return 0, com.HRESULT(ret)
		}
		return hstr, nil
	}

	var hstr HString
	ret, _, _ := procWindowsCreateString.Call(
		uintptr(unsafe.Pointer(&utf16Str[0])),
		uintptr(len(utf16Str)),
		uintptr(unsafe.Pointer(&hstr)),
	)
	if com.HRESULT(ret) != com.S_OK {
		return 0, com.HRESULT(ret)
	}
	return hstr, nil
}

// DeleteHString deletes an HSTRING and releases associated memory
func DeleteHString(hstr HString) {
	if hstr != 0 {
		procWindowsDeleteString.Call(uintptr(hstr))
	}
}

// HStringToString converts an HSTRING to a Go string
func HStringToString(hstr HString) string {
	if hstr == 0 {
		return ""
	}
	var length uint32
	ptr, _, _ := procWindowsGetStringRawBuffer.Call(
		uintptr(hstr),
		uintptr(unsafe.Pointer(&length)),
	)
	if ptr == 0 || length == 0 {
		return ""
	}
	utf16Str := unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), length)
	return string(utf16.Decode(utf16Str))
}

// IsEmpty returns true if the HSTRING is empty or null
func (h HString) IsEmpty() bool {
	return h == 0
}

// String returns the string representation (for debugging)
func (h HString) String() string {
	return HStringToString(h)
}
