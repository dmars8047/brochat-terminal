package ui

import (
	"github.com/dmars8047/broterm/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const APP_SETTINGS_PAGE PageSlug = "app_settings"

// AppSettingsPage is the location where users can configure application level settings.
type AppSettingsPage struct {
	settingsForm *tview.Form
}

// NewAppSettingsPage creates a new instance of the application settings page
func NewAppSettingsPage() *AppSettingsPage {
	return &AppSettingsPage{
		settingsForm: tview.NewForm(),
	}
}

// Setup configures the application settings page and registers it with the page navigator
// The page includes a form which allows the user to set the following settings:
// The host address of the server
// The theme (default, america, matrix, halloween, and morning)
// The log and setting config file storage location
func (page *AppSettingsPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	const title = " BroChat - Application Settings "

	grid := tview.NewGrid()
	grid.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	grid.SetRows(4, 0, 1, 3, 4)
	grid.SetColumns(0, 70, 0)

	page.settingsForm.SetBackgroundColor(ACCENT_BACKGROUND_COLOR)
	page.settingsForm.SetFieldBackgroundColor(ACCENT_COLOR_TWO_COLOR)
	page.settingsForm.SetLabelColor(BROCHAT_YELLOW_COLOR)
	page.settingsForm.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignCenter)
	page.settingsForm.SetButtonStyle(DEFAULT_BUTTON_STYLE)
	page.settingsForm.SetButtonActivatedStyle(ACTIVATED_BUTTON_STYLE)

	// Input field for server address
	page.settingsForm.AddInputField("Server Address", "", 0, nil, nil)

	// Dropdown for theme selection
	page.settingsForm.AddDropDown("Theme", []string{"default", "america", "matrix", "halloween", "morning"}, 0, nil)

	page.settingsForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(WELCOME_PAGE, nil)
			return nil
		}

		return event
	})

	// Add the save and back buttons
	page.settingsForm.AddButton("Save & Apply", func() {
		// TODO: Save the settings in a json file

	})

	page.settingsForm.AddButton("Back", func() {
		nav.NavigateTo(WELCOME_PAGE, nil)
	})

	grid.AddItem(page.settingsForm, 1, 1, 1, 1, 0, 0, true)

	nav.Register(APP_SETTINGS_PAGE, grid, true, false, func(param interface{}) {
		page.settingsForm.SetFocus(0)
	}, func() {
		// Clear the server address field
		serverAddressInput, ok := page.settingsForm.GetFormItemByLabel("Server Address").(*tview.InputField)

		if !ok {
			panic("server address input form access failure")
		}

		serverAddressInput.SetText("")
	})
}
