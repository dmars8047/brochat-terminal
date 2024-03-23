package ui

import (
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/idamlib/idam"
	"github.com/dmars8047/strval"
	"github.com/rivo/tview"
)

const FORGOT_PW_PAGE PageSlug = "forgot_password"

const (
	// FORGOT_PW_TITLE is the title of the forgot password page
	FORGOT_PW_TITLE           = " BroChat - Forgot Password "
	FORGOT_PW_MODAL_INFO      = "auth:forgotpw:alert:info"
	FORGOT_PW_MODAL_ERR       = "auth:forgotpw:alert:err"
	FORGOT_PW_SUCCESS_MESSAGE = "Password Reset Link Sent. Check your email to proceed."
)

// ForgotPasswordPage is the forgot password page
type ForgotPasswordPage struct {
	userAuthClient   *idam.UserAuthClient
	forgotPWForm     *tview.Form
	currentThemeCode string
}

// NewForgotPasswordPage creates a new instance of the forgot password page
func NewForgotPasswordPage(userAuthClient *idam.UserAuthClient) *ForgotPasswordPage {
	return &ForgotPasswordPage{
		userAuthClient:   userAuthClient,
		forgotPWForm:     tview.NewForm(),
		currentThemeCode: "NOT_SET",
	}
}

// Setup sets up the forgot password page
func (page *ForgotPasswordPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	grid := tview.NewGrid()
	grid.SetRows(4, 0, 1, 3, 4)
	grid.SetColumns(0, 70, 0)

	page.forgotPWForm.SetBorder(true).SetTitle(FORGOT_PW_TITLE).SetTitleAlign(tview.AlignCenter)
	page.forgotPWForm.AddInputField("Email", "", 0, nil, nil)

	page.forgotPWForm.AddButton("Submit", func() {
		emailInput, ok := page.forgotPWForm.GetFormItemByLabel("Email").(*tview.InputField)

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
			nav.AlertErrors(FORGOT_PW_MODAL_ERR, "Form Validation Error", valResult.Messages)
			return
		}

		request := &idam.UserPasswordResetInitiationRequest{
			Email: email,
		}

		err := page.userAuthClient.InitiatePasswordReset("brochat", request)

		if err != nil {
			errMessage := err.Error()

			idamErr, ok := err.(*idam.ErrorResponse)

			if ok {
				switch idamErr.Code {
				case idam.RequestValidationFailure:
					errMessage = "Request Validation Error"
					nav.AlertErrors(FORGOT_PW_MODAL_ERR, errMessage, idamErr.Details)
					return
				case idam.UserNotFound:
					errMessage = "User Not Found"
				case idam.UnhandledError:
					errMessage = "An Unexpected Error Occurred"
				}
			}

			nav.Alert(FORGOT_PW_MODAL_ERR, errMessage)
			return
		}

		nav.AlertWithDoneFunc(FORGOT_PW_MODAL_INFO, FORGOT_PW_SUCCESS_MESSAGE, func(buttonIndex int, buttonLabel string) {
			nav.Pages.HidePage(FORGOT_PW_MODAL_INFO).RemovePage(FORGOT_PW_MODAL_INFO)
			nav.NavigateTo(LOGIN_PAGE, nil)
		})
	})

	page.forgotPWForm.AddButton("Back", func() {
		nav.NavigateTo(LOGIN_PAGE, nil)
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetText("Enter your email to recieve a password reset link.")

	grid.AddItem(page.forgotPWForm, 1, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 3, 1, 1, 1, 0, 0, false)

	applyTheme := func() {
		theme := appContext.GetTheme()

		if page.currentThemeCode != theme.Code {
			grid.SetBackgroundColor(theme.BackgroundColor)
			page.forgotPWForm.SetBackgroundColor(theme.AccentColor)
			page.forgotPWForm.SetFieldBackgroundColor(theme.AccentColorTwo)
			page.forgotPWForm.SetLabelColor(theme.HighlightColor)
			page.forgotPWForm.SetButtonStyle(theme.ButtonStyle)
			page.forgotPWForm.SetButtonActivatedStyle(theme.ActivatedButtonStyle)
			page.forgotPWForm.SetBorderColor(theme.BorderColor)
			page.forgotPWForm.SetTitleColor(theme.TitleColor)
			tvInstructions.SetBackgroundColor(theme.BackgroundColor)
			tvInstructions.SetTextColor(theme.InfoColor)
		}
	}

	applyTheme()

	nav.Register(FORGOT_PW_PAGE, grid, true, false,
		func(param interface{}) {
			applyTheme()
			page.onPageLoad()
		},
		func() {
			page.onPageClose()
		})
}

// onPageLoad is called when the forgot password page is navigated to
func (page *ForgotPasswordPage) onPageLoad() {
	page.forgotPWForm.SetFocus(0)
}

// onPageClose is called when the forgot password page is navigated away from
func (page *ForgotPasswordPage) onPageClose() {
	emailInput, ok := page.forgotPWForm.GetFormItemByLabel("Email").(*tview.InputField)

	if !ok {
		panic("email input form clear failure")
	}

	emailInput.SetText("")
}
