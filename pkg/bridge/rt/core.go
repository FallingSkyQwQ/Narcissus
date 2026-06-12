//go:build windows
// +build windows

package rt

import (
	"syscall"

	"github.com/FallingSkyQwQ/Narcissus/pkg/bridge/com"
)

var (
	modcombase = syscall.NewLazyDLL("combase.dll")

	procRoInitialize   = modcombase.NewProc("RoInitialize")
	procRoUninitialize = modcombase.NewProc("RoUninitialize")
)

// ThreadType defines the apartment type for the current thread
type ThreadType uint32

const (
	// RO_INIT_SINGLETHREADED initializes the thread in a single-threaded apartment
	RO_INIT_SINGLETHREADED ThreadType = 0
	// RO_INIT_MULTITHREADED initializes the thread in a multi-threaded apartment
	RO_INIT_MULTITHREADED ThreadType = 1
)

// Initialize initializes the Windows Runtime on the current thread
func Initialize(threadType ThreadType) error {
	ret, _, _ := procRoInitialize.Call(uintptr(threadType))
	// S_FALSE (0x00000001) is also ok (already initialized)
	if com.HRESULT(ret) != com.S_OK && com.HRESULT(ret) != 0x00000001 {
		return com.HRESULT(ret)
	}
	return nil
}

// Uninitialize uninitializes the Windows Runtime on the current thread
func Uninitialize() {
	procRoUninitialize.Call()
}
