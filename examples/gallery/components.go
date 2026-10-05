package main

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

// token resolves a theme colour role against the active theme. Everything in
// the gallery is built from roles, so rebuilding after a theme change restyles
// the whole application.
func token(t ui.ColorToken) ui.Color { return ui.GetThemeColor(t) }

func heading(s string) *ui.Text {
	return ui.NewText(s).FontSize(24).FontWeight(ui.FontWeightBold).TextColor(token(ui.TokenTextPrimary))
}

func subheading(s string) *ui.Text {
	return ui.NewText(s).FontSize(16).FontWeight(ui.FontWeightSemiBold).TextColor(token(ui.TokenTextPrimary))
}

func muted(s string) *ui.Text {
	return ui.NewText(s).FontSize(13).TextColor(token(ui.TokenTextSecondary))
}

func cardStyle() *ui.Style {
	return ui.NewStyle().
		WithBackgroundColor(token(ui.TokenCardBackground)).
		WithBorder(ui.Border{Width: 1, Color: token(ui.TokenBorder), Radius: 10}).
		WithPadding(ui.Insets{Top: 20, Right: 22, Bottom: 20, Left: 22})
}

// card is a padded, bordered surface that stacks its children vertically.
func card(children ...ui.Widget) ui.Widget {
	return ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(16).
		Style(cardStyle()).
		Add(children...)
}

func row(children ...ui.Widget) *ui.Container {
	return ui.NewContainer().Direction(flex.DirectionRow).Gap(18).Align(flex.AlignCenter).Add(children...)
}

func col(children ...ui.Widget) *ui.Container {
	return ui.NewContainer().Direction(flex.DirectionColumn).Gap(16).Add(children...)
}

func statCard(value, label string) ui.Widget {
	return ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(8).
		Style(cardStyle()).
		MinWidth(150).
		Add(
			ui.NewText(value).FontSize(24).FontWeight(ui.FontWeightBold).TextColor(token(ui.TokenPrimary)),
			ui.NewText(label).FontSize(12).TextColor(token(ui.TokenTextSecondary)),
		)
}

// snippetBlock renders a Go source excerpt in a monospaced, bordered card.
func snippetBlock(code string) ui.Widget {
	return ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(8).
		Style(ui.NewStyle().
			WithBackgroundColor(token(ui.TokenSurface)).
			WithBorder(ui.Border{Width: 1, Color: token(ui.TokenBorder), Radius: 10}).
			WithPadding(ui.UniformInsets(14))).
		Add(
			ui.NewText("Source").FontSize(11).FontWeight(ui.FontWeightSemiBold).TextColor(token(ui.TokenTextSecondary)),
			ui.NewText(code).FontFamily("monospace").FontSize(12).TextColor(token(ui.TokenTextPrimary)),
		)
}

// themeButtonLabel is the sidebar theme control's caption.
func themeButtonLabel(dark bool) string {
	if dark {
		return "Dark: on"
	}
	return "Dark: off"
}

func ghostButton(label string) *ui.Button {
	return ui.NewButton().
		Text(label).
		BackgroundColor(token(ui.TokenSurface)).
		TextColor(token(ui.TokenTextPrimary))
}

// navListStyle is the sidebar navigation list.
func navListStyle() *ui.Style {
	return ui.NewStyle().
		WithFlexGrow(1).
		WithBackgroundColor(token(ui.TokenSurface)).
		WithTextColor(token(ui.TokenTextPrimary)).
		WithBorderRadius(8).
		WithPadding(ui.UniformInsets(4))
}

// selectionListStyle is an inline list used inside the Selection page.
func selectionListStyle() *ui.Style {
	return ui.NewStyle().
		WithBackgroundColor(token(ui.TokenCardBackground)).
		WithTextColor(token(ui.TokenTextPrimary)).
		WithBorder(ui.Border{Width: 1, Color: token(ui.TokenBorder), Radius: 8}).
		WithPadding(ui.UniformInsets(4))
}
