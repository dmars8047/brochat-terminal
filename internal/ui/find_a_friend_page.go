package ui

import (
	"fmt"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const FRIENDS_FINDER_PAGE PageSlug = "find_a_friend"

const (
	FIND_A_FRIEND_PAGE_ALERT_INFO = "home:findafriend:alert:info"
	FIND_A_FRIEND_PAGE_ALERT_ERR  = "home:findafriend:alert:err"
	FIND_A_FRIEND_PAGE_CONFIRM    = "home:findafriend:confirm"
)

// FindAFriendPage is the find a friend page
type FindAFriendPage struct {
	brochatClient *chat.BroChatUserClient
	table         *tview.Table
	users         map[uint8]chat.UserInfo
}

// NewFindAFriendPage creates a new find a friend page
func NewFindAFriendPage(brochatClient *chat.BroChatUserClient) *FindAFriendPage {
	return &FindAFriendPage{
		brochatClient: brochatClient,
		table:         tview.NewTable(),
		users:         make(map[uint8]chat.UserInfo, 0),
	}
}

// Setup sets up the find a friend page and registers it with the page navigator
func (page *FindAFriendPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvHeader.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvHeader.SetText("Find Friends")

	page.table.SetBorders(true)
	page.table.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	page.table.SetFixed(1, 1)
	page.table.SetSelectable(true, false)

	page.table.SetSelectedFunc(func(row int, _ int) {
		selectedUser, ok := page.users[uint8(row)]

		if !ok {
			return
		}

		nav.Confirm(FIND_A_FRIEND_PAGE_CONFIRM, fmt.Sprintf("Send Friend Request to %s?", selectedUser.Username), func() {
			err := page.brochatClient.SendFriendRequest(appContext.GetAuthInfo(), &chat.SendFriendRequestRequest{
				RequestedUserId: selectedUser.Id,
			})

			if err != nil {
				if err.Error() == "friend request already exists or users are already a friend" {
					nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("Friend Request Already Sent to %s", selectedUser.Username))
					return
				} else if err.Error() == "user not found" {
					nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("User %s Not Found", selectedUser.Username))
					return
				}

				nav.AlertFatal(app, FIND_A_FRIEND_PAGE_ALERT_ERR, err.Error())
				return
			}

			// Update the relationship type
			for i := 0; i < len(appContext.BrochatUser.Relationships); i++ {
				if selectedUser.Id == appContext.BrochatUser.Relationships[i].UserId {
					existingRelationshipType := appContext.BrochatUser.Relationships[i].Type

					if existingRelationshipType == chat.RELATIONSHIP_TYPE_DEFAULT {
						appContext.BrochatUser.Relationships[i].Type = chat.RELATIONSHIP_TYPE_FRIENDSHIP_REQUESTED
					} else {
						appContext.BrochatUser.Relationships[i].Type = existingRelationshipType | chat.RELATIONSHIP_TYPE_FRIENDSHIP_REQUESTED
					}
					break
				}
			}

			page.table.RemoveRow(row)
			nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("Friend Request Sent to %s", selectedUser.Username))
		})
	})

	page.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(HOME_PAGE, nil)
			page.users = make(map[uint8]chat.UserInfo, 0)
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

	nav.Register(FRIENDS_FINDER_PAGE, grid, true, false,
		func(_ interface{}) {
			page.onPageLoad(app, appContext, nav)
		},
		func() {
			page.onPageClose()
		})
}

// onPageLoad is called when the find a friend page is navigated to
func (page *FindAFriendPage) onPageLoad(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
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

	usrs, err := page.brochatClient.GetUsers(appContext.GetAuthInfo(), true, true, 1, 10, "")

	if err != nil {
		nav.AlertFatal(app, FIND_A_FRIEND_PAGE_ALERT_ERR, err.Error())
		return
	}

	for i, usr := range usrs {
		row := i + 1

		page.table.SetCell(row, 0, tview.NewTableCell(usr.Username).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
		var dateString string = usr.LastOnlineUtc.Local().Format("Jan 2, 2006")
		page.table.SetCell(row, 1, tview.NewTableCell(dateString).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight))

		page.users[uint8(row)] = usr
	}
}

// onPageClose is called when the find a friend page is navigated away from
func (page *FindAFriendPage) onPageClose() {
	page.users = make(map[uint8]chat.UserInfo, 0)
	page.table.Clear()
}
