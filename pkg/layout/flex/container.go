package flex

// Container represents a flex container
type Container struct {
	Items []Item

	// Flex properties
	Direction    Direction
	Wrap         Wrap
	Justify      Justify
	Align        Align
	AlignContent Align

	// Gap properties
	RowGap    float32
	ColumnGap float32
}

// NewContainer creates a new flex container
func NewContainer() *Container {
	return &Container{
		Items:        make([]Item, 0),
		Direction:    DirectionRow,
		Wrap:         WrapNoWrap,
		Justify:      JustifyFlexStart,
		Align:        AlignStretch,
		AlignContent: AlignStretch,
	}
}

// AddItem adds an item to the container
func (c *Container) AddItem(item Item) {
	c.Items = append(c.Items, item)
}

// Measure performs the measure pass
func (c *Container) Measure(constraints Constraint) Size {
	// Simplified implementation - full algorithm in Task 3.3
	var totalWidth, totalHeight float32

	for i := range c.Items {
		item := &c.Items[i]
		if item.Node != nil {
			size := item.Node.Measure(constraints)
			item.SetMeasuredSize(size.Width, size.Height)
			totalWidth += size.Width
			totalHeight = max(totalHeight, size.Height)
		}
	}

	// Apply constraints
	if constraints.MaxWidth > 0 && totalWidth > constraints.MaxWidth {
		totalWidth = constraints.MaxWidth
	}
	if constraints.MaxHeight > 0 && totalHeight > constraints.MaxHeight {
		totalHeight = constraints.MaxHeight
	}

	return Size{Width: totalWidth, Height: totalHeight}
}

// Layout performs the complete flexbox layout pass
func (c *Container) Layout(x, y, width, height float32) {
	// Determine main and cross axis dimensions
	var mainSize, crossSize float32
	var isRow bool

	if c.Direction == DirectionRow || c.Direction == DirectionRowReverse {
		mainSize = width
		crossSize = height
		isRow = true
	} else {
		mainSize = height
		crossSize = width
		isRow = false
	}

	// Calculate lines
	lines := calculateLines(c, mainSize)

	// Calculate cross positions for each line
	totalCrossSize := float32(0)
	for i := range lines {
		line := &lines[i]
		for _, item := range line.Items {
			var itemCrossSize float32
			if isRow {
				_, itemCrossSize = item.GetMeasuredSize()
			} else {
				itemCrossSize, _ = item.GetMeasuredSize()
			}
			if itemCrossSize > line.CrossSize {
				line.CrossSize = itemCrossSize
			}
		}
		totalCrossSize += line.CrossSize
	}

	// Calculate remaining cross space
	remainingCrossSpace := crossSize - totalCrossSize

	// Position lines along cross axis
	var crossPos float32
	if c.AlignContent == AlignCenter {
		crossPos = remainingCrossSpace / 2
	} else if c.AlignContent == AlignFlexEnd {
		crossPos = remainingCrossSpace
	}

	// Layout each line
	for _, line := range lines {
		// Calculate total main size of items in this line
		var totalItemsMainSize float32
		for _, item := range line.Items {
			var itemMainSize float32
			if isRow {
				itemMainSize, _ = item.GetMeasuredSize()
			} else {
				_, itemMainSize = item.GetMeasuredSize()
			}
			totalItemsMainSize += itemMainSize
		}

		// Distribute extra space using flex-grow
		availableMainSpace := mainSize - totalItemsMainSize
		distributeExtraSpace(&line, availableMainSpace, true, c.Direction)

		// Calculate justify offset and gap
		justifyOffset := calculateJustifyOffset(c.Justify, availableMainSpace, len(line.Items))
		justifyGap := calculateJustifyGap(c.Justify, availableMainSpace, len(line.Items))

		// Check if main axis is reversed
		isMainReverse := c.Direction == DirectionRowReverse || c.Direction == DirectionColumnReverse

		// Compute total main occupied space using post-distribution sizes and gaps
		var totalMainOccupied float32
		for idx, item := range line.Items {
			var itemFinalMainSize float32
			if isRow {
				// Use post-distribution width
				if item.Width != 0 {
					itemFinalMainSize = item.Width
				} else {
					w, _ := item.GetMeasuredSize()
					itemFinalMainSize = w
				}
			} else {
				// Use post-distribution height
				if item.Height != 0 {
					itemFinalMainSize = item.Height
				} else {
					_, h := item.GetMeasuredSize()
					itemFinalMainSize = h
				}
			}
			totalMainOccupied += itemFinalMainSize
			// Add gap between items (not after the last item)
			if idx < len(line.Items)-1 {
				if isRow {
					totalMainOccupied += c.ColumnGap
				} else {
					totalMainOccupied += c.RowGap
				}
			}
		}

		// Position items along main axis
		var mainPos float32
		var itemsToIterate []*Item
		if isMainReverse {
			// For reverse direction, start from the end and iterate in reverse
			mainPos = mainSize - justifyOffset - totalMainOccupied - (float32(len(line.Items)-1) * justifyGap)
			// Reverse iteration order
			itemsToIterate = make([]*Item, len(line.Items))
			for i, item := range line.Items {
				itemsToIterate[len(line.Items)-1-i] = item
			}
		} else {
			mainPos = justifyOffset
			itemsToIterate = line.Items
		}

		for _, item := range itemsToIterate {
			mw, mh := item.GetMeasuredSize()

			// Use measured size as base, apply flex-grow adjustments
			var w, h float32
			if isRow {
				w = item.Width // This may have been adjusted by distributeExtraSpace
				h = mh
				if w == 0 {
					w = mw
				}
			} else {
				w = mw
				h = item.Height // This may have been adjusted by distributeExtraSpace
				if h == 0 {
					h = mh
				}
			}

			// Calculate cross-axis position
			align := c.Align
			if item.AlignSelf != AlignAuto {
				align = item.AlignSelf
			}
			itemCrossOffset := calculateAlignment(align, h, line.CrossSize)

			// Set final position and size
			if isRow {
				item.Left = x + mainPos
				item.Top = y + crossPos + itemCrossOffset
				item.Width = w
				item.Height = h
			} else {
				item.Left = x + crossPos + itemCrossOffset
				item.Top = y + mainPos
				item.Width = w
				item.Height = h
			}

			// Advance main position
			if isRow {
				mainPos += w + justifyGap + c.ColumnGap
			} else {
				mainPos += h + justifyGap + c.RowGap
			}
		}

		if isRow {
			crossPos += line.CrossSize + c.RowGap
		} else {
			crossPos += line.CrossSize + c.ColumnGap
		}
	}
}
