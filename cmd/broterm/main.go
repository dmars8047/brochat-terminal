package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dmars8047/broterm/internal/auth"
	"github.com/dmars8047/broterm/internal/bro"
	"github.com/dmars8047/broterm/internal/feed"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/broterm/internal/ui"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

func main() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Setup idam user auth client
	userAuthClient := auth.NewUserAuthClient(httpClient, "https://dev.marshall-labs.com")

	brochatClient := bro.NewBroChatClient(httpClient, "https://dev.marshall-labs.com")

	dialer := &websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	appContext := state.NewApplicationContext(context)

	feedClient := feed.NewFeedClient(dialer, "dev.marshall-labs.com")

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
