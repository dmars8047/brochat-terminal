package ui

import (
	"context"
	"time"

	"github.com/dmars8047/broterm/internal/auth"
	"github.com/dmars8047/broterm/internal/bro"
	"github.com/dmars8047/broterm/internal/feed"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/strval"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AuthModule struct {
	userAuthClient *auth.UserAuthClient
	appContext     *state.ApplicationContext
	brochatClient  *bro.BroChatUserClient
	pageNav        *PageNavigator
	app            *tview.Application
	feedClient     *feed.Client
}

func NewAuthModule(userAuthClient *auth.UserAuthClient,
	brochatClient *bro.BroChatUserClient,
	appContext *state.ApplicationContext,
	pageNavigator *PageNavigator,
	application *tview.Application,
	feedClient *feed.Client) *AuthModule {
	mod := AuthModule{
		userAuthClient: userAuthClient,
		appContext:     appContext,
		brochatClient:  brochatClient,
		pageNav:        pageNavigator,
		app:            application,
		feedClient:     feedClient,
	}

	return &mod
}

func (mod *AuthModule) SetupAuthPages() {
	mod.setupWelcomePage()
	mod.setupLoginPage()
	mod.setupRegistrationPage()
	mod.setupForgotPasswordPage()
}

func (mod *AuthModule) setupWelcomePage() {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(4, 8, 8, 1, 1, 0).
		SetColumns(0, 31, 39, 0)

	logoBro := tview.NewTextView()
	logoBro.SetTextAlign(tview.AlignLeft).
		SetBackgroundColor(DefaultBackgroundColor)
	logoBro.SetTextColor(tcell.ColorWhite)
	logoBro.SetText(
		`BBBBBBB\                      
BB  __BB\                     
BB |  BB | RRRRRR\   OOOOOO\  
BBBBBBB\ |RR  __RR\ OO  __OO\ 
BB  __BB\ RR |  \__|OO /  OO |
BB |  BB |RR |      OO |  OO |
BBBBBBB  |RR |      \OOOOOO  |
\_______/ \__|       \______/ `)

	logoChat := tview.NewTextView()
	logoChat.SetTextAlign(tview.AlignLeft)
	logoChat.SetBackgroundColor(DefaultBackgroundColor)
	logoChat.SetTextColor(BroChatYellowColor)
	logoChat.SetText(
		` CCCCCC\  HH\                  TT\
CC  __CC\ HH |                 TT |
CC /  \__|HHHHHHH\   AAAAAA\ TTTTTT\
CC |      HH  __HH\  \____AA\\_TT  _|
CC |      HH |  HH | AAAAAAA | TT |
CC |  CC\ HH |  HH |AA  __AA | TT |TT\
\CCCCCC  |HH |  HH |\AAAAAAA | \TTTT  |
 \______/ \__|  \__| \_______|  \____/`)

	loginButton := tview.NewButton("Login").SetSelectedFunc(func() {
		mod.pageNav.NavigateTo(LOGIN_PAGE, nil)
	}).SetActivatedStyle(ActivatedButtonStyle).SetStyle(ButtonStyle)

	registrationButton := tview.NewButton("Register").SetSelectedFunc(func() {
		mod.pageNav.NavigateTo(REGISTER_PAGE, nil)
	}).SetActivatedStyle(ActivatedButtonStyle).SetStyle(ButtonStyle)

	exitButton := tview.NewButton("Exit").SetSelectedFunc(func() {
		mod.app.Stop()
	}).SetActivatedStyle(ActivatedButtonStyle).SetStyle(ButtonStyle)

	buttonGrid := tview.NewGrid()
	buttonGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if loginButton.HasFocus() {
				mod.app.SetFocus(registrationButton)
			} else if registrationButton.HasFocus() {
				mod.app.SetFocus(exitButton)
			} else if exitButton.HasFocus() {
				mod.app.SetFocus(loginButton)
			}
		} else if event.Key() == tcell.KeyBacktab {
			if loginButton.HasFocus() {
				mod.app.SetFocus(exitButton)
			} else if registrationButton.HasFocus() {
				mod.app.SetFocus(loginButton)
			} else if exitButton.HasFocus() {
				mod.app.SetFocus(registrationButton)
			}
		}
		return event
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DefaultBackgroundColor)
	tvInstructions.SetText("Navigate with Tab and Shift+Tab")
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))

	tvVersionNumber := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvVersionNumber.SetBackgroundColor(DefaultBackgroundColor)
	tvVersionNumber.SetText("Version - v0.0.12")
	tvVersionNumber.SetTextColor(tcell.NewHexColor(0x777777))

	buttonGrid.SetRows(3, 1, 1).SetColumns(0, 4, 0, 4, 0)

	buttonGrid.AddItem(loginButton, 0, 0, 1, 1, 0, 0, true).
		AddItem(registrationButton, 0, 2, 1, 1, 0, 0, false).
		AddItem(exitButton, 0, 4, 1, 1, 0, 0, false).
		AddItem(tvInstructions, 2, 0, 1, 5, 0, 0, false)

	grid.AddItem(logoBro, 1, 1, 1, 1, 0, 0, false).
		AddItem(logoChat, 1, 2, 1, 1, 0, 0, false).
		AddItem(buttonGrid, 2, 1, 1, 2, 0, 0, true).
		AddItem(tvVersionNumber, 4, 1, 1, 2, 0, 0, false)

	mod.pageNav.Register(WELCOME_PAGE, grid, true, true, nil, nil)
}

