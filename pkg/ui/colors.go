package ui

import "fmt"

// Color 表示 RGBA 颜色
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// RGBA 返回颜色的 RGBA 分量
func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	g = uint32(c.G)
	b = uint32(c.B)
	a = uint32(c.A)
	return
}

// Hex 返回颜色的十六进制字符串表示
func (c Color) Hex() string {
	if c.A == 255 {
		return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}

// WithAlpha 返回具有指定透明度的新颜色
func (c Color) WithAlpha(alpha uint8) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: alpha}
}

// IsTransparent 返回颜色是否完全透明
func (c Color) IsTransparent() bool {
	return c.A == 0
}

// RGB 创建一个不透明的颜色
func RGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, A: 255}
}

// RGBA 创建一个带透明度的颜色
func RGBA(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// 常用颜色定义
var (
	// 基础颜色
	ColorTransparent = Color{R: 0, G: 0, B: 0, A: 0}
	ColorBlack       = RGB(0, 0, 0)
	ColorWhite       = RGB(255, 255, 255)
	ColorRed         = RGB(255, 0, 0)
	ColorGreen       = RGB(0, 255, 0)
	ColorBlue        = RGB(0, 0, 255)
	ColorYellow      = RGB(255, 255, 0)
	ColorCyan        = RGB(0, 255, 255)
	ColorMagenta     = RGB(255, 0, 255)

	// 灰度颜色
	ColorGray      = RGB(128, 128, 128)
	ColorLightGray = RGB(211, 211, 211)
	ColorDarkGray  = RGB(64, 64, 64)
	ColorSilver    = RGB(192, 192, 192)
	ColorDimGray   = RGB(105, 105, 105)

	// 主题颜色 - 浅色主题
	ColorLightBackground     = RGB(255, 255, 255)
	ColorLightSurface        = RGB(250, 250, 250)
	ColorLightCardBackground = RGB(255, 255, 255)
	ColorLightTextPrimary    = RGB(33, 33, 33)
	ColorLightTextSecondary  = RGB(117, 117, 117)
	ColorLightBorder         = RGB(224, 224, 224)
	ColorLightDivider        = RGB(238, 238, 238)

	// 主题颜色 - 深色主题
	ColorDarkBackground     = RGB(18, 18, 18)
	ColorDarkSurface        = RGB(30, 30, 30)
	ColorDarkCardBackground = RGB(42, 42, 42)
	ColorDarkTextPrimary    = RGB(255, 255, 255)
	ColorDarkTextSecondary  = RGB(176, 176, 176)
	ColorDarkBorder         = RGB(64, 64, 64)
	ColorDarkDivider        = RGB(48, 48, 48)

	// 强调色
	ColorPrimary   = RGB(33, 150, 243)
	ColorSecondary = RGB(156, 39, 176)
	ColorAccent    = RGB(255, 64, 129)
	ColorSuccess   = RGB(76, 175, 80)
	ColorWarning   = RGB(255, 152, 0)
	ColorError     = RGB(244, 67, 54)
	ColorInfo      = RGB(33, 150, 243)
)

// ColorToken 是主题颜色标识符
type ColorToken int

const (
	TokenBackground ColorToken = iota
	TokenSurface
	TokenCardBackground
	TokenTextPrimary
	TokenTextSecondary
	TokenBorder
	TokenDivider
	TokenPrimary
	TokenSecondary
	TokenAccent
	TokenSuccess
	TokenWarning
	TokenError
)

// ColorTokenResolver 将 ColorToken 解析为实际颜色
type ColorTokenResolver func(token ColorToken) Color
