package ui

import (
	"fmt"
	"time"

	"github.com/dmars8047/broterm/internal/auth"
	"github.com/dmars8047/broterm/internal/bro"
	"github.com/dmars8047/broterm/internal/feed"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HomeModule struct {
	appContext     *state.ApplicationContext
	userAuthClient *auth.UserAuthClient
	brochatClient  *bro.BroChatUserClient
	pageNav        *PageNavigator
	app            *tview.Application
	feedClient     *feed.Client
}

func NewHomeModule(userAuthClient *auth.UserAuthClient,
	application *tview.Application,
	pageNavigator *PageNavigator,
	brochatClient *bro.BroChatUserClient,
	appContext *state.ApplicationContext,
	feedClient *feed.Client) *HomeModule {
	return &HomeModule{
		appContext:     appContext,
		brochatClient:  brochatClient,
		userAuthClient: userAuthClient,
		app:            application,
		pageNav:        pageNavigator,
		feedClient:     feedClient,
	}
}

func (mod *HomeModule) SetupHomePages() {
	mod.setupMenuPage()
	mod.setupFriendListPage()
	mod.setupFindAFriendPage()
	mod.setupAcceptPendingRequestPage()
	mod.setupChatPage()
}

func (mod *HomeModule) setupMenuPage() {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(4, 8, 8, 0).
		SetColumns(0, 31, 39, 0)

	logoBro := tview.NewTextView()
	logoBro.SetTextAlign(tview.AlignLeft).
		SetBackgroundColor(DefaultBackgroundColor)
	logoBro.SetTextColor(tcell.ColorWhite)
	logoBro.SetText(
		`BBBBBBB\                      
BB  __BB\                     
BB |  BB | RRRRRR\   OOOOOO\  
BBBBBBB\ |RR  __RR\ OO  __OO\ 
BB  __BB\ RR |  \__|OO /  OO |
BB |  BB |RR |      OO |  OO |
BBBBBBB  |RR |      \OOOOOO  |
\_______/ \__|       \______/ `)

	logoChat := tview.NewTextView()
	logoChat.SetTextAlign(tview.AlignLeft)
	logoChat.SetBackgroundColor(DefaultBackgroundColor)
	logoChat.SetTextColor(BroChatYellowColor)
	logoChat.SetText(
		` CCCCCC\  HH\                  TT\
CC  __CC\ HH |                 TT |
CC /  \__|HHHHHHH\   AAAAAA\ TTTTTT\
CC |      HH  __HH\  \____AA\\_TT  _|
CC |      HH |  HH | AAAAAAA | TT |
CC |  CC\ HH |  HH |AA  __AA | TT |TT\
\CCCCCC  |HH |  HH |\AAAAAAA | \TTTT  |
 \______/ \__|  \__| \_______|  \____/`)

	brosButton := tview.NewButton("Bros").
		SetActivatedStyle(ActivatedButtonStyle).
		SetStyle(ButtonStyle)

	brosButton.SetSelectedFunc(func() {
		mod.pageNav.NavigateTo(HOME_FRIENDS_LIST_PAGE, nil)
	})

	chatButton := tview.NewButton("Chat").
		SetActivatedStyle(ActivatedButtonStyle).
		SetStyle(ButtonStyle)

	chatButton.SetSelectedFunc(func() {
		Alert(mod.pageNav.Pages, "home:menu:alert:info", "Chat Servers Not Implemented Yet")
	})

	settingsButton := tview.NewButton("Settings").
		SetActivatedStyle(ActivatedButtonStyle).
		SetStyle(ButtonStyle)

	settingsButton.SetSelectedFunc(func() {
		Alert(mod.pageNav.Pages, "home:menu:alert:info", "Settings Not Implemented Yet")
	})

	logoutButton := tview.NewButton("Logout").
		SetActivatedStyle(ActivatedButtonStyle).
		SetStyle(ButtonStyle)

	logoutButton.SetSelectedFunc(func() {
		err := mod.userAuthClient.Logout(mod.appContext.UserSession.Auth.AccessToken)

		if err != nil {
			AlertFatal(mod.app, mod.pageNav.Pages, "home:menu:alert:err", err.Error())
			return
		}

		mod.appContext.UserSession.CancelFunc()
		mod.appContext.UserSession = nil

		mod.pageNav.NavigateTo(WELCOME_PAGE, nil)
	})

	buttonGrid := tview.NewGrid()

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DefaultBackgroundColor)
	tvInstructions.SetTextColor(tcell.ColorWhite)

	logoutButton.SetFocusFunc(func() {
		tvInstructions.SetText("Sign out of your account.")
	})

	settingsButton.SetFocusFunc(func() {
		tvInstructions.SetText("Change your account settings.")
	})

	chatButton.SetFocusFunc(func() {
		tvInstructions.SetText("Chat in a server or find one to join.")
	})

	brosButton.SetFocusFunc(func() {
		tvInstructions.SetText("Talk to your Bros or find new ones!")
	})

	buttonGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if brosButton.HasFocus() {
				mod.app.SetFocus(chatButton)
			} else if chatButton.HasFocus() {
				mod.app.SetFocus(settingsButton)
			} else if settingsButton.HasFocus() {
				mod.app.SetFocus(logoutButton)
			} else if logoutButton.HasFocus() {
				mod.app.SetFocus(brosButton)
			}
		} else if event.Key() == tcell.KeyBacktab {
			if logoutButton.HasFocus() {
				mod.app.SetFocus(settingsButton)
			} else if settingsButton.HasFocus() {
				mod.app.SetFocus(chatButton)
			} else if chatButton.HasFocus() {
				mod.app.SetFocus(brosButton)
			} else if brosButton.HasFocus() {
				mod.app.SetFocus(logoutButton)
			}
		}
		return event
	})

	buttonGrid.SetRows(3, 1, 1).
		SetColumns(0, 1, 0, 1, 0, 1, 0)

	buttonGrid.AddItem(brosButton, 0, 0, 1, 1, 0, 0, true).
		AddItem(chatButton, 0, 2, 1, 1, 0, 0, true).
		AddItem(settingsButton, 0, 4, 1, 1, 0, 0, true).
		AddItem(logoutButton, 0, 6, 1, 1, 0, 0, true).
		AddItem(tvInstructions, 2, 0, 1, 7, 0, 0, false)

	grid.AddItem(logoBro, 1, 1, 1, 1, 0, 0, false).
		AddItem(logoChat, 1, 2, 1, 1, 0, 0, false).
		AddItem(buttonGrid, 2, 1, 1, 2, 0, 0, true)

	mod.pageNav.Register(HOME_MENU_PAGE, grid, true, false,
		func(param interface{}) {

			// Make sure the session is still valid
			if mod.appContext.UserSession.Auth.TokenExpiration.Before(time.Now()) {
				mod.appContext.UserSession.CancelFunc()
				mod.appContext.UserSession = nil
				mod.pageNav.NavigateTo(LOGIN_PAGE, nil)
			}
		}, nil)
}

