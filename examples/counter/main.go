// Command counter shows reactive state: a Signal drives an Effect that updates
// the label automatically whenever the value changes.
package main

import (
	"fmt"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

func main() {
	count := reactive.NewSignal(0)

	label := ui.NewText("Count: 0").FontSize(22)

	// The effect re-runs on every count change; no manual subscription needed.
	reactive.NewEffect(func() {
		label.Text(fmt.Sprintf("Count: %d", count.Get()))
	})

	increment := ui.NewButton().
		Text("+").
		OnClick(func() { count.Set(count.Get() + 1) })

	decrement := ui.NewButton().
		Text("-").
		BackgroundColor(ui.ColorSecondary).
		OnClick(func() { count.Set(count.Get() - 1) })

	reset := ui.NewButton().
		Text("Reset").
		BackgroundColor(ui.ColorGray).
		OnClick(func() { count.Set(0) })

	controls := ui.NewContainer().
		Direction(flex.DirectionRow).
		Gap(8).
		Add(increment, decrement, reset)

	root := ui.NewContainer().
		Direction(flex.DirectionColumn).
		Justify(flex.JustifyCenter).
		Align(flex.AlignCenter).
		Gap(20).
		Padding(ui.UniformInsets(24)).
		Add(label, controls)

	app := ui.NewApp("Reactive Counter", ui.WithSize(420, 260))
	app.SetContent(root)
	if err := app.Run(); err != nil {
		fmt.Println("run:", err)
	}
}
