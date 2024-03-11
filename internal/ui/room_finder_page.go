package ui

import (
	"fmt"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const ROOM_FINDER_PAGE PageSlug = "room_finder"

const (
	ROOM_FINDER_PAGE_ALERT_INFO = "home:roomfinder:alert:info"
	ROOM_FINDER_PAGE_ALERT_ERR  = "home:roomfinder:alert:err"
	ROOM_FINDER_PAGE_CONFIRM    = "home:roomfinder:confirm"
)

// RoomFinderPage is the room finder page
type RoomFinderPage struct {
	brochatClient *chat.BroChatUserClient
	table         *tview.Table
	publicRooms   map[int]chat.Room
}

// NewRoomFinderPage creates a new room finder page
func NewRoomFinderPage(brochatClient *chat.BroChatUserClient) *RoomFinderPage {
	return &RoomFinderPage{
		brochatClient: brochatClient,
		table:         tview.NewTable(),
		publicRooms:   make(map[int]chat.Room, 0),
	}
}

// Setup sets up the room finder page and registers it with the page navigator
func (page *RoomFinderPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvHeader.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvHeader.SetText("Find Rooms")

	page.table.SetBorders(true)
	page.table.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	page.table.SetFixed(1, 1)
	page.table.SetSelectable(true, false)

	page.table.SetSelectedFunc(func(row int, _ int) {
		room, ok := page.publicRooms[row]

		if !ok {
			return
		}

		nav.Confirm(ROOM_FINDER_PAGE_CONFIRM, fmt.Sprintf("Join %s?", room.Name), func() {
			joinRoomErr := page.brochatClient.JoinRoom(appContext.GetAuthInfo(), room.Id)

			if joinRoomErr != nil {
				nav.Alert(ROOM_FINDER_PAGE_ALERT_ERR, fmt.Sprintf("An error occurred while joining room: %s", joinRoomErr.Error()))
				return
			}

			// Add the room to the user's rooms
			appContext.BrochatUser.Rooms = append(appContext.BrochatUser.Rooms, room)

			nav.AlertWithDoneFunc(ROOM_FINDER_PAGE_ALERT_INFO, fmt.Sprintf("You have successfuly joined the room '%s'.", room.Name), func(buttonIndex int, buttonLabel string) {
				nav.Pages.HidePage(ROOM_FINDER_PAGE_ALERT_INFO).RemovePage(ROOM_FINDER_PAGE_ALERT_INFO)
				nav.NavigateTo(CHAT_PAGE, ChatPageParameters{
					channel_id: room.ChannelId,
				})
			})
		})
	})

	page.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(HOME_PAGE, nil)
			page.publicRooms = make(map[int]chat.Room, 0)
			page.table.Clear()
		}

		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvInstructions.SetText("(enter) Join room -(esc) Quit")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)

	grid.SetRows(2, 1, 1, 0, 1, 1, 2)
	grid.SetColumns(0, 76, 0)

	grid.AddItem(tvHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(page.table, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 5, 1, 1, 1, 0, 0, false)

	nav.Register(ROOM_FINDER_PAGE, grid, true, false,
		func(_ interface{}) {
			page.onPageLoad(appContext, nav)
		},
		func() {
			page.onPageClose()
		})
}

// onPageLoad is called when the room finder page is navigated to
func (page *RoomFinderPage) onPageLoad(appContext *state.ApplicationContext, nav *PageNavigator) {
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

	rooms, err := page.brochatClient.GetRooms(appContext.GetAuthInfo())

	if err != nil {
		nav.Alert(ROOM_FINDER_PAGE_ALERT_ERR, fmt.Sprintf("An error occurred while retrieving public rooms: %s", err.Error()))
		return
	}

	for i, rel := range rooms {
		row := i + 1

		page.table.SetCell(row, 0, tview.NewTableCell(rel.Name).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
		page.table.SetCell(row, 1, tview.NewTableCell(rel.Owner.Username).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))

		page.publicRooms[row] = rel
	}
}

// onPageClose is called when the room finder page is navigated away from
func (page *RoomFinderPage) onPageClose() {
	page.publicRooms = make(map[int]chat.Room, 0)
	page.table.Clear()
}