const (
	FRIENDS_LIST_PAGE_ALERT_INFO = "home:friendlist:alert:info"
	FRIENDS_LIST_PAGE_ALERT_ERR  = "home:friendlist:alert:err"
)

func (mod *HomeModule) setupFriendListPage() {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetBackgroundColor(DefaultBackgroundColor)
	tvHeader.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvHeader.SetText("Friends List")

	userFriends := make(map[uint8]bro.UserRelationship, 0)

	table := tview.NewTable().
		SetBorders(true)
	table.SetBackgroundColor(DefaultBackgroundColor)
	table.SetFixed(1, 1)
	table.SetSelectable(true, false)

	table.SetSelectedFunc(func(row int, _ int) {
		rel, ok := userFriends[uint8(row)]

		if !ok {
			return
		}

		mod.pageNav.NavigateTo(HOME_CHAT_PAGE, ChatParams{
			channel_id: rel.DirectMessageChannelId,
		})
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'p':
				mod.pageNav.NavigateTo(HOME_PENDING_REQUESTS_PAGE, nil)
				userFriends = make(map[uint8]bro.UserRelationship, 0)
				table.Clear()
			case 'f':
				mod.pageNav.NavigateTo(HOME_FRIENDS_FINDER_PAGE, nil)
				userFriends = make(map[uint8]bro.UserRelationship, 0)
				table.Clear()
			}
		} else if event.Key() == tcell.KeyEscape {
			mod.pageNav.NavigateTo(HOME_MENU_PAGE, nil)
			userFriends = make(map[uint8]bro.UserRelationship, 0)
			table.Clear()
		}

		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DefaultBackgroundColor)
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvInstructions.SetText("(f) Find a new Bro - (p) View Pending - (esc) Quit")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(2, 1, 1, 0, 1, 1, 2)
	grid.SetColumns(0, 76, 0)

	grid.AddItem(tvHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(table, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 5, 1, 1, 1, 0, 0, false)

	mod.pageNav.Register(HOME_FRIENDS_LIST_PAGE, grid, true, false,
		func(param interface{}) {
			table.SetCell(0, 0, tview.NewTableCell("Username").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).
				SetExpansion(1).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

			table.SetCell(0, 1, tview.NewTableCell("Status").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

			table.SetCell(0, 2, tview.NewTableCell("Last Active").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignRight).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

			usr, err := mod.brochatClient.GetUser(&bro.AuthInfo{
				AccessToken: mod.appContext.UserSession.Auth.AccessToken,
				TokenType:   DEFAULT_AUTH_TOKEN_TYPE,
			}, mod.appContext.UserSession.Info.Id)

			if err != nil {
				AlertFatal(mod.app, mod.pageNav.Pages, FRIENDS_LIST_PAGE_ALERT_ERR, err.Error())
				return
			}

			mod.appContext.BrochatUser = usr

			countOfPendingFriendRequests := 0

			for _, rel := range usr.Relationships {
				if rel.Type&bro.RELATIONSHIP_TYPE_FRIEND_REQUEST_RECIEVED != 0 {
					countOfPendingFriendRequests++
				}
			}

			tvInstructions.SetText(fmt.Sprintf("(f) Find a new Bro - (p) View Pending [%d] - (esc) Quit", countOfPendingFriendRequests))

			for i, rel := range usr.Relationships {
				row := i + 1

				if rel.Type != bro.RELATIONSHIP_TYPE_FRIEND {
					continue
				}

				table.SetCell(row, 0, tview.NewTableCell(rel.Username).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
				if rel.IsOnline {
					table.SetCell(row, 1, tview.NewTableCell("Online").SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))
				} else {
					table.SetCell(row, 1, tview.NewTableCell("Offline").SetTextColor(tcell.ColorGray).SetAlign(tview.AlignCenter))
				}

				var dateString string = rel.LastOnlineUtc.Local().Format("Jan 2, 2006")

				table.SetCell(row, 2, tview.NewTableCell(dateString).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight))

				userFriends[uint8(row)] = rel
			}
		},
		func() {
			userFriends = make(map[uint8]bro.UserRelationship, 0)
			table.Clear()
		})
}

