package flex

// Line represents a line of flex items in wrapped layout
type Line struct {
	Items     []*Item
	MainSize  float32
	CrossSize float32
}

// calculateLines splits items into lines based on wrap mode
func calculateLines(container *Container, availableMain float32) []Line {
	if container.Wrap == WrapNoWrap {
		line := Line{Items: make([]*Item, len(container.Items))}
		for i := range container.Items {
			line.Items[i] = &container.Items[i]
		}
		return []Line{line}
	}

	var lines []Line
	currentLine := Line{Items: make([]*Item, 0)}
	currentMainSize := float32(0)

	for i := range container.Items {
		item := &container.Items[i]
		w, h := item.GetMeasuredSize()

		var itemMainSize float32
		if container.Direction == DirectionRow || container.Direction == DirectionRowReverse {
			itemMainSize = w
		} else {
			itemMainSize = h
		}

		if container.Wrap != WrapNoWrap && currentMainSize+itemMainSize > availableMain && len(currentLine.Items) > 0 {
			lines = append(lines, currentLine)
			currentLine = Line{Items: make([]*Item, 0)}
			currentMainSize = 0
		}

		currentLine.Items = append(currentLine.Items, item)
		currentMainSize += itemMainSize
	}

	if len(currentLine.Items) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

// distributeExtraSpace handles flex-grow and flex-shrink
func distributeExtraSpace(line *Line, availableSpace float32, isMainAxis bool, direction Direction) {
	if !isMainAxis {
		return
	}

	// Handle flex-shrink when available space is negative (overflow)
	if availableSpace < 0 {
		// Track which items are frozen (clamped to 0)
		frozen := make([]bool, len(line.Items))
		remainingAvailableSpace := availableSpace

		// Iterative redistribution: repeat until space satisfied or all items frozen
		for remainingAvailableSpace < 0 {
			// Compute totalFlexShrink only over unfrozen items
			totalFlexShrink := float32(0)
			unfrozenCount := 0
			for i, item := range line.Items {
				if !frozen[i] {
					totalFlexShrink += item.FlexShrink
					unfrozenCount++
				}
			}

			// If no unfrozen items or no shrink capacity, stop
			if unfrozenCount == 0 || totalFlexShrink == 0 {
				break
			}

			shrinkPerUnit := remainingAvailableSpace / totalFlexShrink
			anyItemClamped := false

			for i, item := range line.Items {
				if frozen[i] {
					continue
				}

				shrink := shrinkPerUnit * item.FlexShrink
				if direction == DirectionRow || direction == DirectionRowReverse {
					measuredBase, _ := item.GetMeasuredSize()
					if item.Width == 0 {
						item.Width = measuredBase
					}
					newWidth := item.Width + shrink
					if newWidth < 0 {
						// Clamp to 0 and mark as frozen
						item.Width = 0
						frozen[i] = true
						anyItemClamped = true
					} else {
						item.Width = newWidth
					}
				} else {
					_, measuredBase := item.GetMeasuredSize()
					if item.Height == 0 {
						item.Height = measuredBase
					}
					newHeight := item.Height + shrink
					if newHeight < 0 {
						// Clamp to 0 and mark as frozen
						item.Height = 0
						frozen[i] = true
						anyItemClamped = true
					} else {
						item.Height = newHeight
					}
				}
			}

			// Recompute remaining available space
			totalItemsSize := float32(0)
			for _, item := range line.Items {
				if direction == DirectionRow || direction == DirectionRowReverse {
					totalItemsSize += item.Width
				} else {
					totalItemsSize += item.Height
				}
			}

			// Compute new remaining space relative to original availableSpace target
			// availableSpace is negative, so we need to check if we've distributed enough shrink
			if direction == DirectionRow || direction == DirectionRowReverse {
				var totalMeasured float32
				for _, item := range line.Items {
					w, _ := item.GetMeasuredSize()
					totalMeasured += w
				}
				remainingAvailableSpace = availableSpace - (totalItemsSize - totalMeasured)
			} else {
				var totalMeasured float32
				for _, item := range line.Items {
					_, h := item.GetMeasuredSize()
					totalMeasured += h
				}
				remainingAvailableSpace = availableSpace - (totalItemsSize - totalMeasured)
			}

			// If no items were clamped this iteration and we still have negative space,
			// we've distributed as much as possible
			if !anyItemClamped {
				break
			}
		}
		return
	}

	// Handle flex-grow when available space is positive
	if availableSpace > 0 {
		totalFlexGrow := float32(0)
		for _, item := range line.Items {
			totalFlexGrow += item.FlexGrow
		}

		if totalFlexGrow > 0 {
			spacePerGrow := availableSpace / totalFlexGrow
			for _, item := range line.Items {
				extra := spacePerGrow * item.FlexGrow
				if direction == DirectionRow || direction == DirectionRowReverse {
					measuredBase, _ := item.GetMeasuredSize()
					if item.Width == 0 {
						item.Width = measuredBase
					}
					item.Width = item.Width + extra
				} else {
					_, measuredBase := item.GetMeasuredSize()
					if item.Height == 0 {
						item.Height = measuredBase
					}
					item.Height = item.Height + extra
				}
			}
		}
	}
}

// calculateJustifyOffset calculates the starting offset based on justify content
func calculateJustifyOffset(justify Justify, availableSpace float32, itemCount int) float32 {
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
	if itemCount <= 1 {
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
