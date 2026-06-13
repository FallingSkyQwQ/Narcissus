package pack

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Options struct {
	ProjectDir string
	OutputDir  string
	AppName    string
	Version    string
	Publisher  string
	ExePath    string
}

func Package(opts Options) error {
	if opts.OutputDir == "" {
		opts.OutputDir = "dist"
	}
	if opts.Version == "" {
		opts.Version = "1.0.0.0"
	}
	if opts.Publisher == "" {
		opts.Publisher = "CN=Developer"
	}

	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	stagingDir := filepath.Join(opts.OutputDir, "msix_staging")
	os.RemoveAll(stagingDir)
	if err := os.MkdirAll(stagingDir, 0755); err != nil {
		return fmt.Errorf("failed to create staging directory: %w", err)
	}
	defer os.RemoveAll(stagingDir)

	if opts.ExePath == "" {
		opts.ExePath = filepath.Join(opts.OutputDir, opts.AppName+".exe")
	}

	destExe := filepath.Join(stagingDir, opts.AppName+".exe")
	if err := copyFile(opts.ExePath, destExe); err != nil {
		return fmt.Errorf("failed to copy executable: %w", err)
	}

	manifestPath := filepath.Join(stagingDir, "AppxManifest.xml")
	if err := generateManifest(manifestPath, opts); err != nil {
		return fmt.Errorf("failed to generate manifest: %w", err)
	}

	msixPath := filepath.Join(opts.OutputDir, opts.AppName+"_"+opts.Version+".msix")

	fmt.Printf("Creating MSIX package...\n")
	fmt.Printf("  Source: %s\n", stagingDir)
	fmt.Printf("  Output: %s\n", msixPath)

	if err := createMSIX(stagingDir, msixPath); err != nil {
		fmt.Println("  MakeAppx not available, creating unsigned package...")
		msixPath = filepath.Join(opts.OutputDir, opts.AppName+"_"+opts.Version+".zip")
		if err := createZip(stagingDir, msixPath); err != nil {
			return fmt.Errorf("failed to create zip: %w", err)
		}
		fmt.Printf("✓ Package created: %s (unsigned, rename to .msix to use)\n", msixPath)
	} else {
		fmt.Printf("✓ MSIX package created: %s\n", msixPath)
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Chmod(0755)
}

func generateManifest(path string, opts Options) error {
	// Escape XML special characters in user-provided values
	appName := xmlEscape(opts.AppName)
	publisher := xmlEscape(opts.Publisher)
	version := xmlEscape(opts.Version)

	manifest := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10"
         xmlns:mp="http://schemas.microsoft.com/appx/2014/phone/manifest"
         xmlns:uap="http://schemas.microsoft.com/appx/manifest/uap/windows10"
         IgnorableNamespaces="uap mp">
  <Identity Name="%s" Publisher="%s" Version="%s" ProcessorArchitecture="x64" />
  <mp:PhoneIdentity PhoneProductId="00000000-0000-0000-0000-000000000000" PhonePublisherId="00000000-0000-0000-0000-000000000000"/>
  <Properties>
    <DisplayName>%s</DisplayName>
    <PublisherDisplayName>%s</PublisherDisplayName>
    <Logo>Assets\StoreLogo.png</Logo>
  </Properties>
  <Dependencies>
    <TargetDeviceFamily Name="Windows.Desktop" MinVersion="10.0.17763.0" MaxVersionTested="10.0.19041.0" />
  </Dependencies>
  <Resources>
    <Resource Language="x-generate"/>
  </Resources>
  <Applications>
    <Application Id="App" Executable="%s.exe" EntryPoint="Windows.FullTrustApplication">
      <uap:VisualElements DisplayName="%s" Description="A Narcissus application"
                          Square150x150Logo="Assets\Square150x150Logo.png"
                          Square44x44Logo="Assets\Square44x44Logo.png"
                          BackgroundColor="transparent" />
    </Application>
  </Applications>
  <Capabilities>
    <Capability Name="internetClient" />
  </Capabilities>
</Package>
`, appName, publisher, version, appName, publisher, appName, appName)
	return os.WriteFile(path, []byte(manifest), 0644)
}

// xmlEscape escapes XML special characters
func xmlEscape(s string) string {
	var buf bytes.Buffer
	xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

func createMSIX(stagingDir, outputPath string) error {
	cmd := exec.Command("MakeAppx.exe", "pack", "/d", stagingDir, "/p", outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createZip(sourceDir, outputPath string) error {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			file.Close()
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			file.Close()
			return err
		}

		_, err = io.Copy(writer, file)
		file.Close()
		return err
	})
}
