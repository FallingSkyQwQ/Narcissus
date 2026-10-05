package ui

import "sync"

// devicePixelRatioProvider is implemented by backends whose window knows the
// screen's device pixel ratio (for example GTK's scale factor or WinUI's
// rasterization scale). The ratio is exposed through the package-level helpers
// below so application code and image loading can adapt to HiDPI displays.
type devicePixelRatioProvider interface {
	DevicePixelRatio() float32
}

var (
	dprMu sync.RWMutex
	// dpr is the current device pixel ratio; 1.0 means one device pixel per
	// logical pixel.
	dpr float32 = 1
)

// SetDevicePixelRatio records the current device pixel ratio. Values <= 0 fall
// back to 1. Backends call this as their window is laid out.
func SetDevicePixelRatio(ratio float32) {
	if ratio <= 0 {
		ratio = 1
	}
	dprMu.Lock()
	dpr = ratio
	dprMu.Unlock()
}

// DevicePixelRatio returns the current device pixel ratio (default 1).
func DevicePixelRatio() float32 {
	dprMu.RLock()
	defer dprMu.RUnlock()
	return dpr
}

// ScaleToDevice converts a logical (density-independent) length to device
// pixels.
func ScaleToDevice(logical float32) float32 {
	return logical * DevicePixelRatio()
}

// ScaleToLogical converts a device-pixel length to logical units.
func ScaleToLogical(device float32) float32 {
	ratio := DevicePixelRatio()
	if ratio <= 0 {
		ratio = 1
	}
	return device / ratio
}

// applyWindowPixelRatio updates the global ratio from a window that reports
// one. Windows that do not implement devicePixelRatioProvider leave it as is.
func applyWindowPixelRatio(window any) {
	if provider, ok := window.(devicePixelRatioProvider); ok && provider != nil {
		SetDevicePixelRatio(provider.DevicePixelRatio())
	}
}
