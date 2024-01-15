package ui

import "github.com/rivo/tview"

type Page string

const (
	WELCOME_PAGE             Page = "auth:welcome"
	LOGIN_PAGE               Page = "auth:login"
	REGISTER_PAGE            Page = "auth:register"
	FORGOT_PW_PAGE           Page = "auth:forgotpw"
	HOME_MENU_PAGE           Page = "home:menu"
	HOME_FRIENDS_LIST_PAGE   Page = "home:friendslist"
	HOME_FRIENDS_FINDER_PAGE Page = "home:findafriend"
)

type PageNavigator struct {
	current   Page
	Pages     *tview.Pages
	callbacks map[Page]func()
}

func NewNavigator(pages *tview.Pages) *PageNavigator {
	return &PageNavigator{
		current:   WELCOME_PAGE,
		Pages:     pages,
		callbacks: make(map[Page]func()),
	}
}

func (nav *PageNavigator) Register(page Page, primitive tview.Primitive, resize, visible bool, callback func()) {
	nav.callbacks[page] = callback
	nav.Pages.AddPage(string(page), primitive, resize, visible)
}

func (nav *PageNavigator) NavigateTo(pageName Page) {
	callback, ok := nav.callbacks[pageName]

	if ok {
		callback()
	}

	nav.Pages.SwitchToPage(string(pageName))
}
