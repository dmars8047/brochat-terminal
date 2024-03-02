package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Confirm(pages *tview.Pages, id string, massage string, yesFunc func()) *tview.Pages {
	return pages.AddPage(
		id,
		tview.NewModal().
			SetText(massage).
			AddButtons([]string{"Yes", "No"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					yesFunc()
				}
				pages.HidePage(id).RemovePage(id)
			}),
		false,
		true,
	)
}

func Alert(pages *tview.Pages, id string, message string) *tview.Pages {
	return pages.AddPage(
		id,
		tview.NewModal().
			SetText(message).
			AddButtons([]string{"Close"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.HidePage(id).RemovePage(id)
			}),
		false,
		true,
	)
}

func AlertWithDoneFunc(pages *tview.Pages, id string, message string, doneFunc func(buttonIndex int, buttonLabel string)) *tview.Pages {
	return pages.AddPage(
		id,
		tview.NewModal().
			SetText(message).
			AddButtons([]string{"Close"}).
			SetDoneFunc(doneFunc),
		false,
		true,
	)
}

func AlertFatal(app *tview.Application, pages *tview.Pages, id string, message string) *tview.Pages {
	return pages.AddPage(
		id,
		tview.NewModal().
			SetText("Fatal Error: "+message).
			AddButtons([]string{"Exit"}).
			SetBackgroundColor(DangerBackgroundColor).
			SetTextColor(tcell.ColorWhite).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.Stop()
			}),
		false,
		true,
	)
}

func AlertErrors(pages *tview.Pages, id, errMessage string, messages []string) {
	added := false

	for _, message := range messages {
		if len(message) > 2 {
			if !added {
				errMessage += "\n"
				added = true
			}
			val := strings.ToUpper(string(message[0])) + message[1:]
			errMessage += fmt.Sprintf("\n- %s", val)
		}
	}

	Alert(pages, id, errMessage)
}
