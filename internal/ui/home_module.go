package ui

import (
	"github.com/dmars8047/brochat-terminal/internal/state"
	"github.com/dmars8047/idam-service/pkg/idam"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HomeModule struct {
	appState       state.ApplicationState
	userAuthClient *idam.UserAuthClient
}

func NewHomeModule(userAuthClient *idam.UserAuthClient, appState state.ApplicationState) *HomeModule {
	return &HomeModule{
		appState:       appState,
		userAuthClient: userAuthClient,
	}
}

func (mod *HomeModule) SetupMenuPage(app *tview.Application, pages *tview.Pages) {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(4, 8, 8, 0).
		SetColumns(0, 31, 39, 0)

	logoBro := tview.NewTextView()
	logoBro.SetTextAlign(tview.AlignLeft).
		SetBackgroundColor(DefaultBackgroundColor)
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
	logoChat.SetBackgroundColor(DefaultBackgroundColor)
	logoChat.SetTextColor(BroChatYellowColor)
	logoChat.SetText(
		` CCCCCC\  HH\                  TT\
CC  __CC\ HH |                 TT |
CC /  \__|HHHHHHH\   AAAAAA\ TTTTTT\
CC |      HH  __HH\  \____AA\\_TT  _|
CC |      HH |  HH | AAAAAAA | TT |
CC |  CC\ HH |  HH |AA  __AA | TT |TT\
\CCCCCC  |HH |  HH |\AAAAAAA | \TTTT  |
 \______/ \__|  \__| \_______|  \____/`)

	list := tview.NewList().
		AddItem("Bros", "", 'b', func() {}).
		AddItem("Chats", "", 'c', nil).
		AddItem("Settings", "", 's', nil).
		AddItem("Logout", "", 'q', func() {
			session, ok := state.Get[state.UserSession](mod.appState, state.UserSessionProp)

			if !ok {
				AlertFatal(app, pages, "home:menu:alert:err", "Application State Error - Could not get user session.")
			}

			err := mod.userAuthClient.Logout(session.Auth.AccessToken)

			if err != nil {
				AlertFatal(app, pages, "home:menu:alert:err", err.Error())
				return
			}

			pages.SwitchToPage("auth:login")
		})

	list.SetBackgroundColor(DefaultBackgroundColor)

	grid.AddItem(logoBro, 1, 1, 1, 1, 0, 0, false).
		AddItem(logoChat, 1, 2, 1, 1, 0, 0, false).
		AddItem(list, 2, 1, 1, 2, 0, 0, true)

	pages.AddPage("home:menu", grid, true, false)
}
