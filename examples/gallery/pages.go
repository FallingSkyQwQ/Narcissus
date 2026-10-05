package main

import (
	"fmt"
	"strings"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

// scaffold wraps a page's content with its title, summary and source snippet.
func (g *gallery) scaffold(p page, content ...ui.Widget) ui.Widget {
	items := []ui.Widget{
		heading(p.title),
		muted(p.subtitle),
	}
	items = append(items, content...)
	items = append(items, snippetBlock(p.snippet))
	return ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(28).
		Padding(ui.Insets{Top: 28, Right: 56, Bottom: 48, Left: 56}).
		Add(items...)
}

func (g *gallery) pageOverview(p page) ui.Widget {
	intro := ui.NewText("Narcissus builds native desktop UI from Go: a declarative widget model, a flexbox layout engine, reactive state and a toolkit-neutral accessibility tree.").
		FontSize(14).
		TextColor(token(ui.TokenTextPrimary))

	howto := card(
		subheading("How to read this gallery"),
		muted("Pick a category on the left. Each page shows live, interactive examples and, at the bottom, the Go source that produced them."),
		muted("Use the theme button at the bottom of the sidebar to restyle the whole application; it starts from your system's light or dark preference."),
	)

	return g.scaffold(p,
		intro,
		row(
			statCard("16", "control kinds"),
			statCard("2", "widget toolkits"),
			statCard("100%", "Go, no XAML compiler"),
		),
		howto,
	)
}

func (g *gallery) pageLayout(p page) ui.Widget {
	demo := ui.NewContainer().
		Direction(flex.DirectionRow).
		Gap(8).
		Style(ui.NewStyle().
			WithBackgroundColor(token(ui.TokenSurface)).
			WithBorderRadius(10).
			WithPadding(ui.UniformInsets(12))).
		MinHeight(140)

	for i, colour := range []ui.Color{ui.ColorPrimary, ui.ColorSecondary, ui.ColorAccent, ui.ColorSuccess} {
		demo.Add(ui.NewContainer().
			Direction(flex.DirectionColumn).
			Justify(flex.JustifyCenter).
			Align(flex.AlignCenter).
			Style(ui.NewStyle().WithBackgroundColor(colour).WithBorderRadius(8)).
			Size(64, 64).
			Add(ui.NewText(fmt.Sprintf("%d", i+1)).TextColor(ui.ColorWhite)))
	}

	state := muted("direction: row   justify: flex-start")

	directionButtons := row(
		ghostButton("Row").OnClick(func() {
			demo.Direction(flex.DirectionRow)
			state.Text("direction: row   justify: flex-start")
			g.refresh()
		}),
		ghostButton("Column").OnClick(func() {
			demo.Direction(flex.DirectionColumn)
			state.Text("direction: column")
			g.refresh()
		}),
	)

	justifies := []flex.Justify{
		flex.JustifyFlexStart,
		flex.JustifyCenter,
		flex.JustifyFlexEnd,
		flex.JustifySpaceBetween,
		flex.JustifySpaceAround,
	}
	justifyIndex := 0
	justifyButton := ghostButton("Cycle justify").OnClick(func() {
		justifyIndex = (justifyIndex + 1) % len(justifies)
		demo.Justify(justifies[justifyIndex])
		state.Text("justify: " + justifyName(justifies[justifyIndex]))
		g.refresh()
	})

	// The gap slider mutates the demo container and relayouts on every tick.
	gapLabel := muted("gap: 8")
	gapSlider := ui.NewSlider().
		Range(0, 32).
		Value(8).
		Width(220).
		OnChange(func(v float32) {
			demo.Gap(v)
			gapLabel.Text(fmt.Sprintf("gap: %.0f", v))
			g.refresh()
		})

	return g.scaffold(p,
		muted("A row of boxes laid out by the framework's flexbox engine. The controls below mutate the live container and relayout it."),
		demo,
		directionButtons,
		row(justifyButton, state),
		row(gapSlider, gapLabel),
	)
}

func (g *gallery) pageTypography(p page) ui.Widget {
	return g.scaffold(p,
		ui.NewText("Display").FontSize(28).FontWeight(ui.FontWeightBold).TextColor(token(ui.TokenTextPrimary)),
		ui.NewText("Heading").FontSize(20).FontWeight(ui.FontWeightSemiBold).TextColor(token(ui.TokenTextPrimary)),
		ui.NewText("Body text uses the theme's primary text colour.").FontSize(14).TextColor(token(ui.TokenTextPrimary)),
		muted("Secondary caption, used for hints and metadata."),
		row(
			ui.NewText("Primary").TextColor(token(ui.TokenPrimary)),
			ui.NewText("Success").TextColor(token(ui.TokenSuccess)),
			ui.NewText("Warning").TextColor(token(ui.TokenWarning)),
			ui.NewText("Error").TextColor(token(ui.TokenError)),
		),
	)
}

func (g *gallery) pageButtons(p page) ui.Widget {
	clicks := 0
	counter := muted("clicks: 0")
	count := func() {
		clicks++
		counter.Text(fmt.Sprintf("clicks: %d", clicks))
	}

	primary := ui.NewButton().Text("Primary").
		BackgroundColor(token(ui.TokenPrimary)).TextColor(ui.ColorWhite).OnClick(count)
	secondary := ui.NewButton().Text("Secondary").
		BackgroundColor(token(ui.TokenSecondary)).TextColor(ui.ColorWhite).OnClick(count)
	danger := ui.NewButton().Text("Danger").
		BackgroundColor(token(ui.TokenError)).TextColor(ui.ColorWhite).OnClick(count)

	disabled := ui.NewButton().Text("Disabled")
	disabled.SetEnabled(false)

	return g.scaffold(p,
		muted("Buttons take the accent colours of the current theme. A disabled button drops out of the focus order."),
		row(primary, secondary, danger, disabled),
		counter,
	)
}

func (g *gallery) pageInputs(p page) ui.Widget {
	echo := muted("value: ")
	single := ui.NewTextInput().
		Placeholder("Your name").
		Width(260).
		BackgroundColor(token(ui.TokenCardBackground)).
		TextColor(token(ui.TokenTextPrimary)).
		OnChange(func(v string) { echo.Text("value: " + v) })
	multiline := ui.NewTextInput().
		InputType(ui.TextInputTypeMultiline).
		Placeholder("Notes…").
		Width(320).
		Height(80)
	readonly := ui.NewTextInput().Value("read-only value").ReadOnly(true).Width(240)

	bar := ui.NewProgressBar().Range(0, 100).Value(30).ShowText(true).Width(260)
	valueLabel := muted("30")
	slider := ui.NewSlider().Range(0, 100).Value(30).Width(220).OnChange(func(v float32) {
		bar.Value(v)
		valueLabel.Text(fmt.Sprintf("%.0f", v))
	})

	return g.scaffold(p,
		card(
			subheading("Text fields"),
			muted("Single line, multi-line and read-only inputs."),
			single,
			echo,
			multiline,
			readonly,
		),
		card(
			subheading("Slider drives the progress bar"),
			row(slider, valueLabel),
			bar,
		),
	)
}

func (g *gallery) pageSelection(p page) ui.Widget {
	note := muted("make a selection")
	status := func(s string) { note.Text(s) }

	check := ui.NewCheckbox().Label("Enable telemetry").
		OnChange(func(v bool) { status(fmt.Sprintf("checkbox: %v", v)) })
	sw := ui.NewSwitch().
		OnToggle(func(v bool) { status(fmt.Sprintf("switch: %v", v)) })
	switchRow := row(ui.NewText("Wi-Fi").TextColor(token(ui.TokenTextPrimary)), sw)

	group := fmt.Sprintf("density-%d", g.gen)
	compact := ui.NewRadioButton().Group(group).Label("Compact").
		OnSelect(func() { status("density: compact") }).Select()
	cozy := ui.NewRadioButton().Group(group).Label("Cozy").
		OnSelect(func() { status("density: cozy") })
	spacious := ui.NewRadioButton().Group(group).Label("Spacious").
		OnSelect(func() { status("density: spacious") })

	combo := ui.NewComboBox().Items("Small", "Medium", "Large").Placeholder("Pick a size").
		Width(220).
		BackgroundColor(token(ui.TokenCardBackground)).TextColor(token(ui.TokenTextPrimary)).
		OnChange(func(_ int, v string) { status("combo: " + v) })

	list := ui.NewList().Items("Inbox", "Drafts", "Sent", "Trash").
		Style(selectionListStyle()).
		Width(200).
		Height(150).
		OnSelect(func(_ int, v string) { status("list: " + v) })

	return g.scaffold(p,
		card(check, switchRow),
		card(subheading("Radio group"), row(compact, cozy, spacious)),
		card(subheading("Combo box"), combo),
		card(subheading("List"), list),
		note,
	)
}

func (g *gallery) pageFeedback(p page) ui.Widget {
	status := muted("idle")

	bar := ui.NewProgressBar().Range(0, 100).Value(65).ShowText(true).Width(280)
	indeterminate := ui.NewProgressBar().Indeterminate(true).Width(280)
	indeterminateSwitch := ui.NewSwitch().
		OnToggle(func(on bool) {
			indeterminate.Indeterminate(on)
			status.Text(fmt.Sprintf("indeterminate: %v", on))
		})
	indeterminateRow := row(ui.NewText("Indeterminate").TextColor(token(ui.TokenTextPrimary)), indeterminateSwitch)

	success := ui.NewButton().Text("Success toast").
		BackgroundColor(token(ui.TokenSuccess)).TextColor(ui.ColorWhite).
		OnClick(func() {
			ui.NewToast("Saved successfully").Severity(ui.ToastSuccess).Show()
			status.Text("showed: success toast")
		})
	failure := ui.NewButton().Text("Error toast").
		BackgroundColor(token(ui.TokenError)).TextColor(ui.ColorWhite).
		OnClick(func() {
			ui.NewToast("Something went wrong").Severity(ui.ToastError).Show()
			status.Text("showed: error toast")
		})

	return g.scaffold(p,
		card(subheading("Progress"), bar, indeterminate, indeterminateRow),
		card(subheading("Toasts"), muted("Presented by the backend, not the layout tree."), row(success, failure)),
		status,
	)
}

func (g *gallery) pageOverlays(p page) ui.Widget {
	status := muted("no action yet")

	menu := ui.NewMenu("File").
		AddItem("New").
		AddItem("Open").
		AddSeparator().
		AddItem("Exit").
		OnSelect(func(_ int, item ui.MenuItem) { status.Text("menu: " + item.Label) })

	openDialog := ui.NewButton().Text("Open dialog").OnClick(func() {
		ui.NewDialog("Delete file", "This action cannot be undone.").
			Buttons("Cancel", "Delete").
			OnResult(func(choice int) {
				if choice == 1 {
					status.Text("dialog: deleted")
				} else {
					status.Text("dialog: cancelled")
				}
			}).
			Show()
	})

	return g.scaffold(p,
		card(subheading("Menu"), menu),
		card(subheading("Modal dialog"), muted("Dialogs are presented by the backend's Presenter capability."), openDialog),
		status,
	)
}

func (g *gallery) pageTabs(p page) ui.Widget {
	status := muted("tab 0")
	tabs := ui.NewTabs().
		AddTab("Overview", ui.NewText("The overview tab.").TextColor(token(ui.TokenTextPrimary))).
		AddTab("Details", ui.NewText("The details tab.").TextColor(token(ui.TokenTextPrimary))).
		AddTab("Settings", ui.NewText("The settings tab.").TextColor(token(ui.TokenTextPrimary))).
		OnChange(func(i int) { status.Text(fmt.Sprintf("tab %d", i)) })

	return g.scaffold(p,
		muted("Tabs are composed from a header row of buttons plus a content area, so they need no backend support."),
		tabs,
		status,
	)
}

func (g *gallery) pageReactive(p page) ui.Widget {
	count := reactive.NewSignal(0)
	counterText := ui.NewText("0").FontSize(40).FontWeight(ui.FontWeightBold).TextColor(token(ui.TokenPrimary))
	doubled := reactive.NewComputed(func() int { return count.Get() * 2 })
	doubledText := muted("doubled: 0")

	// Effects re-run automatically whenever a signal they read changes.
	reactive.NewEffect(func() { counterText.Text(fmt.Sprintf("%d", count.Get())) })
	reactive.NewEffect(func() { doubledText.Text(fmt.Sprintf("doubled: %d", doubled.Get())) })

	increment := ui.NewButton().Text("+").OnClick(func() { count.Set(count.Get() + 1) })
	decrement := ui.NewButton().Text("−").
		BackgroundColor(token(ui.TokenSecondary)).TextColor(ui.ColorWhite).
		OnClick(func() { count.Set(count.Get() - 1) })
	reset := ui.NewButton().Text("Reset").
		BackgroundColor(token(ui.TokenSurface)).TextColor(token(ui.TokenTextPrimary)).
		OnClick(func() { count.Set(0) })

	// Batching demo: an effect that reads two signals.
	a := reactive.NewSignal(0)
	b := reactive.NewSignal(0)
	runs := 0
	runText := muted("effect runs: 0")
	reactive.NewEffect(func() {
		_, _ = a.Get(), b.Get()
		runs++
		runText.Text(fmt.Sprintf("effect runs: %d", runs))
	})
	batched := ui.NewButton().Text("Batched update").OnClick(func() {
		reactive.Batch(func() {
			a.Set(a.Get() + 1)
			b.Set(b.Get() + 1)
		})
	})
	unbatched := ui.NewButton().Text("Two plain updates").OnClick(func() {
		a.Set(a.Get() + 1)
		b.Set(b.Get() + 1)
	})

	return g.scaffold(p,
		card(
			subheading("Signal + Computed + Effect"),
			counterText,
			doubledText,
			row(increment, decrement, reset),
		),
		card(
			subheading("Batching"),
			muted("A batched update runs the effect once; two plain updates run it twice."),
			runText,
			row(batched, unbatched),
		),
	)
}

func (g *gallery) pageAccessibility(p page) ui.Widget {
	email := ui.NewTextInput().Placeholder("you@example.com").Width(260)
	email.SetAccessibleName("Email address")
	email.SetAccessibleDescription("Used to sign in to your account")

	subscribe := ui.NewCheckbox().Label("Subscribe to updates")
	subscribe.SetAccessibleName("Subscribe to updates")

	submit := ui.NewButton().Text("Sign in")
	submit.SetAccessibleRole(ui.RoleButton)
	submit.SetAccessibleDescription("Submits the sign-in form")

	form := ui.NewContainer().Direction(flex.DirectionColumn).Gap(10).Add(email, subscribe, submit)

	output := ui.NewText("Press “Inspect” to dump the accessibility tree.").
		FontFamily("monospace").FontSize(12).TextColor(token(ui.TokenTextSecondary))
	inspect := ui.NewButton().Text("Inspect").OnClick(func() {
		output.Text(dumpAccessibility(ui.AccessibilityTree(form), 0))
	})

	return g.scaffold(p,
		ui.NewText("Every widget carries an accessible name, description and role. Names default to the widget's own content; override them with SetAccessibleName.").
			FontSize(13).TextColor(token(ui.TokenTextPrimary)),
		card(subheading("Form"), form),
		row(inspect),
		card(output),
	)
}

// dumpAccessibility renders an accessibility tree as indented text.
func dumpAccessibility(node *ui.AccessibleNode, depth int) string {
	if node == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(strings.Repeat("  ", depth))
	b.WriteString("- ")
	b.WriteString(node.Role.String())
	if node.Name != "" {
		b.WriteString(": ")
		b.WriteString(node.Name)
	}
	if node.Description != "" {
		b.WriteString("  (")
		b.WriteString(node.Description)
		b.WriteString(")")
	}
	b.WriteByte('\n')
	for _, child := range node.Children {
		b.WriteString(dumpAccessibility(child, depth+1))
	}
	return b.String()
}

// justifyName renders a flex justify value for the layout page readout.
func justifyName(j flex.Justify) string {
	switch j {
	case flex.JustifyCenter:
		return "center"
	case flex.JustifyFlexEnd:
		return "flex-end"
	case flex.JustifySpaceBetween:
		return "space-between"
	case flex.JustifySpaceAround:
		return "space-around"
	case flex.JustifySpaceEvenly:
		return "space-evenly"
	default:
		return "flex-start"
	}
}
