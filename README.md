# Narcissus
Build your WinUI app with GO!

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
```

### Build

Build for production:

```bash
narc build
```

### Package

Create MSIX package:

```bash
narc package
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `narc init [name]` | Create a new project |
| `narc run` | Run with hot reload |
| `narc build` | Build production executable |
| `narc package` | Create MSIX package |
| `narc doctor` | Check environment |
| `narc version` | Show version |

## Requirements

- Windows 10/11
- Go 1.26.4+
- Git
- UPX (optional, for compression)
- Air (optional, for hot reload)

## Documentation

- [Design Document](./DESIGN.md)
- [Roadmap](./ROADMAP.md)
