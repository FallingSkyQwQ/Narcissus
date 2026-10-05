# Examples

Runnable Narcissus applications. Each directory is a `main` package; run one
from the repository root with:

```bash
go run ./examples/<name>
```

On Linux the GTK4 backend needs `libgtk-4-dev` and `pkg-config` (see the
[README](../README.md)).

| Example | Demonstrates |
|---------|--------------|
| [`gallery`](./gallery) | The full showcase: a sidebar-navigated gallery of every component with live examples, source snippets, a light/dark theme switch, a reactive-state demo and an accessibility-tree inspector. |
| [`hello`](./hello) | The smallest app: a `Container`, `Text` and `Button`. |
| [`counter`](./counter) | Reactive state with `reactive.Signal` and an auto-tracking `reactive.Effect`. |
| [`layout`](./layout) | The flexbox engine: direction, justify, align, wrap, gaps and `FlexGrow`. |
| [`widgets`](./widgets) | A single scrolling page of every built-in control. |
| [`forms`](./forms) | Inputs plus `Dialog` and `Toast` presented through the backend. |

### Showcase

`examples/gallery` is the recommended starting point:

```bash
go run ./examples/gallery
```

It pairs a navigation sidebar with a content pane and covers, per component,
multiple variants and their interactive behaviour. Each page ends with the Go
source that produced it. The app starts from the system's light/dark preference
(see `ui.SystemPrefersDark`), the sidebar's theme button restyles it live, and
the Accessibility page dumps the live `ui.AccessibilityTree`.

## Widget reference

The gallery touches the full control set:

- Layout: `Container` (flex), `ScrollView`, `Tabs`
- Text: `Text`, `TextInput`
- Actions: `Button`, `Menu`
- Selection: `Checkbox`, `Switch`, `RadioButton`, `ComboBox`, `List`, `Slider`
- Feedback: `ProgressBar`, `Dialog`, `Toast`

Every widget is configured with chainable methods, for example:

```go
ui.NewContainer().
    Direction(flex.DirectionColumn).
    Gap(8).
    Add(
        ui.NewText("Hello").FontSize(20),
        ui.NewButton().Text("OK").OnClick(func() { /* ... */ }),
    )
```

Widgets expose a toolkit-neutral accessibility view through
`ui.AccessibilityTree`; set names and roles with
`SetAccessibleName` / `SetAccessibleDescription` / `SetAccessibleRole`.
