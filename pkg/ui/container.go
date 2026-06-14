package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// Container 是布局容器组件
type Container struct {
	*BaseWidget

	// Flex 容器属性
	flexContainer *flex.Container
	direction     flex.Direction
	wrap          flex.Wrap
	justify       flex.Justify
	align         flex.Align
	alignContent  flex.Align
	rowGap        float32
	columnGap     float32
}

// NewContainer 创建新的容器组件
func NewContainer() *Container {
	c := &Container{
		BaseWidget:    NewBaseWidget(),
		flexContainer: flex.NewContainer(),
		direction:     flex.DirectionRow,
		wrap:          flex.WrapNoWrap,
		justify:       flex.JustifyFlexStart,
		align:         flex.AlignStretch,
		alignContent:  flex.AlignStretch,
		rowGap:        0,
		columnGap:     0,
	}

	// 设置默认样式
	c.style.BackgroundColor = ColorTransparent

	return c
}

// Direction 设置主轴方向（链式调用）
func (c *Container) Direction(direction flex.Direction) *Container {
	c.direction = direction
	c.flexContainer.Direction = direction
	return c
}

// Wrap 设置换行模式（链式调用）
func (c *Container) Wrap(wrap flex.Wrap) *Container {
	c.wrap = wrap
	c.flexContainer.Wrap = wrap
	return c
}

// Justify 设置主轴对齐方式（链式调用）
func (c *Container) Justify(justify flex.Justify) *Container {
	c.justify = justify
	c.flexContainer.Justify = justify
	return c
}

// Align 设置交叉轴对齐方式（链式调用）
func (c *Container) Align(align flex.Align) *Container {
	c.align = align
	c.flexContainer.Align = align
	return c
}

// AlignContent 设置多行对齐方式（链式调用）
func (c *Container) AlignContent(alignContent flex.Align) *Container {
	c.alignContent = alignContent
	c.flexContainer.AlignContent = alignContent
	return c
}

// RowGap 设置行间距（链式调用）
func (c *Container) RowGap(gap float32) *Container {
	c.rowGap = gap
	c.flexContainer.RowGap = gap
	return c
}

// ColumnGap 设置列间距（链式调用）
func (c *Container) ColumnGap(gap float32) *Container {
	c.columnGap = gap
	c.flexContainer.ColumnGap = gap
	return c
}

// Gap 同时设置行列间距（链式调用）
func (c *Container) Gap(gap float32) *Container {
	c.rowGap = gap
	c.columnGap = gap
	c.flexContainer.RowGap = gap
	c.flexContainer.ColumnGap = gap
	return c
}

// FlexGrow 设置 FlexGrow（链式调用）
func (c *Container) FlexGrow(grow float32) *Container {
	c.style.FlexGrow = grow
	c.flexItem.FlexGrow = grow
	return c
}

// FlexShrink 设置 FlexShrink（链式调用）
func (c *Container) FlexShrink(shrink float32) *Container {
	c.style.FlexShrink = shrink
	c.flexItem.FlexShrink = shrink
	return c
}

// FlexBasis 设置 FlexBasis（链式调用）
func (c *Container) FlexBasis(basis float32) *Container {
	c.style.FlexBasis = basis
	c.flexItem.FlexBasis = basis
	return c
}

// Margin 设置外边距（链式调用）
func (c *Container) Margin(margin Insets) *Container {
	c.style.Margin = margin
	c.flexItem.MarginTop = margin.Top
	c.flexItem.MarginRight = margin.Right
	c.flexItem.MarginBottom = margin.Bottom
	c.flexItem.MarginLeft = margin.Left
	return c
}

// Padding 设置内边距（链式调用）
func (c *Container) Padding(padding Insets) *Container {
	c.style.Padding = padding
	c.flexItem.PaddingTop = padding.Top
	c.flexItem.PaddingRight = padding.Right
	c.flexItem.PaddingBottom = padding.Bottom
	c.flexItem.PaddingLeft = padding.Left
	return c
}

// BackgroundColor 设置背景色（链式调用）
func (c *Container) BackgroundColor(color Color) *Container {
	c.style.BackgroundColor = color
	return c
}

// BackgroundColorToken 使用主题颜色设置背景色（链式调用）
func (c *Container) BackgroundColorToken(token ColorToken) *Container {
	c.style.BackgroundColor = GetThemeColor(token)
	return c
}

// Width 设置宽度（链式调用）
func (c *Container) Width(width float32) *Container {
	c.style.Width = width
	return c
}

// Height 设置高度（链式调用）
func (c *Container) Height(height float32) *Container {
	c.style.Height = height
	return c
}

