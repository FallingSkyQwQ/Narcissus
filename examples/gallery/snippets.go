package main

// Source excerpts shown at the bottom of each page. They are intentionally
// short and mirror the real code on the page rather than reproducing it
// verbatim.

const snippetOverview = `app := ui.NewApp("My App", ui.WithSize(800, 600))
app.SetContent(
    ui.NewContainer().
        Direction(flex.DirectionColumn).
        Gap(12).
        Padding(ui.UniformInsets(24)).
        Add(
            ui.NewText("Hello").FontSize(20),
            ui.NewButton().Text("OK").OnClick(save),
        ),
)
app.Run()`

const snippetLayout = `ui.NewContainer().
    Direction(flex.DirectionRow).
    Justify(flex.JustifySpaceBetween).
    Align(flex.AlignCenter).
    Gap(8).
    Add(box1, box2, box3)

// Flexible children fill the leftover space.
ui.NewContainer().FlexGrow(1)`

const snippetTypography = `ui.NewText("Display").
    FontSize(28).
    FontWeight(ui.FontWeightBold).
    TextColor(ui.GetThemeColor(ui.TokenTextPrimary))`

const snippetButtons = `ui.NewButton().
    Text("Primary").
    BackgroundColor(ui.GetThemeColor(ui.TokenPrimary)).
    TextColor(ui.ColorWhite).
    OnClick(func() { /* ... */ })

// Disable interaction and skip it in the focus order.
btn.SetEnabled(false)`

const snippetInputs = `ui.NewTextInput().
    Placeholder("Your name").
    OnChange(func(v string) { echo.Text(v) })

ui.NewSlider().
    Range(0, 100).
    Value(30).
    OnChange(func(v float32) { bar.Value(v) })`

const snippetSelection = `ui.NewCheckbox().Label("Enable").OnChange(onCheck)
ui.NewSwitch().Label("Wi-Fi").OnToggle(onToggle)

ui.NewRadioButton().Group("size").Label("Small").Select()
ui.NewComboBox().Items("Small", "Medium").OnChange(onPick)

ui.NewList().Items("Inbox", "Sent").OnSelect(onRow)`

const snippetFeedback = `ui.NewProgressBar().Range(0, 100).Value(65).ShowText(true)
ui.NewProgressBar().Indeterminate(true)

ui.NewToast("Saved").Severity(ui.ToastSuccess).Show()`

const snippetOverlays = `ui.NewMenu("File").
    AddItem("New").
    AddSeparator().
    OnSelect(func(i int, item ui.MenuItem) { /* ... */ })

ui.NewDialog("Delete file", "This cannot be undone.").
    Buttons("Cancel", "Delete").
    OnResult(func(choice int) { /* ... */ }).
    Show()`

const snippetTabs = `ui.NewTabs().
    AddTab("Overview", overviewView).
    AddTab("Details", detailsView).
    OnChange(func(i int) { /* ... */ })`

const snippetReactive = `count := reactive.NewSignal(0)
doubled := reactive.NewComputed(func() int { return count.Get() * 2 })

// Re-runs automatically when a dependency changes.
reactive.NewEffect(func() {
    label.Text(fmt.Sprintf("%d", doubled.Get()))
})

// Coalesce several updates into one run.
reactive.Batch(func() {
    a.Set(a.Get() + 1)
    b.Set(b.Get() + 1)
})`

const snippetAccessibility = `email := ui.NewTextInput().Placeholder("you@example.com")
email.SetAccessibleName("Email address")
email.SetAccessibleDescription("Used to sign in")

// The toolkit-neutral tree, in document order.
node := ui.AccessibilityTree(root)`
