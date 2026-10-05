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

// clampAxis applies non-zero min/max bounds to a value.
func clampAxis(value, minValue, maxValue float32) float32 {
	if minValue > 0 && value < minValue {
		value = minValue
	}
	if maxValue > 0 && value > maxValue {
		value = maxValue
	}
	return value
}

// reverseItems returns the items in reverse order without mutating the line.
func reverseItems(items []*Item) []*Item {
	out := make([]*Item, len(items))
	for i, item := range items {
		out[len(items)-1-i] = item
	}
	return out
}

// Measure performs the intrinsic size pass. It refreshes each item's measured
// size from its node (when one is attached), breaks the items into lines and
// returns the container's natural size, clamped to the constraints.
func (c *Container) Measure(constraints Constraint) Size {
	for i := range c.Items {
		item := &c.Items[i]
		if item.Node != nil {
			size := item.Node.Measure(constraints)
			item.SetMeasuredSize(size.Width, size.Height)
		}
	}

	row := isRow(c.Direction)
	lines := calculateLines(c, mainExtent(row, constraints))

	var maxMain, totalCross float32
	for li := range lines {
		line := &lines[li]

		var lineMain float32
		for i, item := range line.Items {
			if i > 0 {
				lineMain += mainGap(c, row)
			}
			lineMain += itemBaseMain(item, row) + mainOuterMargin(item, row)
		}
		line.MainSize = lineMain
		if lineMain > maxMain {
			maxMain = lineMain
		}

		var lineCross float32
		for _, item := range line.Items {
			if v := itemBaseCross(item, row) + crossOuterMargin(item, row); v > lineCross {
				lineCross = v
			}
		}
		line.CrossSize = lineCross
	}

	for i := range lines {
		if i > 0 {
			totalCross += crossGap(c, row)
		}
		totalCross += lines[i].CrossSize
	}

	var width, height float32
	if row {
		width, height = maxMain, totalCross
	} else {
		width, height = totalCross, maxMain
	}

	width = clampAxis(width, constraints.MinWidth, constraints.MaxWidth)
	height = clampAxis(height, constraints.MinHeight, constraints.MaxHeight)

	return Size{Width: width, Height: height}
}

// Layout performs the complete flexbox layout pass. It resolves flexible
// lengths, distributes leftover space through justify-content and
// align-content, applies per-item cross-axis alignment (including stretch of
// items without a definite cross size) and writes the final rectangles back to
// the items.
func (c *Container) Layout(x, y, width, height float32) {
	row := isRow(c.Direction)

	var mainSize, crossSize float32
	if row {
		mainSize, crossSize = width, height
	} else {
		mainSize, crossSize = height, width
	}

	lines := calculateLines(c, mainSize)

	// Each line's cross size is driven by its tallest/widest outer item.
	lineCross := make([]float32, len(lines))
	for li := range lines {
		var lc float32
		for _, item := range lines[li].Items {
			if v := itemBaseCross(item, row) + crossOuterMargin(item, row); v > lc {
				lc = v
			}
		}
		lineCross[li] = lc
	}

	var totalCross float32
	for i := range lineCross {
		if i > 0 {
			totalCross += crossGap(c, row)
		}
		totalCross += lineCross[i]
	}
	freeCross := crossSize - totalCross
	if freeCross < 0 {
		freeCross = 0
	}

	// align-content: how leftover cross space is distributed across lines.
	var crossOffset, extraLineGap float32
	switch c.AlignContent {
	case AlignFlexEnd:
		crossOffset = freeCross
	case AlignCenter:
		crossOffset = freeCross / 2
	case AlignStretch:
		if n := len(lines); n > 0 {
			each := freeCross / float32(n)
			for i := range lineCross {
				lineCross[i] += each
			}
		}
	case AlignSpaceBetween:
		if n := len(lines); n > 1 {
			extraLineGap = freeCross / float32(n-1)
		}
	case AlignSpaceAround:
		if n := len(lines); n > 0 {
			extraLineGap = freeCross / float32(n)
			crossOffset = extraLineGap / 2
		}
	case AlignSpaceEvenly:
		if n := len(lines); n > 0 {
			extraLineGap = freeCross / float32(n+1)
			crossOffset = extraLineGap
		}
	}

	crossPos := crossOffset
	for li := range lines {
		line := &lines[li]
		lc := lineCross[li]

		// Resolve flex-grow / flex-shrink over the line's free main space.
		usedMain := float32(0)
		for i, item := range line.Items {
			if i > 0 {
				usedMain += mainGap(c, row)
			}
			usedMain += itemBaseMain(item, row) + mainOuterMargin(item, row)
		}
		resolveFlexibleLengths(line, mainSize-usedMain, row)

		// Usage after resolution drives the justify-content leftover space.
		usedMain = 0
		for i, item := range line.Items {
			if i > 0 {
				usedMain += mainGap(c, row)
			}
			usedMain += resolvedMain(item, row) + mainOuterMargin(item, row)
		}
		remaining := mainSize - usedMain
		justifyOffset := calculateJustifyOffset(c.Justify, remaining, len(line.Items))
		justifyGap := calculateJustifyGap(c.Justify, remaining, len(line.Items))

		items := line.Items
		if isReverse(c.Direction) {
			items = reverseItems(items)
		}

		cursor := justifyOffset
		for _, item := range items {
			itemMain := resolvedMain(item, row)

			align := c.Align
			if item.AlignSelf != AlignAuto {
				align = item.AlignSelf
			}
			baseCross := itemBaseCross(item, row)
			itemCross := baseCross
			crossStart := float32(0)
			switch align {
			case AlignFlexEnd:
				crossStart = lc - (itemCross + crossOuterMargin(item, row))
			case AlignCenter:
				crossStart = (lc - (itemCross + crossOuterMargin(item, row))) / 2
			case AlignStretch:
				// Items whose cross size is indefinite (auto) fill the line.
				if baseCross == 0 || itemCrossAuto(item, row) {
					itemCross = lc - crossOuterMargin(item, row)
					if itemCross < 0 {
						itemCross = 0
					}
				}
			}
			crossStart += crossPos + leadingCrossMargin(item, row)

			mainFromStart := cursor + leadingMainMargin(item, row)
			if isReverse(c.Direction) {
				mainFromStart = mainSize - mainFromStart - itemMain
			}

			var dx, dy float32
			if row {
				dx, dy = mainFromStart, crossStart
				item.Height = itemCross
			} else {
				dx, dy = crossStart, mainFromStart
				item.Width = itemCross
			}
			item.Left = x + dx
			item.Top = y + dy

			cursor += leadingMainMargin(item, row) + itemMain + trailingMainMargin(item, row) +
				justifyGap + mainGap(c, row)
		}

		crossPos += lc + crossGap(c, row) + extraLineGap
	}
}

// resolvedMain reads an item's final main size after resolveFlexibleLengths.
func resolvedMain(item *Item, row bool) float32 {
	if row {
		return item.Width
	}
	return item.Height
}
