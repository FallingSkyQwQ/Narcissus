package ui

import "testing"

type fakeDPRWindow struct{ ratio float32 }

func (w *fakeDPRWindow) DevicePixelRatio() float32 { return w.ratio }

func TestDevicePixelRatioDefaultAndClamp(t *testing.T) {
	SetDevicePixelRatio(0)
	if got := DevicePixelRatio(); got != 1 {
		t.Errorf("ratio for invalid input = %v, want 1", got)
	}
	SetDevicePixelRatio(-2)
	if got := DevicePixelRatio(); got != 1 {
		t.Errorf("ratio for negative input = %v, want 1", got)
	}
	t.Cleanup(func() { SetDevicePixelRatio(1) })
}

func TestDevicePixelRatioScaling(t *testing.T) {
	SetDevicePixelRatio(2)
	t.Cleanup(func() { SetDevicePixelRatio(1) })

	if got := ScaleToDevice(10); got != 20 {
		t.Errorf("ScaleToDevice(10) = %v, want 20", got)
	}
	if got := ScaleToLogical(20); got != 10 {
		t.Errorf("ScaleToLogical(20) = %v, want 10", got)
	}
}

func TestApplyWindowPixelRatio(t *testing.T) {
	t.Cleanup(func() { SetDevicePixelRatio(1) })

	applyWindowPixelRatio(&fakeDPRWindow{ratio: 1.5})
	if got := DevicePixelRatio(); got != 1.5 {
		t.Errorf("ratio = %v, want 1.5", got)
	}

	// Windows without the capability leave the ratio untouched.
	applyWindowPixelRatio(struct{}{})
	if got := DevicePixelRatio(); got != 1.5 {
		t.Errorf("ratio changed for non-provider window: %v", got)
	}
	applyWindowPixelRatio(nil)
	if got := DevicePixelRatio(); got != 1.5 {
		t.Errorf("ratio changed for nil window: %v", got)
	}
}
