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

	// Setup idam user auth client
	userAuthClient := idam.NewUserAuthClient(httpClient, "http://localhost:8083")

	brochatClient := chat.NewBroChatClient(httpClient, "http://localhost:8083")

	dialer := &websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	appContext := state.NewApplicationContext(context)

	feedClient := feed.NewFeedClient(dialer, "localhost:8083")

	app := tview.NewApplication()

	pages := tview.NewPages()
	pages.SetBackgroundColor(ui.DefaultBackgroundColor)

	pageNav := ui.NewNavigator(pages)

	homeModule := ui.NewHomeModule(userAuthClient, app, pageNav, brochatClient, appContext, feedClient)
	homeModule.SetupHomePages()

	authModule := ui.NewAuthModule(userAuthClient, brochatClient, appContext, pageNav, app, feedClient)
	authModule.SetupAuthPages()

	// Start the application.
	err := app.SetRoot(pages, true).Run()

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
