# Narcissus
Build native UI apps with Go.

Narcissus provides a declarative widget model, a flexbox layout engine,
reactive state and a toolkit-neutral accessibility tree. The platform-specific
rendering lives behind a single backend interface, so the same application code
runs on multiple toolkits:

| Platform | Toolkit | Status |
|----------|---------|--------|
| Linux    | GTK4    | Supported |
| Windows  | WinUI 3 | Implemented, compile-verified, not yet run on Windows |

## Quick Start

### Installation

```bash
go install github.com/FallingSkyQwQ/Narcissus/cmd/narc@latest
```

### Create a New Project

```bash
narc init myapp
cd myapp
```

### Development

Run with hot reload:

```bash
narc run
narc run --target linux
```

### Build

Build for the host platform, or pick a target explicitly:

```bash
narc build
narc build --target linux
narc build --target windows
narc build --target linux --arch arm64
narc build --target windows --no-windows-app-sdk
```

`--target` selects `GOOS` (the output extension and linker flags follow
automatically). When omitted, the host platform is used.

For a Windows target, `narc build` also prepares the Windows App SDK runtime
files a WinUI 3 app needs beside its executable: the bootstrapper DLL and a
`resources.pri`. The bootstrapper is downloaded from NuGet; `resources.pri`
needs `makepri` and the installed framework package, so it is only generated
when building on Windows (on another host the remaining command is printed).
Pass `--no-windows-app-sdk` to skip both steps.

## CLI Commands

| Command | Description |
|---------|-------------|
| `narc init [name]` | Create a new project |
| `narc run` | Run with hot reload |
| `narc build` | Build production executable |
| `narc package` | Create MSIX package (Windows) |
| `narc doctor` | Check environment |
| `narc version` | Show version |

## Examples

The [`examples/`](./examples) directory contains runnable applications covering
the widget set, the flexbox layout engine, reactive state and modal overlays:

```bash
go run ./examples/hello
go run ./examples/counter
go run ./examples/widgets
```

## Requirements

Common:

- Go 1.27+
- Git

Linux (GTK4 backend):

- GTK4 development files, e.g. `sudo apt install libgtk-4-dev`
- `pkg-config`

Windows (WinUI 3 backend):

- Windows 10/11, `amd64` (the bindings are generated for `windows/amd64` only)
- Windows App SDK 2.x runtime installed
- The bootstrapper DLL and `resources.pri` beside the executable (produced by
  `narc build`, or by hand)
- `UPX` (optional, for compression)
- `Air` (optional, for hot reload)

Run `narc doctor` to verify that your environment has everything the backends
need.

### Windows (WinUI 3) notes

The Windows backend renders through
[`go-bindings-windowsappsdk`](https://github.com/deploymenttheory/go-bindings-windowsappsdk),
a pure-Go WinUI 3 projection (no cgo, no .NET SDK, no XAML compiler). Startup
order — lock the UI thread, enter a single-threaded apartment, bootstrap the
Windows App SDK, `Application.Start`, then create the window — is handled by
that library, mirroring how the GTK backend defers all widget creation to
the application's `activate` handler.

Layout follows the same rule as Linux: a container maps to a `Canvas` and every
control is placed at the coordinates the framework's flex engine computed, so
windows lay out identically across platforms.

**Status:** the code is compile-verified (`GOOS=windows go build ./...`) but has
not been executed on Windows. Treat the runtime behaviour as untested.

## Architecture

Application code only imports `github.com/FallingSkyQwQ/Narcissus/pkg/ui`.
Each widget exposes its state through a toolkit-neutral property snapshot, and
the platform backend (`backend_gtk_linux.go`, `backend_windows.go`, ...)
translates the laid-out widget tree into native controls. The framework's flex
engine remains authoritative for layout on every platform, so positions are
identical across backends.

## Documentation

- [Design Document](./DESIGN.md)
- [Roadmap](./ROADMAP.md)
