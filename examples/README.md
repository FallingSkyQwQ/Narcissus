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
| [`hello`](./hello) | The smallest app: a `Container`, `Text` and `Button`. |
| [`counter`](./counter) | Reactive state with `reactive.Signal` and an auto-tracking `reactive.Effect`. |
| [`layout`](./layout) | The flexbox engine: direction, justify, align, wrap, gaps and `FlexGrow`. |
| [`widgets`](./widgets) | A gallery of every built-in control inside a `ScrollView`. |
| [`forms`](./forms) | Inputs plus `Dialog` and `Toast` presented through the backend. |

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
