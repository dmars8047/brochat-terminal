package ui

import (
	"fmt"
	"log"

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
	brochatClient *chat.BroChatClient
	table         *tview.Table
	users         map[uint8]chat.UserInfo
	themeCode     string
}

// NewFindAFriendPage creates a new find a friend page
func NewFindAFriendPage(brochatClient *chat.BroChatClient) *FindAFriendPage {
	return &FindAFriendPage{
		brochatClient: brochatClient,
		table:         tview.NewTable(),
		users:         make(map[uint8]chat.UserInfo, 0),
		themeCode:     "NOT_SET",
	}
}

// Setup sets up the find a friend page and registers it with the page navigator
func (page *FindAFriendPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetText("Find Friends")

	page.table.SetBorders(true)
	page.table.SetFixed(1, 1)
	page.table.SetSelectable(true, false)

	page.table.SetSelectedFunc(func(row int, _ int) {
		selectedUser, ok := page.users[uint8(row)]

		if !ok {
			return
		}

		accessToken, ok := appContext.GetAccessToken()

		if !ok {
			log.Printf("Valid user authentication information not found. Redirecting to login page.")
			nav.NavigateTo(LOGIN_PAGE, nil)
			return
		}

		nav.Confirm(FIND_A_FRIEND_PAGE_CONFIRM, fmt.Sprintf("Send Friend Request to %s?", selectedUser.Username), func() {
			sendFriendRequestResult := page.brochatClient.SendFriendRequest(accessToken, chat.SendFriendRequestRequest{
				RequestedUserId: selectedUser.Id,
			})

			err := sendFriendRequestResult.Err()

			if err != nil {
				if len(sendFriendRequestResult.ErrorDetails) > 0 {
					nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, sendFriendRequestResult.ErrorDetails[0])
					return
				}

				if sendFriendRequestResult.ResponseCode == chat.BROCHAT_RESPONSE_CODE_FORBIDDEN_ERROR {
					nav.Alert(FIND_A_FRIEND_PAGE_ALERT_ERR, FORBIDDEN_OPERATION_ERROR_MESSAGE)
					return
				}

				nav.AlertFatal(app, FIND_A_FRIEND_PAGE_ALERT_ERR, err.Error())
				return
			}

			page.table.RemoveRow(row)
			nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("Friend Request Sent to %s", selectedUser.Username))
		})
	})

	page.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(FRIENDS_LIST_PAGE, nil)
			page.users = make(map[uint8]chat.UserInfo, 0)
			page.table.Clear()
		} else if event.Key() == tcell.KeyTab {
			// Change the selected row to the next row
			row, _ := page.table.GetSelection()
			if row+1 >= page.table.GetRowCount() {
				row = 1
			} else {
				row++
			}

			page.table.Select(row, 0)
		} else if event.Key() == tcell.KeyBacktab {
			// Change the selected row to the previous row
			row, _ := page.table.GetSelection()

			if row-1 < 1 {
				row = page.table.GetRowCount() - 1
			} else {
				row--
			}

			page.table.Select(row, 0)
		}

		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetText("(esc) Quit")

	grid := tview.NewGrid()

	grid.SetRows(2, 1, 1, 0, 1, 1, 2)
	grid.SetColumns(0, 76, 0)

	grid.AddItem(tvHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(page.table, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 5, 1, 1, 1, 0, 0, false)

	applyTheme := func() {
		theme := appContext.GetTheme()

		if page.themeCode != theme.Code {
			page.themeCode = theme.Code
			grid.SetBackgroundColor(theme.BackgroundColor)
			tvHeader.SetBackgroundColor(theme.BackgroundColor)
			tvHeader.SetTextColor(theme.TitleColor)
			page.table.SetBordersColor(theme.BorderColor)
			page.table.SetBorderColor(theme.BorderColor)
			page.table.SetTitleColor(theme.TitleColor)
			page.table.SetBackgroundColor(theme.BackgroundColor)
			page.table.SetSelectedStyle(theme.DropdownListSelectedStyle)
			tvInstructions.SetBackgroundColor(theme.BackgroundColor)
			tvInstructions.SetTextColor(theme.InfoColor)
		}
	}

	applyTheme()

	nav.Register(FRIENDS_FINDER_PAGE, grid, true, false,
		func(_ interface{}) {
			applyTheme()
			page.onPageLoad(app, appContext, nav)
		},
		func() {
			page.onPageClose()
		})
}

// onPageLoad is called when the find a friend page is navigated to
func (page *FindAFriendPage) onPageLoad(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	accessToken, ok := appContext.GetAccessToken()

	if !ok {
		log.Printf("Valid user authentication information not found. Redirecting to login page.")
		nav.NavigateTo(LOGIN_PAGE, nil)
		return
	}

	thm := appContext.GetTheme()

	page.table.SetCell(0, 0, tview.NewTableCell("Username").
		SetTextColor(thm.ForgroundColor).
		SetAlign(tview.AlignCenter).
		SetExpansion(1).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	page.table.SetCell(0, 1, tview.NewTableCell("Last Active").
		SetTextColor(thm.ForgroundColor).
		SetAlign(tview.AlignRight).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

	getUsersResult := page.brochatClient.GetUsers(accessToken, chat.GetUsersOption_ExcludeSelf(),
		chat.GetUsersOption_ExcludeFriends(),
		chat.GetUsersOption_Page(1),
		chat.GetUsersOption_PageSize(10))

	err := getUsersResult.Err()

	if err != nil {
		if len(getUsersResult.ErrorDetails) > 0 {
			nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, getUsersResult.ErrorDetails[0])
			return
		}

		if getUsersResult.ResponseCode == chat.BROCHAT_RESPONSE_CODE_FORBIDDEN_ERROR {
			nav.Alert(FIND_A_FRIEND_PAGE_ALERT_INFO, FORBIDDEN_OPERATION_ERROR_MESSAGE)
			return
		}

		nav.AlertFatal(app, FIND_A_FRIEND_PAGE_ALERT_ERR, err.Error())
		return
	}

	usrs := getUsersResult.Content

	for i, usr := range usrs {
		row := i + 1

		page.table.SetCell(row, 0, tview.NewTableCell(usr.Username).SetTextColor(thm.ForgroundColor).SetAlign(tview.AlignCenter))
		var dateString string = usr.LastOnlineUtc.Local().Format("Jan 2, 2006")
		page.table.SetCell(row, 1, tview.NewTableCell(dateString).SetTextColor(thm.ForgroundColor).SetAlign(tview.AlignRight))

		page.users[uint8(row)] = usr
	}
}

// onPageClose is called when the find a friend page is navigated away from
func (page *FindAFriendPage) onPageClose() {
	page.users = make(map[uint8]chat.UserInfo, 0)
	page.table.Clear()
}