func (mod *AuthModule) setupLoginPage() {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)
	grid.SetRows(4, 0, 1, 3, 4)
	grid.SetColumns(0, 70, 0)

	loginForm := tview.NewForm()
	loginForm.SetBackgroundColor(AccentBackgroundColor)
	loginForm.SetFieldBackgroundColor(AccentColorTwoColor)
	loginForm.SetLabelColor(BroChatYellowColor)
	loginForm.SetBorder(true).SetTitle(" BroChat - Login ").SetTitleAlign(tview.AlignCenter)
	loginForm.SetButtonStyle(ButtonStyle)
	loginForm.SetButtonActivatedStyle(ActivatedButtonStyle)
	loginForm.AddInputField("Email", "", 0, nil, nil)
	loginForm.AddPasswordField("Password", "", 0, '*', nil)

	loginForm.AddButton("Login", func() {
		emailInput, ok := loginForm.GetFormItemByLabel("Email").(*tview.InputField)

		formValidationErrors := make([]string, 0)

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

		passwordInput, ok := loginForm.GetFormItemByLabel("Password").(*tview.InputField)

		if !ok {
			panic("password input form clear failure")
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
			alertErrors(mod.pageNav.Pages, "auth:login:alert:err", "Login Failed - Form Validation Error", formValidationErrors)
			return
		}

		request := &auth.UserLoginRequest{
			Email:    email,
			Password: password,
		}

		loginResponse, err := mod.userAuthClient.Login("brochat", request)

		if err != nil {
			errMessage := err.Error()

			idamErr, ok := err.(*auth.ErrorResponse)

			if ok {
				switch idamErr.Code {
				case auth.RequestValidationFailure:
					errMessage = "Login Failed - Request Validation Error"
					alertErrors(mod.pageNav.Pages, "auth:login:alert:err", errMessage, idamErr.Details)
					return
				case auth.UserNotFound:
					errMessage = "Login Failed - User Not Found"
				case auth.InvalidCredentials:
					errMessage = "Login Failed - Invalid Credentials"
				case auth.UnhandledError:
					errMessage = "Login Failed - An Unexpected Error Occurred"
				}
			}

			Alert(mod.pageNav.Pages, "auth:login:alert:err", errMessage)
			return
		}

		sessionContext, sessionCancelFunc := context.WithCancel(mod.appContext.Context)

		session := &state.UserSession{
			Auth: state.UserAuth{
				AccessToken:     loginResponse.Token,
				TokenExpiration: time.Now().Add(time.Duration(loginResponse.ExpiresIn * int64(time.Second))),
			},
			Info: state.UserInfo{
				Id:       loginResponse.UserId,
				Username: loginResponse.Username,
			},
			Context:    sessionContext,
			CancelFunc: sessionCancelFunc,
		}

		mod.appContext.UserSession = session

		passwordInput.SetText("")
		emailInput.SetText("")

		brochatUser, err := mod.brochatClient.GetUser(&bro.AuthInfo{
			AccessToken: session.Auth.AccessToken,
			TokenType:   DEFAULT_AUTH_TOKEN_TYPE,
		}, session.Info.Id)

		if err != nil {
			Alert(mod.pageNav.Pages, "auth:login:alert:err", err.Error())
			return
		}

		mod.appContext.BrochatUser = brochatUser

		err = mod.feedClient.Connect(session.Auth, session.Context)

		if err != nil {
			Alert(mod.pageNav.Pages, "auth:login:alert:err", err.Error())
			return
		}

		go func(ctx *state.ApplicationContext, pageNav *PageNavigator) {
			select {
			case <-ctx.UserSession.Context.Done():
				return
			case <-time.After(time.Until(ctx.UserSession.Auth.TokenExpiration)):
				ctx.UserSession.CancelFunc()
				pageNav.NavigateTo(LOGIN_PAGE, nil)
				return
			}
		}(mod.appContext, mod.pageNav)

		mod.pageNav.NavigateTo(HOME_MENU_PAGE, nil)
	})

	loginForm.AddButton("Back", func() {
		mod.pageNav.NavigateTo(WELCOME_PAGE, nil)
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DefaultBackgroundColor)
	tvInstructions.SetText("(CTRL + F) Forgot Password?")
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))

	loginForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlF {
			mod.pageNav.NavigateTo(FORGOT_PW_PAGE, nil)
		}

		return event
	})

	grid.AddItem(loginForm, 1, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 3, 1, 1, 1, 0, 0, false)

	mod.pageNav.Register(LOGIN_PAGE, grid, true, false, func(param interface{}) {
		loginForm.SetFocus(0)
	}, func() {
		emailInput, ok := loginForm.GetFormItemByLabel("Email").(*tview.InputField)

		if !ok {
			panic("email input form clear failure")
		}

		emailInput.SetText("")

		pwInput, ok := loginForm.GetFormItemByLabel("Password").(*tview.InputField)

		if !ok {
			panic("password input form clear failure")
		}

		pwInput.SetText("")
	})
}

