package state

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/theme"
)

// ApplicationContext is the context for the application
// It manages the context (lifetime) for the user and chat sessions
// It also contains references to the logged in user and their authentication information.
type ApplicationContext struct {
	// The context for the application
	Context           context.Context
	brochatUser       *chat.User
	userSession       *UserSession
	mut               sync.RWMutex
	monitoringContext context.Context
	cancelMonitoring  context.CancelFunc
	theme             *theme.Theme
}

func NewApplicationContext(context context.Context) *ApplicationContext {
	return &ApplicationContext{
		Context: context,
		theme:   theme.NewTheme("satanic"),
	}
}

func (appContext *ApplicationContext) GetTheme() theme.Theme {
	appContext.mut.RLock()
	defer appContext.mut.RUnlock()
	return *appContext.theme
}

func (appContext *ApplicationContext) SetTheme(themeName string) {
	appContext.mut.Lock()
	defer appContext.mut.Unlock()
	appContext.theme = theme.NewTheme(themeName)
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
	appContext.mut.RLock()
	defer appContext.mut.RUnlock()

	if appContext.userSession == nil {
		return "", false
	}

	// Check if the token has expired
	if time.Now().After(appContext.userSession.Auth.TokenExpiration) {
		return "", false
	}

	return appContext.userSession.Auth.AccessToken, true
}

// SetUserSession sets the user session.
// This will cancel the previous user session if it exists.
// It will also create a new context for the user session.
// The user session will be cancelled after the token expires.
// The redirect function will be called when the token expires in a separate goroutine.
func (appContext *ApplicationContext) SetUserSession(auth UserAuth, redirect func()) {
	appContext.mut.Lock()
	defer appContext.mut.Unlock()

	if appContext.userSession != nil {
		appContext.cancelMonitoring()
	}

	userSessionContext, cancelUserSession := context.WithCancel(appContext.Context)

	appContext.userSession = &UserSession{
		Auth:    auth,
		context: userSessionContext,
		cancel:  cancelUserSession,
	}

	appContext.monitoringContext, appContext.cancelMonitoring = context.WithCancel(userSessionContext)

	go func() {
		select {
		case <-appContext.monitoringContext.Done():
			return
		case <-time.After(time.Until(appContext.userSession.Auth.TokenExpiration)):
			redirect()
			appContext.CancelUserSession()
			return
		}
	}()
}

// GenerateUserSessionBoundContextWithCancel generates a new context with cancel function that is bound to the lifetime of the user session.
// If the user session is not set, this will panic.
func (appContext *ApplicationContext) GenerateUserSessionBoundContextWithCancel() (context.Context, context.CancelFunc) {
	appContext.mut.RLock()
	defer appContext.mut.RUnlock()

	if appContext.userSession == nil {
		panic(errors.New("user session is not set"))
	}

	ctx, cancel := context.WithCancel(appContext.userSession.context)

	return ctx, cancel
}

// CancelUserSession cancels the user session
// This will cancel the context and set the user session to nil
// Calling this method before the user session is set will do nothing
func (appContext *ApplicationContext) CancelUserSession() {
	appContext.mut.Lock()
	defer appContext.mut.Unlock()

	if appContext.userSession == nil {
		return
	}

	cancel := appContext.userSession.cancel
	appContext.userSession = nil
	cancel()
}

type UserAuth struct {
	AccessToken     string
	TokenExpiration time.Time
}

type UserSession struct {
	Auth    UserAuth
	context context.Context
	cancel  context.CancelFunc
}
