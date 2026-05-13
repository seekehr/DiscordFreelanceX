package gui

import (
	"image/color"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/seekehr/DiscordFreelanceX/internal"
)

// TabbedDisplay holds the per-guild and "New" tab RichText widgets
// so callers can direct output to the appropriate tab.
type TabbedDisplay struct {
	NewRT    *widget.RichText
	GuildRTs map[string]*widget.RichText
	Tabs     *container.AppTabs
	guildTab map[string]*container.TabItem
}

// forcedDarkTheme wraps Fyne's default theme but always returns dark-variant
// colours, with the foreground forced to pure white for readability.
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

// CreateApp initialises a new Fyne application with the forced dark theme.
func CreateApp() fyne.App {
	a := app.NewWithID("com.seekehr.discordfreelancex")
	a.Settings().SetTheme(&forcedDarkTheme{})
	return a
}

// CreateTabbedWindow builds the main window with a tab bar: a "New" tab for
// incoming real-time messages, plus one tab per configured guild for historical
// analysis results.
func CreateTabbedWindow(a fyne.App, cfg *internal.Config) (fyne.Window, *TabbedDisplay) {
	w := a.NewWindow("DiscordFreelanceX — Made by @seekehr")
	w.Resize(fyne.NewSize(1000, 700))
	w.CenterOnScreen()

	td := &TabbedDisplay{
		GuildRTs: make(map[string]*widget.RichText, len(cfg.ReceiveServers)),
		guildTab: make(map[string]*container.TabItem, len(cfg.ReceiveServers)),
	}

	newRT := newRichText("Waiting for new messages...\n")
	td.NewRT = newRT
	newTab := container.NewTabItem("New", container.NewVScroll(newRT))

	tabs := container.NewAppTabs(newTab)
	tabs.SetTabLocation(container.TabLocationTop)

	for _, server := range cfg.ReceiveServers {
		rt := newRichText("Connecting to Discord...\n")
		td.GuildRTs[server.GuildID] = rt
		tab := container.NewTabItem(server.GuildID, container.NewVScroll(rt))
		td.guildTab[server.GuildID] = tab
		tabs.Append(tab)
	}

	td.Tabs = tabs
	w.SetContent(tabs)
	return w, td
}

// RenameGuildTab updates a guild tab's title (e.g. once the guild name is resolved).
func (td *TabbedDisplay) RenameGuildTab(guildID, name string) {
	fyne.Do(func() {
		if tab, ok := td.guildTab[guildID]; ok {
			tab.Text = name
			td.Tabs.Refresh()
		}
	})
}

func newRichText(initialText string) *widget.RichText {
	rt := widget.NewRichText(
		&widget.TextSegment{
			Text: initialText,
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameForeground,
				TextStyle: fyne.TextStyle{Monospace: true},
			},
		},
	)
	rt.Wrapping = fyne.TextWrapWord
	return rt
}

// AppendAnalysisText adds a plain white monospace line to the RichText widget.
// Safe to call from any goroutine (uses fyne.Do).
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

// parseStyledSegments splits text on *bold* markers and returns a mix of
// normal and bold RichText segments. Unmatched asterisks are left as-is.
func parseStyledSegments(text string) []widget.RichTextSegment {
	var segments []widget.RichTextSegment
	remaining := text
	for {
		start := strings.Index(remaining, "*")
		if start == -1 {
			break
		}
		end := strings.Index(remaining[start+1:], "*")
		if end == -1 {
			break
		}
		end += start + 1

		if start > 0 {
			segments = append(segments, &widget.TextSegment{
				Text: remaining[:start],
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameForeground,
					TextStyle: fyne.TextStyle{Monospace: true},
					Inline:    true,
				},
			})
		}

		segments = append(segments, &widget.TextSegment{
			Text: remaining[start+1 : end],
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameForeground,
				TextStyle: fyne.TextStyle{Monospace: true, Bold: true},
				Inline:    true,
			},
		})

		remaining = remaining[end+1:]
	}

	if remaining != "" {
		segments = append(segments, &widget.TextSegment{
			Text: remaining,
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameForeground,
				TextStyle: fyne.TextStyle{Monospace: true},
				Inline:    true,
			},
		})
	}

	return segments
}

// AppendAnalysisEntries renders structured analysis data into the RichText widget.
// Each entry's text is parsed for *bold* markers, and entries with a MessageURL
// get a clickable "Go to message" hyperlink appended.
// Safe to call from any goroutine (uses fyne.Do).
func AppendAnalysisEntries(rt *widget.RichText, entries []internal.AnalysisEntry) {
	fyne.Do(func() {
		for _, e := range entries {
			rt.Segments = append(rt.Segments, parseStyledSegments(e.Text+" ")...)

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
