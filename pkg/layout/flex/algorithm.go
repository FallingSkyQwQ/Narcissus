package flex

// Line represents a line of flex items in wrapped layout
type Line struct {
	Items     []*Item
	MainSize  float32
	CrossSize float32
}

// isRow reports whether direction lays items out horizontally.
func isRow(direction Direction) bool {
	return direction == DirectionRow || direction == DirectionRowReverse
}

// isReverse reports whether direction lays items out from the far end.
func isReverse(direction Direction) bool {
	return direction == DirectionRowReverse || direction == DirectionColumnReverse
}

// mainGap returns the physical gap between items along the main axis.
func mainGap(container *Container, row bool) float32 {
	if row {
		return container.ColumnGap
	}
	return container.RowGap
}

// crossGap returns the physical gap between lines along the cross axis.
func crossGap(container *Container, row bool) float32 {
	if row {
		return container.RowGap
	}
	return container.ColumnGap
}

// mainExtent returns the definite main-axis bound of a constraint, or 0 when
// the main axis is unbounded.
func mainExtent(row bool, c Constraint) float32 {
	if row {
		return c.MaxWidth
	}
	return c.MaxHeight
}

// itemBaseMain resolves an item's main-axis base size: flex-basis wins over the
// measured size, matching the CSS flex-basis:auto fallback order.
func itemBaseMain(item *Item, row bool) float32 {
	if item.FlexBasis > 0 {
		return item.FlexBasis
	}
	if row {
		w, _ := item.GetMeasuredSize()
		if w == 0 {
			w = item.Width
		}
		return w
	}
	_, h := item.GetMeasuredSize()
	if h == 0 {
		h = item.Height
	}
	return h
}

// itemBaseCross resolves an item's cross-axis base size.
func itemBaseCross(item *Item, row bool) float32 {
	if row {
		_, h := item.GetMeasuredSize()
		if h == 0 {
			h = item.Height
		}
		return h
	}
	w, _ := item.GetMeasuredSize()
	if w == 0 {
		w = item.Width
	}
	return w
}

// mainOuterMargin returns the sum of the main-axis margins of an item.
func mainOuterMargin(item *Item, row bool) float32 {
	if row {
		return item.MarginLeft + item.MarginRight
	}
	return item.MarginTop + item.MarginBottom
}

// crossOuterMargin returns the sum of the cross-axis margins of an item.
func crossOuterMargin(item *Item, row bool) float32 {
	if row {
		return item.MarginTop + item.MarginBottom
	}
	return item.MarginLeft + item.MarginRight
}

// leadingMainMargin returns the main-axis margin before an item along the axis.
func leadingMainMargin(item *Item, row bool) float32 {
	if row {
		return item.MarginLeft
	}
	return item.MarginTop
}

// trailingMainMargin returns the main-axis margin after an item along the axis.
func trailingMainMargin(item *Item, row bool) float32 {
	if row {
		return item.MarginRight
	}
	return item.MarginBottom
}

// leadingCrossMargin returns the cross-axis margin before an item.
func leadingCrossMargin(item *Item, row bool) float32 {
	if row {
		return item.MarginTop
	}
	return item.MarginLeft
}

func itemMainMin(item *Item, row bool) float32 {
	if row {
		return item.MinWidth
	}
	return item.MinHeight
}

func itemMainMax(item *Item, row bool) float32 {
	if row {
		return item.MaxWidth
	}
	return item.MaxHeight
}

