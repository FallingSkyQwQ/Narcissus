package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Options struct {
	ProjectDir string
	OutputDir  string
	AppName    string
	Compress   bool
	LDFlags    string
}

func Build(opts Options) error {
	if opts.OutputDir == "" {
		opts.OutputDir = "dist"
	}
	if opts.AppName == "" {
		opts.AppName = "app"
	}

	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPath := filepath.Join(opts.OutputDir, opts.AppName+".exe")

	ldflags := opts.LDFlags
	if ldflags == "" {
		ldflags = "-s -w -H=windowsgui"
	}

	fmt.Printf("Building %s...\n", opts.AppName)
	fmt.Printf("  Output: %s\n", outputPath)
	fmt.Printf("  LDFlags: %s\n", ldflags)

	args := []string{"build", "-ldflags", ldflags, "-o", outputPath, "."}
	cmd := exec.Command("go", args...)
	cmd.Dir = opts.ProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	info, _ := os.Stat(outputPath)
	originalSize := info.Size()
	fmt.Printf("✓ Build complete: %s (%.2f MB)\n", outputPath, float64(originalSize)/(1024*1024))

	if opts.Compress {
		if err := compressWithUPX(outputPath); err != nil {
			fmt.Printf("⚠ Compression failed: %v\n", err)
			fmt.Println("  Output is still usable without compression.")
		} else {
			info, _ = os.Stat(outputPath)
			compressedSize := info.Size()
			ratio := float64(originalSize-compressedSize) / float64(originalSize) * 100
			fmt.Printf("✓ Compressed: %.2f MB → %.2f MB (%.1f%% reduction)\n",
				float64(originalSize)/(1024*1024),
				float64(compressedSize)/(1024*1024),
				ratio)
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
