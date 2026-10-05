# Design

Narcissus builds native desktop UI from Go. This document describes how the
pieces fit together. For status and planned work see [ROADMAP.md](./ROADMAP.md).

## Goals

- A declarative, chainable widget model that reads well in Go.
- A layout engine owned by the framework, so windows lay out identically on
  every platform. Backends translate already-computed rectangles rather than
  delegating layout to the toolkit.
- A single backend interface so application code imports only `pkg/ui`.

## Layers

```
pkg/ui           widget model, styling/theme, events, backend interface
pkg/layout/flex  the flexbox layout engine
pkg/reactive     signals and UI-thread dispatch
pkg/bridge       low-level COM/WinRT helpers for the WinUI 3 backend
cmd/narc         the CLI (init / run / build / package / doctor / version)
internal/        CLI internals (build, doctor, pack, runner, scaffold)
```

`pkg/layout/flex` and `pkg/reactive` have no dependency on `pkg/ui`, so they can
be tested and reused on their own.

## Widget model

Every widget embeds `BaseWidget`, which implements the `Widget` interface:
identity, flex properties, style, parent/children, events, visibility and
enabled state. `BaseWidget` holds a `sync.RWMutex` and a cached
`flex.Item`.

Widgets are configured with chainable methods, for example:

```go
ui.NewContainer().
    Direction(flex.DirectionColumn).
    Gap(8).
    Add(
        ui.NewText("Hello").FontSize(20),
        ui.NewButton().Text("OK").OnClick(func() { /* ... */ }),
    )
```

`pkg/ui/render.go` maps each concrete widget to a toolkit-neutral
`ControlKind` and a `ControlProps` snapshot. Adding a widget means adding a
kind, a props extractor and backend construction.

## Layout pipeline

The flex engine is authoritative on every platform:

1. `layoutAndMount` (in `pkg/ui/render.go`) measures the content tree against
   the window size, then lays the root out at the full window rectangle.
2. `Container.Layout` runs the flexbox algorithm: it breaks items into lines,
   resolves flexible lengths, distributes leftover space through justify /
   align-content, and writes final rectangles back to the items.
3. Each widget applies its laid-out rectangle to its native control via
   `NativeControl.SetBounds`, converting absolute coordinates to
   parent-relative ones.

Text is measured through the `TextMeasurer` interface (`pkg/ui/textmetrics.go`).
The GTK4 backend installs a Pango-backed measurer so layout uses the same font
metrics the toolkit renders with; other platforms fall back to a documented
estimator until they provide their own.

## Backend interface

`pkg/ui/backend.go` defines the seam between the framework and a toolkit:

- `Backend` — `Init`, `CreateWindow`, `CreateControl`, `Run`, `Quit`, `Post`,
  `OnUIThread`.
- `NativeWindow` — the top-level window and its root surface.
- `NativeControl` — `SetBounds`, `SetVisible`, `SetEnabled`, `SetStyle`,
  `Update`, `AttachTo`, `Destroy`.
- `ControlKind` / `ControlProps` — the toolkit-neutral description of a control.

A backend wires toolkit events back into widgets; for example a GTK button
click becomes an `EventClick` dispatched through `Widget.HandleEvent`.

`CurrentBackend` returns the platform default installed by a build-tagged file
(`backend_gtk_linux.go`, `backend_windows.go`, `backend_unsupported.go`), and
`SetBackend` lets tests inject a fake.

## Focus and keyboard

`pkg/ui/focus.go` owns a toolkit-neutral focus manager. Widgets opt into
focus via `BaseWidget.SetFocusable` (Button, Checkbox, ComboBox, Slider and
TextInput do; Text, Image and Container do not).

- `RequestFocus(w)` moves focus, dispatching `EventBlur` to the previous
  widget and `EventFocus` to the new one, then asking the backend for native
  focus through the optional `focusRequester` interface on the control.
- `FocusNext` / `FocusPrevious` walk the focusable widgets in document order,
  wrapping around; `DispatchKey` turns an unhandled TAB into a traversal step.
- Key events from the backend enter through `DispatchKey`, which routes them to
  the focused widget and bubbles up its ancestors until one consumes them.

The mount path calls `SetFocusRoot` with the window content so traversal knows
which tree to walk. Backends report toolkit focus changes with
`notifyNativeFocus` / `notifyNativeBlur`.

## Styling and theming

`Style` (`pkg/ui/style.go`) carries layout, background, border, font, shadow
and opacity properties. `Theme` maps `ColorToken` values to colors, and a global
`ThemeManager` exposes the current theme and a subscription signal. The GTK
backend turns a `Style` into a per-widget CSS rule.

## Reactivity

`pkg/reactive` provides `Signal` (observable state), `Computed` and `Effect`.
Updates are marshaled onto the toolkit UI thread through a `Dispatcher` that
`App.Run` installs. Dependency tracking is currently manual; see the roadmap.

## Platform support

| Platform | Backend | Status |
|----------|---------|--------|
| Linux | GTK4 | Supported (requires `libgtk-4-dev`) |
| Windows | WinUI 3 | Compile-verified, not yet run on Windows |
| Other (`darwin`, ...) | unsupported | Returns `ErrBackendUnavailable` |

## CLI

`narc` scaffolds projects (`init`), runs them with hot reload (`run`), builds
executables for a chosen `GOOS`/`GOARCH` (`build`, optionally preparing Windows
App SDK runtime files), packages MSIX (`package`), checks the environment
(`doctor`) and prints its version (`version`).
