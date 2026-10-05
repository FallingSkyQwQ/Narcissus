// Command hello is the smallest Narcissus application: a window with a label
// and a button.
package main

import (
	"fmt"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

func main() {
	message := ui.NewText("Hello, Narcissus!").FontSize(20)

	greet := ui.NewButton().
		Text("Greet").
		OnClick(func() {
			message.Text("Hello from a Go callback!")
		})

	root := ui.NewContainer().
		Direction(flex.DirectionColumn).
		Justify(flex.JustifyCenter).
		Align(flex.AlignCenter).
		Gap(16).
		Padding(ui.UniformInsets(24)).
		Add(message, greet)

	app := ui.NewApp("Hello Narcissus", ui.WithSize(420, 240))
	app.SetContent(root)
	if err := app.Run(); err != nil {
		fmt.Println("run:", err)
	}
}
