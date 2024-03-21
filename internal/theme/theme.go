package theme

import (
	"github.com/gdamore/tcell/v2"
)

type Theme struct {
	Name                        string
	BackgroundColor             tcell.Color
	ForgroundColor              tcell.Color
	HighlightColor              tcell.Color
	AccentColor                 tcell.Color
	AccentColorTwo              tcell.Color
	ButtonStyle                 tcell.Style
	ActivatedButtonStyle        tcell.Style
	DropdownListUnselectedStyle tcell.Style
	DropdownListSelectedStyle   tcell.Style
	TextAreaStyle               tcell.Style
}

func NewTheme(themeName string) *Theme {
	getDefault := func() *Theme {

		// tview.Styles.BorderColor = tcell.ColorWhite
		// tview.Styles.TitleColor = tcell.ColorWhite
		// tview.Styles.PrimaryTextColor = tcell.ColorWhite

		return &Theme{
			Name:                        "default",
			BackgroundColor:             tcell.NewHexColor(0x111111),
			ForgroundColor:              tcell.ColorWhite,
			HighlightColor:              tcell.NewHexColor(0xFFC300),
			AccentColor:                 tcell.NewHexColor(0x444444),
			AccentColorTwo:              tcell.NewHexColor(0x222222),
			ButtonStyle:                 tcell.StyleDefault.Background(tcell.NewHexColor(0x222222)).Foreground(tcell.ColorWhite),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(tcell.NewHexColor(0xFFC300)).Foreground(tcell.ColorBlack),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(tcell.NewHexColor(0xFFC300)).Foreground(tcell.ColorBlack),
			TextAreaStyle:               tcell.StyleDefault.Background(0x111111).Foreground(tcell.ColorWhite),
		}
	}

	switch themeName {
	case "default":
		return getDefault()
	case "america":
		// tview.Styles.BorderColor = tcell.ColorRed
		// tview.Styles.TitleColor = tcell.ColorRed
		// tview.Styles.PrimaryTextColor = tcell.ColorWhite

		return &Theme{
			Name:                        "america",
			BackgroundColor:             tcell.ColorBlue,
			ForgroundColor:              tcell.ColorWhite,
			HighlightColor:              tcell.ColorRed,
			AccentColor:                 tcell.ColorWhite,
			AccentColorTwo:              tcell.ColorRed,
			ButtonStyle:                 tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorRed),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite),
			TextAreaStyle:               tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite),
		}
	case "matrix":
		// tview.Styles.BorderColor = tcell.ColorGreen
		// tview.Styles.TitleColor = tcell.ColorGreen
		// tview.Styles.PrimaryTextColor = tcell.ColorGreen

		return &Theme{
			Name:                        "matrix",
			BackgroundColor:             tcell.ColorBlack,
			ForgroundColor:              tcell.ColorGreen,
			HighlightColor:              tcell.ColorGreen,
			AccentColor:                 tcell.NewHexColor(0x111111),
			AccentColorTwo:              tcell.ColorBlack,
			ButtonStyle:                 tcell.StyleDefault.Background(tcell.NewHexColor(0x111111)).Foreground(tcell.ColorGreen),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.NewHexColor(0x111111)).Foreground(tcell.ColorGreen),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack),
			TextAreaStyle:               tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen),
		}
	case "halloween":
		orange := tcell.NewHexColor(0xDD863A)

		// tview.Styles.BorderColor = tcell.ColorBlack
		// tview.Styles.TitleColor = tcell.ColorBlack
		// tview.Styles.PrimaryTextColor = tcell.ColorBlack

		return &Theme{
			Name:                        "halloween",
			BackgroundColor:             tcell.ColorDarkOrange,
			ForgroundColor:              tcell.ColorBlack,
			HighlightColor:              orange,
			AccentColor:                 tcell.ColorOrangeRed,
			AccentColorTwo:              tcell.ColorBlack,
			ButtonStyle:                 tcell.StyleDefault.Background(orange).Foreground(tcell.ColorBlack),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(orange),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.NewHexColor(0x111111)).Foreground(orange),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(orange).Foreground(tcell.ColorBlack),
			TextAreaStyle:               tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(orange),
		}
	case "christmas":
		// tview.Styles.BorderColor = tcell.ColorRed
		// tview.Styles.TitleColor = tcell.ColorRed
		// tview.Styles.PrimaryTextColor = tcell.ColorWhite

		lightGreen := tcell.NewHexColor(0x00FF00)

		return &Theme{
			Name:                        "christmas",
			BackgroundColor:             tcell.ColorDarkGreen,
			ForgroundColor:              tcell.ColorWhite,
			HighlightColor:              tcell.ColorRed,
			AccentColor:                 lightGreen,
			AccentColorTwo:              tcell.ColorRed,
			ButtonStyle:                 tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorWhite),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorWhite),
			TextAreaStyle:               tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorWhite),
		}
	case "satanic":
		// tview.Styles.BorderColor = tcell.ColorRed
		// tview.Styles.TitleColor = tcell.ColorRed
		// tview.Styles.PrimaryTextColor = tcell.ColorBlack

		return &Theme{
			Name:                        "satanic",
			BackgroundColor:             tcell.ColorBlack,
			ForgroundColor:              tcell.ColorRed,
			HighlightColor:              tcell.ColorRed,
			AccentColor:                 tcell.ColorDarkRed,
			AccentColorTwo:              tcell.NewHexColor(0x111111),
			ButtonStyle:                 tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.NewHexColor(0x111111)),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(tcell.NewHexColor(0x222222)).Foreground(tcell.ColorRed),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.ColorDarkRed).Foreground(tcell.NewHexColor(0x111111)),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorBlack),
			TextAreaStyle:               tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed),
		}

	default:
		return getDefault()
	}
}
