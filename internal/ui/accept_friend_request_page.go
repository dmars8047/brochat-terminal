package ui

import (
	"fmt"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const ACCEPT_FRIEND_REQUEST_PAGE PageSlug = "accept_friend_request"

// AcceptFriendRequestPage is the page for accepting friend requests
type AcceptFriendRequestPage struct {
	brochatClient       *chat.BroChatUserClient
	userPendingRequests map[uint8]chat.UserRelationship
	table               *tview.Table
}

// NewAcceptFriendRequestPage creates a new accept friend request page
func NewAcceptFriendRequestPage(brochatClient *chat.BroChatUserClient) *AcceptFriendRequestPage {
	return &AcceptFriendRequestPage{
		brochatClient:       brochatClient,
		userPendingRequests: make(map[uint8]chat.UserRelationship, 0),
		table:               tview.NewTable(),
	}
}

// Setup sets up the accept friend request page and registers it with the page navigator
func (page *AcceptFriendRequestPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvHeader.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvHeader.SetText("Pending Friend Requests")

	page.table.SetBorders(true)
	page.table.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	page.table.SetFixed(1, 1)
	page.table.SetSelectable(true, false)

	page.table.SetSelectedFunc(func(row int, _ int) {
		selectedUser, ok := page.userPendingRequests[uint8(row)]

		if !ok {
			return
		}

		nav.Confirm(FIND_A_FRIEND_PAGE_CONFIRM, fmt.Sprintf("Accept Friend Request from %s?", selectedUser.Username), func() {
			err := page.brochatClient.AcceptFriendRequest(appContext.GetAuthInfo(), &chat.AcceptFriendRequestRequest{
				InitiatingUserId: selectedUser.UserId,
			})

			if err != nil {
				if err.Error() == "user not found or friend request not found" {
					nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("Friend Request from %s Not Found", selectedUser.Username))
					return
				} else if err.Error() == "bad request" {
					nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, "Friend Request Acceptance Not Processable")
					return
				}

				nav.AlertFatal(app, FIND_A_FRIEND_PAGE_ALERT_ERR, err.Error())
				return
			}

			// Update the relationship type
			for i := 0; i < len(appContext.BrochatUser.Relationships); i++ {
				if selectedUser.UserId == appContext.BrochatUser.Relationships[i].UserId {
					appContext.BrochatUser.Relationships[i].Type = chat.RELATIONSHIP_TYPE_FRIEND
					break
				}
			}

			page.table.RemoveRow(row)
			nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("Accepted Friend Request from %s", selectedUser.Username))
		})
	})

	page.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(HOME_PAGE, nil)
			page.userPendingRequests = make(map[uint8]chat.UserRelationship, 0)
			page.table.Clear()
		}

		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvInstructions.SetText("(esc) Quit")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)

	grid.SetRows(2, 1, 1, 0, 1, 1, 2)
	grid.SetColumns(0, 76, 0)

	grid.AddItem(tvHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(page.table, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 5, 1, 1, 1, 0, 0, false)

	nav.Register(ACCEPT_FRIEND_REQUEST_PAGE, grid, true, false,
		func(param interface{}) {
			page.onPageLoad(appContext)
		},
		func() {
			page.onPageClose()
		})
}

// onPageLoad is called when the page is navigated to
func (page *AcceptFriendRequestPage) onPageLoad(appContext *state.ApplicationContext) {
	page.table.SetCell(0, 0, tview.NewTableCell("Username").
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignCenter).
		SetExpansion(1).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	page.table.SetCell(0, 1, tview.NewTableCell("Last Active").
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignRight).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	row := 1

	for _, rel := range appContext.BrochatUser.Relationships {
		if rel.Type&chat.RELATIONSHIP_TYPE_FRIEND_REQUEST_RECIEVED != 0 {
			page.table.SetCell(row, 0, tview.NewTableCell(rel.Username).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
			var dateString string = rel.LastOnlineUtc.Local().Format("Jan 2, 2006")
			page.table.SetCell(row, 1, tview.NewTableCell(dateString).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight))

			page.userPendingRequests[uint8(row)] = rel
			row++
		}
	}
}

// onPageClose is called when the page is navigated away from
func (page *AcceptFriendRequestPage) onPageClose() {
	page.userPendingRequests = make(map[uint8]chat.UserRelationship)
	page.table.Clear()
}
