package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CyberpunkTheme struct{}

var _ fyne.Theme = (*CyberpunkTheme)(nil)

func NewCyberpunkTheme() fyne.Theme {
	return &CyberpunkTheme{}
}

func (m CyberpunkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		// Dark background, slightly transparent (not natively supported by all OS window backends without extra config,
		// but we simulate the "app background").
		// Actually, Fyne window background transparency requires `window.SetTransparent(true)` and specific colored canvas.
		// For the theme itself, we stick to a solid dark or slightly alpha if we want "glass" widgets.
		return color.RGBA{R: 20, G: 20, B: 30, A: 240} // Deep dark blue-grey
	case theme.ColorNameForeground:
		return color.RGBA{R: 220, G: 230, B: 240, A: 255} // Off-white
	case theme.ColorNamePrimary:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255} // Neutral Grey/White
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 40, G: 40, B: 50, A: 200}
	case theme.ColorNameOverlayBackground:
		return color.RGBA{R: 30, G: 30, B: 40, A: 230}
	case theme.ColorNameButton: // Fallback or override standard button color
		return color.RGBA{R: 100, G: 100, B: 110, A: 255} // Neutral Slate
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m CyberpunkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m CyberpunkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m CyberpunkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 24
	case theme.SizeNameScrollBar:
		return 16
	}
	return theme.DefaultTheme().Size(name)
}
