package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type forcedDarkTheme struct{}

func (t *forcedDarkTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (t *forcedDarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *forcedDarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *forcedDarkTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func CreateApp() fyne.App {
	a := app.New()
	a.Settings().SetTheme(&forcedDarkTheme{})
	return a
}

func CreateAnalysisWindow(a fyne.App) (fyne.Window, *widget.Entry) {
	w := a.NewWindow("DiscordFreelanceX — Message Analysis")
	w.Resize(fyne.NewSize(1000, 700))
	w.CenterOnScreen()

	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.SetText("Connecting to Discord and loading messages...")
	entry.Disable()

	w.SetContent(container.NewPadded(entry))
	return w, entry
}

func AppendAnalysisText(entry *widget.Entry, text string) {
	fyne.Do(func() {
		current := entry.Text
		if current != "" {
			current += "\n"
		}
		entry.SetText(current + text)
	})
}
