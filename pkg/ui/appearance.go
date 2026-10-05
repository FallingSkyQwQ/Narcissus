package ui

// SystemThemeProvider is implemented by backends that can report the operating
// system's light/dark preference. GTK4 implements it by inspecting the toolkit
// settings; backends without support fall back to light.
type SystemThemeProvider interface {
	// SystemPrefersDark reports whether the OS is set to a dark appearance.
	SystemPrefersDark() bool
}

// ThemeApplier is implemented by backends that can push a light/dark preference
// into the native toolkit, so the chrome the toolkit draws (entries, buttons,
// scrollbars) matches the framework theme. GTK4 implements it through
// GtkSettings:gtk-application-prefer-dark-theme.
type ThemeApplier interface {
	// ApplyTheme asks the toolkit to render in a dark or light appearance.
	ApplyTheme(dark bool)
}

// SystemPrefersDark reports whether the operating system is currently set to a
// dark appearance. It returns false when the backend cannot tell or none is
// installed.
func SystemPrefersDark() bool {
	if backend := CurrentBackend(); backend != nil {
		if provider, ok := backend.(SystemThemeProvider); ok {
			return provider.SystemPrefersDark()
		}
	}
	return false
}

// SystemTheme returns the framework theme that matches the OS appearance.
//
// The backend may not be able to answer before the toolkit is initialized, so
// prefer calling it from an OnReady callback.
func SystemTheme() Theme {
	if SystemPrefersDark() {
		return NewDarkTheme()
	}
	return NewLightTheme()
}

// applyThemeToBackend forwards a theme to the toolkit when the backend supports
// it. SetTheme calls it so that toggling the framework theme keeps native
// controls in step.
func applyThemeToBackend(theme Theme) {
	if theme == nil {
		return
	}
	if backend := CurrentBackend(); backend != nil {
		if applier, ok := backend.(ThemeApplier); ok {
			applier.ApplyTheme(theme.GetType() == ThemeDark)
		}
	}
}
