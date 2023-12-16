package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AuthModule struct {
}

func NewAuthModule() *AuthModule {
	mod := AuthModule{}

	return &mod
}

func (mod *AuthModule) SetupAuthPages(app *tview.Application, pages *tview.Pages) {
	mod.setupWelcomePage(app, pages)
	mod.setupLoginPage(app, pages)
	mod.setupRegistrationPage(app, pages)
}

func (mod *AuthModule) setupWelcomePage(app *tview.Application, pages *tview.Pages) {
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

	loginButton := tview.NewButton("Login").SetSelectedFunc(func() {
		pages.SwitchToPage("auth:login")
	}).SetActivatedStyle(ActivatedButtonStyle).SetStyle(ButtonStyle)

	registrationButton := tview.NewButton("Register").SetSelectedFunc(func() {
		pages.SwitchToPage("auth:registration")
	}).SetActivatedStyle(ActivatedButtonStyle).SetStyle(ButtonStyle)

	exitButton := tview.NewButton("Exit").SetSelectedFunc(func() {
		app.Stop()
	}).SetActivatedStyle(ActivatedButtonStyle).SetStyle(ButtonStyle)

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
	tvInstructions.SetBackgroundColor(DefaultBackgroundColor)
	tvInstructions.SetText("Navigate with Tab and Shift+Tab")
	tvInstructions.SetTextColor(tcell.NewHexColor(0x444444))

	buttonGrid.SetRows(3, 1, 1).SetColumns(0, 4, 0, 4, 0)

	buttonGrid.AddItem(loginButton, 0, 0, 1, 1, 0, 0, true).
		AddItem(registrationButton, 0, 2, 1, 1, 0, 0, false).
		AddItem(exitButton, 0, 4, 1, 1, 0, 0, false).
		AddItem(tvInstructions, 2, 0, 1, 5, 0, 0, false)

	// grid.SetRows(4, 0, 5, 0).
	// SetColumns(0, 31, 39, 0)

	grid.AddItem(logoBro, 1, 1, 1, 1, 0, 0, false).
		AddItem(logoChat, 1, 2, 1, 1, 0, 0, false).
		AddItem(buttonGrid, 2, 1, 1, 2, 0, 0, true)

	pages.AddPage("auth:welcome", grid, true, true)
}

func (mod *AuthModule) setupLoginPage(app *tview.Application, pages *tview.Pages) {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(4, 0, 10)
	grid.SetColumns(0, 70, 0)

	loginForm := tview.NewForm()
	loginForm.SetBackgroundColor(AccentBackgroundColor)
	loginForm.SetFieldBackgroundColor(AccentColorTwoColorCode)
	loginForm.SetLabelColor(BroChatYellowColor)
	loginForm.SetBorder(true).SetTitle(" BroChat - Login ").SetTitleAlign(tview.AlignCenter)
	loginForm.SetButtonStyle(ButtonStyle)
	loginForm.SetButtonActivatedStyle(ActivatedButtonStyle)
	loginForm.AddInputField("Email", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddButton("Login", func() {
			fmt.Println("Login pressed!")
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("auth:welcome")
		})

	grid.AddItem(loginForm, 1, 1, 1, 1, 0, 0, true)

	pages.AddPage("auth:login", grid, true, false)
}

func (mod *AuthModule) setupRegistrationPage(app *tview.Application, pages *tview.Pages) {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(4, 0, 6)
	grid.SetColumns(0, 70, 0)

	registrationForm := tview.NewForm()
	registrationForm.SetBorder(true).SetTitle(" BroChat - Register ").SetTitleAlign(tview.AlignCenter)
	registrationForm.SetBackgroundColor(AccentBackgroundColor)
	registrationForm.SetFieldBackgroundColor(AccentColorTwoColorCode)
	registrationForm.SetButtonStyle(ButtonStyle)
	registrationForm.SetButtonActivatedStyle(ActivatedButtonStyle)

	registrationForm.AddInputField("Username", "", 0, nil, nil).
		AddInputField("Email", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddPasswordField("Confirm Password", "", 0, '*', nil).
		AddButton("Register", func() {}).
		AddButton("Back", func() { pages.SwitchToPage("auth:welcome") })

	grid.AddItem(registrationForm, 1, 1, 1, 1, 0, 0, true)

	pages.AddPage("auth:registration", grid, true, false)
}
