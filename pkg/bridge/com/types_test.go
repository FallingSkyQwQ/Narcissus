package com

import (
	"testing"
)

func TestGUID(t *testing.T) {
	guid := GUID{
		Data1: 0x12345678,
		Data2: 0x1234,
		Data3: 0x5678,
		Data4: [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}
	if guid.Data1 != 0x12345678 {
		t.Errorf("expected Data1 = 0x12345678, got 0x%x", guid.Data1)
	}
	if guid.Data2 != 0x1234 {
		t.Errorf("expected Data2 = 0x1234, got 0x%x", guid.Data2)
	}
	if guid.Data3 != 0x5678 {
		t.Errorf("expected Data3 = 0x5678, got 0x%x", guid.Data3)
	}
	expectedData4 := [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	if guid.Data4 != expectedData4 {
		t.Errorf("expected Data4 = %v, got %v", expectedData4, guid.Data4)
	}
}

func TestNewGUID(t *testing.T) {
	data4 := [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	guid := NewGUID(0x12345678, 0x1234, 0x5678, data4)

	if guid.Data1 != 0x12345678 {
		t.Errorf("expected Data1 = 0x12345678, got 0x%x", guid.Data1)
	}
}

func TestGUIDToBytes(t *testing.T) {
	guid := GUID{
		Data1: 0x12345678,
		Data2: 0x1234,
		Data3: 0x5678,
		Data4: [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	bytes := guid.ToBytes()
	if len(bytes) != 16 {
		t.Errorf("expected 16 bytes, got %d", len(bytes))
	}
}

func TestHRESULT(t *testing.T) {
	if S_OK != 0 {
		t.Error("S_OK should be 0")
	}
	if E_FAIL != -2147467259 {
		t.Errorf("E_FAIL should be -2147467259, got %d", E_FAIL)
	}
	if E_INVALIDARG != -2147024809 {
		t.Errorf("E_INVALIDARG should be -2147024809, got %d", E_INVALIDARG)
	}
	if E_NOINTERFACE != -2147467262 {
		t.Errorf("E_NOINTERFACE should be -2147467262, got %d", E_NOINTERFACE)
	}
}

func TestHRESULTError(t *testing.T) {
	tests := []struct {
		hr       HRESULT
		expected string
	}{
		{S_OK, "S_OK"},
		{E_FAIL, "E_FAIL"},
		{E_INVALIDARG, "E_INVALIDARG"},
		{E_NOINTERFACE, "E_NOINTERFACE"},
		{E_OUTOFMEMORY, "E_OUTOFMEMORY"},
		{E_NOTIMPL, "E_NOTIMPL"},
		{E_ACCESSDENIED, "E_ACCESSDENIED"},
		{E_UNEXPECTED, "E_UNEXPECTED"},
	}

	for _, tt := range tests {
		result := tt.hr.Error()
		if result != tt.expected {
			t.Errorf("HRESULT(0x%x).Error() = %q, expected %q", tt.hr, result, tt.expected)
		}
	}
}

func TestPredefinedGUIDs(t *testing.T) {
	// Test IID_IInspectable
	if IID_IInspectable.Data1 != 0xAF86E2E0 {
		t.Errorf("IID_IInspectable.Data1 incorrect: 0x%x", IID_IInspectable.Data1)
	}

	// Test IID_IUnknown
	if IID_IUnknown.Data1 != 0x00000000 {
		t.Errorf("IID_IUnknown.Data1 incorrect: 0x%x", IID_IUnknown.Data1)
	}
}
