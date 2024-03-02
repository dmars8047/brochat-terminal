package ui

import "github.com/gdamore/tcell/v2"

const (
	BroChatYellowColorCode     = 0xFFC300
	AccentColorTwoColorCode    = 0x222222
	AccentColorFourColorCode   = 0x444444
	DefaultBackgroundColorCode = 0x111111
	DangerBackgroundColorCode  = 0xFF0E0E

	DEFAULT_AUTH_TOKEN_TYPE = "Bearer"
)

var (
	DangerBackgroundColor    = tcell.NewHexColor(DangerBackgroundColorCode)
	BROCHAT_YELLOW_COLOR     = tcell.NewHexColor(BroChatYellowColorCode)
	DEFAULT_BACKGROUND_COLOR = tcell.NewHexColor(DefaultBackgroundColorCode)
	ACCENT_BACKGROUND_COLOR  = tcell.NewHexColor(AccentColorFourColorCode)
	ACCENT_COLOR_TWO_COLOR   = tcell.NewHexColor(AccentColorTwoColorCode)
	DEFAULT_BUTTON_STYLE     = tcell.StyleDefault.Background(tcell.NewHexColor(AccentColorTwoColorCode)).Foreground(tcell.ColorWhite)
	ACTIVATED_BUTTON_STYLE   = tcell.StyleDefault.Background(tcell.NewHexColor(0xFFC300)).Foreground(tcell.ColorBlack)
)
