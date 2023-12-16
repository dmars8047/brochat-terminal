package main

import (
	"log"

	"github.com/dmars8047/brochat-terminal/internal/ui"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	pages := tview.NewPages()
	pages.SetBackgroundColor(ui.DefaultBackgroundColor)

	authModule := ui.NewAuthModule()
	authModule.SetupAuthPages(app, pages)

	// Start the application.
	err := app.SetRoot(pages, true).Run()

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
