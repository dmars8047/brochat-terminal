package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dmars8047/brochat-service/pkg/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/broterm/internal/ui"
	"github.com/dmars8047/idam-service/pkg/idam"
	"github.com/rivo/tview"
)

func main() {
	// Setup idam user auth client
	userAuthClient := idam.NewUserAuthClient("https://dev.marshall-labs.com")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	brochatClient := chat.NewBroChatClient(httpClient, "https://dev.marshall-labs.com")
	appState := state.NewApplicationState()

	app := tview.NewApplication()

	pages := tview.NewPages()
	pages.SetBackgroundColor(ui.DefaultBackgroundColor)

	pageNav := ui.NewNavigator(pages)

	homeModule := ui.NewHomeModule(userAuthClient, app, pageNav, brochatClient, appState)
	homeModule.SetupHomePages()

	authModule := ui.NewAuthModule(userAuthClient, brochatClient, appState, pageNav, app)
	authModule.SetupAuthPages()

	// Start the application.
	err := app.SetRoot(pages, true).Run()

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
