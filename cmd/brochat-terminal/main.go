package main

import (
	"log"

	"github.com/dmars8047/brochat-terminal/internal/state"
	"github.com/dmars8047/brochat-terminal/internal/ui"
	"github.com/dmars8047/idam-service/pkg/idam"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	pages := tview.NewPages()
	pages.SetBackgroundColor(ui.DefaultBackgroundColor)

	// Setup idam user auth client
	userAuthClient := idam.NewUserAuthClient("https://dev.marshall-labs.com")

	appState := state.NewApplicationState()

	homeModule := ui.NewHomeModule(userAuthClient, appState)
	homeModule.SetupMenuPage(app, pages)

	authModule := ui.NewAuthModule(userAuthClient, appState)
	authModule.SetupAuthPages(app, pages)

	// Start the application.
	err := app.SetRoot(pages, true).Run()

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
