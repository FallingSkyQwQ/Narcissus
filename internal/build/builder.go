package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Supported target operating systems.
const (
	TargetWindows = "windows"
	TargetLinux   = "linux"
	TargetDarwin  = "darwin"
)

type Options struct {
	ProjectDir string
	OutputDir  string
	AppName    string
	Compress   bool
	LDFlags    string
	// Target is the GOOS to build for. An empty value means the host OS.
	Target string
	// Arch is the GOARCH to build for. An empty value means the host GOARCH.
	Arch string
	// WindowsAppSDK fetches the Windows App SDK runtime files (bootstrapper DLL
	// and resources.pri) beside the executable for Windows targets. It is
	// ignored for other targets.
	WindowsAppSDK bool
}

// ResolveTarget returns the effective GOOS and GOARCH for the given options,
// falling back to the host runtime when they are not specified.
func ResolveTarget(opts Options) (goos, goarch string, err error) {
	goos = strings.ToLower(strings.TrimSpace(opts.Target))
	if goos == "" {
		goos = runtime.GOOS
	}
	switch goos {
	case TargetWindows, TargetLinux, TargetDarwin:
	default:
		return "", "", fmt.Errorf("unsupported target OS %q (supported: windows, linux, darwin)", goos)
	}

	goarch = strings.ToLower(strings.TrimSpace(opts.Arch))
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	return goos, goarch, nil
}

// ExecutableExtension returns the executable suffix for a GOOS.
func ExecutableExtension(goos string) string {
	if goos == TargetWindows {
		return ".exe"
	}
	return ""
}

// DefaultLDFlags returns the platform-appropriate default linker flags.
func DefaultLDFlags(goos string) string {
	if goos == TargetWindows {
		// Windows GUI subsystem so no console window appears.
		return "-s -w -H=windowsgui"
	}
	return "-s -w"
}

func Build(opts Options) error {
	if opts.OutputDir == "" {
		opts.OutputDir = "dist"
	}
	if opts.AppName == "" {
		opts.AppName = "app"
	}
	if opts.ProjectDir == "" {
		opts.ProjectDir = "."
	}

	goos, goarch, err := ResolveTarget(opts)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPath := filepath.Join(opts.OutputDir, opts.AppName+ExecutableExtension(goos))
	// The go command runs with its working directory set to ProjectDir, so a
	// relative -o would be resolved there. Resolve it up front instead.
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	ldflags := opts.LDFlags
	if ldflags == "" {
		ldflags = DefaultLDFlags(goos)
	}

	crossCompiling := goos != runtime.GOOS || goarch != runtime.GOARCH

	fmt.Printf("Building %s...\n", opts.AppName)
	fmt.Printf("  Target: %s/%s\n", goos, goarch)
	if crossCompiling {
		fmt.Printf("  Cross-compiling from %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}
	fmt.Printf("  Output: %s\n", outputPath)
	fmt.Printf("  LDFlags: %s\n", ldflags)

	args := []string{"build", "-ldflags", ldflags, "-o", absOutputPath, "."}
	cmd := exec.Command("go", args...)
	cmd.Dir = opts.ProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if crossCompiling {
		cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	info, err := os.Stat(absOutputPath)
	if err != nil {
		return fmt.Errorf("failed to stat output file: %w", err)
	}
	originalSize := info.Size()
	fmt.Printf("✓ Build complete: %s (%.2f MB)\n", outputPath, float64(originalSize)/(1024*1024))

	if opts.Compress {
		if err := compressWithUPX(absOutputPath); err != nil {
			fmt.Printf("⚠ Compression failed: %v\n", err)
			fmt.Println("  Output is still usable without compression.")
		} else {
			info, err := os.Stat(absOutputPath)
			if err != nil {
				fmt.Printf("⚠ Failed to stat compressed file: %v\n", err)
			} else if originalSize > 0 {
				compressedSize := info.Size()
				ratio := float64(originalSize-compressedSize) / float64(originalSize) * 100
				fmt.Printf("✓ Compressed: %.2f MB → %.2f MB (%.1f%% reduction)\n",
					float64(originalSize)/(1024*1024),
					float64(compressedSize)/(1024*1024),
					ratio)
			}
		}
	}

	if goos == TargetWindows && opts.WindowsAppSDK {
		if err := PrepareWindowsAppSDK(WASDKOptions{
			OutputDir: opts.OutputDir,
			AppName:   opts.AppName,
			GOARCH:    goarch,
		}); err != nil {
			fmt.Printf("⚠ Windows App SDK setup incomplete: %v\n", err)
			fmt.Println("  The executable was built, but a WinUI app also needs the bootstrapper")
			fmt.Println("  DLL and resources.pri beside it; see the README for the manual steps.")
		}
	}

	return nil
}

func compressWithUPX(exePath string) error {
	if _, err := exec.LookPath("upx"); err != nil {
		return fmt.Errorf("UPX not found in PATH")
	}

	fmt.Println("Compressing with UPX...")
	cmd := exec.Command("upx", "--best", "--lzma", exePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GetAppName(projectDir string) string {
	modFile := filepath.Join(projectDir, "go.mod")
	if data, err := os.ReadFile(modFile); err == nil {
		content := string(data)
		if idx := strings.Index(content, "module "); idx != -1 {
			start := idx + 7
			end := strings.Index(content[start:], "\n")
			if end == -1 {
				end = len(content) - start
			}
			moduleName := strings.TrimSpace(content[start : start+end])
			parts := strings.Split(moduleName, "/")
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
		}
	}

	absPath, _ := filepath.Abs(projectDir)
	return filepath.Base(absPath)
}
