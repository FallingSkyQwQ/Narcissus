// Command gallery is a showcase application for the Narcissus UI library.
//
// It is organised as a sidebar of component categories plus a content pane:
// selecting a category shows live, interactive examples together with the Go
// source that produced them. The sidebar also carries a light/dark theme
// switch that restyles the whole application.
//
// Run it from the repository root:
//
//	go run ./examples/gallery
package main

import (
	"fmt"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

// page describes one entry in the gallery: its sidebar title, a one-line
// summary, the builder that produces its content, and the source snippet
// shown at the bottom of the page.
type page struct {
	title    string
	subtitle string
	build    func(*gallery, page) ui.Widget
	snippet  string
}

// gallery holds the mutable application state: the pages, the current
// selection and the theme.
type gallery struct {
	app     *ui.App
	root    ui.Widget
	pages   []page
	pageUI  []ui.Widget
	current int
	dark    bool

	// building guards callbacks that fire while the tree is being assembled.
	building bool
	// gen makes radio-group names unique across rebuilds.
	gen int
}

func main() {
	g := &gallery{}
	g.pages = g.definePages()
	g.app = ui.NewApp("Narcissus Gallery", ui.WithSize(1180, 800))
	g.root = g.build()
	g.app.SetContent(g.root)

	// The OS preference is only readable once the toolkit is initialized, so
	// adopt it from the ready callback and rebuild if it differs.
	ui.OnReady(func() {
		if dark := ui.SystemPrefersDark(); dark != g.dark {
			g.setDark(dark)
		}
	})

	if err := g.app.Run(); err != nil {
		fmt.Println("gallery:", err)
	}
}

// definePages wires the categories shown in the sidebar.
func (g *gallery) definePages() []page {
	return []page{
		{"Overview", "What Narcissus is and how this gallery is organised.", (*gallery).pageOverview, snippetOverview},
		{"Layout", "The flexbox engine: direction, justify, gaps and growth.", (*gallery).pageLayout, snippetLayout},
		{"Typography", "Text sizes, weights and the theme's colour roles.", (*gallery).pageTypography, snippetTypography},
		{"Buttons", "Accent-coloured actions and disabled states.", (*gallery).pageButtons, snippetButtons},
		{"Inputs", "Text fields, sliders and progress reporting.", (*gallery).pageInputs, snippetInputs},
		{"Selection", "Checkbox, switch, radio group, combo box and list.", (*gallery).pageSelection, snippetSelection},
		{"Feedback", "Progress bars and transient toasts.", (*gallery).pageFeedback, snippetFeedback},
		{"Overlays", "Menus and modal dialogs.", (*gallery).pageOverlays, snippetOverlays},
		{"Tabs", "Composed tabs that work on every backend.", (*gallery).pageTabs, snippetTabs},
		{"Reactive", "Signals, computed values, effects and batching.", (*gallery).pageReactive, snippetReactive},
		{"Accessibility", "The toolkit-neutral accessibility tree.", (*gallery).pageAccessibility, snippetAccessibility},
	}
}

// build assembles the whole widget tree for the current theme and selection.
// It is called once at start-up and again on every theme change.
func (g *gallery) build() ui.Widget {
	g.building = true
	g.gen++
	defer func() { g.building = false }()

	g.pageUI = g.pageUI[:0]
	titles := make([]string, 0, len(g.pages))
	for _, p := range g.pages {
		titles = append(titles, p.title)
		g.pageUI = append(g.pageUI, p.build(g, p))
	}
	if g.current < 0 || g.current >= len(g.pageUI) {
		g.current = 0
	}
	for i, w := range g.pageUI {
		w.SetVisible(i == g.current)
	}

	pages := ui.NewContainer().Direction(flex.DirectionColumn).Add(g.pageUI...)
	content := ui.NewScrollView().
		SetContent(pages).
		Style(ui.NewStyle().WithFlexGrow(1))

	nav := ui.NewList().Items(titles...).Style(navListStyle())
	nav.OnSelect(func(index int, _ string) { g.show(index) })

	brand := ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(2).
		Add(
			ui.NewText("Narcissus").FontSize(18).FontWeight(ui.FontWeightBold).TextColor(token(ui.TokenTextPrimary)),
			ui.NewText("Widget gallery").FontSize(12).TextColor(token(ui.TokenTextSecondary)),
		)

	// Capture this tree's generation: a control left over from a previous,
	// torn-down tree must not drive the theme. A button is used rather than a
	// switch because the toolkit's switch drops an on-state set before the
	// widget is realized.
	gen := g.gen
	themeButton := ghostButton(themeButtonLabel(g.dark)).OnClick(func() {
		if gen != g.gen {
			return
		}
		g.setDark(!g.dark)
	})

	themeRow := ui.NewContainer().
		Direction(flex.DirectionRow).
		Gap(8).
		Align(flex.AlignCenter).
		Add(
			ui.NewText("Appearance").FontSize(13).TextColor(token(ui.TokenTextSecondary)),
			themeButton,
		)

	sidebar := ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(20).
		Style(ui.NewStyle().
			WithBackgroundColor(token(ui.TokenSurface)).
			WithPadding(ui.Insets{Top: 20, Right: 16, Bottom: 20, Left: 16})).
		Width(248).
		FlexShrink(0).
		Add(brand, nav, themeRow)

	root := ui.NewContainer().
		Direction(flex.DirectionRow).
		Style(ui.NewStyle().WithBackgroundColor(token(ui.TokenBackground))).
		Add(sidebar, content)

	nav.Select(g.current)
	return root
}

// show selects a page, flipping visibility and relaying the tree out. It is a
// no-op while the tree is still being built.
func (g *gallery) show(index int) {
	if index < 0 || index >= len(g.pageUI) {
		return
	}
	changed := index != g.current
	g.current = index
	for i, w := range g.pageUI {
		w.SetVisible(i == index)
	}
	if !g.building && changed {
		g.refresh()
	}
}

// refresh re-lays out the existing tree. The widget instances are stable, so
// the backend reuses the native controls and only repositions them.
func (g *gallery) refresh() {
	if g.app != nil && g.root != nil {
		g.app.SetContent(g.root)
	}
}

// setDark switches the active theme and rebuilds the tree so every widget
// picks up the new colour roles.
func (g *gallery) setDark(dark bool) {
	if g.dark == dark {
		return
	}
	g.dark = dark
	if dark {
		ui.SetTheme(ui.NewDarkTheme())
	} else {
		ui.SetTheme(ui.NewLightTheme())
	}
	g.rebuild()
}

// rebuild replaces the tree with a freshly themed one. The old tree is torn
// down on an idle callback so the widget that triggered the change is not
// destroyed from inside its own event handler.
func (g *gallery) rebuild() {
	old := g.root
	apply := func() {
		teardown(old)
		g.root = g.build()
		g.app.SetContent(g.root)
	}
	if backend := ui.CurrentBackend(); backend != nil {
		if err := backend.Post(apply); err == nil {
			return
		}
	}
	apply()
}

// nativable is implemented by every widget through the embedded BaseWidget; it
// exposes the native control handle for teardown.
type nativable interface {
	NativeControl() ui.NativeControl
}

// teardown destroys the native controls of a widget tree bottom-up. Children
// are released before their parent so a container is never finalized while it
// still owns live children.
func teardown(w ui.Widget) {
	if w == nil {
		return
	}
	for _, child := range w.GetChildren() {
		teardown(child)
	}
	if n, ok := w.(nativable); ok {
		if c := n.NativeControl(); c != nil {
			c.Destroy()
		}
	}
}
