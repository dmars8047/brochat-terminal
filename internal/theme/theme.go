package theme

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Theme struct {
	Code                        string
	BackgroundColor             tcell.Color
	ForgroundColor              tcell.Color
	HighlightColor              tcell.Color
	AccentColor                 tcell.Color
	AccentColorTwo              tcell.Color
	ButtonStyle                 tcell.Style
	ActivatedButtonStyle        tcell.Style
	DropdownListUnselectedStyle tcell.Style
	DropdownListSelectedStyle   tcell.Style
	TextAreaTextStyle           tcell.Style
	BorderColor                 tcell.Color
	TitleColor                  tcell.Color
	InfoColor                   tcell.Color
	InfoColorTwo                tcell.Color
	ChatTextColor               tcell.Color
	ChatLabelColors             []string
}

func NewTheme(themeName string) *Theme {
	getDefault := func() *Theme {
		return &Theme{
			Code:                        "default",
			BackgroundColor:             tcell.NewHexColor(0x111111),
			ForgroundColor:              tcell.ColorWhite,
			HighlightColor:              tcell.NewHexColor(0xFFC300),
			AccentColor:                 tcell.NewHexColor(0x444444),
			AccentColorTwo:              tcell.NewHexColor(0x222222),
			ButtonStyle:                 tcell.StyleDefault.Background(tcell.NewHexColor(0x222222)).Foreground(tcell.ColorWhite),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(tcell.NewHexColor(0xFFC300)).Foreground(tcell.ColorBlack),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(tcell.NewHexColor(0xFFC300)).Foreground(tcell.ColorBlack),
			TextAreaTextStyle:           tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.NewHexColor(0x111111)),
			BorderColor:                 tcell.ColorWhite,
			TitleColor:                  tcell.ColorWhite,
			InfoColor:                   tcell.ColorWhite,
			InfoColorTwo:                tcell.NewHexColor(0x777777),
			ChatTextColor:               tcell.ColorWhite,
			ChatLabelColors: []string{
				"#33DA7A", // Light Green
				"#C061CB", // Lilac
				"#FF6B30", // Orange
				"#5928ED", // Purple
				"#00FFFF", // Cyan
				"#FF5555", // Light Red
				"#FAEC34", // Yellow
				"#FFAAFF", // Light Pink
			},
		}
	}

	switch themeName {
	case "default":
		return getDefault()
	case "america":
		return &Theme{
			Code:                        "america",
			BackgroundColor:             tcell.ColorBlue,
			ForgroundColor:              tcell.ColorWhite,
			HighlightColor:              tcell.ColorRed,
			AccentColor:                 tcell.ColorWhite,
			AccentColorTwo:              tcell.ColorRed,
			ButtonStyle:                 tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(tcell.ColorCornflowerBlue).Foreground(tcell.ColorRed),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite),
			TextAreaTextStyle:           tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite),
			BorderColor:                 tcell.ColorRed,
			TitleColor:                  tcell.ColorRed,
			InfoColor:                   tcell.ColorWhite,
			InfoColorTwo:                tcell.ColorGhostWhite,
			ChatTextColor:               tcell.ColorWhite,
			ChatLabelColors: []string{
				tcell.ColorRed.CSS(),
				tcell.ColorGold.CSS(),
				tcell.ColorDarkBlue.CSS(),
				tcell.ColorGreenYellow.CSS(),
				tcell.ColorDarkRed.CSS(),
				"#222222",
			},
		}
	case "matrix":
		trueBlack := tcell.NewHexColor(0x000000)
		black := tcell.NewHexColor(0x111111)
		brightGreen := tcell.NewHexColor(0x00FF00)
		darkerGreen := tcell.NewHexColor(0x00CC00)

		return &Theme{
			Code:                        "matrix",
			BackgroundColor:             trueBlack,
			ForgroundColor:              brightGreen,
			HighlightColor:              brightGreen,
			AccentColor:                 black,
			AccentColorTwo:              trueBlack,
			ButtonStyle:                 tcell.StyleDefault.Background(black).Foreground(darkerGreen),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(darkerGreen).Foreground(trueBlack),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.NewHexColor(0x222222)).Foreground(tcell.ColorGreen),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(brightGreen).Foreground(trueBlack),
			TextAreaTextStyle:           tcell.StyleDefault.Background(trueBlack).Foreground(brightGreen),
			BorderColor:                 darkerGreen,
			TitleColor:                  brightGreen,
			InfoColor:                   darkerGreen,
			InfoColorTwo:                tcell.ColorDarkGreen,
			ChatTextColor:               tcell.ColorWhite,
			ChatLabelColors: []string{
				tcell.ColorFuchsia.CSS(),
				tcell.ColorAqua.CSS(),
				tcell.ColorYellow.CSS(),
				tcell.ColorPink.CSS(),
				tcell.ColorLavender.CSS(),
				tcell.ColorMintCream.CSS(),
				tcell.ColorRed.CSS(),
				tcell.ColorLightSkyBlue.CSS(),
			},
		}
	case "halloween":
		orange := tcell.ColorOrange
		black := tcell.NewHexColor(0x111111)
		trueBlack := tcell.NewHexColor(0x000000)

		return &Theme{
			Code:                        "halloween",
			BackgroundColor:             tcell.ColorDarkOrange,
			ForgroundColor:              trueBlack,
			HighlightColor:              trueBlack,
			AccentColor:                 tcell.ColorOrangeRed,
			AccentColorTwo:              tcell.ColorYellow,
			ButtonStyle:                 tcell.StyleDefault.Background(orange).Foreground(trueBlack),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(trueBlack).Foreground(orange),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(black).Foreground(orange),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(orange).Foreground(tcell.ColorDarkOrange),
			TextAreaTextStyle:           tcell.StyleDefault.Background(trueBlack).Foreground(orange),
			BorderColor:                 trueBlack,
			TitleColor:                  trueBlack,
			InfoColor:                   trueBlack,
			InfoColorTwo:                tcell.NewHexColor(0x444444),
			ChatTextColor:               trueBlack,
			ChatLabelColors: []string{
				tcell.ColorOrangeRed.CSS(),
				tcell.ColorYellow.CSS(),
				tcell.ColorDarkOrange.CSS(),
				tcell.ColorOrange.CSS(),
				tcell.ColorBrown.CSS(),
				tcell.ColorWhite.CSS()},
		}
	case "christmas":
		lightGreen := tcell.NewHexColor(0x00FF00)
		trueRed := tcell.NewHexColor(0xFF0000)

		return &Theme{
			Code:                        "christmas",
			BackgroundColor:             tcell.ColorDarkGreen,
			ForgroundColor:              tcell.ColorWhite,
			HighlightColor:              trueRed,
			AccentColor:                 lightGreen,
			AccentColorTwo:              tcell.ColorRed,
			ButtonStyle:                 tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorWhite),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(trueRed).Foreground(tcell.ColorWhite),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(trueRed).Foreground(tcell.ColorWhite),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorWhite),
			TextAreaTextStyle:           tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorWhite),
			BorderColor:                 trueRed,
			TitleColor:                  trueRed,
			InfoColor:                   tcell.ColorWhite,
			InfoColorTwo:                tcell.ColorAntiqueWhite,
			ChatTextColor:               tcell.ColorWhite,
			ChatLabelColors:             []string{tcell.ColorGold.CSS(), tcell.ColorYellow.CSS(), tcell.ColorRed.CSS(), lightGreen.CSS(), tcell.ColorGreen.CSS()},
		}
	case "satanic":
		trueBlack := tcell.NewHexColor(0x000000)
		black := tcell.NewHexColor(0x111111)
		darkRed := tcell.NewHexColor(0x660000)
		red := tcell.NewHexColor(0xFF0000)
		mediumRed := tcell.NewHexColor(0xCC0000)

		return &Theme{
			Code:                        "satanic",
			BackgroundColor:             trueBlack,
			ForgroundColor:              red,
			HighlightColor:              red,
			AccentColor:                 black,
			AccentColorTwo:              trueBlack,
			ButtonStyle:                 tcell.StyleDefault.Background(tcell.ColorDarkRed).Foreground(tcell.NewHexColor(0x111111)),
			ActivatedButtonStyle:        tcell.StyleDefault.Background(tcell.NewHexColor(0x222222)).Foreground(red),
			DropdownListUnselectedStyle: tcell.StyleDefault.Background(tcell.ColorDarkRed).Foreground(tcell.NewHexColor(0x111111)),
			DropdownListSelectedStyle:   tcell.StyleDefault.Background(red).Foreground(black),
			TextAreaTextStyle:           tcell.StyleDefault.Background(trueBlack).Foreground(red),
			BorderColor:                 darkRed,
			TitleColor:                  red,
			InfoColor:                   mediumRed,
			InfoColorTwo:                darkRed,
			ChatTextColor:               tcell.ColorWhite,
			ChatLabelColors: []string{
				tcell.ColorYellow.CSS(),
				tcell.ColorDarkOrange.CSS(),
				tcell.ColorOrange.CSS(),
				tcell.ColorDarkGoldenrod.CSS(),
				tcell.ColorGold.CSS(),
				"#C061CB",
				tcell.ColorPink.CSS()},
		}

	default:
		return getDefault()
	}
}

func (theme Theme) ApplyGlobals() {
	tview.Styles.BorderColor = theme.BorderColor
	tview.Styles.TitleColor = theme.BorderColor
	tview.Styles.PrimaryTextColor = theme.ForgroundColor
}
