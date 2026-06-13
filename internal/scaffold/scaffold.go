package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FallingSkyQwQ/Narcissus/cmd/narc/templates"
)

// Options 包含创建项目的选项
type Options struct {
	AppName    string
	ModuleName string
	TargetDir  string
	Template   string
}

// CreateProject 根据选项创建新项目
func CreateProject(opts Options) error {
	// 验证选项
	if opts.AppName == "" {
		return fmt.Errorf("app name is required")
	}
	if opts.ModuleName == "" {
		return fmt.Errorf("module name is required")
	}
	if opts.TargetDir == "" {
		// 默认使用当前目录下的应用名称作为目标目录
		opts.TargetDir = opts.AppName
	}
	if opts.Template == "" {
		// 默认使用 basic 模板
		opts.Template = "basic"
	}

	// 验证模板是否可用
	availableTemplates := templates.ListTemplates()
	found := false
	for _, t := range availableTemplates {
		if t == opts.Template {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("template %q not found, available templates: %v", opts.Template, availableTemplates)
	}

	// 解析目标目录为绝对路径
	absTargetDir, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve target directory: %w", err)
	}
	opts.TargetDir = absTargetDir

	// 检查目标目录是否已存在
	if _, err := os.Stat(opts.TargetDir); err == nil {
		// 目录已存在，检查是否为空
		entries, err := os.ReadDir(opts.TargetDir)
		if err != nil {
			return fmt.Errorf("failed to read target directory: %w", err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("target directory %q already exists and is not empty", opts.TargetDir)
		}
	}

	// 根据模板类型调用相应的脚手架函数
	if err := templates.ScaffoldProject(opts.TargetDir, opts.AppName, opts.ModuleName, opts.Template); err != nil {
		return fmt.Errorf("failed to scaffold project: %w", err)
	}

	fmt.Printf("✓ Created project %q at %s\n", opts.AppName, opts.TargetDir)
	fmt.Printf("✓ Using template: %s\n", opts.Template)
	fmt.Printf("✓ Module name: %s\n", opts.ModuleName)

	return nil
}
