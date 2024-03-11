package state

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/gorilla/websocket"
)

const (
	feedSuffix = "/api/brochat/connect"
	feedScheme = "ws"
)

type FeedClient struct {
	brochatUserClient  *chat.BroChatUserClient
	dialer             *websocket.Dialer
	url                url.URL
	conn               *websocket.Conn
	ChatMessageChannel chan chat.ChatMessage
	Closed             bool
}

func NewFeedClient(dialer *websocket.Dialer, baseUrl string, brochatUserClient *chat.BroChatUserClient) *FeedClient {
	return &FeedClient{
		brochatUserClient:  brochatUserClient,
		dialer:             dialer,
		url:                url.URL{Scheme: feedScheme, Host: baseUrl, Path: feedSuffix},
		ChatMessageChannel: make(chan chat.ChatMessage, 1),
	}
}

func (c *FeedClient) Connect(appContext *ApplicationContext) error {
	c.ChatMessageChannel = make(chan chat.ChatMessage, 1)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+appContext.GetAuthInfo().AccessToken)

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
					log.Printf("Error unmarshalling feed message: %s", msgErr.Error())
					continue
				}

				switch feedMessage.Type {
				case chat.FEED_MESSAGE_TYPE_CHAT_MESSAGE:
					var chatMessage chat.ChatMessage

					chtMsgErr := json.Unmarshal(feedMessage.Content, &chatMessage)

					if chtMsgErr != nil {
						// TODO: figure out what to do with these errors.
						continue
					}

					c.ChatMessageChannel <- chatMessage
				case chat.FEED_MESSAGE_TYPE_USER_PROFILE_UPDATED:
					usrProfile, err := c.brochatUserClient.GetUser(appContext.GetAuthInfo(), appContext.BrochatUser.Id)

					if err != nil {
						log.Printf("An error occurred during the processing of a user profile updated event. "+
							"The call to retrieve user data resulted in the following error: %s", err.Error())
						continue
					}

					appContext.BrochatUser = usrProfile
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-appContext.userSession.context.Done():
				close(c.ChatMessageChannel)

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
				ctx, cancel := context.WithTimeout(appContext.Context, 30*time.Second)
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
				c.conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

	return nil
}

func (c *FeedClient) SendFeedMessage(messageType chat.FeedMessageType, content interface{}) error {
	if c.Closed || c.conn == nil {
		return errors.New("feed connection failure")
	}

	feedMessage, err := chat.NewFeedMessageJSON(messageType, content)

	if err != nil {
		return err
	}

	return c.conn.WriteJSON(feedMessage)
}
