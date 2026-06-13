package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed basic/*
var basicFS embed.FS

// TemplateData 包含模板渲染所需的数据
type TemplateData struct {
	AppName    string
	ModuleName string
}

// ScaffoldProject 根据模板生成项目脚手架
func ScaffoldProject(targetDir, appName, moduleName string) error {
	data := TemplateData{
		AppName:    appName,
		ModuleName: moduleName,
	}

	// 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// 遍历嵌入的文件系统
	return fs.WalkDir(basicFS, "basic", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if d.IsDir() {
			return nil
		}

		// 读取模板文件内容
		content, err := basicFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		// 解析模板
		tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		// 生成目标文件名（去掉 .tmpl 后缀）
		relPath, err := filepath.Rel("basic", path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		targetFileName := strings.TrimSuffix(relPath, ".tmpl")
		targetFilePath := filepath.Join(targetDir, targetFileName)

		// 确保目标文件的父目录存在
		targetFileDir := filepath.Dir(targetFilePath)
		if err := os.MkdirAll(targetFileDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", targetFilePath, err)
		}

		// 创建目标文件
		file, err := os.Create(targetFilePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", targetFilePath, err)
		}
		defer file.Close()

		// 执行模板
		if err := tmpl.Execute(file, data); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", path, err)
		}

		return nil
	})
}

// ListTemplates 返回可用的模板列表
func ListTemplates() []string {
	return []string{"basic"}
}
