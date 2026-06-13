package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// ImageFit 定义图片适应模式
type ImageFit int

const (
	// ImageFitContain 保持比例，完整显示在容器内
	ImageFitContain ImageFit = iota
	// ImageFitCover 保持比例，填满容器（可能裁剪）
	ImageFitCover
	// ImageFitFill 拉伸填满容器（不保持比例）
	ImageFitFill
	// ImageFitNone 原始尺寸
	ImageFitNone
	// ImageFitScaleDown 类似 contain，但不会放大超过原始尺寸
	ImageFitScaleDown
)

// Image 是图片显示组件
type Image struct {
	*BaseWidget

	// 图片源路径或 URL
	source string

	// 图片适应模式
	fit ImageFit

	// 图片原始尺寸（加载后设置）
	naturalWidth  float32
	naturalHeight float32

	// 是否已加载
	loaded bool
}

// NewImage 创建新的图片组件
func NewImage() *Image {
	img := &Image{
		BaseWidget: NewBaseWidget(),
		fit:        ImageFitContain,
	}

	// 设置默认样式
	img.style.BackgroundColor = ColorTransparent

	return img
}

// Source 设置图片源（链式调用）
func (img *Image) Source(source string) *Image {
	img.source = source
	img.loaded = false
	img.measured = false
	return img
}

// GetSource 获取图片源
func (img *Image) GetSource() string {
	return img.source
}

// Fit 设置图片适应模式（链式调用）
func (img *Image) Fit(fit ImageFit) *Image {
	img.fit = fit
	return img
}

// GetFit 获取图片适应模式
func (img *Image) GetFit() ImageFit {
	return img.fit
}

// Width 设置宽度（链式调用）
func (img *Image) Width(width float32) *Image {
	img.style.Width = width
	return img
}

// Height 设置高度（链式调用）
func (img *Image) Height(height float32) *Image {
	img.style.Height = height
	return img
}

// Size 设置尺寸（链式调用）
func (img *Image) Size(width, height float32) *Image {
	img.style.Width = width
	img.style.Height = height
	return img
}

// FlexGrow 设置 FlexGrow（链式调用）
func (img *Image) FlexGrow(grow float32) *Image {
	img.style.FlexGrow = grow
	img.flexItem.FlexGrow = grow
	return img
}

// FlexShrink 设置 FlexShrink（链式调用）
func (img *Image) FlexShrink(shrink float32) *Image {
	img.style.FlexShrink = shrink
	img.flexItem.FlexShrink = shrink
	return img
}

// Margin 设置外边距（链式调用）
func (img *Image) Margin(margin Insets) *Image {
	img.style.Margin = margin
	img.flexItem.MarginTop = margin.Top
	img.flexItem.MarginRight = margin.Right
	img.flexItem.MarginBottom = margin.Bottom
	img.flexItem.MarginLeft = margin.Left
	return img
}

// BorderRadius 设置边框圆角（链式调用）
func (img *Image) BorderRadius(radius float32) *Image {
	img.style.Border.Radius = radius
	return img
}

// Style 设置样式（链式调用）
func (img *Image) Style(style *Style) *Image {
	img.SetStyle(style)
	return img
}

// calculateSize 根据适应模式计算图片显示尺寸
func (img *Image) calculateSize(containerWidth, containerHeight float32) (width, height float32) {
	// 如果图片未加载或没有原始尺寸，使用容器尺寸
	if !img.loaded || img.naturalWidth == 0 || img.naturalHeight == 0 {
		return containerWidth, containerHeight
	}

	// 如果没有设置容器尺寸，使用原始尺寸
	if containerWidth == 0 && containerHeight == 0 {
		return img.naturalWidth, img.naturalHeight
	}

	switch img.fit {
	case ImageFitContain:
		return img.calculateContainSize(containerWidth, containerHeight)
	case ImageFitCover:
		return img.calculateCoverSize(containerWidth, containerHeight)
	case ImageFitFill:
		return containerWidth, containerHeight
	case ImageFitNone:
		return img.naturalWidth, img.naturalHeight
	case ImageFitScaleDown:
		w, h := img.calculateContainSize(containerWidth, containerHeight)
		if w > img.naturalWidth || h > img.naturalHeight {
			return img.naturalWidth, img.naturalHeight
		}
		return w, h
	default:
		return containerWidth, containerHeight
	}
}

