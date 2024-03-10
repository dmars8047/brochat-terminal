package ui

import (
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
	table         *tview.Table
	userRooms     map[int]chat.Room
}

func NewRoomListPage(brochatClient *chat.BroChatUserClient) *RoomListPage {
	return &RoomListPage{
		brochatClient: brochatClient,
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

	nav.Register(ROOM_LIST_PAGE, grid, true, false,
		func(_ interface{}) {
			page.onPageLoad(appContext)
		},
		func() {
			page.onPageClose()
		})
}

func (page *RoomListPage) onPageLoad(appContext *state.ApplicationContext) {
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

	for i, rel := range appContext.BrochatUser.Rooms {
		row := i + 1

		page.table.SetCell(row, 0, tview.NewTableCell(rel.Name).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
		page.table.SetCell(row, 1, tview.NewTableCell(rel.Owner.Username).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))

		page.userRooms[row] = rel
	}
}

func (page *RoomListPage) onPageClose() {
	page.userRooms = make(map[int]chat.Room, 0)
	page.table.Clear()
}
