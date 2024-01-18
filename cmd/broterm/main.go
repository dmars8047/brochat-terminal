package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dmars8047/brochat-service/pkg/chat"
	"github.com/dmars8047/broterm/internal/feed"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/broterm/internal/ui"
	"github.com/dmars8047/idam-service/pkg/idam"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

func main() {
	// Setup idam user auth client
	userAuthClient := idam.NewUserAuthClient("https://dev.marshall-labs.com")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	brochatClient := chat.NewBroChatClient(httpClient, "https://dev.marshall-labs.com")

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
