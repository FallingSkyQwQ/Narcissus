package build

import (
	"runtime"
	"testing"
)

func TestResolveTargetDefaultsToHost(t *testing.T) {
	goos, goarch, err := ResolveTarget(Options{})
	if err != nil {
		t.Fatalf("ResolveTarget: %v", err)
	}
	if goos != runtime.GOOS {
		t.Errorf("goos = %q, want host %q", goos, runtime.GOOS)
	}
	if goarch != runtime.GOARCH {
		t.Errorf("goarch = %q, want host %q", goarch, runtime.GOARCH)
	}
}

func TestResolveTargetExplicit(t *testing.T) {
	goos, goarch, err := ResolveTarget(Options{Target: "Windows", Arch: "ARM64"})
	if err != nil {
		t.Fatalf("ResolveTarget: %v", err)
	}
	if goos != TargetWindows {
		t.Errorf("goos = %q, want %q", goos, TargetWindows)
	}
	if goarch != "arm64" {
		t.Errorf("goarch = %q, want arm64", goarch)
	}
}

func TestResolveTargetUnsupported(t *testing.T) {
	if _, _, err := ResolveTarget(Options{Target: "plan9"}); err == nil {
		t.Fatal("expected an error for an unsupported target")
	}
}

func TestExecutableExtension(t *testing.T) {
	if got := ExecutableExtension(TargetWindows); got != ".exe" {
		t.Errorf("windows extension = %q, want .exe", got)
	}
	if got := ExecutableExtension(TargetLinux); got != "" {
		t.Errorf("linux extension = %q, want empty", got)
	}
}

func TestWindowsArch(t *testing.T) {
	tests := map[string]string{
		"amd64": "x64",
		"386":   "x86",
		"arm64": "arm64",
		"AMD64": "x64",
	}
	for goarch, want := range tests {
		if got := windowsArch(goarch); got != want {
			t.Errorf("windowsArch(%q) = %q, want %q", goarch, got, want)
		}
	}
}

func TestDefaultLDFlags(t *testing.T) {
	if got := DefaultLDFlags(TargetWindows); got != "-s -w -H=windowsgui" {
		t.Errorf("windows ldflags = %q", got)
	}
	if got := DefaultLDFlags(TargetLinux); got != "-s -w" {
		t.Errorf("linux ldflags = %q", got)
	}
}