// calculateContainSize 计算 contain 模式的尺寸
func (img *Image) calculateContainSize(containerWidth, containerHeight float32) (width, height float32) {
	if containerWidth == 0 || containerHeight == 0 {
		return img.naturalWidth, img.naturalHeight
	}

	containerRatio := containerWidth / containerHeight
	imageRatio := img.naturalWidth / img.naturalHeight

	if imageRatio > containerRatio {
		// 图片更宽，以宽度为基准
		width = containerWidth
		height = containerWidth / imageRatio
	} else {
		// 图片更高，以高度为基准
		height = containerHeight
		width = containerHeight * imageRatio
	}

	return width, height
}

// calculateCoverSize 计算 cover 模式的尺寸
func (img *Image) calculateCoverSize(containerWidth, containerHeight float32) (width, height float32) {
	if containerWidth == 0 || containerHeight == 0 {
		return img.naturalWidth, img.naturalHeight
	}

	containerRatio := containerWidth / containerHeight
	imageRatio := img.naturalWidth / img.naturalHeight

	if imageRatio > containerRatio {
		// 图片更宽，以高度为基准
		height = containerHeight
		width = containerHeight * imageRatio
	} else {
		// 图片更高，以宽度为基准
		width = containerWidth
		height = containerWidth / imageRatio
	}

	return width, height
}

// Measure 测量图片组件大小
func (img *Image) Measure(constraints flex.Constraint) flex.Size {
	img.mu.Lock()
	defer img.mu.Unlock()

	var width, height float32

	// 优先使用样式中设置的固定尺寸
	if img.style.Width > 0 && img.style.Height > 0 {
		width = img.style.Width
		height = img.style.Height
	} else if img.style.Width > 0 {
		// 只设置了宽度，根据比例计算高度
		width = img.style.Width
		if img.loaded && img.naturalWidth > 0 {
			height = width * (img.naturalHeight / img.naturalWidth)
		} else {
			height = width // 默认正方形
		}
	} else if img.style.Height > 0 {
		// 只设置了高度，根据比例计算宽度
		height = img.style.Height
		if img.loaded && img.naturalHeight > 0 {
			width = height * (img.naturalWidth / img.naturalHeight)
		} else {
			width = height // 默认正方形
		}
	} else {
		// 没有设置尺寸，使用原始尺寸或默认尺寸
		if img.loaded && img.naturalWidth > 0 {
			width = img.naturalWidth
			height = img.naturalHeight
		} else {
			// 默认尺寸
			width = 100
			height = 100
		}
	}

	// 应用约束
	if constraints.MaxWidth > 0 && width > constraints.MaxWidth {
		width = constraints.MaxWidth
	}
	if constraints.MinWidth > 0 && width < constraints.MinWidth {
		width = constraints.MinWidth
	}
	if constraints.MaxHeight > 0 && height > constraints.MaxHeight {
		height = constraints.MaxHeight
	}
	if constraints.MinHeight > 0 && height < constraints.MinHeight {
		height = constraints.MinHeight
	}

	img.flexItem.SetMeasuredSize(width, height)
	img.measured = true

	return flex.Size{Width: width, Height: height}
}

// Layout 布局图片组件
func (img *Image) Layout(x, y, width, height float32) {
	img.mu.Lock()
	defer img.mu.Unlock()

	img.rect.Position.X = x
	img.rect.Position.Y = y
	img.rect.Size.Width = width
	img.rect.Size.Height = height

	img.flexItem.Left = x
	img.flexItem.Top = y
	img.flexItem.Width = width
	img.flexItem.Height = height
}

// Render 渲染图片组件
func (img *Image) Render() error {
	if !img.GetVisible() {
		return nil
	}

	// 实际渲染应该调用底层 WinUI API 创建 Image
	// 这里只是占位实现
	return nil
}

// SetNaturalSize 设置图片原始尺寸（加载完成后调用）
func (img *Image) SetNaturalSize(width, height float32) {
	img.mu.Lock()
	defer img.mu.Unlock()

	img.naturalWidth = width
	img.naturalHeight = height
	img.loaded = true
	img.measured = false // 需要重新测量
}

// IsLoaded 返回图片是否已加载
func (img *Image) IsLoaded() bool {
	img.mu.RLock()
	defer img.mu.RUnlock()
	return img.loaded
}
