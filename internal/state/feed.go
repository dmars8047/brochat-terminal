package state

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	feedSuffix = "/api/brochat/connect"
	feedScheme = "ws"
)

type FeedClient struct {
	appContext                *ApplicationContext
	broChatClient             *chat.BroChatClient
	dialer                    *websocket.Dialer
	url                       url.URL
	conn                      *websocket.Conn
	chatMessageChannels       map[string]chan chat.ChatMessage
	userProfileUpdateChannels map[string]chan chat.UserProfileUpdateCode
	channelUpdateChannels     map[string]chan string
	Closed                    bool
	mu                        sync.RWMutex
}

// NewFeedClient creates a new instance of the feed client.
func NewFeedClient(dialer *websocket.Dialer, baseUrl string, broChatClient *chat.BroChatClient, appContext *ApplicationContext) *FeedClient {
	return &FeedClient{
		broChatClient:             broChatClient,
		dialer:                    dialer,
		url:                       url.URL{Scheme: feedScheme, Host: baseUrl, Path: feedSuffix},
		chatMessageChannels:       make(map[string]chan chat.ChatMessage, 0),
		userProfileUpdateChannels: make(map[string]chan chat.UserProfileUpdateCode, 0),
		channelUpdateChannels:     make(map[string]chan string, 0),
		mu:                        sync.RWMutex{},
		appContext:                appContext,
	}
}

// SubscribeToChatMessages subscribes to chat messages and returns a channel to receive messages on.
// The returned string is the subscription ID and is used to unsubscribe from chat messages.
// The returned channel will be closed when the subscription is removed. Suggested usage is to defer the call to UnsubscribeFromChatMessages.
func (c *FeedClient) SubscribeToChatMessages() (string, <-chan chat.ChatMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := uuid.NewString()
	ch := make(chan chat.ChatMessage)
	c.chatMessageChannels[id] = ch
	return id, ch
}

// UnsubscribeFromChatMessages unsubscribes from chat messages.
func (c *FeedClient) UnsubscribeFromChatMessages(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.chatMessageChannels[id]

	if !ok {
		return
	}

	close(ch)
	delete(c.chatMessageChannels, id)
}

// SubscribeToUserProfileUpdates subscribes to user profile updates and returns a channel to receive updates on.
// The returned string is the subscription ID and is used to unsubscribe from user profile updates.
// The returned channel will be closed when the subscription is removed. Suggested usage is to defer the call to UnsubscribeFromUserProfileUpdates.
func (c *FeedClient) SubscribeToUserProfileUpdates() (string, <-chan chat.UserProfileUpdateCode) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := uuid.NewString()
	ch := make(chan chat.UserProfileUpdateCode)
	c.userProfileUpdateChannels[id] = ch
	return id, ch
}

// UnsubscribeFromUserProfileUpdates unsubscribes from user profile updates.
func (c *FeedClient) UnsubscribeFromUserProfileUpdates(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.userProfileUpdateChannels[id]

	if !ok {
		return
	}

	close(ch)
	delete(c.userProfileUpdateChannels, id)
}

// SubscribeToChannelUpdates subscribes to channel updates and returns a channel to receive updates on.
// The returned string is the subscription ID and is used to unsubscribe from channel updates.
// The returned channel will be closed when the subscription is removed. Suggested usage is to defer the call to UnsubscribeFromChannelUpdates.
func (c *FeedClient) SubscribeToChannelUpdates() (string, <-chan string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := uuid.NewString()
	ch := make(chan string)

	c.channelUpdateChannels[id] = ch

	return id, ch
}

// UnsubscribeFromChannelUpdates unsubscribes from channel updates.
func (c *FeedClient) UnsubscribeFromChannelUpdates(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.channelUpdateChannels[id]

	if !ok {
		return
	}

	close(ch)
	delete(c.channelUpdateChannels, id)
}

