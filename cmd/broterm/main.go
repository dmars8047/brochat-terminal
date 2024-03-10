package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/feed"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/broterm/internal/ui"
	"github.com/dmars8047/idamlib/idam"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

func main() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Setup dependencies
	userAuthClient := idam.NewUserAuthClient(httpClient, "http://localhost:8083")

	brochatClient := chat.NewBroChatClient(httpClient, "http://localhost:8083")

	dialer := &websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	feedClient := feed.NewFeedClient(dialer, "localhost:8083")

	// Setup the application context
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	appContext := state.NewApplicationContext(context)

	// Configure the application
	app := tview.NewApplication()

	// Setup the page navigator
	nav := ui.NewNavigator()

	// Setup the welcome page
	welcomePage := ui.NewWelcomePage()
	welcomePage.Setup(app, appContext, nav)

	// Setup the registration page
	registrationPage := ui.NewRegistrationPage(userAuthClient)
	registrationPage.Setup(app, appContext, nav)

	// Setup the login page
	loginPage := ui.NewLoginPage(userAuthClient, brochatClient, feedClient)
	loginPage.Setup(app, appContext, nav)

	// Setup the forgot password page
	forgotPasswordPage := ui.NewForgotPasswordPage(userAuthClient)
	forgotPasswordPage.Setup(app, appContext, nav)

	// Setup the chat page
	chatPage := ui.NewChatPage(brochatClient, feedClient)
	chatPage.Setup(app, appContext, nav)

	// Setup the home page
	homePage := ui.NewHomePage(userAuthClient)
	homePage.Setup(app, appContext, nav)

	// Setup the friends list page
	friendsListPage := ui.NewFriendsListPage(brochatClient)
	friendsListPage.Setup(app, appContext, nav)

	// Setup the find a friend page
	findAFriendPage := ui.NewFindAFriendPage(brochatClient)
	findAFriendPage.Setup(app, appContext, nav)

	// Setup the accept friend request page
	acceptFriendRequestPage := ui.NewAcceptFriendRequestPage(brochatClient)
	acceptFriendRequestPage.Setup(app, appContext, nav)

	// Setup the room list page
	roomListPage := ui.NewRoomListPage(brochatClient)
	roomListPage.Setup(app, appContext, nav)

	// Setup the room editor page
	roomEditorPage := ui.NewRoomEditorPage(brochatClient)
	roomEditorPage.Setup(app, appContext, nav)

	// Setup the room finder page
	roomFinderPage := ui.NewRoomFinderPage(brochatClient)
	roomFinderPage.Setup(app, appContext, nav)

	// Start the application.
	err := app.SetRoot(nav.Pages, true).Run()

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