const (
	FIND_A_FRIEND_PAGE_ALERT_INFO = "home:findafriend:alert:info"
	FIND_A_FRIEND_PAGE_ALERT_ERR  = "home:findafriend:alert:err"
	FIND_A_FRIEND_PAGE_CONFIRM    = "home:findafriend:confirm"
)

func (mod *HomeModule) setupFindAFriendPage() {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetBackgroundColor(DefaultBackgroundColor)
	tvHeader.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvHeader.SetText("Find Friends")

	users := make(map[uint8]bro.UserInfo, 0)

	table := tview.NewTable().
		SetBorders(true)
	table.SetBackgroundColor(DefaultBackgroundColor)
	table.SetFixed(1, 1)
	table.SetSelectable(true, false)

	table.SetSelectedFunc(func(row int, _ int) {
		uInfo, ok := users[uint8(row)]

		if !ok {
			return
		}

		Confirm(mod.pageNav.Pages, FIND_A_FRIEND_PAGE_CONFIRM, fmt.Sprintf("Send Friend Request to %s?", uInfo.Username), func() {
			err := mod.brochatClient.SendFriendRequest(&bro.AuthInfo{
				AccessToken: mod.appContext.UserSession.Auth.AccessToken,
				TokenType:   DEFAULT_AUTH_TOKEN_TYPE,
			}, &bro.SendFriendRequestRequest{
				RequestedUserId: uInfo.ID,
			})

			if err != nil {
				if err.Error() == "friend request already exists or users are already a friend" {
					Alert(mod.pageNav.Pages, FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("Friend Request Already Sent to %s", uInfo.Username))
					return
				} else if err.Error() == "user not found" {
					Alert(mod.pageNav.Pages, FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("User %s Not Found", uInfo.Username))
					return
				}

				AlertFatal(mod.app, mod.pageNav.Pages, FIND_A_FRIEND_PAGE_ALERT_ERR, err.Error())
				return
			}

			Alert(mod.pageNav.Pages, FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("Friend Request Sent to %s", uInfo.Username))
		})
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mod.pageNav.NavigateTo(HOME_MENU_PAGE, nil)
			users = make(map[uint8]bro.UserInfo, 0)
			table.Clear()
		}

		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DefaultBackgroundColor)
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvInstructions.SetText("(esc) Quit")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(2, 1, 1, 0, 1, 1, 2)
	grid.SetColumns(0, 76, 0)

	grid.AddItem(tvHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(table, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 5, 1, 1, 1, 0, 0, false)

	mod.pageNav.Register(HOME_FRIENDS_FINDER_PAGE, grid, true, false,
		func(param interface{}) {
			table.SetCell(0, 0, tview.NewTableCell("Username").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).
				SetExpansion(1).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

			table.SetCell(0, 1, tview.NewTableCell("Last Active").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignRight).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

			usrs, err := mod.brochatClient.GetUsers(&bro.AuthInfo{
				AccessToken: mod.appContext.UserSession.Auth.AccessToken,
				TokenType:   DEFAULT_AUTH_TOKEN_TYPE,
			}, true, true, 1, 10, "")

			if err != nil {
				AlertFatal(mod.app, mod.pageNav.Pages, FIND_A_FRIEND_PAGE_ALERT_ERR, err.Error())
				return
			}

			for i, usr := range usrs {
				row := i + 1

				table.SetCell(row, 0, tview.NewTableCell(usr.Username).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
				var dateString string = usr.LastOnlineUtc.Local().Format("Jan 2, 2006")
				table.SetCell(row, 1, tview.NewTableCell(dateString).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight))

				users[uint8(row)] = usr
			}
		},
		func() {
			users = make(map[uint8]bro.UserInfo, 0)
			table.Clear()
		})
}

func (mod *HomeModule) setupAcceptPendingRequestPage() {
	tvHeader := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvHeader.SetBackgroundColor(DefaultBackgroundColor)
	tvHeader.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvHeader.SetText("Pending Friend Requests")

	table := tview.NewTable().
		SetBorders(true)
	table.SetBackgroundColor(DefaultBackgroundColor)
	table.SetFixed(1, 1)
	table.SetSelectable(true, false)

	userPendingRequests := make(map[uint8]bro.UserRelationship, 0)

	table.SetSelectedFunc(func(row int, _ int) {
		rel, ok := userPendingRequests[uint8(row)]

		if !ok {
			return
		}

		Confirm(mod.pageNav.Pages, FIND_A_FRIEND_PAGE_CONFIRM, fmt.Sprintf("Accept Friend Request from %s?", rel.Username), func() {
			err := mod.brochatClient.AcceptFriendRequest(&bro.AuthInfo{
				AccessToken: mod.appContext.UserSession.Auth.AccessToken,
				TokenType:   DEFAULT_AUTH_TOKEN_TYPE,
			}, &bro.AcceptFriendRequestRequest{
				InitiatingUserId: rel.UserId,
			})

			if err != nil {
				if err.Error() == "user not found or friend request not found" {
					Alert(mod.pageNav.Pages, FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("Friend Request from %s Not Found", rel.Username))
					return
				} else if err.Error() == "bad request" {
					Alert(mod.pageNav.Pages, FIND_A_FRIEND_PAGE_ALERT_INFO, "Friend Request Acceptance Not Processable")
					return
				}

				AlertFatal(mod.app, mod.pageNav.Pages, FIND_A_FRIEND_PAGE_ALERT_ERR, err.Error())
				return
			}

			Alert(mod.pageNav.Pages, FIND_A_FRIEND_PAGE_ALERT_INFO, fmt.Sprintf("Accepted Friend Request from %s", rel.Username))
		})
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mod.pageNav.NavigateTo(HOME_MENU_PAGE, nil)
			userPendingRequests = make(map[uint8]bro.UserRelationship, 0)
			table.Clear()
		}

		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DefaultBackgroundColor)
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))
	tvInstructions.SetText("(esc) Quit")

	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(2, 1, 1, 0, 1, 1, 2)
	grid.SetColumns(0, 76, 0)

	grid.AddItem(tvHeader, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(table, 3, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 5, 1, 1, 1, 0, 0, false)

	mod.pageNav.Register(HOME_PENDING_REQUESTS_PAGE, grid, true, false,
		func(param interface{}) {
			table.SetCell(0, 0, tview.NewTableCell("Username").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).
				SetExpansion(1).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

			table.SetCell(0, 1, tview.NewTableCell("Last Active").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignRight).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold|tcell.AttrUnderline))

			for i, rel := range mod.appContext.BrochatUser.Relationships {
				row := i + 1

				if rel.Type&bro.RELATIONSHIP_TYPE_FRIEND_REQUEST_RECIEVED != 0 {
					table.SetCell(row, 0, tview.NewTableCell(rel.Username).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
					var dateString string = rel.LastOnlineUtc.Local().Format("Jan 2, 2006")
					table.SetCell(row, 1, tview.NewTableCell(dateString).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight))

					userPendingRequests[uint8(row)] = rel
				}
			}
		},
		func() {
			userPendingRequests = make(map[uint8]bro.UserRelationship)
			table.Clear()
		})
}

