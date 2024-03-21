package ui

import (
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/idamlib/idam"
	"github.com/dmars8047/strval"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const REGISTER_PAGE PageSlug = "register"

const (
	REGISTRATION_MODAL_INFO      = "auth:register:alert:info"
	REGISTRATION_MODAL_ERR       = "auth:register:alert:err"
	REGISTRATION_SUCCESS_MESSAGE = "A verification email has been sent to the email address provided. Please check your email to verify your account."
)

type RegistrationPage struct {
	userAuthClient   *idam.UserAuthClient
	registrationForm *tview.Form
}

// NewRegistrationPage creates a new instance of the registration page
func NewRegistrationPage(userAuthClient *idam.UserAuthClient) *RegistrationPage {
	return &RegistrationPage{
		userAuthClient:   userAuthClient,
		registrationForm: tview.NewForm(),
	}
}

func (page *RegistrationPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {

	theme := appContext.GetTheme()

	grid := tview.NewGrid()
	grid.SetBackgroundColor(theme.BackgroundColor)

	grid.SetRows(4, 0, 6)
	grid.SetColumns(0, 70, 0)

	page.registrationForm.SetBorder(true).SetTitle(" BroChat - Register ").SetTitleAlign(tview.AlignCenter)
	page.registrationForm.SetBackgroundColor(theme.AccentColor)
	page.registrationForm.SetFieldBackgroundColor(theme.AccentColorTwo)
	page.registrationForm.SetFieldTextColor(theme.ForgroundColor)
	page.registrationForm.SetLabelColor(theme.HighlightColor)
	page.registrationForm.SetButtonStyle(theme.ButtonStyle)
	page.registrationForm.SetButtonActivatedStyle(theme.ActivatedButtonStyle)

	// If the user presses the escape key, navigate back to the welcome page
	page.registrationForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(WELCOME_PAGE, nil)
			return nil
		}

		return event
	})

	page.registrationForm.
		AddInputField("Email", "", 0, nil, nil).
		AddInputField("Username", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddPasswordField("Confirm Password", "", 0, '*', nil).
		AddButton("Register", func() {
			formValidationErrors := make([]string, 0)

			emailInput, ok := page.registrationForm.GetFormItemByLabel("Email").(*tview.InputField)

			if !ok {
				panic("email input form clear failure")
			}

			email := emailInput.GetText()

			valResult := strval.ValidateStringWithName(email,
				"Email",
				strval.MustNotBeEmpty(),
				strval.MustBeValidEmailFormat(),
			)

			if !valResult.Valid {
				formValidationErrors = append(formValidationErrors, valResult.Messages...)
			}

			passwordInput, ok := page.registrationForm.GetFormItemByLabel("Password").(*tview.InputField)

			if !ok {
				panic("password input form clear failure")
			}

			password := passwordInput.GetText()

			valResult = strval.ValidateStringWithName(password,
				"Password",
				strval.MustNotBeEmpty(),
				strval.MustHaveMinLengthOf(idam.MinPasswordLength),
				strval.MustHaveMaxLengthOf(idam.MaxPasswordLength),
				strval.MustContainAtLeastOne([]rune(idam.AllowablePasswordSpecialCharacters)),
				strval.MustNotContainAnyOf([]rune(idam.DisallowedPassowrdSpecialCharacters)),
				strval.MustContainNumbers(),
				strval.MustContainUppercaseLetter(),
				strval.MustContainLowercaseLetter(),
				strval.MustOnlyContainPrintableCharacters(),
				strval.MustOnlyContainASCIICharacters(),
			)

			if !valResult.Valid {
				formValidationErrors = append(formValidationErrors, valResult.Messages...)
			}

			confirmPasswordInput, ok := page.registrationForm.GetFormItemByLabel("Confirm Password").(*tview.InputField)

			if !ok {
				panic("confirm password input form clear failure")
			}

			confirmPassword := confirmPasswordInput.GetText()

			if password != confirmPassword {
				formValidationErrors = append(formValidationErrors, "Passwords do not match")
			}

			usernameInput, ok := page.registrationForm.GetFormItemByLabel("Username").(*tview.InputField)

			if !ok {
				panic("username input form clear failure")
			}

			username := usernameInput.GetText()

			valResult = strval.ValidateStringWithName(username,
				"Username",
				strval.MustNotBeEmpty(),
				strval.MustBeAlphaNumeric(),
				strval.MustHaveMinLengthOf(3),
				strval.MustHaveMaxLengthOf(20),
			)

			if !valResult.Valid {
				formValidationErrors = append(formValidationErrors, valResult.Messages...)
			}

			if len(formValidationErrors) > 0 {
				nav.AlertErrors(REGISTRATION_MODAL_ERR, "Login Failed - Form Validation Error", formValidationErrors)
				return
			}

			request := &idam.UserRegistrationRequest{
				Email:    email,
				Password: password,
				Username: username,
			}

			_, err := page.userAuthClient.Register("brochat", request)

			if err != nil {
				errMessage := err.Error()

				idamErr, ok := err.(*idam.ErrorResponse)

				if ok {
					switch idamErr.Code {
					case idam.RequestValidationFailure:
						errMessage = "Registration Failed - Request Validation Error"
						nav.AlertErrors(REGISTRATION_MODAL_ERR, errMessage, idamErr.Details)
						return
					case idam.InvalidCredentials:
						errMessage = "Registration Failed - Invalid Credentials"
					case idam.UnhandledError:
						errMessage = "Registration Failed - An Unexpected Error Occurred"
					}
				}

				nav.Alert(REGISTRATION_MODAL_ERR, errMessage)
				return
			}

			nav.AlertWithDoneFunc(REGISTRATION_MODAL_INFO, REGISTRATION_SUCCESS_MESSAGE, func(buttonIndex int, buttonLabel string) {
				nav.Pages.HidePage(REGISTRATION_MODAL_INFO).RemovePage(REGISTRATION_MODAL_INFO)
				nav.NavigateTo(WELCOME_PAGE, nil)
			})
		}).
		AddButton("Back", func() { nav.NavigateTo(WELCOME_PAGE, nil) })

	grid.AddItem(page.registrationForm, 1, 1, 1, 1, 0, 0, true)

	nav.Register(REGISTER_PAGE, grid, true, false,
		func(param interface{}) {
			page.onPageLoad()
		},
		func() {
			page.onPageClose()
		})
}

func (page *RegistrationPage) onPageLoad() {
	page.registrationForm.SetFocus(0)
}

func (page *RegistrationPage) onPageClose() {
	emailInput, ok := page.registrationForm.GetFormItemByLabel("Email").(*tview.InputField)

	if !ok {
		panic("email input form clear failure")
	}

	emailInput.SetText("")

	passwordInput, ok := page.registrationForm.GetFormItemByLabel("Password").(*tview.InputField)

	if !ok {
		panic("password input form clear failure")
	}

	passwordInput.SetText("")

	confirmPasswordInput, ok := page.registrationForm.GetFormItemByLabel("Confirm Password").(*tview.InputField)

	if !ok {
		panic("confirm password input form clear failure")
	}

	confirmPasswordInput.SetText("")

	usernameInput, ok := page.registrationForm.GetFormItemByLabel("Username").(*tview.InputField)

	if !ok {
		panic("username input form clear failure")
	}

	usernameInput.SetText("")
}
