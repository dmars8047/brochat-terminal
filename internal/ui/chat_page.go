package ui

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const CHAT_PAGE PageSlug = "chat"

// ChatPage is the chat page
type ChatPage struct {
	brochatClient *chat.BroChatClient
	feedClient    *state.FeedClient
	textView      *tview.TextView
	textArea      *tview.TextArea
}

// NewChatPage creates a new chat page
func NewChatPage(brochatClient *chat.BroChatClient, feedClient *state.FeedClient) *ChatPage {
	return &ChatPage{
		brochatClient: brochatClient,
		feedClient:    feedClient,
		textView:      tview.NewTextView(),
		textArea:      tview.NewTextArea(),
	}
}

// Setup configures the chat page and registers it with the page navigator
func (page *ChatPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {

	page.textView.SetDynamicColors(true)
	page.textView.SetBorder(true)
	page.textView.SetScrollable(true)
	page.textView.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)

	page.textView.SetChangedFunc(func() {
		app.Draw()
	})

	page.textArea.SetTextStyle(TEXT_AREA_STYLE)
	page.textArea.SetBorder(true)
	page.textArea.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvInstructions.SetTextColor(tcell.ColorWhite)

	tvInstructions.SetText("(enter) Send - (esc) Back")

	grid := tview.NewGrid()

	grid.SetRows(0, 6, 2)
	grid.SetColumns(0)

	grid.AddItem(page.textView, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(page.textArea, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 2, 0, 1, 1, 0, 0, false)

	pageContext, cancel := context.WithCancel(appContext.Context)

	nav.Register(CHAT_PAGE, grid, true, false,
		func(param interface{}) {
			pageContext, cancel = context.WithCancel(appContext.Context)
			page.onPageLoad(param, app, appContext, nav, pageContext)
		},
		func() {
			cancel()
			page.onPageClose(appContext)
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

	// Get the color manifest
	colorManifest := getColorManifest(channel.Users)

	// Get the channel messages
	getChannelMessagesResult := page.brochatClient.GetChannelMessages(accessToken, chatParam.channel_id)

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

		msgString := fmt.Sprintf("[%s]%s [%s][%s]: %s", color, senderUsername, dateString, "#FFFFFF", msg.Content)
		fmt.Fprintln(w, msgString)
	}

	page.textView.ScrollToEnd()

	brochatUser := appContext.GetBrochatUser()

	// Set the chat context
	appContext.SetChatSession(&channel)

	page.feedClient.SendFeedMessage(chat.FEED_MESSAGE_TYPE_SET_ACTIVE_CHANNEL_REQUEST, &chat.SetActiveChannelRequest{
		ChannelId: channel.Id,
	})

	page.textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			text := page.textArea.GetText()

			if len(text) > 0 {
				page.feedClient.SendFeedMessage(chat.FEED_MESSAGE_TYPE_CHAT_MESSAGE_REQUEST, chat.ChatMessage{
					ChannelId:    channel.Id,
					Content:      text,
					SenderUserId: brochatUser.Id,
				})

				page.textArea.SetText("", false)
			}

			return nil
		} else if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(chatParam.returnPage, nil)
		}

		return event
	})

	mux := sync.RWMutex{}

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

					mux.Lock()
					colorManifest = getColorManifest(usersForManifest)
					mux.Unlock()

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
						var senderUsername string

						mux.RLock()
						color := colorManifest[msg.SenderUserId]
						mux.RUnlock()

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

						msgString := fmt.Sprintf("[%s]%s [%s][%s]: %s", color, senderUsername, dateString, "#FFFFFF", msg.Content)
						tv.Write([]byte(msgString + "\n"))
						tv.ScrollToEnd()
					})
				}
			}
		}
	}(&channel, app, page.textView)
}

// onPageClose is called when the chat page is navigated away from
func (page *ChatPage) onPageClose(
	appContext *state.ApplicationContext) {

	page.textView.Clear()
	page.textArea.SetText("", false)

	page.feedClient.SendFeedMessage(chat.FEED_MESSAGE_TYPE_SET_ACTIVE_CHANNEL_REQUEST, &chat.SetActiveChannelRequest{
		ChannelId: "NONE",
	})

	appContext.CancelChatSession()
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
func getColorManifest(users []chat.UserInfo) map[string]string {
	var possibleColors = []string{
		"#33DA7A", // Light Green
		"#C061CB", // Lilac
		"#FF6B30", // Orange
		"#5928ED", // Purple
		"#00FFFF", // Cyan
		"#FF5555", // Light Red
		"#FAEC34", // Yellow
		"#FFAAFF", // Light Pink
	}

	colorManifest := make(map[string]string)

	for i, user := range users {
		if i >= len(possibleColors) {
			i = i % len(possibleColors)
		}

		colorManifest[user.Id] = possibleColors[i]
	}

	return colorManifest
}