const (
	FORGOT_PW_TITLE           = " BroChat - Forgot Password "
	FORGOT_PW_MODAL_INFO      = "auth:forgotpw:alert:info"
	FORGOT_PW_MODAL_ERR       = "auth:forgotpw:alert:err"
	FORGOT_PW_SUCCESS_MESSAGE = "Password Reset Link Sent. Check your email to proceed."
)

func (mod *AuthModule) setupForgotPasswordPage() {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)
	grid.SetRows(4, 0, 1, 3, 4)
	grid.SetColumns(0, 70, 0)

	forgotPWForm := tview.NewForm()
	forgotPWForm.SetBackgroundColor(AccentBackgroundColor)
	forgotPWForm.SetFieldBackgroundColor(AccentColorTwoColor)
	forgotPWForm.SetLabelColor(BroChatYellowColor)
	forgotPWForm.SetBorder(true).SetTitle(FORGOT_PW_TITLE).SetTitleAlign(tview.AlignCenter)
	forgotPWForm.SetButtonStyle(ButtonStyle)
	forgotPWForm.SetButtonActivatedStyle(ActivatedButtonStyle)
	forgotPWForm.AddInputField("Email", "", 0, nil, nil)

	forgotPWForm.AddButton("Submit", func() {
		emailInput, ok := forgotPWForm.GetFormItemByLabel("Email").(*tview.InputField)

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
			alertErrors(mod.pageNav.Pages, FORGOT_PW_MODAL_ERR, "Form Validation Error", valResult.Messages)
			return
		}

		request := &auth.UserPasswordResetInitiationRequest{
			Email: email,
		}

		err := mod.userAuthClient.InitiatePasswordReset("brochat", request)

		if err != nil {
			errMessage := err.Error()

			idamErr, ok := err.(*auth.ErrorResponse)

			if ok {
				switch idamErr.Code {
				case auth.RequestValidationFailure:
					errMessage = "Request Validation Error"
					alertErrors(mod.pageNav.Pages, FORGOT_PW_MODAL_ERR, errMessage, idamErr.Details)
					return
				case auth.UserNotFound:
					errMessage = "User Not Found"
				case auth.UnhandledError:
					errMessage = "An Unexpected Error Occurred"
				}
			}

			Alert(mod.pageNav.Pages, FORGOT_PW_MODAL_ERR, errMessage)
			return
		}

		AlertWithDoneFunc(mod.pageNav.Pages, FORGOT_PW_MODAL_INFO, FORGOT_PW_SUCCESS_MESSAGE, func(buttonIndex int, buttonLabel string) {
			mod.pageNav.Pages.HidePage(FORGOT_PW_MODAL_INFO).RemovePage(FORGOT_PW_MODAL_INFO)
			mod.pageNav.NavigateTo(LOGIN_PAGE, nil)
		})
	})

	forgotPWForm.AddButton("Back", func() {
		mod.pageNav.NavigateTo(LOGIN_PAGE, nil)
	})

	tvInstructions := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	tvInstructions.SetBackgroundColor(DefaultBackgroundColor)
	tvInstructions.SetText("Enter your email to recieve a password reset link.")
	tvInstructions.SetTextColor(tcell.NewHexColor(0xFFFFFF))

	grid.AddItem(forgotPWForm, 1, 1, 1, 1, 0, 0, true)
	grid.AddItem(tvInstructions, 3, 1, 1, 1, 0, 0, false)

	mod.pageNav.Register(FORGOT_PW_PAGE, grid, true, false,
		func(param interface{}) {
			forgotPWForm.SetFocus(0)
		},
		func() {
			emailInput, ok := forgotPWForm.GetFormItemByLabel("Email").(*tview.InputField)

			if !ok {
				panic("email input form clear failure")
			}

			emailInput.SetText("")
		})
}

