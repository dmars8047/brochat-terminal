package ui

import (
	"fmt"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const CHAT_PAGE PageSlug = "chat"

// ChatPage is the chat page
type ChatPage struct {
	brochatClient *chat.BroChatUserClient
	feedClient    *state.FeedClient
	textView      *tview.TextView
	textArea      *tview.TextArea
}

// NewChatPage creates a new chat page
func NewChatPage(brochatClient *chat.BroChatUserClient, feedClient *state.FeedClient) *ChatPage {
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

	nav.Register(CHAT_PAGE, grid, true, false,
		func(param interface{}) {
			page.onPageLoad(param, app, appContext, nav)
		},
		func() {
			page.onPageClose(appContext)
		})
}

// onPageLoad is called when the chat page is navigated to
func (page *ChatPage) onPageLoad(param interface{},
	app *tview.Application,
	appContext *state.ApplicationContext,
	nav *PageNavigator) {

	// The param should be a ChatParams struct
	chatParam, ok := param.(ChatPageParameters)

	if !ok {
		nav.AlertFatal(app, "home:chat:alert:err", "Application State Error - Could not get chat params.")
		return
	}

	// Get the channel
	channel, err := page.brochatClient.GetChannelManifest(appContext.GetAuthInfo(), chatParam.channel_id)

	if err != nil {
		nav.Alert("home:chat:alert:err", err.Error())
		return
	}

	if channel.Type == chat.CHANNEL_TYPE_DIRECT_MESSAGE {
		page.textView.SetTitle(fmt.Sprintf(" %s - %s ", channel.Users[0].Username, channel.Users[1].Username))
	} else if chatParam.title != "" {
		page.textView.SetTitle(fmt.Sprintf(" %s ", chatParam.title))
	}

	// Get the channel messages
	messages, err := page.brochatClient.GetChannelMessages(appContext.GetAuthInfo(), chatParam.channel_id)

	if err != nil {
		nav.AlertFatal(app, "home:chat:alert:err", err.Error())
		return
	}

	// Write the messages to the text view
	w := page.textView.BatchWriter()
	defer w.Close()
	w.Clear()

	for i := len(messages) - 1; i >= 0; i-- {
		// Write the messages to the text view
		var senderUsername string
		var msg = messages[i]

		color := CHAT_COLOR_TWO

		if msg.SenderUserId == appContext.BrochatUser.Id {
			color = CHAT_COLOR_ONE
		}

		for _, u := range channel.Users {
			if u.Id == msg.SenderUserId {
				senderUsername = u.Username
				break
			}
		}

		if senderUsername == "" {
			// Make a request to get the channel manifest
			// This is a fallback in case the user info is not in the channel manifest for some reason (maybe the just joined the channel)
			newChannel, getChanErr := page.brochatClient.GetChannelManifest(appContext.GetAuthInfo(), chatParam.channel_id)

			if getChanErr != nil {
				nav.Alert("home:chat:alert:err", getChanErr.Error())
				return
			}

			channel = newChannel
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

	// Set the chat context
	appContext.SetChatSession(channel)

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
					SenderUserId: appContext.BrochatUser.Id,
				})

				page.textArea.SetText("", false)
			}

			return nil
		} else if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(HOME_PAGE, nil)
		}

		return event
	})

	// Start the chat message listener
	go func(ch *chat.Channel, a *tview.Application, tv *tview.TextView) {
		for msg := range page.feedClient.ChatMessageChannel {
			if msg.ChannelId == ch.Id {
				a.QueueUpdateDraw(func() {
					var senderUsername string
					color := "#C061CB"

					if msg.SenderUserId == appContext.BrochatUser.Id {
						color = "#33DA7A"
					}

					for _, u := range ch.Users {
						if u.Id == msg.SenderUserId {
							senderUsername = u.Username
							break
						}
					}

					if senderUsername == "" {
						// Make a request to get the channel manifest
						// This is a fallback in case the user info is not in the channel manifest for some reason (maybe the just joined the channel)
						newChannel, getChanErr := page.brochatClient.GetChannelManifest(appContext.GetAuthInfo(), ch.Id)

						if getChanErr != nil {
							nav.Alert("home:chat:alert:err", getChanErr.Error())
							return
						}

						channel = newChannel
					}

					dateString := msg.RecievedAtUtc.Local().Format(time.Kitchen)

					msgString := fmt.Sprintf("[%s]%s [%s][%s]: %s", color, senderUsername, dateString, "#FFFFFF", msg.Content)
					tv.Write([]byte(msgString + "\n"))
					tv.ScrollToEnd()
				})
			}
		}
	}(channel, app, page.textView)
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
}
