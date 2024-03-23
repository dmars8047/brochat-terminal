package ui

import (
	"context"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/broterm/internal/theme"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const ROOM_LIST_PAGE PageSlug = "roomlist"

const (
	ROOM_LIST_PAGE_ALERT_INFO = "home:roomlist:alert:info"
	ROOM_LIST_PAGE_ALERT_ERR  = "home:roomlist:alert:err"
)

type RoomListPage struct {
	brochatClient    *chat.BroChatClient
	feedClient       *state.FeedClient
	table            *tview.Table
	userRooms        map[int]chat.Room
	currentThemeCode string
}

func NewRoomListPage(brochatClient *chat.BroChatClient, feedClient *state.FeedClient) *RoomListPage {
	return &RoomListPage{
		brochatClient:    brochatClient,
		feedClient:       feedClient,
		table:            tview.NewTable(),
		userRooms:        make(map[int]chat.Room, 0),
		currentThemeCode: "NOT_SET",
	}
}

func (page *RoomListPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetText("Room List")

	page.table.SetBorders(true)
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
			returnPage: ROOM_LIST_PAGE,
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
	tvInstructions.SetText("(n) Create a Room - (f) Find a Room - (esc) Quit")

	grid := tview.NewGrid()

	grid.SetRows(2, 1, 1, 0, 1, 1, 2)
	grid.SetColumns(0, 76, 0)

	grid.AddItem(tvHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(page.table, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 5, 1, 1, 1, 0, 0, false)

	var pageContext context.Context
	var cancel context.CancelFunc

	applyTheme := func() {
		theme := appContext.GetTheme()

		if page.currentThemeCode != theme.Code {
			page.currentThemeCode = theme.Code
			grid.SetBackgroundColor(theme.BackgroundColor)
			page.table.SetBordersColor(theme.BorderColor)
			page.table.SetBorderColor(theme.BorderColor)
			page.table.SetTitleColor(theme.TitleColor)
			page.table.SetBackgroundColor(theme.BackgroundColor)
			page.table.SetSelectedStyle(theme.DropdownListSelectedStyle)
			tvHeader.SetBackgroundColor(theme.BackgroundColor)
			tvHeader.SetTextColor(theme.TitleColor)
			tvInstructions.SetBackgroundColor(theme.BackgroundColor)
			tvInstructions.SetTextColor(theme.InfoColor)
		}
	}

	applyTheme()

	nav.Register(ROOM_LIST_PAGE, grid, true, false,
		func(_ interface{}) {
			applyTheme()
			pageContext, cancel = appContext.GenerateUserSessionBoundContextWithCancel()
			page.onPageLoad(app, appContext, pageContext)
		},
		func() {
			cancel()
			page.onPageClose()
		})
}

func (page *RoomListPage) onPageLoad(app *tview.Application, appContext *state.ApplicationContext, pageContext context.Context) {
	page.populateTable(appContext.GetBrochatUser(), appContext.GetTheme())

	// Create a go routine to monitor for changes to the user's rooms via a user profile update event
	go func() {
		subId, userUpdatedChannel := page.feedClient.SubscribeToUserProfileUpdates()
		defer page.feedClient.UnsubscribeFromUserProfileUpdates(subId)

		for {
			select {
			case <-pageContext.Done():
				return
			case eventCode := <-userUpdatedChannel:
				if eventCode == chat.USER_PROFILE_UPDATE_CODE_ROOM_UPDATE {
					app.QueueUpdateDraw(func() {
						page.populateTable(appContext.GetBrochatUser(), appContext.GetTheme())
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

func (page *RoomListPage) populateTable(brochatUser chat.User, thm theme.Theme) {
	page.table.SetCell(0, 0, tview.NewTableCell("Name").
		SetTextColor(thm.ForgroundColor).
		SetAlign(tview.AlignCenter).
		SetExpansion(1).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	page.table.SetCell(0, 1, tview.NewTableCell("Owner").
		SetTextColor(thm.ForgroundColor).
		SetAlign(tview.AlignCenter).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	for i, rel := range brochatUser.Rooms {
		row := i + 1

		page.table.SetCell(row, 0, tview.NewTableCell(rel.Name).SetTextColor(thm.ForgroundColor).SetAlign(tview.AlignCenter))
		page.table.SetCell(row, 1, tview.NewTableCell(rel.Owner.Username).SetTextColor(thm.ForgroundColor).SetAlign(tview.AlignCenter))

		page.userRooms[row] = rel
	}
}
