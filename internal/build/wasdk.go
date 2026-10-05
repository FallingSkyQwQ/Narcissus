package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// A WinUI 3 application built with the pure-Go Windows App SDK bindings needs
// two things beside its executable that `go build` does not produce:
//
//  1. Microsoft.WindowsAppRuntime.Bootstrap.dll -- an unpackaged app cannot
//     activate any Microsoft.UI.* type until the Windows App SDK framework
//     package is added to its process, and that DLL is only shipped in the
//     NuGet package.
//  2. resources.pri -- WinUI resolves its default control styles through
//     ms-appx:/// URIs, which an unpackaged app resolves against its own
//     resource index. Without it the controls whose templates come from that
//     index cannot be created.
//
// Both are produced by the tooling in the bindings module rather than
// reimplemented here. Generating resources.pri additionally needs makepri and
// the installed framework package, so it only runs on Windows; on another host
// the bootstrapper is fetched and the remaining steps are printed.
const (
	windowsAppSDKGenerator       = "github.com/deploymenttheory/go-bindings-windowsappsdk/cmd/generate"
	windowsAppSDKBindingsVersion = "v0.1.0"
	bootstrapDLLName             = "Microsoft.WindowsAppRuntime.Bootstrap.dll"
	resourcesPRIName             = "resources.pri"
)

// WASDKOptions describes where the Windows App SDK artifacts should be placed.
type WASDKOptions struct {
	// OutputDir is the directory holding the built executable.
	OutputDir string
	// AppName is the executable base name, used as the resource map name.
	AppName string
	// GOARCH is the target architecture (amd64, 386, arm64).
	GOARCH string
}

// PrepareWindowsAppSDK fetches the bootstrapper DLL and, on Windows, generates
// resources.pri into the output directory.
func PrepareWindowsAppSDK(opts WASDKOptions) error {
	if opts.OutputDir == "" {
		opts.OutputDir = "dist"
	}
	if opts.GOARCH == "" {
		opts.GOARCH = runtime.GOARCH
	}

	absOut, err := filepath.Abs(opts.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output directory: %w", err)
	}
	if err := os.MkdirAll(absOut, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	arch := windowsArch(opts.GOARCH)

	fmt.Println("Preparing Windows App SDK runtime files...")
	if err := runGoTool(
		"fetch-bootstrap", "--out", absOut, "--arch", arch,
	); err != nil {
		return fmt.Errorf("fetching %s: %w", bootstrapDLLName, err)
	}

	if runtime.GOOS != "windows" {
		fmt.Println("  resources.pri was not generated: it needs makepri and the installed")
		fmt.Println("  Windows App SDK framework package, so run these on Windows:")
		fmt.Printf("    go run %s@%s app-resources --out %s --name %s\n",
			windowsAppSDKGenerator, windowsAppSDKBindingsVersion, opts.OutputDir, opts.AppName)
		return nil
	}

	if err := runGoTool(
		"app-resources", "--out", absOut, "--name", opts.AppName,
	); err != nil {
		return fmt.Errorf("generating %s: %w", resourcesPRIName, err)
	}

	fmt.Printf("✓ Windows App SDK runtime files ready in %s\n", opts.OutputDir)
	return nil
}

// runGoTool invokes a command from the bindings generator at the pinned
// version, so it does not depend on the target project's module graph.
func runGoTool(command string, args ...string) error {
	fullArgs := append([]string{
		"run", windowsAppSDKGenerator + "@" + windowsAppSDKBindingsVersion, command,
	}, args...)
	cmd := exec.Command("go", fullArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// windowsArch maps a Go architecture to the naming the Windows App SDK
// redistributable uses.
func windowsArch(goarch string) string {
	switch strings.ToLower(goarch) {
	case "386":
		return "x86"
	case "arm64":
		return "arm64"
	default:
		return "x64"
	}
}
