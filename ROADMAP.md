# Roadmap

Narcissus is an early-stage UI framework. This roadmap tracks the work needed
to move it from a working skeleton toward a mature, dependable UI library. The
order reflects impact: each item unblocks the ones below it.

Status legend: ✅ done · 🚧 in progress · ⬜ not started

| # | Item | Status | Notes |
|---|------|--------|-------|
| 1 | Real text measurement | ✅ | `TextMeasurer` abstraction with a Pango-backed implementation for GTK4 and a documented fallback estimator. Text, Button, Checkbox and ComboBox measure through it. |
| 2 | Complete the flex layout algorithm | ✅ | Intrinsic sizing for node-less items, margins/gaps, `flex-basis`, bounded grow/shrink, `align-content`, per-item stretch and min/max size constraints. |
| 3 | Repository and documentation cleanup | 🚧 | `DESIGN.md` and `ROADMAP.md` are now tracked (previously ignored yet linked from the README); local build artifacts removed. A usage example gallery is still missing. |
| 4 | Interaction layer: focus, keyboard, tab order | ⬜ | Key/mouse-enter/leave/focus events are declared but not wired in the backends. `TextInput.Focus`/`Blur` are stubs. |
| 5 | Essential widgets | ⬜ | Missing ScrollView, List/ListBox, Dialog/Modal, Menu, ProgressBar, Radio/Switch, Tabs and Toast. Only 8 controls exist today. |
| 6 | Accessibility and HiDPI | ⬜ | No accessibility tree or screen-reader support, and no device-pixel-ratio handling. |
| 7 | Reactive auto dependency tracking | ⬜ | `Computed` needs a manual `Recompute` and `Effect` runs once; there is no automatic dependency collection or batching. |

## Beyond the list

These are not required for a "usable" library but are part of reaching
maturity:

- macOS/Android/Web backends (`darwin` currently resolves to the unsupported
  backend).
- Custom drawing / canvas, animation and transitions.
- Scrolling, virtualization and z-ordering.
- API stability guarantees, a changelog and a component documentation site.
- Headless/Golden rendering tests and backend mocks.
