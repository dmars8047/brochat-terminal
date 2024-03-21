package ui

import (
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/idamlib/idam"
	"github.com/dmars8047/strval"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const LOGIN_PAGE PageSlug = "login"

// LoginPage is the login page
type LoginPage struct {
	userAuthClient *idam.UserAuthClient
	brochatClient  *chat.BroChatClient
	feedClient     *state.FeedClient
	loginForm      *tview.Form
}

// NewLoginPage creates a new instance of the login page
func NewLoginPage(userAuthClient *idam.UserAuthClient, brochatClient *chat.BroChatClient, feedClient *state.FeedClient) *LoginPage {
	return &LoginPage{
		userAuthClient: userAuthClient,
		brochatClient:  brochatClient,
		feedClient:     feedClient,
		loginForm:      tview.NewForm(),
	}
}

func (page *LoginPage) Setup(app *tview.Application, appContext *state.ApplicationContext, nav *PageNavigator) {

	theme := appContext.GetTheme()

	const title = " BroChat - Login "

	grid := tview.NewGrid()
	grid.SetBackgroundColor(theme.BackgroundColor)
	grid.SetRows(4, 0, 1, 3, 4)
	grid.SetColumns(0, 70, 0)

	page.loginForm.SetBackgroundColor(theme.AccentColor)
	page.loginForm.SetFieldBackgroundColor(theme.AccentColorTwo)
	page.loginForm.SetFieldTextColor(theme.ForgroundColor)
	page.loginForm.SetLabelColor(theme.HighlightColor)
	page.loginForm.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignCenter)
	page.loginForm.SetButtonStyle(theme.ButtonStyle)
	page.loginForm.SetButtonActivatedStyle(theme.ActivatedButtonStyle)
	page.loginForm.AddInputField("Email", "", 0, nil, nil)
	page.loginForm.AddPasswordField("Password", "", 0, '*', nil)

	page.loginForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			nav.NavigateTo(WELCOME_PAGE, nil)
			return nil
		} else if event.Key() == tcell.KeyCtrlF {
			nav.NavigateTo(FORGOT_PW_PAGE, nil)
		}

		return event
	})

	page.loginForm.AddButton("Login", func() {
		emailInput, ok := page.loginForm.GetFormItemByLabel("Email").(*tview.InputField)

		formValidationErrors := make([]string, 0)

		if !ok {
			panic("email input form access failure")
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

		passwordInput, ok := page.loginForm.GetFormItemByLabel("Password").(*tview.InputField)

		if !ok {
			panic("password input form access failure")
		}

		password := passwordInput.GetText()

		valResult = strval.ValidateStringWithName(password,
			"Password",
			strval.MustNotBeEmpty(),
		)

		if !valResult.Valid {
			formValidationErrors = append(formValidationErrors, valResult.Messages...)
		}

		if len(formValidationErrors) > 0 {
			nav.AlertErrors("auth:login:alert:err", "Login Failed - Form Validation Error", formValidationErrors)
			return
		}

		request := &idam.UserLoginRequest{
			Email:    email,
			Password: password,
		}

		loginResponse, err := page.userAuthClient.Login("brochat", request)

		if err != nil {
			errMessage := err.Error()

			idamErr, ok := err.(*idam.ErrorResponse)

			if ok {
				switch idamErr.Code {
				case idam.RequestValidationFailure:
					errMessage = "Login Failed - Request Validation Error"
					nav.AlertErrors("auth:login:alert:err", errMessage, idamErr.Details)
					return
				case idam.UserNotFound:
					errMessage = "Login Failed - User Not Found"
				case idam.InvalidCredentials:
					errMessage = "Login Failed - Invalid Credentials"
				case idam.UserAccountLockout:
					errMessage = "User Account Lockout - Too Many Failed Login Requests"
				case idam.UnhandledError:
					errMessage = "Login Failed - An Unexpected Error Occurred"
				}
			}

			nav.Alert("auth:login:alert:err", errMessage)
			return
		}

		userAuth := state.UserAuth{
			AccessToken:     loginResponse.Token,
			TokenExpiration: time.Now().Add(time.Duration(loginResponse.ExpiresIn * int64(time.Second))),
		}

		appContext.SetUserSession(userAuth, func() {
			app.QueueUpdateDraw(
				func() {
					nav.NavigateTo(WELCOME_PAGE, WelcomePageParams{isRedirect: true, redirectMessage: "Your session has expired. Please login again."})
				},
			)
		})

		passwordInput.SetText("")
		emailInput.SetText("")

		getUserResult := page.brochatClient.GetUser(loginResponse.Token, loginResponse.UserId)

		err = getUserResult.Err()

		if err != nil {
			if len(getUserResult.ErrorDetails) > 0 {
				nav.Alert("auth:login:alert:err", getUserResult.ErrorDetails[0])
				return
			}

			if getUserResult.ResponseCode == chat.BROCHAT_RESPONSE_CODE_FORBIDDEN_ERROR {
				nav.Alert("auth:login:alert:err", FORBIDDEN_OPERATION_ERROR_MESSAGE)
				return
			}

			nav.Alert("auth:login:alert:err", err.Error())
			return
		}

		brochatUser := getUserResult.Content

		appContext.SetBrochatUser(brochatUser)

		err = page.feedClient.Connect()

		if err != nil {
			nav.Alert("auth:login:alert:err", err.Error())
			return
		}

		nav.NavigateTo(HOME_PAGE, nil)
	})

	page.loginForm.AddButton("Back", func() {
		nav.NavigateTo(WELCOME_PAGE, nil)
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(theme.BackgroundColor)
	tvInstructions.SetText("(CTRL + F) Forgot Password?")
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))

	grid.AddItem(page.loginForm, 1, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 3, 1, 1, 1, 0, 0, false)

	nav.Register(LOGIN_PAGE, grid, true, false, func(param interface{}) {
		page.onPageLoad(appContext)
	}, func() {
		page.onPageClose()
	})
}

func (page *LoginPage) onPageLoad(appContext *state.ApplicationContext) {
	appContext.CancelUserSession()
	page.loginForm.SetFocus(0)
}

func (page *LoginPage) onPageClose() {
	emailInput, ok := page.loginForm.GetFormItemByLabel("Email").(*tview.InputField)

	if !ok {
		panic("email input form clear failure")
	}

	emailInput.SetText("")

	pwInput, ok := page.loginForm.GetFormItemByLabel("Password").(*tview.InputField)

	if !ok {
		panic("password input form clear failure")
	}

	pwInput.SetText("")
}
