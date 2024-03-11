package ui

import "github.com/gdamore/tcell/v2"

const (
	BroChatYellowColorCode     = 0xFFC300
	AccentColorTwoColorCode    = 0x222222
	AccentColorFourColorCode   = 0x444444
	DefaultBackgroundColorCode = 0x111111
	DangerBackgroundColorCode  = 0xFF0E0E
)

var (
	DangerBackgroundColor    = tcell.NewHexColor(DangerBackgroundColorCode)
	BROCHAT_YELLOW_COLOR     = tcell.NewHexColor(BroChatYellowColorCode)
	DEFAULT_BACKGROUND_COLOR = tcell.NewHexColor(DefaultBackgroundColorCode)
	ACCENT_BACKGROUND_COLOR  = tcell.NewHexColor(AccentColorFourColorCode)
	ACCENT_COLOR_TWO_COLOR   = tcell.NewHexColor(AccentColorTwoColorCode)
	DEFAULT_BUTTON_STYLE     = tcell.StyleDefault.Background(tcell.NewHexColor(AccentColorTwoColorCode)).Foreground(tcell.ColorWhite)
	ACTIVATED_BUTTON_STYLE   = tcell.StyleDefault.Background(tcell.NewHexColor(0xFFC300)).Foreground(tcell.ColorBlack)
	TEXT_AREA_STYLE          = tcell.StyleDefault.Reverse(true).Background(tcell.NewHexColor(0x222222)).Foreground(tcell.ColorWhite)
)

// Chat colors
const (
	CHAT_COLOR_ONE   = "#33DA7A"
	CHAT_COLOR_TWO   = "#C061CB"
	CHAT_COLOR_THREE = "#900C3F"
	CHAT_COLOR_FOUR  = "#FF5733"
	CHAT_COLOR_FIVE  = "#3498DB"
	CHAT_COLOR_SIX   = "#117a65"
	CHAT_COLOR_SEVEN = "#F0B27A"
	CHAT_COLOR_EIGHT = "#ABB2B9"
)