// calculateLines splits items into lines based on wrap mode, accounting for
// main-axis margins and gaps so a line does not overflow the available space.
func calculateLines(container *Container, availableMain float32) []Line {
	row := isRow(container.Direction)
	gap := mainGap(container, row)

	if container.Wrap == WrapNoWrap {
		line := Line{Items: make([]*Item, len(container.Items))}
		for i := range container.Items {
			line.Items[i] = &container.Items[i]
		}
		return []Line{line}
	}

	var lines []Line
	currentLine := Line{Items: make([]*Item, 0)}
	var currentMain float32

	for i := range container.Items {
		item := &container.Items[i]
		itemMain := itemBaseMain(item, row) + mainOuterMargin(item, row)

		add := itemMain
		if len(currentLine.Items) > 0 {
			add += gap
		}

		if len(currentLine.Items) > 0 && availableMain > 0 && currentMain+add > availableMain {
			lines = append(lines, currentLine)
			currentLine = Line{Items: make([]*Item, 0)}
			currentMain = 0
			add = itemMain
		}

		currentLine.Items = append(currentLine.Items, item)
		currentMain += add
	}

	if len(currentLine.Items) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

// resolveFlexibleLengths distributes free main-axis space among the items of a
// line via flex-grow (positive space) or flex-shrink (negative space), clamping
// each item to its min/max main size and re-running while any item is frozen.
// The resolved sizes are written back to the items' main axis.
func resolveFlexibleLengths(line *Line, free float32, row bool) {
	n := len(line.Items)
	if n == 0 {
		return
	}

	sizes := make([]float32, n)
	for i, item := range line.Items {
		sizes[i] = itemBaseMain(item, row)
	}

	if free > 0 {
		frozen := make([]bool, n)
		remaining := free
		for remaining > 0.0001 {
			var totalGrow float32
			for i, item := range line.Items {
				if !frozen[i] && item.FlexGrow > 0 {
					totalGrow += item.FlexGrow
				}
			}
			if totalGrow <= 0 {
				break
			}
			perGrow := remaining / totalGrow
			consumed := float32(0)
			for i, item := range line.Items {
				if frozen[i] || item.FlexGrow <= 0 {
					continue
				}
				delta := perGrow * item.FlexGrow
				if maxSize := itemMainMax(item, row); maxSize > 0 && sizes[i]+delta >= maxSize {
					consumed += maxSize - sizes[i]
					sizes[i] = maxSize
					frozen[i] = true
					continue
				}
				sizes[i] += delta
				consumed += delta
			}
			remaining -= consumed
			if consumed <= 0 {
				break
			}
		}
	} else if free < 0 {
		frozen := make([]bool, n)
		remaining := -free
		for remaining > 0.0001 {
			var totalScaled float32
			for i, item := range line.Items {
				if !frozen[i] {
					totalScaled += item.FlexShrink * sizes[i]
				}
			}
			if totalScaled <= 0 {
				break
			}
			consumed := float32(0)
			for i, item := range line.Items {
				if frozen[i] {
					continue
				}
				scaled := item.FlexShrink * sizes[i]
				if scaled <= 0 {
					continue
				}
				delta := remaining * scaled / totalScaled
				floor := itemMainMin(item, row)
				if floor < 0 {
					floor = 0
				}
				if sizes[i]-delta <= floor {
					consumed += sizes[i] - floor
					sizes[i] = floor
					frozen[i] = true
					continue
				}
				sizes[i] -= delta
				consumed += delta
			}
			remaining -= consumed
			if consumed <= 0 {
				break
			}
		}
	}

	for i, item := range line.Items {
		if row {
			item.Width = sizes[i]
		} else {
			item.Height = sizes[i]
		}
	}
}

// distributeExtraSpace handles flex-grow and flex-shrink for a line. It is kept
// as the public-ish entry point used by tests; it delegates to the bounds-aware
// resolver.
func distributeExtraSpace(line *Line, availableSpace float32, isMainAxis bool, direction Direction) {
	if !isMainAxis {
		return
	}
	resolveFlexibleLengths(line, availableSpace, isRow(direction))
}

// calculateJustifyOffset calculates the starting offset based on justify content
func calculateJustifyOffset(justify Justify, availableSpace float32, itemCount int) float32 {
	if availableSpace < 0 {
		availableSpace = 0
	}
	switch justify {
	case JustifyFlexStart:
		return 0
	case JustifyFlexEnd:
		return availableSpace
	case JustifyCenter:
		return availableSpace / 2
	case JustifySpaceBetween:
		return 0
	case JustifySpaceAround:
		if itemCount == 0 {
			return 0.0
		}
		gap := availableSpace / float32(itemCount)
		return gap / 2
	case JustifySpaceEvenly:
		if itemCount == 0 {
			return 0.0
		}
		gap := availableSpace / float32(itemCount+1)
		return gap
	default:
		return 0
	}
}

// calculateJustifyGap calculates the gap between items based on justify content
func calculateJustifyGap(justify Justify, availableSpace float32, itemCount int) float32 {
	if itemCount <= 1 || availableSpace < 0 {
		return 0
	}

	switch justify {
	case JustifySpaceBetween:
		return availableSpace / float32(itemCount-1)
	case JustifySpaceAround:
		return availableSpace / float32(itemCount)
	case JustifySpaceEvenly:
		return availableSpace / float32(itemCount+1)
	default:
		return 0
	}
}

// calculateAlignment calculates cross-axis alignment position
func calculateAlignment(align Align, itemSize float32, lineSize float32) float32 {
	switch align {
	case AlignFlexStart:
		return 0
	case AlignFlexEnd:
		return lineSize - itemSize
	case AlignCenter:
		return (lineSize - itemSize) / 2
	case AlignStretch:
		return 0
	case AlignBaseline:
		return 0
	default:
		return 0
	}
}

// max returns the maximum of two float32 values
func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
