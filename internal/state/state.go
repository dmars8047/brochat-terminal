package state

import (
	"context"
	"sync"
	"time"

	"github.com/dmars8047/brolib/chat"
)

// ApplicationContext is the context for the application
// It manages the context (lifetime) for the user and chat sessions
// It also contains references to the logged in user and their authentication information.
type ApplicationContext struct {
	// The context for the application
	Context     context.Context
	brochatUser *chat.User
	chatSession *ChatSession
	userSession *UserSession
	mut         sync.RWMutex
}

func NewApplicationContext(context context.Context) *ApplicationContext {
	return &ApplicationContext{
		Context: context,
	}
}

func (appContext *ApplicationContext) GetBrochatUser() chat.User {
	return *appContext.brochatUser
}

func (appContext *ApplicationContext) SetBrochatUser(user chat.User) {
	appContext.mut.Lock()
	defer appContext.mut.Unlock()
	appContext.brochatUser = &user
}

func (appContext *ApplicationContext) GetUserAuth() UserAuth {
	appContext.mut.RLock()
	defer appContext.mut.RUnlock()

	if appContext.userSession == nil {
		return UserAuth{}
	}

	return appContext.userSession.Auth
}

// GetAccessToken returns the auth info for the user session. This is used to authenticate calls made using the brochat client.
func (appContext *ApplicationContext) GetAccessToken() (string, bool) {
	if appContext.userSession == nil {
		return "", false
	}

	// Check if the token has expired
	if time.Now().After(appContext.userSession.Auth.TokenExpiration) {
		return "", false
	}

	return appContext.userSession.Auth.AccessToken, true
}

// SetUserSession sets the user session
// This will cancel the previous user session if it exists
// It will also create a new context for the user session
// The user session will be cancelled after the token expires
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
	context, cancel := context.WithCancel(appContext.userSession.context)

	appContext.chatSession = &ChatSession{
		channel:    channel,
		context:    context,
		cancelFunc: cancel,
	}
}

// CancelUserSession cancels the user session
// This will cancel the context and set the user session to nil
// Calling this method before the user session is set will do nothing
func (appContext *ApplicationContext) CancelUserSession() {
	if appContext.userSession == nil {
		return
	}

	appContext.userSession.cancelFunc()
	appContext.userSession = nil
}

// CancelChatSession cancels the chat session
// This will cancel the context and set the chat session to nil
// Calling this method before the chat session is set will do nothing
func (appContext *ApplicationContext) CancelChatSession() {
	if appContext.chatSession == nil {
		return
	}

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
