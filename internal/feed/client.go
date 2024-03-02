package feed

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gorilla/websocket"
)

const (
	feedSuffix = "/api/brochat/connect"
	feedScheme = "ws"
)

type Client struct {
	dialer             *websocket.Dialer
	url                url.URL
	conn               *websocket.Conn
	ChatMessageChannel chan chat.ChatMessage
	Connected          bool
}

func NewFeedClient(dialer *websocket.Dialer, baseUrl string) *Client {
	return &Client{
		dialer:             dialer,
		url:                url.URL{Scheme: feedScheme, Host: baseUrl, Path: feedSuffix},
		ChatMessageChannel: make(chan chat.ChatMessage, 1),
	}
}

func (c *Client) Connect(userAuth state.UserAuth, ctx context.Context) error {
	c.ChatMessageChannel = make(chan chat.ChatMessage, 1)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+userAuth.AccessToken)

	conn, _, err := c.dialer.Dial(c.url.String(), headers)

	if err != nil {
		return err
	}

	c.Connected = true
	c.conn = conn

	go func() {
		for {
			messageType, message, err := c.conn.ReadMessage()

			if err != nil {
				c.Connected = false

				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return
				}

				return
			}

			switch messageType {
			case websocket.TextMessage:
				var feedMessage chat.FeedMessage
				msgErr := json.Unmarshal(message, &feedMessage)

				if msgErr != nil {
					// TODO: figure out what to do with these errors.
					continue
				}

				switch feedMessage.Type {
				case chat.FEED_MESSAGE_TYPE_CHAT_MESSAGE:
					var chatMessage chat.ChatMessage

					chtMsgErr := json.Unmarshal(feedMessage.Content, &chatMessage)

					if chtMsgErr != nil {
						continue
					}

					c.ChatMessageChannel <- chatMessage
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(c.ChatMessageChannel)
				c.conn.Close()
				return
			case <-time.After(30 * time.Second):
				c.conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

	return nil
}

func (c *Client) SendFeedMessage(messageType chat.FeedMessageType, content interface{}) error {
	if !c.Connected || c.conn == nil {
		return errors.New("feed connection failure")
	}

	feedMessage, err := chat.NewFeedMessageJSON(messageType, content)

	if err != nil {
		return err
	}

	return c.conn.WriteJSON(feedMessage)
}