func (c *FeedClient) Connect() error {

	accessToken, ok := c.appContext.GetAccessToken()

	if !ok {
		return errors.New("no valid authentication information available for feed connection")
	}

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+accessToken)

	conn, _, err := c.dialer.Dial(c.url.String(), headers)

	if err != nil {
		return err
	}

	c.Closed = false
	c.conn = conn

	// The default close handler will automatically respond to close messages from the server
	defaultCloseHandler := c.conn.CloseHandler()

	c.conn.SetCloseHandler(func(code int, text string) error {
		c.Closed = true
		return defaultCloseHandler(code, text)
	})

	go func() {
		for {
			messageType, message, err := c.conn.ReadMessage()

			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) || err == websocket.ErrCloseSent {
					log.Printf("Websocket connection to %s closed", c.url.String())
				} else {
					log.Printf("Error reading message from websocket: %s", err.Error())
				}
				return
			}

			switch messageType {
			case websocket.TextMessage:
				var feedMessage chat.FeedMessage
				msgErr := json.Unmarshal(message, &feedMessage)

				if msgErr != nil {
					log.Printf("Error unmarshaling feed message: %s", msgErr.Error())
					continue
				}

				switch feedMessage.Type {
				case chat.FEED_MESSAGE_TYPE_CHANNEL_UPDATED:
					var channelUpdatedEvent chat.ChannelUpdatedEvent

					chtMsgErr := json.Unmarshal(feedMessage.Content, &channelUpdatedEvent)

					if chtMsgErr != nil {
						log.Printf("Error unmarshaling channel updated event during channel updated event processing: %s", chtMsgErr.Error())
						continue
					}

					c.mu.RLock()

					for ch := range c.channelUpdateChannels {
						c.channelUpdateChannels[ch] <- channelUpdatedEvent.ChannelId
					}

					c.mu.RUnlock()
				case chat.FEED_MESSAGE_TYPE_CHAT_MESSAGE:
					var chatMessage chat.ChatMessage

					chtMsgErr := json.Unmarshal(feedMessage.Content, &chatMessage)

					if chtMsgErr != nil {
						log.Printf("Error unmarshaling chat message during chat message event processing: %s", chtMsgErr.Error())
						continue
					}

					c.mu.RLock()

					for ch := range c.chatMessageChannels {
						c.chatMessageChannels[ch] <- chatMessage
					}

					c.mu.RUnlock()
				case chat.FEED_MESSAGE_TYPE_USER_PROFILE_UPDATED:
					brochatUser := c.appContext.GetBrochatUser()

					accessToken, ok := c.appContext.GetAccessToken()

					if !ok {
						log.Println("No valid authentication information available for user profile updated event processing")
						c.appContext.CancelUserSession()
						return
					}

					result := c.broChatClient.GetUser(accessToken, brochatUser.Id)

					err = result.Err()

					if err != nil {
						log.Printf("An error occurred during the processing of a user profile updated event. "+
							"The call to retrieve user data resulted in the following error: %s", err.Error())

						continue
					}

					usrProfile := result.Content

					c.appContext.SetBrochatUser(usrProfile)

					var userProfileUpdatedEvent chat.UserProfileUpdatedEvent

					chtMsgErr := json.Unmarshal(feedMessage.Content, &userProfileUpdatedEvent)

					if chtMsgErr != nil {
						log.Printf("Error unmarshaling user profile updated event during user profile updated event processing: %s", chtMsgErr.Error())
						continue
					}

					c.mu.RLock()

					for ch := range c.userProfileUpdateChannels {
						c.userProfileUpdateChannels[ch] <- userProfileUpdatedEvent.UpdateCode
					}

					c.mu.RUnlock()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-c.appContext.userSession.context.Done():
				c.mu.Lock()
				defer c.mu.Unlock()

				// Close all chat message channels
				for ch := range c.chatMessageChannels {
					close(c.chatMessageChannels[ch])
				}

				clear(c.chatMessageChannels)

				// Close all user profile update channels
				for ch := range c.userProfileUpdateChannels {
					close(c.userProfileUpdateChannels[ch])
				}

				clear(c.userProfileUpdateChannels)

				// Close the connection
				defer func() {
					log.Printf("Closing websocket connection to %s", c.url.String())
					if err := c.conn.Close(); err != nil {
						log.Println("websocket close error:", err)
					}
				}()

				// Create a close message
				msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Client closed connection.")
				// Write the close message to the server
				if err := c.conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
					return
				}

				// Wait for a close message from the server or timeout after 30 seconds
				ctx, cancel := context.WithTimeout(c.appContext.Context, 30*time.Second)
				defer cancel()

				done := make(chan struct{})

				go func(ctx context.Context) {
					defer close(done)
					for {
						select {
						case <-ctx.Done():
							// Context was cancelled or timeout occurred, stop the goroutine
							return
						default:
							if c.Closed {
								return
							}
						}
					}
				}(ctx)

				select {
				case <-done:
				case <-ctx.Done():
					// Timeout occurred or context was cancelled
				}

				return
			case <-time.After(30 * time.Second):
				if c.conn != nil && !c.Closed && c.appContext.userSession != nil {
					c.conn.WriteMessage(websocket.PingMessage, nil)
				}
			}
		}
	}()

	return nil
}

func (c *FeedClient) SendFeedMessage(messageType chat.FeedMessageType, content interface{}) error {
	if c.Closed || c.conn == nil || c.appContext.userSession == nil {
		return errors.New("feed connection failure")
	}

	feedMessage, err := chat.NewFeedMessageJSON(messageType, content)

	if err != nil {
		return err
	}

	return c.conn.WriteJSON(feedMessage)
}
