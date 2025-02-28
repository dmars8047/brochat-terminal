package ui

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/broterm/internal/theme"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const CHAT_PAGE PageSlug = "chat"

// ChatPage is the chat page
type ChatPage struct {
	brochatClient    *chat.BroChatClient
	feedClient       *state.FeedClient
	textView         *tview.TextView
	textArea         *tview.TextArea
	mu               sync.Mutex
	currentThemeCode string
}

// NewChatPage creates a new chat page
func NewChatPage(brochatClient *chat.BroChatClient, feedClient *state.FeedClient) *ChatPage {
	return &ChatPage{
		brochatClient:    brochatClient,
		feedClient:       feedClient,
		textView:         tview.NewTextView(),
		textArea:         tview.NewTextArea(),
		currentThemeCode: "NOT_SET",
	}
}

// Setup configures the chat page and registers it with the page navigator
func (page *ChatPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	page.textView.SetDynamicColors(true)
	page.textView.SetBorder(true)
	page.textView.SetScrollable(true)

	page.textView.SetChangedFunc(func() {
		app.Draw()
	})

	page.textArea.SetBorder(true)

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	tvInstructions.SetText("(enter) Send - (pgup/pgdn) Scroll - (esc) Back")

	grid := tview.NewGrid()

	grid.SetRows(0, 6, 2)
	grid.SetColumns(0)

	grid.AddItem(page.textView, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(page.textArea, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 2, 0, 1, 1, 0, 0, false)

	var pageContext context.Context
	var cancel context.CancelFunc

	applyTheme := func() {
		theme := appContext.GetTheme()

		if page.currentThemeCode != theme.Code {
			page.currentThemeCode = theme.Code
			grid.SetBackgroundColor(theme.BackgroundColor)
			page.textView.SetBackgroundColor(theme.BackgroundColor)
			page.textView.SetBorderColor(theme.BorderColor)
			page.textView.SetTitleColor(theme.TitleColor)

			page.textArea.SetTextStyle(theme.TextAreaTextStyle)
			page.textArea.SetBorderColor(theme.BorderColor)
			page.textArea.SetTitleColor(theme.TitleColor)
			page.textArea.SetBorderStyle(theme.TextAreaTextStyle)

			tvInstructions.SetBackgroundColor(theme.BackgroundColor)
			tvInstructions.SetTextColor(theme.InfoColor)
		}
	}

	applyTheme()

	nav.Register(CHAT_PAGE, grid, true, false,
		func(param interface{}) {
			applyTheme()
			pageContext, cancel = appContext.GenerateUserSessionBoundContextWithCancel()
			page.onPageLoad(param, app, appContext, nav, pageContext)
		},
		func() {
			cancel()
			page.onPageClose()
		})
}

// onPageLoad is called when the chat page is navigated to
func (page *ChatPage) onPageLoad(param interface{},
	app *tview.Application,
	appContext *state.ApplicationContext,
	nav *PageNavigator,
	pageContext context.Context) {

	// The param should be a ChatParams struct
	chatParam, ok := param.(ChatPageParameters)

	if !ok {
		nav.AlertFatal(app, "home:chat:alert:err", "Application State Error - Could not get chat params.")
		return
	}

	accessToken, ok := appContext.GetAccessToken()

	if !ok {
		log.Printf("Valid user authentication information not found. Redirecting to login page.")
		nav.NavigateTo(LOGIN_PAGE, nil)
		return
	}

	// Get the channel
	getChannelResult := page.brochatClient.GetChannel(accessToken, chatParam.channel_id)

	err := getChannelResult.Err()

	if err != nil {
		if len(getChannelResult.ErrorDetails) > 0 {
			nav.Alert("home:chat:alert:err", getChannelResult.ErrorDetails[0])
			return
		}

		if getChannelResult.ResponseCode == chat.BROCHAT_RESPONSE_CODE_FORBIDDEN_ERROR {
			nav.Alert("home:chat:alert:err", FORBIDDEN_OPERATION_ERROR_MESSAGE)
			return
		}

		nav.Alert("home:chat:alert:err", err.Error())
		return
	}

	channel := getChannelResult.Content

	if channel.Type == chat.CHANNEL_TYPE_DIRECT_MESSAGE {
		page.textView.SetTitle(fmt.Sprintf(" %s - %s ", channel.Users[0].Username, channel.Users[1].Username))
	} else if chatParam.title != "" {
		page.textView.SetTitle(fmt.Sprintf(" %s ", chatParam.title))
	}

	theme := appContext.GetTheme()

	// Get the color manifest
	colorManifest := getColorManifest(channel.Users, theme)

	const pageSize = 100
	entireConversationLoaded := false
	oldestMessageId := ""

	// Get the channel messages
	getChannelMessagesResult := page.brochatClient.GetChannelMessages(accessToken, chatParam.channel_id, chat.GetChannelMessages_Page(1), chat.GetChannelMessages_PageSize(pageSize))

	err = getChannelMessagesResult.Err()

	if err != nil {
		if len(getChannelMessagesResult.ErrorDetails) > 0 {
			nav.Alert("home:chat:alert:err", getChannelMessagesResult.ErrorDetails[0])
			return
		}

		if getChannelMessagesResult.ResponseCode == chat.BROCHAT_RESPONSE_CODE_FORBIDDEN_ERROR {
			nav.Alert("home:chat:alert:err", FORBIDDEN_OPERATION_ERROR_MESSAGE)
			return
		}

		nav.AlertFatal(app, "home:chat:alert:err", err.Error())
		return
	}

	messages := getChannelMessagesResult.Content

	if len(messages) < pageSize {
		entireConversationLoaded = true
	} else {
		oldestMessageId = messages[len(messages)-1].Id
	}

	// Write the messages to the text view
	w := page.textView.BatchWriter()
	defer w.Close()
	w.Clear()

	for i := len(messages) - 1; i >= 0; i-- {
		// Write the messages to the text view
		var senderUsername string
		var msg = messages[i]

		color := colorManifest[msg.SenderUserId]

		for _, u := range channel.Users {
			if u.Id == msg.SenderUserId {
				senderUsername = u.Username
				break
			}
		}

		// If for some reason the user info is not found just make the username "Unknown User"
		if senderUsername == "" {
			senderUsername = "Unknown User"
		}

		// If the color is not found then just make it red
		if color == "" {
			color = "#FF0000"
		}

		var dateString string

		// If the message is from a date in the past (not today) then format the date string differently
		if msg.RecievedAtUtc.Local().Day() == time.Now().Day() {
			dateString = msg.RecievedAtUtc.Local().Format(time.Kitchen)
		} else {
			dateString = msg.RecievedAtUtc.Local().Format("Jan 2, 2006 3:04 PM")
		}

		msgString := fmt.Sprintf("[%s]%s [%s][%s]: %s", color, senderUsername, dateString, theme.ChatTextColor.CSS(), msg.Content)
		fmt.Fprintln(w, msgString)
	}

	page.textView.ScrollToEnd()

	// Tell the server that this is the active channel
	page.feedClient.SendFeedMessage(chat.FEED_MESSAGE_TYPE_SET_ACTIVE_CHANNEL_REQUEST, &chat.SetActiveChannelRequest{
		ChannelId: channel.Id,
	})

	page.textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyPgUp {
			// scroll up 10 lines
			r, _ := page.textView.GetScrollOffset()

			if r == 0 {
				if !entireConversationLoaded {
					getChannelMessagesResult := page.brochatClient.GetChannelMessages(accessToken, chatParam.channel_id,
						chat.GetChannelMessages_Page(1),
						chat.GetChannelMessages_PageSize(pageSize),
						chat.GetChannelMessages_BeforeMessage(oldestMessageId))

					err = getChannelMessagesResult.Err()

					if err != nil {
						if len(getChannelMessagesResult.ErrorDetails) > 0 {
							nav.Alert("home:chat:alert:err", getChannelMessagesResult.ErrorDetails[0])
							return nil
						}

						if getChannelMessagesResult.ResponseCode == chat.BROCHAT_RESPONSE_CODE_FORBIDDEN_ERROR {
							nav.Alert("home:chat:alert:err", FORBIDDEN_OPERATION_ERROR_MESSAGE)
							return nil
						}

						nav.AlertFatal(app, "home:chat:alert:err", err.Error())
						return nil
					}

					messages := getChannelMessagesResult.Content

					page.mu.Lock()
					defer page.mu.Unlock()

					if len(messages) < pageSize {
						entireConversationLoaded = true
					} else {
						oldestMessageId = messages[len(messages)-1].Id
					}

					oldText := []byte(page.textView.GetText(false))

					// Prepend the messages to the text view
					writer := page.textView.BatchWriter()
					defer writer.Close()

					writer.Clear()

					for i := len(messages) - 1; i >= 0; i-- {
						// Write the messages to the text view
						var senderUsername string
						var msg = messages[i]

						color := colorManifest[msg.SenderUserId]

						for _, u := range channel.Users {
							if u.Id == msg.SenderUserId {
								senderUsername = u.Username
								break
							}
						}

						// If for some reason the user info is not found just make the username "Unknown User"
						if senderUsername == "" {
							senderUsername = "Unknown User"
						}

						// If the color is not found then just make it red
						if color == "" {
							color = "#FF0000"
						}

						var dateString string

						// If the message is from a date in the past (not today) then format the date string differently
						if msg.RecievedAtUtc.Local().Day() == time.Now().Day() {
							dateString = msg.RecievedAtUtc.Local().Format(time.Kitchen)
						} else {
							dateString = msg.RecievedAtUtc.Local().Format("Jan 2, 2006 3:04 PM")
						}

						msgString := fmt.Sprintf("[%s]%s [%s][%s]: %s", color, senderUsername, dateString, theme.ChatTextColor.CSS(), msg.Content)
						fmt.Fprintln(writer, msgString)
					}

					writer.Write(oldText)

					// Scroll to the top if there are less than 10 messages otherwise scroll up the normal 10 lines
					if len(messages) > 10 {
						page.textView.ScrollTo(len(messages)-10, 0)
					} else {
						page.textView.ScrollToBeginning()
					}
				}

				return nil
			}

			page.textView.ScrollTo(r-10, 0)
			return nil
		} else if event.Key() == tcell.KeyPgDn {
			r, _ := page.textView.GetScrollOffset()
			page.textView.ScrollTo(r+10, 0)
			return nil
		} else if event.Key() == tcell.KeyEnter {
			text := page.textArea.GetText()

			if len(text) > 0 {

				isMacro, macroType := chat.IsMacro(text)

				if isMacro {
					page.feedClient.SendFeedMessage(chat.FEED_MESSAGE_TYPE_MACRO_REQUEST, chat.MacroRequest{
						Type:      macroType,
						Body:      text,
						ChannelId: channel.Id,
					})
				} else {
					page.feedClient.SendFeedMessage(chat.FEED_MESSAGE_TYPE_CHAT_MESSAGE_REQUEST, chat.ChatMessageRequest{
						ChannelId: channel.Id,
						Content:   text,
					})
				}

				page.textArea.SetText("", false)
			}

			return nil
		} else if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(chatParam.returnPage, nil)
		}

		return event
	})

	// Start the listener for channel updates
	go func() {
		subscriptionId, channelUpdateChannel := page.feedClient.SubscribeToChannelUpdates()

		defer page.feedClient.UnsubscribeFromChannelUpdates(subscriptionId)

		for {
			select {
			case <-pageContext.Done():
				return
			case eventChannelId := <-channelUpdateChannel:
				if eventChannelId == channel.Id {
					page.mu.Lock()
					defer page.mu.Unlock()

					accessToken, ok := appContext.GetAccessToken()

					if !ok {
						log.Println("No valid authentication information available for channel update event processing")
						appContext.CancelUserSession()
						return
					}

					getChannelResult := page.brochatClient.GetChannel(accessToken, channel.Id)

					err := getChannelResult.Err()

					if err != nil {
						log.Printf("Error getting channel during channel update event processing: %s", err.Error())
						return
					}

					newChannel := getChannelResult.Content

					usersForManifest := newChannel.Users

					for _, u := range newChannel.Users {
						// if the user is not in the manifest then add them
						if _, ok := colorManifest[u.Id]; !ok {
							usersForManifest = append(usersForManifest, u)
						}
					}

					colorManifest = getColorManifest(usersForManifest, theme)

					channel = newChannel
				}
			}
		}
	}()

	// Start the chat message listener
	go func(ch *chat.Channel, a *tview.Application, tv *tview.TextView) {
		subscriptionId, chatMsgChannel := page.feedClient.SubscribeToChatMessages()
		defer page.feedClient.UnsubscribeFromChatMessages(subscriptionId)

		for {
			select {
			case <-pageContext.Done():
				return
			case msg := <-chatMsgChannel:
				if msg.ChannelId == ch.Id {
					a.QueueUpdateDraw(func() {
						page.mu.Lock()
						defer page.mu.Unlock()

						var senderUsername string
						color := colorManifest[msg.SenderUserId]

						for _, u := range ch.Users {
							if u.Id == msg.SenderUserId {
								senderUsername = u.Username
								break
							}
						}

						// If for some reason the user info is not found just make the username "Unknown User"
						if senderUsername == "" {
							senderUsername = "Unknown User"
						}

						// If the color is not found then just make it red
						if color == "" {
							color = "#FF0000"
						}

						dateString := msg.RecievedAtUtc.Local().Format(time.Kitchen)

						msgString := fmt.Sprintf("[%s]%s [%s][%s]: %s", color, senderUsername, dateString, theme.ChatTextColor.CSS(), msg.Content)
						tv.Write([]byte(msgString + "\n"))
						tv.ScrollToEnd()
					})
				}
			}
		}
	}(&channel, app, page.textView)
}

// onPageClose is called when the chat page is navigated away from
func (page *ChatPage) onPageClose() {
	page.textView.Clear()
	page.textArea.SetText("", false)

	page.feedClient.SendFeedMessage(chat.FEED_MESSAGE_TYPE_SET_ACTIVE_CHANNEL_REQUEST, &chat.SetActiveChannelRequest{
		ChannelId: "NONE",
	})
}

// ChatPageParameters is load time parameters for the chat page
type ChatPageParameters struct {
	channel_id string
	title      string
	returnPage PageSlug
}

// getColorManifest takes in a slice of users and assigns each users and a color.
// The color manifest is a map of user ids to hex colors.
// The colors are assigned based upon the users index (position) in the slice.
func getColorManifest(users []chat.UserInfo, thm theme.Theme) map[string]string {

	var possibleColors = thm.ChatLabelColors

	colorManifest := make(map[string]string)

	for i, user := range users {
		if i >= len(possibleColors) {
			i = i % len(possibleColors)
		}

		colorManifest[user.Id] = possibleColors[i]
	}

	return colorManifest
}
