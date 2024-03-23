package ui

import (
	"log"

	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/broterm/internal/theme"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const APP_SETTINGS_PAGE PageSlug = "app_settings"

// AppSettingsPage is the location where users can configure application level settings.
type AppSettingsPage struct {
	settingsForm *tview.Form
	currentTheme string
}

// NewAppSettingsPage creates a new instance of the application settings page
func NewAppSettingsPage() *AppSettingsPage {
	return &AppSettingsPage{
		settingsForm: tview.NewForm(),
		currentTheme: "NOT_SET",
	}
}

// Setup configures the application settings page and registers it with the page navigator
// The page includes a form which allows the user to set the following settings:
// The host address of the server
// The theme (default, america, matrix, halloween, and morning)
// The log and setting config file storage location
func (page *AppSettingsPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	const title = " BroChat - Application Settings "

	page.currentTheme = appContext.GetTheme().Code

	grid := tview.NewGrid()
	grid.SetRows(4, 0, 1, 3, 4)
	grid.SetColumns(0, 70, 0)

	page.settingsForm.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignCenter)

	// Input field for server address
	page.settingsForm.AddInputField("Server Address", "", 0, nil, nil)

	// Dropdown for theme selection
	page.settingsForm.AddDropDown("Theme", []string{"default", "america", "matrix", "halloween", "christmas", "satanic"}, 0, nil)

	page.settingsForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(WELCOME_PAGE, nil)
			return nil
		}

		return event
	})

	applyTheme := func(previewTheme *theme.Theme) {
		var theme theme.Theme

		isPreview := previewTheme != nil

		if isPreview {
			theme = *previewTheme
		} else {
			theme = appContext.GetTheme()
		}

		nav.Pages.SetBackgroundColor(theme.BackgroundColor)
		theme.ApplyGlobals()

		grid.SetBackgroundColor(theme.BackgroundColor)
		page.settingsForm.SetBackgroundColor(theme.AccentColor)
		page.settingsForm.SetFieldBackgroundColor(theme.AccentColorTwo)
		page.settingsForm.SetFieldTextColor(theme.ForgroundColor)
		page.settingsForm.SetLabelColor(theme.HighlightColor)
		page.settingsForm.SetButtonStyle(theme.ButtonStyle)
		page.settingsForm.SetButtonActivatedStyle(theme.ActivatedButtonStyle)
		page.settingsForm.SetBorderColor(theme.BorderColor)
		page.settingsForm.SetTitleColor(theme.TitleColor)

		// Get the dropdown
		themeDropdown, ok := page.settingsForm.GetFormItemByLabel("Theme").(*tview.DropDown)

		if ok {
			themeDropdown.SetListStyles(theme.DropdownListUnselectedStyle, theme.DropdownListSelectedStyle)
		} else {
			log.Printf("Theme dropdown form access failure on theme change for settings page")
		}

		page.currentTheme = theme.Code
	}

	// Get the dropdown
	themeDropdown, ok := page.settingsForm.GetFormItemByLabel("Theme").(*tview.DropDown)

	if !ok {
		log.Printf("Theme dropdown form access failure on setup for settings page")
	} else {
		themeDropdown.SetSelectedFunc(func(text string, index int) {
			preview := theme.NewTheme(text)
			applyTheme(preview)
		})
	}

	// Add the save and back buttons
	page.settingsForm.AddButton("Save & Apply", func() {

	})

	page.settingsForm.AddButton("Back", func() {
		nav.NavigateTo(WELCOME_PAGE, nil)
	})

	grid.AddItem(page.settingsForm, 1, 1, 1, 1, 0, 0, true)

	nav.Register(APP_SETTINGS_PAGE, grid, true, false, func(param interface{}) {
		applyTheme(nil)
		page.settingsForm.SetFocus(0)

		// Set the dropdown to the current theme
		switch page.currentTheme {
		case "default":
			themeDropdown.SetCurrentOption(0)
		case "america":
			themeDropdown.SetCurrentOption(1)
		case "matrix":
			themeDropdown.SetCurrentOption(2)
		case "halloween":
			themeDropdown.SetCurrentOption(3)
		case "christmas":
			themeDropdown.SetCurrentOption(4)
		case "satanic":
			themeDropdown.SetCurrentOption(5)
		default:
			themeDropdown.SetCurrentOption(0)
		}

	}, func() {
		applyTheme(nil)

		// Clear the server address field
		serverAddressInput, ok := page.settingsForm.GetFormItemByLabel("Server Address").(*tview.InputField)

		if !ok {
			panic("server address input form access failure")
		}

		serverAddressInput.SetText("")
	})
}
