# Roadmap

Narcissus is an early-stage UI framework. This roadmap tracks the work needed
to move it from a working skeleton toward a mature, dependable UI library. The
order reflects impact: each item unblocks the ones below it.

Status legend: ✅ done · 🚧 in progress · ⬜ not started

| # | Item | Status | Notes |
|---|------|--------|-------|
| 1 | Real text measurement | ✅ | `TextMeasurer` abstraction with a Pango-backed implementation for GTK4 and a documented fallback estimator. Text, Button, Checkbox and ComboBox measure through it. |
| 2 | Complete the flex layout algorithm | ✅ | Intrinsic sizing for node-less items, margins/gaps, `flex-basis`, bounded grow/shrink, `align-content`, per-item stretch and min/max size constraints. |
| 3 | Repository and documentation cleanup | ✅ | `DESIGN.md` and `ROADMAP.md` are tracked (previously ignored yet linked from the README); local build artifacts removed. A runnable usage gallery now lives under [`examples/`](./examples) with an index and a widget reference. |
| 4 | Interaction layer: focus, keyboard, tab order | ✅ | Framework focus manager with `RequestFocus`, TAB traversal and focus/blur bubbling; key events routed to the focused widget and bubbled up. GTK4 wired (key controller + focus controllers). WinUI 3 focus/key wiring is deferred until that backend is verified. |
| 5 | Essential widgets | ✅ | ScrollView, List/ListBox, Dialog/Modal, Menu, ProgressBar, Radio/Switch, Tabs and Toast added; 16 control kinds in total. Dialog and Toast present through the backend's optional `Presenter` capability. |
| 6 | Accessibility and HiDPI | ✅ | Widgets carry an accessible name, description and role, exposed as a toolkit-neutral `AccessibilityTree`; GTK publishes accessible label/description and WinUI publishes `AutomationProperties`. The device pixel ratio is plumbed from the GTK scale factor and WinUI rasterization scale through `DevicePixelRatio` / `ScaleToDevice`. |
| 7 | Reactive auto dependency tracking | ✅ | Dependency collection is automatic via a per-goroutine observer context: `Computed` re-evaluates lazily on read (eagerly while subscribed), `Effect` re-runs when a tracked source changes, and `Batch` coalesces a burst of updates into a single run. |

## Beyond the list

These are not required for a "usable" library but are part of reaching
maturity:

- macOS/Android/Web backends (`darwin` currently resolves to the unsupported
  backend).
- Custom drawing / canvas, animation and transitions.
- Scrolling, virtualization and z-ordering.
- API stability guarantees, a changelog and a component documentation site.
- Headless/Golden rendering tests and backend mocks.
