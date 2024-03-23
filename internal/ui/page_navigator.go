package ui

import (
	"fmt"
	"strings"

	"github.com/dmars8047/broterm/internal/state"
	"github.com/rivo/tview"
)

type PageSlug string

// PageNavigator is a page navigator
type PageNavigator struct {
	current    PageSlug
	Pages      *tview.Pages
	appContext *state.ApplicationContext
	openFuncs  map[PageSlug]func(interface{})
	closeFuncs map[PageSlug]func()
}

// NewNavigator creates a new page navigator
func NewNavigator(appContext *state.ApplicationContext) *PageNavigator {
	pages := tview.NewPages()

	return &PageNavigator{
		appContext: appContext,
		current:    WELCOME_PAGE,
		Pages:      pages,
		openFuncs:  make(map[PageSlug]func(interface{})),
		closeFuncs: make(map[PageSlug]func()),
	}
}

// Register registers a page with the page navigator
func (nav *PageNavigator) Register(page PageSlug,
	primitive tview.Primitive,
	resize, visible bool,
	openFunc func(interface{}),
	closeFunc func()) {

	nav.Pages.AddPage(string(page), primitive, resize, visible)

	if openFunc != nil {
		nav.openFuncs[page] = openFunc
	}

	if closeFunc != nil {
		nav.closeFuncs[page] = closeFunc
	}
}

// NavigateTo navigates to a page
func (nav *PageNavigator) NavigateTo(pageName PageSlug, param interface{}) {
	close, ok := nav.closeFuncs[nav.current]

	if ok {
		close()
	}

	open, ok := nav.openFuncs[pageName]

	if ok {
		open(param)
	}

	nav.Pages.SwitchToPage(string(pageName))

	nav.current = pageName
}

// Confirm creates a confirmation modal
func (nav *PageNavigator) Confirm(id string, massage string, yesFunc func()) *tview.Pages {
	theme := nav.appContext.GetTheme()

	modal := tview.NewModal().
		SetText(massage).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				yesFunc()
			}
			nav.Pages.HidePage(id).RemovePage(id)
		})

	modal.SetBackgroundColor(theme.BackgroundColor)
	modal.SetTextColor(theme.ForgroundColor)
	modal.SetButtonStyle(theme.ButtonStyle)
	modal.SetButtonActivatedStyle(theme.ActivatedButtonStyle)
	modal.SetBorderColor(theme.BorderColor)
	modal.SetBorderStyle(theme.TextAreaTextStyle)
	modal.SetTitleColor(theme.TitleColor)

	return nav.Pages.AddPage(
		id,
		modal,
		false,
		true,
	)
}

// Alert creates an alert modal
func (nav *PageNavigator) Alert(id string, message string) *tview.Pages {
	theme := nav.appContext.GetTheme()

	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			nav.Pages.HidePage(id).RemovePage(id)
		})

	modal.SetBackgroundColor(theme.BackgroundColor)
	modal.SetTextColor(theme.ForgroundColor)
	modal.SetButtonStyle(theme.ButtonStyle)
	modal.SetButtonActivatedStyle(theme.ActivatedButtonStyle)
	modal.SetBorderColor(theme.BorderColor)
	modal.SetBorderStyle(theme.TextAreaTextStyle)
	modal.SetTitleColor(theme.TitleColor)

	return nav.Pages.AddPage(
		id,
		modal,
		false,
		true,
	)
}

// AlertWithDoneFunc creates an alert modal with a done function
func (nav *PageNavigator) AlertWithDoneFunc(id string, message string, doneFunc func(buttonIndex int, buttonLabel string)) *tview.Pages {
	theme := nav.appContext.GetTheme()

	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Close"}).
		SetDoneFunc(doneFunc)

	modal.SetBackgroundColor(theme.BackgroundColor)
	modal.SetTextColor(theme.ForgroundColor)
	modal.SetButtonStyle(theme.ButtonStyle)
	modal.SetButtonActivatedStyle(theme.ActivatedButtonStyle)
	modal.SetBorderColor(theme.BorderColor)
	modal.SetBorderStyle(theme.TextAreaTextStyle)
	modal.SetTitleColor(theme.TitleColor)

	return nav.Pages.AddPage(
		id,
		modal,
		false,
		true,
	)
}

// AlertFatal creates a fatal alert modal
func (nav *PageNavigator) AlertFatal(app *tview.Application, id string, message string) *tview.Pages {
	theme := nav.appContext.GetTheme()

	modal := tview.NewModal().
		SetText("Fatal Error: " + message).
		AddButtons([]string{"Exit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.Stop()
		})

	modal.SetBackgroundColor(theme.BackgroundColor)
	modal.SetTextColor(theme.ForgroundColor)
	modal.SetButtonStyle(theme.ButtonStyle)
	modal.SetButtonActivatedStyle(theme.ActivatedButtonStyle)
	modal.SetBorderColor(theme.BorderColor)
	modal.SetBorderStyle(theme.TextAreaTextStyle)
	modal.SetTitleColor(theme.TitleColor)

	return nav.Pages.AddPage(
		id,
		modal,
		false,
		true,
	)
}

// AlertErrors creates an alert modal with a list of errors
func (nav *PageNavigator) AlertErrors(id, errMessage string, messages []string) {
	added := false

	for _, message := range messages {
		if len(message) > 2 {
			if !added {
				errMessage += "\n"
				added = true
			}
			val := strings.ToUpper(string(message[0])) + message[1:]
			errMessage += fmt.Sprintf("\n- %s", val)
		}
	}

	nav.Alert(id, errMessage)
}
