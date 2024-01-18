package ui

import "github.com/rivo/tview"

type Page string

const (
	WELCOME_PAGE               Page = "auth:welcome"
	LOGIN_PAGE                 Page = "auth:login"
	REGISTER_PAGE              Page = "auth:register"
	FORGOT_PW_PAGE             Page = "auth:forgotpw"
	HOME_MENU_PAGE             Page = "home:menu"
	HOME_FRIENDS_LIST_PAGE     Page = "home:friendslist"
	HOME_FRIENDS_FINDER_PAGE   Page = "home:findafriend"
	HOME_PENDING_REQUESTS_PAGE Page = "home:pendingrequests"
	HOME_CHAT_PAGE             Page = "home:chat"
)

type PageNavigator struct {
	current    Page
	Pages      *tview.Pages
	openFuncs  map[Page]func(interface{})
	closeFuncs map[Page]func()
}

func NewNavigator(pages *tview.Pages) *PageNavigator {
	return &PageNavigator{
		current:    WELCOME_PAGE,
		Pages:      pages,
		openFuncs:  make(map[Page]func(interface{})),
		closeFuncs: make(map[Page]func()),
	}
}

func (nav *PageNavigator) Register(page Page,
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

func (nav *PageNavigator) NavigateTo(pageName Page, param interface{}) {
	close, ok := nav.closeFuncs[nav.current]

	if ok {
		close()
	}

	nav.current = pageName

	open, ok := nav.openFuncs[pageName]

	if ok {
		open(param)
	}

	nav.Pages.SwitchToPage(string(pageName))
}