// Size 设置尺寸（链式调用）
func (c *Container) Size(width, height float32) *Container {
	c.style.Width = width
	c.style.Height = height
	return c
}

// MinWidth 设置最小宽度（链式调用）
func (c *Container) MinWidth(width float32) *Container {
	c.style.MinWidth = width
	return c
}

// MinHeight 设置最小高度（链式调用）
func (c *Container) MinHeight(height float32) *Container {
	c.style.MinHeight = height
	return c
}

// MaxWidth 设置最大宽度（链式调用）
func (c *Container) MaxWidth(width float32) *Container {
	c.style.MaxWidth = width
	return c
}

// MaxHeight 设置最大高度（链式调用）
func (c *Container) MaxHeight(height float32) *Container {
	c.style.MaxHeight = height
	return c
}

// Style 设置样式（链式调用）
func (c *Container) Style(style *Style) *Container {
	c.SetStyle(style)
	return c
}

// Add 添加子组件（链式调用）
func (c *Container) Add(children ...Widget) *Container {
	for _, child := range children {
		c.AddChild(child)
	}
	return c
}

// Measure 测量容器大小
func (c *Container) Measure(constraints flex.Constraint) flex.Size {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 计算内边距
	padding := c.style.Padding
	paddingHorizontal := padding.Left + padding.Right
	paddingVertical := padding.Top + padding.Bottom

	// 计算内部约束（减去 padding）
	innerConstraints := constraints
	if innerConstraints.MaxWidth > 0 {
		innerConstraints.MaxWidth -= paddingHorizontal
		if innerConstraints.MaxWidth < 0 {
			innerConstraints.MaxWidth = 0
		}
	}
	if innerConstraints.MinWidth > 0 {
		innerConstraints.MinWidth -= paddingHorizontal
		if innerConstraints.MinWidth < 0 {
			innerConstraints.MinWidth = 0
		}
	}
	if innerConstraints.MaxHeight > 0 {
		innerConstraints.MaxHeight -= paddingVertical
		if innerConstraints.MaxHeight < 0 {
			innerConstraints.MaxHeight = 0
		}
	}
	if innerConstraints.MinHeight > 0 {
		innerConstraints.MinHeight -= paddingVertical
		if innerConstraints.MinHeight < 0 {
			innerConstraints.MinHeight = 0
		}
	}

	// 应用容器样式中的显式尺寸约束
	if c.style.Width > 0 {
		innerConstraints.MaxWidth = c.style.Width - paddingHorizontal
		innerConstraints.MinWidth = c.style.Width - paddingHorizontal
	}
	if c.style.Height > 0 {
		innerConstraints.MaxHeight = c.style.Height - paddingVertical
		innerConstraints.MinHeight = c.style.Height - paddingVertical
	}
	if c.style.MinWidth > 0 {
		minW := c.style.MinWidth - paddingHorizontal
		if minW > innerConstraints.MinWidth {
			innerConstraints.MinWidth = minW
		}
	}
	if c.style.MinHeight > 0 {
		minH := c.style.MinHeight - paddingVertical
		if minH > innerConstraints.MinHeight {
			innerConstraints.MinHeight = minH
		}
	}
	if c.style.MaxWidth > 0 {
		maxW := c.style.MaxWidth - paddingHorizontal
		if innerConstraints.MaxWidth == 0 || maxW < innerConstraints.MaxWidth {
			innerConstraints.MaxWidth = maxW
		}
	}
	if c.style.MaxHeight > 0 {
		maxH := c.style.MaxHeight - paddingVertical
		if innerConstraints.MaxHeight == 0 || maxH < innerConstraints.MaxHeight {
			innerConstraints.MaxHeight = maxH
		}
	}

	// 复用已有的 items 切片
	c.flexContainer.Items = c.flexContainer.Items[:0]

	// 测量所有子组件并添加到 flex 容器
	for _, child := range c.children {
		if !child.GetVisible() {
			continue
		}

		childSize := child.Measure(innerConstraints)
		childStyle := child.GetStyle()
		item := flex.Item{
			Width:         childSize.Width,
			Height:        childSize.Height,
			FlexGrow:      child.GetFlexProperties().FlexGrow,
			FlexShrink:    child.GetFlexProperties().FlexShrink,
			FlexBasis:     child.GetFlexProperties().FlexBasis,
			AlignSelf:     child.GetFlexProperties().AlignSelf,
			MarginTop:     childStyle.Margin.Top,
			MarginRight:   childStyle.Margin.Right,
			MarginBottom:  childStyle.Margin.Bottom,
			MarginLeft:    childStyle.Margin.Left,
			PaddingTop:    childStyle.Padding.Top,
			PaddingRight:  childStyle.Padding.Right,
			PaddingBottom: childStyle.Padding.Bottom,
			PaddingLeft:   childStyle.Padding.Left,
		}
		c.flexContainer.AddItem(item)
	}

	// 测量容器自身
	innerSize := c.flexContainer.Measure(innerConstraints)

	// 加回 padding
	size := flex.Size{
		Width:  innerSize.Width + paddingHorizontal,
		Height: innerSize.Height + paddingVertical,
	}

	// 应用容器的显式尺寸
	if c.style.Width > 0 {
		size.Width = c.style.Width
	}
	if c.style.Height > 0 {
		size.Height = c.style.Height
	}

	// 应用容器的 Min/Max 约束
	if c.style.MinWidth > 0 && size.Width < c.style.MinWidth {
		size.Width = c.style.MinWidth
	}
	if c.style.MaxWidth > 0 && size.Width > c.style.MaxWidth {
		size.Width = c.style.MaxWidth
	}
	if c.style.MinHeight > 0 && size.Height < c.style.MinHeight {
		size.Height = c.style.MinHeight
	}
	if c.style.MaxHeight > 0 && size.Height > c.style.MaxHeight {
		size.Height = c.style.MaxHeight
	}

	// 应用外部约束
	if constraints.MaxWidth > 0 && size.Width > constraints.MaxWidth {
		size.Width = constraints.MaxWidth
	}
	if constraints.MinWidth > 0 && size.Width < constraints.MinWidth {
		size.Width = constraints.MinWidth
	}
	if constraints.MaxHeight > 0 && size.Height > constraints.MaxHeight {
		size.Height = constraints.MaxHeight
	}
	if constraints.MinHeight > 0 && size.Height < constraints.MinHeight {
		size.Height = constraints.MinHeight
	}

	c.flexItem.SetMeasuredSize(size.Width, size.Height)
	c.measured = true

	return size
}

