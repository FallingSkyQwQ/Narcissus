package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const AirConfig = `# Air configuration for Narcissus development
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main.exe ."
bin = "tmp/main.exe"
full_bin = "./tmp/main.exe"
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
`

func SetupAir(projectDir string) error {
	configPath := filepath.Join(projectDir, ".air.toml")
	if _, err := os.Stat(configPath); err == nil {
		return nil
	}
	return os.WriteFile(configPath, []byte(AirConfig), 0644)
}

func RunWithAir(projectDir string) error {
	if err := SetupAir(projectDir); err != nil {
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
