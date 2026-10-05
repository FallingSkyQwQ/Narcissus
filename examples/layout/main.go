// Command layout demonstrates the flexbox engine: direction, justify, align,
// wrapping, gaps and flexible growth.
package main

import (
	"fmt"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

// card builds a fixed-size, coloured box with a caption.
func card(title string, color ui.Color) ui.Widget {
	return ui.NewContainer().
		Direction(flex.DirectionColumn).
		Justify(flex.JustifyCenter).
		Align(flex.AlignCenter).
		Style(ui.NewStyle().WithBackgroundColor(color).WithBorderRadius(8)).
		Size(96, 72).
		Add(ui.NewText(title).TextColor(ui.ColorWhite))
}

func main() {
	row := ui.NewContainer().
		Direction(flex.DirectionRow).
		Justify(flex.JustifySpaceBetween).
		Align(flex.AlignCenter).
		Gap(12).
		Add(
			card("One", ui.ColorPrimary),
			card("Two", ui.ColorSecondary),
			card("Three", ui.ColorAccent),
		)

	growing := ui.NewContainer().
		Direction(flex.DirectionColumn).
		Style(ui.NewStyle().WithBackgroundColor(ui.ColorLightSurface).WithBorderRadius(8)).
		Padding(ui.UniformInsets(12)).
		FlexGrow(1).
		Add(ui.NewText("This region grows to fill the leftover height."))

	wrapped := ui.NewContainer().
		Direction(flex.DirectionRow).
		Wrap(flex.WrapWrap).
		Gap(8).
		Add(
			card("A", ui.ColorGray), card("B", ui.ColorGray), card("C", ui.ColorGray),
			card("D", ui.ColorGray), card("E", ui.ColorGray), card("F", ui.ColorGray),
		)

	root := ui.NewContainer().
		Direction(flex.DirectionColumn).
		Gap(16).
		Padding(ui.UniformInsets(20)).
		Add(
			ui.NewText("Layout gallery").FontSize(22).FontWeight(ui.FontWeightBold),
			row,
			growing,
			wrapped,
		)

	app := ui.NewApp("Layout Gallery", ui.WithSize(560, 460))
	app.SetContent(root)
	if err := app.Run(); err != nil {
		fmt.Println("run:", err)
	}
}