type ChatParams struct {
	channel_id string
}

func (mod *HomeModule) setupChatPage() {
	textView := tview.NewTextView().
		SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetScrollable(true)
	textView.SetBackgroundColor(DefaultBackgroundColor)

	textView.SetChangedFunc(func() {
		mod.app.Draw()
	})

	textArea := tview.NewTextArea()
	textArea.SetTextStyle(tcell.StyleDefault.Background(DefaultBackgroundColor))
	textArea.SetBorder(true)

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DefaultBackgroundColor)
	tvInstructions.SetTextColor(tcell.ColorWhite)

	tvInstructions.SetText("(enter) Send - (esc) Back")

	grid := tview.NewGrid()

	grid.SetRows(0, 6, 2)
	grid.SetColumns(0)

	grid.AddItem(textView, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(textArea, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 2, 0, 1, 1, 0, 0, false)

	mod.pageNav.Register(HOME_CHAT_PAGE, grid, true, false,
		func(param interface{}) {

			// The param should be a ChatParams struct
			chatParams, ok := param.(ChatParams)

			if !ok {
				AlertFatal(mod.app, mod.pageNav.Pages, "home:chat:alert:err", "Application State Error - Could not get chat params.")
				return
			}

			// Get the channel
			channel, err := mod.brochatClient.GetChannelManifest(&bro.AuthInfo{
				AccessToken: mod.appContext.UserSession.Auth.AccessToken,
				TokenType:   DEFAULT_AUTH_TOKEN_TYPE,
			}, chatParams.channel_id)

			if err != nil {
				AlertFatal(mod.app, mod.pageNav.Pages, "home:chat:alert:err", err.Error())
				return
			}

			textView.SetTitle(fmt.Sprintf(" %s - %s ", channel.Users[0].Username, channel.Users[1].Username))

			// Get the channel messages
			messages, err := mod.brochatClient.GetChannelMessages(&bro.AuthInfo{
				AccessToken: mod.appContext.UserSession.Auth.AccessToken,
				TokenType:   DEFAULT_AUTH_TOKEN_TYPE,
			}, chatParams.channel_id)

			if err != nil {
				AlertFatal(mod.app, mod.pageNav.Pages, "home:chat:alert:err", err.Error())
				return
			}

			// Write the messages to the text view
			w := textView.BatchWriter()
			defer w.Close()
			w.Clear()

			for i := len(messages) - 1; i >= 0; i-- {
				// Write the messages to teh text view
				var senderUsername string
				var msg = messages[i]

				for _, u := range channel.Users {
					if u.ID == msg.SenderUserId {
						senderUsername = u.Username
						break
					}
				}

				var dateString string

				// If the message is from a date in the past (not today) then format the date string differently
				if msg.RecievedAtUtc.Local().Day() == time.Now().Day() {
					dateString = msg.RecievedAtUtc.Local().Format(time.Kitchen)
				} else {
					dateString = msg.RecievedAtUtc.Local().Format("Jan 2, 2006 3:04 PM")
				}

				msgString := fmt.Sprintf("%s [%s]: %s", senderUsername, dateString, msg.Content)
				fmt.Fprintln(w, msgString)
			}

			textView.ScrollToEnd()

			// Set the chat context
			mod.appContext.ChatSession = state.NewChatSession(channel, mod.appContext.Context)

			textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEnter {
					text := textArea.GetText()

					if len(text) > 0 {
						mod.feedClient.Send(bro.ChatMessage{
							ChannelId:    channel.ID,
							Content:      text,
							SenderUserId: mod.appContext.UserSession.Info.Id,
						})

						textArea.SetText("", false)
					}

					return nil
				} else if event.Key() == tcell.KeyEscape {
					mod.pageNav.NavigateTo(HOME_MENU_PAGE, nil)
				}

				return event
			})

			go func(ch bro.ChannelManifest, cs *state.ChatSession, a *tview.Application, tv *tview.TextView) {
				for {
					select {
					case <-cs.Context.Done():
						return
					case msg := <-mod.feedClient.ChatMessageChannel:
						if msg.ChannelId == ch.ID {
							a.QueueUpdateDraw(func() {
								var senderUsername string

								for _, u := range ch.Users {
									if u.ID == msg.SenderUserId {
										senderUsername = u.Username
										break
									}
								}

								dateString := msg.RecievedAtUtc.Local().Format(time.Kitchen)

								msgString := fmt.Sprintf("%s [%s]: %s", senderUsername, dateString, msg.Content)
								tv.Write([]byte(msgString + "\n"))
								tv.ScrollToEnd()
							})
						}
					}
				}
			}(*channel, mod.appContext.ChatSession, mod.app, textView)
		},
		func() {
			textView.Clear()
			textArea.SetText("", false)
			if mod.appContext.ChatSession != nil {
				mod.appContext.ChatSession.CancelFunc()
				mod.appContext.ChatSession = nil
			}
		})
}
