package state

import (
	"context"
	"time"

	"github.com/dmars8047/broterm/internal/bro"
)

type Property = string

type ApplicationContext struct {
	// The context for the application
	Context     context.Context
	UserSession *UserSession
	BrochatUser *bro.User
	ChatSession *ChatSession
}

func NewApplicationContext(context context.Context) *ApplicationContext {
	return &ApplicationContext{
		Context: context,
	}
}

type UserSession struct {
	Auth       UserAuth
	Info       UserInfo
	Context    context.Context
	CancelFunc context.CancelFunc
}

type UserInfo struct {
	Username string
	Id       string
}

type UserAuth struct {
	AccessToken     string
	TokenExpiration time.Time
}

type ChatSession struct {
	Channel    *bro.ChannelManifest
	Context    context.Context
	CancelFunc context.CancelFunc
}

func NewChatSession(channel *bro.ChannelManifest, ctx context.Context) *ChatSession {
	context, cancel := context.WithCancel(ctx)

	return &ChatSession{
		Channel:    channel,
		Context:    context,
		CancelFunc: cancel,
	}
}