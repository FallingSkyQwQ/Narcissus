// Command forms combines inputs with overlays: a Dialog and a Toast are shown
// through the backend's presenter capability rather than the layout tree.
package main

import (
	"fmt"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

func main() {
	name := reactive.NewSignal("")
	email := reactive.NewSignal("")
	agree := reactive.NewSignal(false)

	validation := ui.NewText("").TextColor(ui.ColorError)

	reactive.NewEffect(func() {
		switch {
		case name.Get() == "":
			validation.Text("Enter your name")
		case email.Get() == "":
			validation.Text("Enter your email")
		case !agree.Get():
			validation.Text("Accept the terms to continue")
		default:
			validation.Text("Looks good!")
		}
	})

	form := ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(12).
		Add(
			ui.NewText("Name").FontSize(12).TextColor(ui.GetThemeColor(ui.TokenTextSecondary)),
			ui.NewTextInput().Placeholder("Jane Doe").OnChange(func(v string) { name.Set(v) }),
			ui.NewText("Email").FontSize(12).TextColor(ui.GetThemeColor(ui.TokenTextSecondary)),
			ui.NewTextInput().Placeholder("jane@example.com").OnChange(func(v string) { email.Set(v) }),
			ui.NewCheckbox().Label("I accept the terms").OnChange(func(v bool) { agree.Set(v) }),
			validation,
		)

	submit := ui.NewButton().
		Text("Submit").
		OnClick(func() {
			dialog := ui.NewDialog("Confirm", fmt.Sprintf("Create the account for %s?", name.Get())).
				Buttons("Cancel", "Create").
				OnResult(func(choice int) {
					if choice != 1 {
						return
					}
					ui.NewToast("Account created").
						Severity(ui.ToastSuccess).
						Show()
				})
			_ = dialog.Show()
		})

	root := ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(20).
		Padding(ui.UniformInsets(24)).
		Add(
			ui.NewText("Create account").FontSize(22).FontWeight(ui.FontWeightBold),
			form,
			submit,
		)

	app := ui.NewApp("Form Example", ui.WithSize(460, 420))
	app.SetContent(root)
	if err := app.Run(); err != nil {
		fmt.Println("run:", err)
	}
}
