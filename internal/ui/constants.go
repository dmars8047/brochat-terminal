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
	DangerBackgroundColor  = tcell.NewHexColor(DangerBackgroundColorCode)
	BroChatYellowColor     = tcell.NewHexColor(BroChatYellowColorCode)
	DefaultBackgroundColor = tcell.NewHexColor(DefaultBackgroundColorCode)
	AccentBackgroundColor  = tcell.NewHexColor(AccentColorFourColorCode)
	ButtonStyle            = tcell.StyleDefault.Background(tcell.NewHexColor(AccentColorTwoColorCode)).Foreground(tcell.ColorWhite)
	ActivatedButtonStyle   = tcell.StyleDefault.Background(tcell.NewHexColor(0xFFC300)).Foreground(tcell.ColorBlack)
)
