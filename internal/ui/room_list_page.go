package ui

import (
	"context"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const ROOM_LIST_PAGE PageSlug = "roomlist"

const (
	ROOM_LIST_PAGE_ALERT_INFO = "home:roomlist:alert:info"
	ROOM_LIST_PAGE_ALERT_ERR  = "home:roomlist:alert:err"
)

type RoomListPage struct {
	brochatClient *chat.BroChatUserClient
	feedClient    *state.FeedClient
	table         *tview.Table
	userRooms     map[int]chat.Room
}

func NewRoomListPage(brochatClient *chat.BroChatUserClient, feedClient *state.FeedClient) *RoomListPage {
	return &RoomListPage{
		brochatClient: brochatClient,
		feedClient:    feedClient,
		table:         tview.NewTable(),
		userRooms:     make(map[int]chat.Room, 0),
	}
}

func (page *RoomListPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvHeader.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvHeader.SetText("Room List")

	page.table.SetBorders(true)
	page.table.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	page.table.SetFixed(1, 1)
	page.table.SetSelectable(true, false)

	page.table.SetSelectedFunc(func(row int, _ int) {
		room, ok := page.userRooms[row]

		if !ok {
			return
		}

		nav.NavigateTo(CHAT_PAGE, ChatPageParameters{
			channel_id: room.ChannelId,
			title:      room.Name,
		})
	})

	page.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'f':
				nav.NavigateTo(ROOM_FINDER_PAGE, nil)
				page.userRooms = make(map[int]chat.Room, 0)
				page.table.Clear()
			case 'n':
				nav.NavigateTo(ROOM_EDITOR_PAGE, nil)
				page.userRooms = make(map[int]chat.Room, 0)
				page.table.Clear()
			}
		} else if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(HOME_PAGE, nil)
			page.userRooms = make(map[int]chat.Room, 0)
			page.table.Clear()
		}

		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvInstructions.SetText("(n) Create a Room - (f) Find a Room - (esc) Quit")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)

	grid.SetRows(2, 1, 1, 0, 1, 1, 2)
	grid.SetColumns(0, 76, 0)

	grid.AddItem(tvHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(page.table, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 5, 1, 1, 1, 0, 0, false)

	pageContext, cancel := context.WithCancel(appContext.Context)

	nav.Register(ROOM_LIST_PAGE, grid, true, false,
		func(_ interface{}) {
			pageContext, cancel = context.WithCancel(appContext.Context)
			page.onPageLoad(app, appContext, pageContext)
		},
		func() {
			page.onPageClose()
			cancel()
		})
}

func (page *RoomListPage) onPageLoad(app *tview.Application, appContext *state.ApplicationContext, pageContext context.Context) {
	page.populateTable(appContext.GetBrochatUser())

	// Create a go routine to monitor for changes to the user's rooms via a user profile update event
	go func() {
		subId, channel := page.feedClient.SubscribeToUserProfileUpdates()
		defer page.feedClient.UnsubscribeFromUserProfileUpdates(subId)

		for {
			select {
			case <-pageContext.Done():
				return
			case eventCode := <-channel:
				if eventCode == chat.USER_PROFILE_UPDATE_CODE_ROOM_UPDATE {
					app.QueueUpdateDraw(func() {
						page.populateTable(appContext.GetBrochatUser())
					})
				}
			}
		}
	}()
}

func (page *RoomListPage) onPageClose() {
	page.userRooms = make(map[int]chat.Room, 0)
	page.table.Clear()
}

func (page *RoomListPage) populateTable(brochatUser chat.User) {
	page.table.SetCell(0, 0, tview.NewTableCell("Name").
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignCenter).
		SetExpansion(1).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	page.table.SetCell(0, 1, tview.NewTableCell("Owner").
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignCenter).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	for i, rel := range brochatUser.Rooms {
		row := i + 1

		page.table.SetCell(row, 0, tview.NewTableCell(rel.Name).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
		page.table.SetCell(row, 1, tview.NewTableCell(rel.Owner.Username).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))

		page.userRooms[row] = rel
	}
}
