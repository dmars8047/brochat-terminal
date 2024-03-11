package state

import (
	"context"
	"time"

	"github.com/dmars8047/brolib/chat"
)

type ApplicationContext struct {
	// The context for the application
	Context     context.Context
	BrochatUser *chat.User
	chatSession *ChatSession
	userSession *UserSession
}

func NewApplicationContext(context context.Context) *ApplicationContext {
	return &ApplicationContext{
		Context: context,
	}
}

func (appContext *ApplicationContext) GetUserAuth() UserAuth {
	if appContext.userSession == nil {
		return UserAuth{}
	}

	return appContext.userSession.Auth
}

// GetAuthInfo returns the auth info for the user session. This is used to authenticate calls made using the brochat client.
func (appContext *ApplicationContext) GetAuthInfo() *chat.AuthInfo {
	if appContext.userSession == nil {
		return nil
	}

	return &chat.AuthInfo{
		AccessToken: appContext.userSession.Auth.AccessToken,
		TokenType:   "Bearer",
	}
}

func (appContext *ApplicationContext) SetUserSession(auth UserAuth) {
	context, cancel := context.WithCancel(appContext.Context)

	appContext.userSession = &UserSession{
		Auth:       auth,
		context:    context,
		cancelFunc: cancel,
	}

	go func(ctx *ApplicationContext) {
		select {
		case <-ctx.userSession.context.Done():
			return
		case <-time.After(time.Until(ctx.userSession.Auth.TokenExpiration)):
			ctx.CancelUserSession()
			return
		}
	}(appContext)
}

func (appContext *ApplicationContext) SetChatSession(channel *chat.Channel) {
	context, cancel := context.WithCancel(appContext.Context)

	appContext.chatSession = &ChatSession{
		channel:    channel,
		context:    context,
		cancelFunc: cancel,
	}
}

func (appContext *ApplicationContext) CancelUserSession() {
	appContext.userSession.cancelFunc()
}

func (appContext *ApplicationContext) CancelChatSession() {
	appContext.chatSession.cancelFunc()
	appContext.chatSession = nil
}

type UserAuth struct {
	AccessToken     string
	TokenExpiration time.Time
}

type UserSession struct {
	Auth       UserAuth
	context    context.Context
	cancelFunc context.CancelFunc
}

type ChatSession struct {
	channel    *chat.Channel
	context    context.Context
	cancelFunc context.CancelFunc
}
