package rt

import (
	"fmt"
	"unsafe"

	"github.com/FallingSkyQwQ/Narcissus/pkg/bridge/com"
)

var (
	procRoActivateInstance = modcombase.NewProc("RoActivateInstance")
)

// HString represents a Windows Runtime HSTRING handle
type HString uintptr

// ActivateInstance creates an instance of the specified runtime class
func ActivateInstance(classID *com.GUID) (*com.IInspectable, error) {
	// Convert GUID to string format for HSTRING
	className := guidToString(classID)

	hstr, err := NewHStringFromString(className)
	if err != nil {
		return nil, err
	}
	defer DeleteHString(hstr)

	var instance *com.IInspectable
	ret, _, _ := procRoActivateInstance.Call(
		uintptr(hstr),
		uintptr(unsafe.Pointer(&instance)),
	)
	if com.HRESULT(ret) != com.S_OK {
		return nil, com.HRESULT(ret)
	}
	return instance, nil
}

// guidToString converts a GUID to the Windows Runtime class name format
// Format: "{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}"
func guidToString(guid *com.GUID) string {
	if guid == nil {
		return ""
	}
	return fmt.Sprintf("{%08X-%04X-%04X-%02X%02X-%02X%02X%02X%02X%02X%02X}",
		guid.Data1,
		guid.Data2,
		guid.Data3,
		guid.Data4[0], guid.Data4[1],
		guid.Data4[2], guid.Data4[3], guid.Data4[4], guid.Data4[5], guid.Data4[6], guid.Data4[7],
	)
}
