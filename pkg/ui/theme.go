package ui

import (
	"sync"

	"github.com/FallingSkyQwQ/Narcissus/pkg/reactive"
)

// ThemeType 主题类型
type ThemeType int

const (
	ThemeLight ThemeType = iota
	ThemeDark
	ThemeAuto
)

// Theme 定义主题接口
type Theme interface {
	GetType() ThemeType
	GetColor(token ColorToken) Color
	GetStyle(componentType string) *Style
}

// BaseTheme 主题基础实现
type BaseTheme struct {
	ThemeType ThemeType
	Colors    map[ColorToken]Color
	Styles    map[string]*Style
}

// GetType 返回主题类型
func (t *BaseTheme) GetType() ThemeType {
	return t.ThemeType
}

// GetColor 获取主题颜色
func (t *BaseTheme) GetColor(token ColorToken) Color {
	if color, ok := t.Colors[token]; ok {
		return color
	}
	return ColorBlack
}

// GetStyle 获取组件样式
func (t *BaseTheme) GetStyle(componentType string) *Style {
	if style, ok := t.Styles[componentType]; ok {
		return style.Clone()
	}
	return NewStyle()
}

// 预定义主题

// NewLightTheme 创建浅色主题
func NewLightTheme() *BaseTheme {
	return &BaseTheme{
		ThemeType: ThemeLight,
		Colors: map[ColorToken]Color{
			TokenBackground:     ColorLightBackground,
			TokenSurface:        ColorLightSurface,
			TokenCardBackground: ColorLightCardBackground,
			TokenTextPrimary:    ColorLightTextPrimary,
			TokenTextSecondary:  ColorLightTextSecondary,
			TokenBorder:         ColorLightBorder,
			TokenDivider:        ColorLightDivider,
			TokenPrimary:        ColorPrimary,
			TokenSecondary:      ColorSecondary,
			TokenAccent:         ColorAccent,
			TokenSuccess:        ColorSuccess,
			TokenWarning:        ColorWarning,
			TokenError:          ColorError,
		},
		Styles: make(map[string]*Style),
	}
}

// NewDarkTheme 创建深色主题
func NewDarkTheme() *BaseTheme {
	return &BaseTheme{
		ThemeType: ThemeDark,
		Colors: map[ColorToken]Color{
			TokenBackground:     ColorDarkBackground,
			TokenSurface:        ColorDarkSurface,
			TokenCardBackground: ColorDarkCardBackground,
			TokenTextPrimary:    ColorDarkTextPrimary,
			TokenTextSecondary:  ColorDarkTextSecondary,
			TokenBorder:         ColorDarkBorder,
			TokenDivider:        ColorDarkDivider,
			TokenPrimary:        ColorPrimary,
			TokenSecondary:      ColorSecondary,
			TokenAccent:         ColorAccent,
			TokenSuccess:        ColorSuccess,
			TokenWarning:        ColorWarning,
			TokenError:          ColorError,
		},
		Styles: make(map[string]*Style),
	}
}

// ThemeManager 主题管理器
type ThemeManager struct {
	mu          sync.RWMutex
	current     Theme
	themeSignal *reactive.Signal[Theme]
}

var (
	globalThemeManager *ThemeManager
	themeManagerOnce   sync.Once
)

// GetThemeManager 获取全局主题管理器
func GetThemeManager() *ThemeManager {
	themeManagerOnce.Do(func() {
		globalThemeManager = &ThemeManager{
			current:     NewLightTheme(),
			themeSignal: reactive.NewSignal[Theme](NewLightTheme()),
		}
	})
	return globalThemeManager
}

// SetTheme 设置当前主题
func (tm *ThemeManager) SetTheme(theme Theme) {
	if theme == nil {
		theme = NewLightTheme()
	}
	tm.mu.Lock()
	tm.current = theme
	tm.mu.Unlock()
	tm.themeSignal.Set(theme)
}

// GetTheme 获取当前主题
func (tm *ThemeManager) GetTheme() Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.current
}

// Subscribe 订阅主题变更
func (tm *ThemeManager) Subscribe(observer func(Theme)) func() {
	return tm.themeSignal.Subscribe(func(theme Theme) {
		observer(theme)
	})
}

// GetColor 获取当前主题的颜色
func (tm *ThemeManager) GetColor(token ColorToken) Color {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if tm.current == nil {
		return NewLightTheme().GetColor(token)
	}
	return tm.current.GetColor(token)
}

// 便捷函数

// SetTheme 设置全局主题
func SetTheme(theme Theme) {
	GetThemeManager().SetTheme(theme)
}

// GetCurrentTheme 获取当前主题
func GetCurrentTheme() Theme {
	return GetThemeManager().GetTheme()
}

// GetThemeColor 获取当前主题的颜色
func GetThemeColor(token ColorToken) Color {
	return GetThemeManager().GetColor(token)
}

// SubscribeTheme 订阅主题变更
func SubscribeTheme(observer func(Theme)) func() {
	return GetThemeManager().Subscribe(observer)
}
