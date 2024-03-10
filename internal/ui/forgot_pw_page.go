package ui

import (
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/idamlib/idam"
	"github.com/dmars8047/strval"
	"github.com/gdamore/tcell/v2"
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
	userAuthClient *idam.UserAuthClient
	forgotPWForm   *tview.Form
}

// NewForgotPasswordPage creates a new instance of the forgot password page
func NewForgotPasswordPage(userAuthClient *idam.UserAuthClient) *ForgotPasswordPage {
	return &ForgotPasswordPage{
		userAuthClient: userAuthClient,
		forgotPWForm:   tview.NewForm(),
	}
}

// Setup sets up the forgot password page
func (page *ForgotPasswordPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	grid.SetRows(4, 0, 1, 3, 4)
	grid.SetColumns(0, 70, 0)

	page.forgotPWForm.SetBackgroundColor(ACCENT_BACKGROUND_COLOR)
	page.forgotPWForm.SetFieldBackgroundColor(ACCENT_COLOR_TWO_COLOR)
	page.forgotPWForm.SetLabelColor(BROCHAT_YELLOW_COLOR)
	page.forgotPWForm.SetBorder(true).SetTitle(FORGOT_PW_TITLE).SetTitleAlign(tview.AlignCenter)
	page.forgotPWForm.SetButtonStyle(DEFAULT_BUTTON_STYLE)
	page.forgotPWForm.SetButtonActivatedStyle(ACTIVATED_BUTTON_STYLE)
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
	tvInstructions.SetBackgroundColor(DEFAULT_BACKGROUND_COLOR)
	tvInstructions.SetText("Enter your email to recieve a password reset link.")
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))

	grid.AddItem(page.forgotPWForm, 1, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 3, 1, 1, 1, 0, 0, false)

	nav.Register(FORGOT_PW_PAGE, grid, true, false,
		func(param interface{}) {
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
