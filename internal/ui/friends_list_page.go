package ui

import (
	"fmt"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const FRIENDS_LIST_PAGE PageSlug = "friends_list"

const (
	FRIENDS_LIST_PAGE_ALERT_INFO = "home:friendlist:alert:info"
	FRIENDS_LIST_PAGE_ALERT_ERR  = "home:friendlist:alert:err"
)

type FriendsListPage struct {
	brochatClient  *chat.BroChatUserClient
	table          *tview.Table
	tvInstructions *tview.TextView
	userFriends    map[uint8]chat.UserRelationship
}

func NewFriendsListPage(brochatClient *chat.BroChatUserClient) *FriendsListPage {
	return &FriendsListPage{
		brochatClient:  brochatClient,
		table:          tview.NewTable(),
		tvInstructions: tview.NewTextView(),
		userFriends:    make(map[uint8]chat.UserRelationship, 0),
	}
}

func (page *FriendsListPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvHeader.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvHeader.SetText("Friends List")

	page.table.SetBorders(true)
	page.table.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	page.table.SetFixed(1, 1)
	page.table.SetSelectable(true, false)

	page.table.SetSelectedFunc(func(row int, _ int) {
		rel, ok := page.userFriends[uint8(row)]

		if !ok {
			return
		}

		nav.NavigateTo(CHAT_PAGE, ChatPageParameters{
			channel_id: rel.DirectMessageChannelId,
		})
	})

	page.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'p':
				nav.NavigateTo(ACCEPT_FRIEND_REQUEST_PAGE, nil)
				page.userFriends = make(map[uint8]chat.UserRelationship, 0)
				page.table.Clear()
			case 'f':
				nav.NavigateTo(FRIENDS_FINDER_PAGE, nil)
				page.userFriends = make(map[uint8]chat.UserRelationship, 0)
				page.table.Clear()
			}
		} else if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(HOME_PAGE, nil)
			page.userFriends = make(map[uint8]chat.UserRelationship, 0)
			page.table.Clear()
		}

		return event
	})

	page.tvInstructions.SetTextAlign(tview.AlignCenter)
	page.tvInstructions.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	page.tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	page.tvInstructions.SetText("(f) Find a new Bro - (p) View Pending - (esc) Quit")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)

	grid.SetRows(2, 1, 1, 0, 1, 1, 2)
	grid.SetColumns(0, 76, 0)

	grid.AddItem(tvHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(page.table, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(page.tvInstructions, 5, 1, 1, 1, 0, 0, false)

	nav.Register(FRIENDS_LIST_PAGE, grid, true, false,
		func(_ interface{}) {
			page.onPageLoad(app, appContext, nav)
		},
		func() {
			page.onPageClose()
		})
}

func (page *FriendsListPage) onPageLoad(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	page.table.SetCell(0, 0, tview.NewTableCell("Username").
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignCenter).
		SetExpansion(1).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	page.table.SetCell(0, 1, tview.NewTableCell("Status").
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignCenter).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	page.table.SetCell(0, 2, tview.NewTableCell("Last Active").
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignRight).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	usr, err := page.brochatClient.GetUser(&chat.AuthInfo{
		AccessToken: appContext.UserSession.Auth.AccessToken,
		TokenType:   DEFAULT_AUTH_TOKEN_TYPE,
	}, appContext.UserSession.Info.Id)

	if err != nil {
		nav.AlertFatal(app, FRIENDS_LIST_PAGE_ALERT_ERR, err.Error())
		return
	}

	appContext.BrochatUser = usr

	countOfPendingFriendRequests := 0

	for _, rel := range usr.Relationships {
		if rel.Type&chat.RELATIONSHIP_TYPE_FRIEND_REQUEST_RECIEVED != 0 {
			countOfPendingFriendRequests++
		}
	}

	page.tvInstructions.SetText(fmt.Sprintf("(f) Find a new Bro - (p) View Pending [%d] - (esc) Quit", countOfPendingFriendRequests))

	for i, rel := range usr.Relationships {
		row := i + 1

		if rel.Type != chat.RELATIONSHIP_TYPE_FRIEND {
			continue
		}

		page.table.SetCell(row, 0, tview.NewTableCell(rel.Username).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
		if rel.IsOnline {
			page.table.SetCell(row, 1, tview.NewTableCell("Online").SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))
		} else {
			page.table.SetCell(row, 1, tview.NewTableCell("Offline").SetTextColor(tcell.ColorGray).SetAlign(tview.AlignCenter))
		}

		var dateString string = rel.LastOnlineUtc.Local().Format("Jan 2, 2006")

		page.table.SetCell(row, 2, tview.NewTableCell(dateString).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight))

		page.userFriends[uint8(row)] = rel
	}
}

func (page *FriendsListPage) onPageClose() {
	page.userFriends = make(map[uint8]chat.UserRelationship, 0)
	page.table.Clear()
}
