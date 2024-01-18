package feed

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/dmars8047/broterm/internal/bro"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gorilla/websocket"
)

const (
	feedSuffix = "/api/brochat/connect"
	feedScheme = "wss"
)

type Client struct {
	dialer             *websocket.Dialer
	url                url.URL
	conn               *websocket.Conn
	ChatMessageChannel chan bro.ChatMessage
	Connected          bool
}

func NewFeedClient(dialer *websocket.Dialer, baseUrl string) *Client {
	return &Client{
		dialer:             dialer,
		url:                url.URL{Scheme: feedScheme, Host: baseUrl, Path: feedSuffix},
		ChatMessageChannel: make(chan bro.ChatMessage, 1),
	}
}

func (c *Client) Connect(userAuth state.UserAuth, ctx context.Context) error {
	c.ChatMessageChannel = make(chan bro.ChatMessage, 1)

	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+userAuth.AccessToken)

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
				var chatMessage bro.ChatMessage
				msgErr := json.Unmarshal(message, &chatMessage)

				if msgErr != nil {
					continue
				}

				c.ChatMessageChannel <- chatMessage
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

func (c *Client) Send(message bro.ChatMessage) error {
	if !c.Connected || c.conn == nil {
		return errors.New("not connected")
	}

	return c.conn.WriteJSON(message)
}
