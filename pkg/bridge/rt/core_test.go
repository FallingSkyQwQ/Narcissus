//go:build windows
// +build windows

package rt

import (
	"runtime"
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/bridge/com"
)

func TestInitialize(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := Initialize(RO_INIT_SINGLETHREADED)
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
		return
	}
	defer Uninitialize()
}

func TestInitializeMultithreaded(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := Initialize(RO_INIT_MULTITHREADED)
	if err != nil {
		t.Errorf("Initialize with MULTITHREADED failed: %v", err)
		return
	}
	defer Uninitialize()
}

func TestDoubleInitialize(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// First initialization
	err := Initialize(RO_INIT_SINGLETHREADED)
	if err != nil {
		t.Errorf("First Initialize failed: %v", err)
		return
	}
	// Defer first Uninitialize
	defer Uninitialize()

	// Second initialization should return S_FALSE (which we treat as success)
	err = Initialize(RO_INIT_SINGLETHREADED)
	if err != nil {
		t.Errorf("Second Initialize should succeed (S_FALSE is ok): %v", err)
		return
	}
	// Must call Uninitialize a second time to balance the second Initialize
	defer Uninitialize()
}

func TestThreadTypeConstants(t *testing.T) {
	if RO_INIT_SINGLETHREADED != 0 {
		t.Errorf("RO_INIT_SINGLETHREADED should be 0, got %d", RO_INIT_SINGLETHREADED)
	}
	if RO_INIT_MULTITHREADED != 1 {
		t.Errorf("RO_INIT_MULTITHREADED should be 1, got %d", RO_INIT_MULTITHREADED)
	}
}

func TestGuidToString(t *testing.T) {
	guid := &com.GUID{
		Data1: 0x12345678,
		Data2: 0x1234,
		Data3: 0x5678,
		Data4: [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	result := guidToString(guid)
	expected := "{12345678-1234-5678-0102-030405060708}"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestGuidToStringNil(t *testing.T) {
	result := guidToString(nil)
	if result != "" {
		t.Errorf("expected empty string for nil GUID, got %q", result)
	}
}
