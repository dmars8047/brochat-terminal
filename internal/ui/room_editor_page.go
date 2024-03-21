package ui

import (
	"fmt"
	"log"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/strval"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const ROOM_EDITOR_PAGE PageSlug = "room_editor"

const (
	ROOM_EDITOR_PAGE_ALERT_INFO = "home:roomeditor:alert:info"
	ROOM_EDITOR_PAGE_ALERT_ERR  = "home:roomeditor:alert:err"
	ROOM_EDITOR_PAGE_CONFIRM    = "home:roomeditor:confirm"
)

// RoomEditorPage is the room editor page
type RoomEditorPage struct {
	brochatClient *chat.BroChatClient
	form          *tview.Form
}

// NewRoomEditorPage creates a new room editor page
func NewRoomEditorPage(brochatClient *chat.BroChatClient) *RoomEditorPage {
	return &RoomEditorPage{
		brochatClient: brochatClient,
		form:          tview.NewForm(),
	}
}

// Setup sets up the room editor page and registers it with the page navigator
func (page *RoomEditorPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {

	theme := appContext.GetTheme()

	grid := tview.NewGrid()
	grid.SetBackgroundColor(theme.BackgroundColor)
	grid.SetRows(4, 0, 1, 3, 4)
	grid.SetColumns(0, 70, 0)

	page.form.SetBackgroundColor(theme.AccentColor)
	page.form.SetFieldBackgroundColor(theme.AccentColorTwo)
	page.form.SetLabelColor(theme.HighlightColor)
	page.form.SetBorder(true)
	page.form.SetTitle(" BroChat - Room Creation Editor ")
	page.form.SetTitleAlign(tview.AlignCenter)
	page.form.SetButtonStyle(theme.ButtonStyle)
	page.form.SetButtonActivatedStyle(theme.ActivatedButtonStyle)

	//Add forms
	page.form.AddInputField("Room Name", "", 0, nil, nil)
	page.form.AddDropDown("Membership Model", []string{string(chat.PUBLIC_MEMBERSHIP_MODEL), string(chat.FRIENDS_MEMBERSHIP_MODEL)}, -1, nil)

	page.form.AddButton("Submit", func() {
		accessToken, ok := appContext.GetAccessToken()

		if !ok {
			log.Printf("Valid user authentication information not found. Redirecting to login page.")
			nav.NavigateTo(LOGIN_PAGE, nil)
			return
		}

		nameInput, ok := page.form.GetFormItemByLabel("Room Name").(*tview.InputField)

		if !ok {
			panic("room name input form access failure")
		}

		name := nameInput.GetText()

		valResult := strval.ValidateStringWithName(name, "Room Name",
			strval.MustNotBeEmpty(),
			strval.MustHaveMinLengthOf(3),
			strval.MustHaveMaxLengthOf(32),
		)

		if !valResult.Valid {
			nav.AlertErrors(ROOM_EDITOR_PAGE_ALERT_ERR, "Room Creation Failed - Form Validation Error", valResult.Messages)
			return
		}

		membershipModelDropdown, ok := page.form.GetFormItemByLabel("Membership Model").(*tview.DropDown)

		if !ok {
			panic("membership model dropdown form access failure")
		}

		optIndex, optstr := membershipModelDropdown.GetCurrentOption()

		if optIndex < 0 || optstr == "" {
			nav.Alert(ROOM_EDITOR_PAGE_ALERT_ERR, "Room Creation Failed - Membership Model Selection Invalid")
			return
		}

		request := chat.CreateRoomRequest{
			Name:            name,
			MembershipModel: optstr,
		}

		createRoomResult := page.brochatClient.CreateRoom(accessToken, request)

		createRoomErr := createRoomResult.Err()

		if createRoomErr != nil {
			if len(createRoomResult.ErrorDetails) > 0 {
				nav.Alert(ROOM_EDITOR_PAGE_ALERT_ERR, createRoomResult.ErrorDetails[0])
				return
			}

			if createRoomResult.ResponseCode == chat.BROCHAT_RESPONSE_CODE_FORBIDDEN_ERROR {
				nav.Alert(ROOM_EDITOR_PAGE_ALERT_ERR, FORBIDDEN_OPERATION_ERROR_MESSAGE)
				return
			}

			nav.Alert(ROOM_EDITOR_PAGE_ALERT_ERR, fmt.Sprintf("An error occurred while creating user room: %s", createRoomErr.Error()))
			return
		}

		nav.AlertWithDoneFunc(ROOM_EDITOR_PAGE_ALERT_INFO, "Room creation successful!", func(buttonIndex int, buttonLabel string) {
			nav.NavigateTo(ROOM_LIST_PAGE, nil)
		})
	})

	page.form.AddButton("Back", func() {
		nav.NavigateTo(ROOM_LIST_PAGE, nil)
	})

	page.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(ROOM_LIST_PAGE, nil)
		}

		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(theme.BackgroundColor)
	tvInstructions.SetText("Enter a name and membership model for your new room.")
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))

	grid.AddItem(page.form, 1, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 3, 1, 1, 1, 0, 0, false)

	nav.Register(ROOM_EDITOR_PAGE, grid, true, false, func(_ interface{}) {
		page.onPageLoad()
	}, func() {
		page.onPageClose()
	})
}

// onPageLoad is called when the page is navigated to
func (page *RoomEditorPage) onPageLoad() {
	page.form.SetFocus(0)
}

// onPageClose is called when the page is navigated away from
func (page *RoomEditorPage) onPageClose() {
	roomNameInput, ok := page.form.GetFormItemByLabel("Room Name").(*tview.InputField)

	if !ok {
		panic("room name input form clear failure")
	}

	roomNameInput.SetText("")

	membModelDropdown, ok := page.form.GetFormItemByLabel("Membership Model").(*tview.DropDown)

	if !ok {
		panic("membership model dropdown form clear failure")
	}

	membModelDropdown.SetCurrentOption(-1)
}
