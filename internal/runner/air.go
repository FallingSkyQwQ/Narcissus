package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// AirConfig returns an air configuration for the given target GOOS. The
// executable suffix follows the target, so hot reload works on Linux as well
// as Windows.
func AirConfig(goos string) string {
	if goos == "" {
		goos = runtime.GOOS
	}
	ext := ""
	if goos == "windows" {
		ext = ".exe"
	}

	return fmt.Sprintf(`# Air configuration for Narcissus development
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main%s ."
bin = "tmp/main%s"
full_bin = "./tmp/main%s"
include_ext = ["go", "mod", "sum"]
exclude_dir = ["assets", "tmp", "vendor", "dist", "build"]
delay = 1000
stop_on_error = true

[log]
time = true

[color]
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"

[misc]
clean_on_exit = true
`, ext, ext, ext)
}

func SetupAir(projectDir, goos string) error {
	configPath := filepath.Join(projectDir, ".air.toml")
	return os.WriteFile(configPath, []byte(AirConfig(goos)), 0644)
}

func RunWithAir(projectDir, goos string) error {
	if err := SetupAir(projectDir, goos); err != nil {
		return fmt.Errorf("failed to setup air config: %w", err)
	}
	if _, err := exec.LookPath("air"); err != nil {
		return fmt.Errorf("air is not installed. Run: go install github.com/air-verse/air@latest")
	}
	cmd := exec.Command("air", "-c", ".air.toml")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	fmt.Println("Starting development server with hot reload...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()
	return cmd.Run()
}