const (
	REGISTRATION_MODAL_INFO      = "auth:register:alert:info"
	REGISTRATION_MODAL_ERR       = "auth:register:alert:err"
	REGISTRATION_SUCCESS_MESSAGE = "A verification email has been sent to the email address provided. Please check your email to verify your account."
)

func (mod *AuthModule) setupRegistrationPage() {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(4, 0, 6)
	grid.SetColumns(0, 70, 0)

	registrationForm := tview.NewForm()
	registrationForm.SetBorder(true).SetTitle(" BroChat - Register ").SetTitleAlign(tview.AlignCenter)
	registrationForm.SetBackgroundColor(AccentBackgroundColor)
	registrationForm.SetFieldBackgroundColor(AccentColorTwoColor)
	registrationForm.SetButtonStyle(ButtonStyle)
	registrationForm.SetButtonActivatedStyle(ActivatedButtonStyle)

	registrationForm.
		AddInputField("Email", "", 0, nil, nil).
		AddInputField("Username", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddPasswordField("Confirm Password", "", 0, '*', nil).
		AddButton("Register", func() {
			formValidationErrors := make([]string, 0)

			emailInput, ok := registrationForm.GetFormItemByLabel("Email").(*tview.InputField)

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

			passwordInput, ok := registrationForm.GetFormItemByLabel("Password").(*tview.InputField)

			if !ok {
				panic("password input form clear failure")
			}

			password := passwordInput.GetText()

			valResult = strval.ValidateStringWithName(password,
				"Password",
				strval.MustNotBeEmpty(),
				strval.MustHaveMinLengthOf(auth.MinPasswordLength),
				strval.MustHaveMaxLengthOf(auth.MaxPasswordLength),
				strval.MustContainAtLeastOne([]rune(auth.AllowablePasswordSpecialCharacters)),
				strval.MustNotContainAnyOf([]rune(auth.DisallowedPassowrdSpecialCharacters)),
				strval.MustContainNumbers(),
				strval.MustContainUppercaseLetter(),
				strval.MustContainLowercaseLetter(),
				strval.MustOnlyContainPrintableCharacters(),
				strval.MustOnlyContainASCIICharacters(),
			)

			if !valResult.Valid {
				formValidationErrors = append(formValidationErrors, valResult.Messages...)
			}

			confirmPasswordInput, ok := registrationForm.GetFormItemByLabel("Confirm Password").(*tview.InputField)

			if !ok {
				panic("confirm password input form clear failure")
			}

			confirmPassword := confirmPasswordInput.GetText()

			if password != confirmPassword {
				formValidationErrors = append(formValidationErrors, "Passwords do not match")
			}

			usernameInput, ok := registrationForm.GetFormItemByLabel("Username").(*tview.InputField)

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
				alertErrors(mod.pageNav.Pages, REGISTRATION_MODAL_ERR, "Login Failed - Form Validation Error", formValidationErrors)
				return
			}

			request := &auth.UserRegistrationRequest{
				Email:    email,
				Password: password,
				Username: username,
			}

			_, err := mod.userAuthClient.Register("brochat", request)

			if err != nil {
				errMessage := err.Error()

				idamErr, ok := err.(*auth.ErrorResponse)

				if ok {
					switch idamErr.Code {
					case auth.RequestValidationFailure:
						errMessage = "Registration Failed - Request Validation Error"
						alertErrors(mod.pageNav.Pages, REGISTRATION_MODAL_ERR, errMessage, idamErr.Details)
						return
					case auth.InvalidCredentials:
						errMessage = "Registration Failed - Invalid Credentials"
					case auth.UnhandledError:
						errMessage = "Registration Failed - An Unexpected Error Occurred"
					}
				}

				Alert(mod.pageNav.Pages, REGISTRATION_MODAL_ERR, errMessage)
				return
			}

			AlertWithDoneFunc(mod.pageNav.Pages, REGISTRATION_MODAL_INFO, REGISTRATION_SUCCESS_MESSAGE, func(buttonIndex int, buttonLabel string) {
				mod.pageNav.Pages.HidePage(REGISTRATION_MODAL_INFO).RemovePage(REGISTRATION_MODAL_INFO)
				mod.pageNav.NavigateTo(WELCOME_PAGE, nil)
			})
		}).
		AddButton("Back", func() { mod.pageNav.NavigateTo(WELCOME_PAGE, nil) })

	grid.AddItem(registrationForm, 1, 1, 1, 1, 0, 0, true)

	mod.pageNav.Register(REGISTER_PAGE, grid, true, false,
		func(param interface{}) {
			registrationForm.SetFocus(0)
		},
		func() {
			emailInput, ok := registrationForm.GetFormItemByLabel("Email").(*tview.InputField)

			if !ok {
				panic("email input form clear failure")
			}

			emailInput.SetText("")

			passwordInput, ok := registrationForm.GetFormItemByLabel("Password").(*tview.InputField)

			if !ok {
				panic("password input form clear failure")
			}

			passwordInput.SetText("")

			confirmPasswordInput, ok := registrationForm.GetFormItemByLabel("Confirm Password").(*tview.InputField)

			if !ok {
				panic("confirm password input form clear failure")
			}

			confirmPasswordInput.SetText("")

			usernameInput, ok := registrationForm.GetFormItemByLabel("Username").(*tview.InputField)

			if !ok {
				panic("username input form clear failure")
			}

			usernameInput.SetText("")
		})
}