// Layout 布局容器及其子组件
func (c *Container) Layout(x, y, width, height float32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 应用内边距
	padding := c.style.Padding
	contentX := x + padding.Left
	contentY := y + padding.Top
	contentWidth := width - padding.Left - padding.Right
	contentHeight := height - padding.Top - padding.Bottom

	if contentWidth < 0 {
		contentWidth = 0
	}
	if contentHeight < 0 {
		contentHeight = 0
	}

	// 更新容器自身的位置和大小
	c.rect.Position.X = x
	c.rect.Position.Y = y
	c.rect.Size.Width = width
	c.rect.Size.Height = height

	c.flexItem.Left = x
	c.flexItem.Top = y
	c.flexItem.Width = width
	c.flexItem.Height = height

	// 复用已有的 items 切片并重新填充
	c.flexContainer.Items = c.flexContainer.Items[:0]
	visibleChildren := make([]Widget, 0, len(c.children))

	for _, child := range c.children {
		if !child.GetVisible() {
			continue
		}
		visibleChildren = append(visibleChildren, child)

		style := child.GetStyle()
		childFlexItem := child.GetFlexItem()
		measuredWidth, measuredHeight := childFlexItem.GetMeasuredSize()
		item := flex.Item{
			Width:         measuredWidth,
			Height:        measuredHeight,
			FlexGrow:      child.GetFlexProperties().FlexGrow,
			FlexShrink:    child.GetFlexProperties().FlexShrink,
			FlexBasis:     child.GetFlexProperties().FlexBasis,
			AlignSelf:     child.GetFlexProperties().AlignSelf,
			MarginTop:     style.Margin.Top,
			MarginRight:   style.Margin.Right,
			MarginBottom:  style.Margin.Bottom,
			MarginLeft:    style.Margin.Left,
			PaddingTop:    style.Padding.Top,
			PaddingRight:  style.Padding.Right,
			PaddingBottom: style.Padding.Bottom,
			PaddingLeft:   style.Padding.Left,
		}
		c.flexContainer.AddItem(item)
	}

	// 执行 flex 布局
	c.flexContainer.Layout(contentX, contentY, contentWidth, contentHeight)

	// 将布局结果应用到子组件
	for i, item := range c.flexContainer.Items {
		if i < len(visibleChildren) {
			visibleChildren[i].Layout(item.Left, item.Top, item.Width, item.Height)
		}
	}
}

// Render 渲染容器
func (c *Container) Render() error {
	if !c.GetVisible() {
		return nil
	}

	// 渲染子组件
	for _, child := range c.children {
		if err := child.Render(); err != nil {
			return err
		}
	}

	return nil
}
