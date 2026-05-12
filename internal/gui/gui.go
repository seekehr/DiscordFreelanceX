package gui

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/seekehr/DiscordFreelanceX/internal"
)

type forcedDarkTheme struct{}

func (t *forcedDarkTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
		return color.White
	}
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

func CreateAnalysisWindow(a fyne.App) (fyne.Window, *widget.RichText) {
	w := a.NewWindow("DiscordFreelanceX — Message Analysis")
	w.Resize(fyne.NewSize(1000, 700))
	w.CenterOnScreen()

	rt := widget.NewRichText(
		&widget.TextSegment{
			Text: "Connecting to Discord and loading messages...\n",
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameForeground,
				TextStyle: fyne.TextStyle{Monospace: true},
			},
		},
	)
	rt.Wrapping = fyne.TextWrapWord

	scroll := container.NewVScroll(rt)
	w.SetContent(scroll)
	return w, rt
}

func AppendAnalysisText(rt *widget.RichText, text string) {
	fyne.Do(func() {
		rt.Segments = append(rt.Segments, &widget.TextSegment{
			Text: text + "\n",
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameForeground,
				TextStyle: fyne.TextStyle{Monospace: true},
			},
		})
		rt.Refresh()
	})
}

func AppendAnalysisEntries(rt *widget.RichText, entries []internal.AnalysisEntry) {
	fyne.Do(func() {
		for _, e := range entries {
			rt.Segments = append(rt.Segments, &widget.TextSegment{
				Text: e.Text + " ",
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameForeground,
					TextStyle: fyne.TextStyle{Monospace: true},
					Inline:    true,
				},
			})

			if e.MessageURL != "" {
				u, _ := url.Parse(e.MessageURL)
				rt.Segments = append(rt.Segments, &widget.HyperlinkSegment{
					Text: "Go to message",
					URL:  u,
				})
			}

			rt.Segments = append(rt.Segments, &widget.TextSegment{
				Text:  "\n",
				Style: widget.RichTextStyle{},
			})
		}
		rt.Refresh()
	})
}
