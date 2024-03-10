package ui

import (
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const WELCOME_PAGE PageSlug = "welcome"

// WelcomePage is the welcome page
type WelcomePage struct {
}

// NewWelcomePage creates a new instance of the welcome page
func NewWelcomePage() *WelcomePage {
	return &WelcomePage{}
}

// Setup configures the welcome page and registers it with the page navigator
func (page *WelcomePage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)

	grid.SetRows(4, 8, 8, 1, 1, 0).
		SetColumns(0, 31, 39, 0)

	logoBro := tview.NewTextView()
	logoBro.SetTextAlign(tview.AlignLeft).
		SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	logoBro.SetTextColor(tcell.ColorWhite)
	logoBro.SetText(
		`BBBBBBB\                      
BB  __BB\                     
BB |  BB | RRRRRR\   OOOOOO\  
BBBBBBB\ |RR  __RR\ OO  __OO\ 
BB  __BB\ RR |  \__|OO /  OO |
BB |  BB |RR |      OO |  OO |
BBBBBBB  |RR |      \OOOOOO  |
\_______/ \__|       \______/ `)

	logoChat := tview.NewTextView()
	logoChat.SetTextAlign(tview.AlignLeft)
	logoChat.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	logoChat.SetTextColor(BROCHAT_YELLOW_COLOR)
	logoChat.SetText(
		` CCCCCC\  HH\                  TT\
CC  __CC\ HH |                 TT |
CC /  \__|HHHHHHH\   AAAAAA\ TTTTTT\
CC |      HH  __HH\  \____AA\\_TT  _|
CC |      HH |  HH | AAAAAAA | TT |
CC |  CC\ HH |  HH |AA  __AA | TT |TT\
\CCCCCC  |HH |  HH |\AAAAAAA | \TTTT  |
 \______/ \__|  \__| \_______|  \____/`)

	loginButton := tview.NewButton("Login").SetSelectedFunc(func() {
		nav.NavigateTo(LOGIN_PAGE, nil)
	}).SetActivatedStyle(ACTIVATED_BUTTON_STYLE).SetStyle(DEFAULT_BUTTON_STYLE)

	registrationButton := tview.NewButton("Register").SetSelectedFunc(func() {
		nav.NavigateTo(REGISTER_PAGE, nil)
	}).SetActivatedStyle(ACTIVATED_BUTTON_STYLE).SetStyle(DEFAULT_BUTTON_STYLE)

	exitButton := tview.NewButton("Exit").SetSelectedFunc(func() {
		app.Stop()
	}).SetActivatedStyle(ACTIVATED_BUTTON_STYLE).SetStyle(DEFAULT_BUTTON_STYLE)

	buttonGrid := tview.NewGrid()
	buttonGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if loginButton.HasFocus() {
				app.SetFocus(registrationButton)
			} else if registrationButton.HasFocus() {
				app.SetFocus(exitButton)
			} else if exitButton.HasFocus() {
				app.SetFocus(loginButton)
			}
		} else if event.Key() == tcell.KeyBacktab {
			if loginButton.HasFocus() {
				app.SetFocus(exitButton)
			} else if registrationButton.HasFocus() {
				app.SetFocus(loginButton)
			} else if exitButton.HasFocus() {
				app.SetFocus(registrationButton)
			}
		}
		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvInstructions.SetText("Navigate with Tab and Shift+Tab")
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))

	tvVersionNumber := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvVersionNumber.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvVersionNumber.SetText("Version - v0.0.12")
	tvVersionNumber.SetTextColor(tcell.NewHexColor(0x777777))

	buttonGrid.SetRows(3, 1, 1).SetColumns(0, 4, 0, 4, 0)

	buttonGrid.AddItem(loginButton, 0, 0, 1, 1, 0, 0, true).
		AddItem(registrationButton, 0, 2, 1, 1, 0, 0, false).
		AddItem(exitButton, 0, 4, 1, 1, 0, 0, false).
		AddItem(tvInstructions, 2, 0, 1, 5, 0, 0, false)

	grid.AddItem(logoBro, 1, 1, 1, 1, 0, 0, false).
		AddItem(logoChat, 1, 2, 1, 1, 0, 0, false).
		AddItem(buttonGrid, 2, 1, 1, 2, 0, 0, true).
		AddItem(tvVersionNumber, 4, 1, 1, 2, 0, 0, false)

	nav.Register(WELCOME_PAGE, grid, true, true, nil, nil)
}
