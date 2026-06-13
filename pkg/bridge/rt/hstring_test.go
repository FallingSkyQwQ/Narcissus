//go:build windows
// +build windows

package rt

import (
	"testing"
)

func TestHStringRoundTrip(t *testing.T) {
	original := "Hello, WinRT!"
	hstr, err := NewHStringFromString(original)
	if err != nil {
		t.Fatalf("NewHStringFromString failed: %v", err)
	}
	defer DeleteHString(hstr)

	result := HStringToString(hstr)
	if result != original {
		t.Errorf("expected %q, got %q", original, result)
	}
}

func TestHStringEmpty(t *testing.T) {
	original := ""
	hstr, err := NewHStringFromString(original)
	if err != nil {
		t.Fatalf("NewHStringFromString failed: %v", err)
	}
	defer DeleteHString(hstr)

	result := HStringToString(hstr)
	if result != original {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestHStringUnicode(t *testing.T) {
	original := "Hello, 世界! 🌍"
	hstr, err := NewHStringFromString(original)
	if err != nil {
		t.Fatalf("NewHStringFromString failed: %v", err)
	}
	defer DeleteHString(hstr)

	result := HStringToString(hstr)
	if result != original {
		t.Errorf("expected %q, got %q", original, result)
	}
}

func TestHStringIsEmpty(t *testing.T) {
	var hstr HString = 0
	if !hstr.IsEmpty() {
		t.Error("HString(0) should be empty")
	}

	hstr, err := NewHStringFromString("test")
	if err != nil {
		t.Fatalf("NewHStringFromString failed: %v", err)
	}
	defer DeleteHString(hstr)

	if hstr.IsEmpty() {
		t.Error("Non-zero HString should not be empty")
	}
}

func TestHStringToStringWithZero(t *testing.T) {
	result := HStringToString(0)
	if result != "" {
		t.Errorf("expected empty string for HString(0), got %q", result)
	}
}
