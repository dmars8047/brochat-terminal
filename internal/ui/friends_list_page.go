package ui

import (
	"context"
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
	brochatClient  *chat.BroChatClient
	feedClient     *state.FeedClient
	table          *tview.Table
	tvInstructions *tview.TextView
	userFriends    map[uint8]chat.UserRelationship
}

func NewFriendsListPage(brochatClient *chat.BroChatClient, feedClient *state.FeedClient) *FriendsListPage {
	return &FriendsListPage{
		brochatClient:  brochatClient,
		feedClient:     feedClient,
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
			returnPage: FRIENDS_LIST_PAGE,
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

	var pageContext context.Context
	var cancel context.CancelFunc

	nav.Register(FRIENDS_LIST_PAGE, grid, true, false,
		func(_ interface{}) {
			pageContext, cancel = appContext.GenerateUserSessionBoundContextWithCancel()
			page.onPageLoad(app, appContext, pageContext)
		},
		func() {
			cancel()
			page.onPageClose()
		})
}

func (page *FriendsListPage) onPageLoad(app *tview.Application, appContext *state.ApplicationContext, pageContext context.Context) {
	page.populateTable(appContext.GetBrochatUser())

	// Create a goroutine to listen for updates to the user's relationships
	// If one is recieved then redraw the table
	go func() {
		subId, userProfileUpdatesChannel := page.feedClient.SubscribeToUserProfileUpdates()

		defer page.feedClient.UnsubscribeFromUserProfileUpdates(subId)

		for {
			select {
			case <-pageContext.Done():
				return
			case updateCode := <-userProfileUpdatesChannel:
				if updateCode == chat.USER_PROFILE_UPDATE_REASON_RELATIONSHIP_UPDATE {
					page.table.Clear()
					app.QueueUpdateDraw(func() {
						page.populateTable(appContext.GetBrochatUser())
					})
				}
			}
		}
	}()

}

func (page *FriendsListPage) onPageClose() {
	page.userFriends = make(map[uint8]chat.UserRelationship, 0)
	page.table.Clear()
}

func (page *FriendsListPage) populateTable(brochatUser chat.User) {
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

	countOfPendingFriendRequests := 0

	for _, rel := range brochatUser.Relationships {
		if rel.Type&chat.RELATIONSHIP_TYPE_FRIEND_REQUEST_RECIEVED != 0 {
			countOfPendingFriendRequests++
		}
	}

	page.tvInstructions.SetText(fmt.Sprintf("(f) Find a new Bro - (p) View Pending [%d] - (esc) Quit", countOfPendingFriendRequests))

	for i, rel := range brochatUser.Relationships {
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
