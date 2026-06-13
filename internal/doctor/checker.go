package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type CheckResult struct {
	Name    string
	Status  Status
	Message string
	Fix     string
}

type Status int

const (
	StatusOK Status = iota
	StatusWarning
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "✓"
	case StatusWarning:
		return "⚠"
	case StatusError:
		return "✗"
	}
	return "?"
}

func RunChecks() []CheckResult {
	var results []CheckResult
	results = append(results, checkOS())
	results = append(results, checkGoVersion())
	results = append(results, checkGit())
	results = append(results, checkUPX())
	results = append(results, checkAir())
	results = append(results, checkWindowsSDK())
	return results
}

func checkOS() CheckResult {
	os := runtime.GOOS
	switch os {
	case "windows":
		return CheckResult{
			Name:    "Operating System",
			Status:  StatusOK,
			Message: fmt.Sprintf("Windows (%s)", runtime.GOARCH),
		}
	case "linux", "darwin":
		return CheckResult{
			Name:    "Operating System",
			Status:  StatusWarning,
			Message: fmt.Sprintf("%s (%s) - Narcissus is designed for Windows", strings.Title(os), runtime.GOARCH),
			Fix:     "Run on Windows for full compatibility",
		}
	default:
		return CheckResult{
			Name:    "Operating System",
			Status:  StatusWarning,
			Message: fmt.Sprintf("%s (%s) - unsupported platform", os, runtime.GOARCH),
			Fix:     "Use Windows 10/11 for best experience",
		}
	}
}

func checkGoVersion() CheckResult {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "Go",
			Status:  StatusError,
			Message: "Not installed or not in PATH",
			Fix:     "Install Go from https://golang.org/dl/",
		}
	}

	versionStr := strings.TrimSpace(string(output))
	// Extract version number from "go version go1.21.0 windows/amd64"
	re := regexp.MustCompile(`go(\d+\.\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(versionStr)

	if len(matches) < 2 {
		return CheckResult{
			Name:    "Go",
			Status:  StatusWarning,
			Message: versionStr,
		}
	}

	version := matches[1]
	return CheckResult{
		Name:    "Go",
		Status:  StatusOK,
		Message: fmt.Sprintf("Version %s", version),
	}
}

func checkGit() CheckResult {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "Git",
			Status:  StatusError,
			Message: "Not installed or not in PATH",
			Fix:     "Install Git from https://git-scm.com/",
		}
	}

	versionStr := strings.TrimSpace(string(output))
	// Extract version from "git version 2.42.0.windows.1"
	re := regexp.MustCompile(`git version (\d+\.\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(versionStr)

	if len(matches) < 2 {
		return CheckResult{
			Name:    "Git",
			Status:  StatusOK,
			Message: "Installed",
		}
	}

	version := matches[1]
	return CheckResult{
		Name:    "Git",
		Status:  StatusOK,
		Message: fmt.Sprintf("Version %s", version),
	}
}

func checkUPX() CheckResult {
	cmd := exec.Command("upx", "--version")
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "UPX",
			Status:  StatusWarning,
			Message: "Not installed or not in PATH",
			Fix:     "Install UPX from https://upx.github.io/ (optional, for binary compression)",
		}
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return CheckResult{
			Name:    "UPX",
			Status:  StatusOK,
			Message: strings.TrimSpace(lines[0]),
		}
	}

	return CheckResult{
		Name:    "UPX",
		Status:  StatusOK,
		Message: "Installed",
	}
}

func checkAir() CheckResult {
	cmd := exec.Command("air", "-v")
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "Air",
			Status:  StatusWarning,
			Message: "Not installed or not in PATH",
			Fix:     "Install with: go install github.com/cosmtrek/air@latest (optional, for hot reload)",
		}
	}

	versionStr := strings.TrimSpace(string(output))
	return CheckResult{
		Name:    "Air",
		Status:  StatusOK,
		Message: versionStr,
	}
}

func checkWindowsSDK() CheckResult {
	// Check for Windows SDK by looking for rc.exe (resource compiler)
	cmd := exec.Command("where", "rc.exe")
	_, err := cmd.Output()
	if err != nil {
		// Try alternative: check for Windows SDK environment variable
		sdkPath := os.Getenv("WindowsSdkDir")
		if sdkPath == "" {
			return CheckResult{
				Name:    "Windows SDK",
				Status:  StatusWarning,
				Message: "Not detected (rc.exe not found)",
				Fix:     "Install Windows SDK via Visual Studio Installer (optional, for building .syso files)",
			}
		}
		return CheckResult{
			Name:    "Windows SDK",
			Status:  StatusOK,
			Message: fmt.Sprintf("Found at %s", sdkPath),
		}
	}

	return CheckResult{
		Name:    "Windows SDK",
		Status:  StatusOK,
		Message: "rc.exe found in PATH",
	}
}
