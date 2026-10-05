// Command widgets is a gallery of the built-in controls, laid out inside a
// ScrollView so every widget is reachable in one window.
package main

import (
	"fmt"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

func section(title string, children ...ui.Widget) ui.Widget {
	return ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(8).
		Add(append([]ui.Widget{
			ui.NewText(title).FontSize(16).FontWeight(ui.FontWeightBold),
		}, children...)...)
}

func main() {
	status := reactive.NewSignal("Ready")

	statusLabel := ui.NewText("Ready").TextColor(ui.GetThemeColor(ui.TokenTextSecondary))
	reactive.NewEffect(func() { statusLabel.Text(status.Get()) })

	button := ui.NewButton().
		Text("Click me").
		OnClick(func() { status.Set("Button clicked") })

	checkbox := ui.NewCheckbox().
		Label("Enable feature").
		OnChange(func(checked bool) {
			if checked {
				status.Set("Feature enabled")
			} else {
				status.Set("Feature disabled")
			}
		})

	toggle := ui.NewSwitch().
		Label("Dark mode").
		OnToggle(func(on bool) { status.Set(fmt.Sprintf("Dark mode: %v", on)) })

	radioGroup := ui.NewContainer().
		Direction(flex.DirectionRow).
		Gap(12).
		Add(
			ui.NewRadioButton().Group("size").Label("Small").Select().OnSelect(func() { status.Set("Small") }),
			ui.NewRadioButton().Group("size").Label("Medium").OnSelect(func() { status.Set("Medium") }),
			ui.NewRadioButton().Group("size").Label("Large").OnSelect(func() { status.Set("Large") }),
		)

	combo := ui.NewComboBox().
		Items("Alpha", "Beta", "Gamma").
		Placeholder("Pick a value").
		OnChange(func(index int, value string) { status.Set("Picked " + value) })

	progress := ui.NewProgressBar().Range(0, 100).Value(40).ShowText(true)

	slider := ui.NewSlider().
		Range(0, 100).
		Value(40).
		OnChange(func(value float32) { progress.Value(value) })

	input := ui.NewTextInput().
		Placeholder("Type a note").
		OnChange(func(value string) { status.Set("Typing: " + value) })

	list := ui.NewList().
		Items("Inbox", "Drafts", "Sent", "Trash").
		OnSelect(func(index int, value string) { status.Set("Selected " + value) })

	menu := ui.NewMenu("Actions").
		AddItem("New").
		AddItem("Open").
		AddSeparator().
		AddItem("Quit").
		OnSelect(func(index int, item ui.MenuItem) { status.Set("Menu: " + item.Label) })

	tabs := ui.NewTabs().
		AddTab("Overview", ui.NewText("The overview tab.")).
		AddTab("Details", ui.NewText("The details tab.")).
		OnChange(func(index int) { status.Set(fmt.Sprintf("Tab %d", index)) })

	content := ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(24).
		Padding(ui.UniformInsets(20)).
		Add(
			ui.NewText("Widget gallery").FontSize(24).FontWeight(ui.FontWeightBold),
			statusLabel,
			section("Buttons & toggles", button, checkbox, toggle),
			section("Radio group", radioGroup),
			section("Selection", combo, menu),
			section("Inputs", input, slider, progress),
			section("List", list),
			section("Tabs", tabs),
		)

	scroll := ui.NewScrollView().SetContent(content)

	app := ui.NewApp("Widget Gallery", ui.WithSize(600, 640))
	app.SetContent(scroll)
	if err := app.Run(); err != nil {
		fmt.Println("run:", err)
	}
}
