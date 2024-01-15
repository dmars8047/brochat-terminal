package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/dmars8047/brochat-service/pkg/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/idam-service/pkg/idam"
	"github.com/dmars8047/strval"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AuthModule struct {
	userAuthClient *idam.UserAuthClient
	appState       state.ApplicationState
	brochatClient  *chat.BroChatUserClient
	pageNav        *PageNavigator
	app            *tview.Application
}

func NewAuthModule(userAuthClient *idam.UserAuthClient,
	brochatClient *chat.BroChatUserClient,
	appState state.ApplicationState,
	pageNavigator *PageNavigator,
	application *tview.Application) *AuthModule {
	mod := AuthModule{
		userAuthClient: userAuthClient,
		appState:       appState,
		brochatClient:  brochatClient,
		pageNav:        pageNavigator,
		app:            application,
	}

	return &mod
}

func (mod *AuthModule) SetupAuthPages() {
	mod.setupWelcomePage()
	mod.setupLoginPage()
	mod.setupRegistrationPage()
}

func (mod *AuthModule) setupWelcomePage() {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(4, 8, 8, 0).
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
		mod.pageNav.NavigateTo(LOGIN_PAGE)
	}).SetActivatedStyle(ActivatedButtonStyle).SetStyle(ButtonStyle)

	registrationButton := tview.NewButton("Register").SetSelectedFunc(func() {
		mod.pageNav.NavigateTo(REGISTER_PAGE)
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

	buttonGrid.SetRows(3, 1, 1).SetColumns(0, 4, 0, 4, 0)

	buttonGrid.AddItem(loginButton, 0, 0, 1, 1, 0, 0, true).
		AddItem(registrationButton, 0, 2, 1, 1, 0, 0, false).
		AddItem(exitButton, 0, 4, 1, 1, 0, 0, false).
		AddItem(tvInstructions, 2, 0, 1, 5, 0, 0, false)

	grid.AddItem(logoBro, 1, 1, 1, 1, 0, 0, false).
		AddItem(logoChat, 1, 2, 1, 1, 0, 0, false).
		AddItem(buttonGrid, 2, 1, 1, 2, 0, 0, true)

	mod.pageNav.Register(WELCOME_PAGE, grid, true, true, func() {})
}

func (mod *AuthModule) setupLoginPage() {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(4, 0, 10)
	grid.SetColumns(0, 70, 0)

	loginForm := tview.NewForm()
	loginForm.SetBackgroundColor(AccentBackgroundColor)
	loginForm.SetFieldBackgroundColor(AccentColorTwoColorCode)
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

		request := &idam.UserLoginRequest{
			Email:    email,
			Password: password,
		}

		loginResponse, err := mod.userAuthClient.Login("brochat", request)

		if err != nil {
			errMessage := err.Error()

			idamErr, ok := err.(*idam.ErrorResponse)

			if ok {
				switch idamErr.Code {
				case idam.RequestValidationFailure:
					errMessage = "Login Failed - Request Validation Error"
					alertErrors(mod.pageNav.Pages, "auth:login:alert:err", errMessage, idamErr.Details)
					return
				case idam.InvalidCredentials:
					errMessage = "Login Failed - Invalid Credentials"
				case idam.UnhandledError:
					errMessage = "Login Failed - An Unexpected Error Occurred"
				}
			}

			Alert(mod.pageNav.Pages, "auth:login:alert:err", errMessage)
			return
		}

		session := &state.UserSession{
			Auth: state.UserAuth{
				AccessToken:     loginResponse.Token,
				TokenExpiration: time.Now().Add(time.Duration(loginResponse.ExpiresIn * int64(time.Second))),
			},
			Info: state.UserInfo{
				Id:       loginResponse.UserId,
				Username: loginResponse.Username,
			},
		}

		state.Set(mod.appState, state.UserSessionProp, session)

		ses, ok := state.Get[state.UserSession](mod.appState, state.UserSessionProp)

		if !ok {
			Alert(mod.pageNav.Pages, "auth:login:alert:err", "State error")
			return
		}

		state.Set(mod.appState, state.UserSessionProp, ses)

		passwordInput.SetText("")
		emailInput.SetText("")

		brochatUser, err := mod.brochatClient.GetUser(&chat.AuthInfo{
			AccessToken: ses.Auth.AccessToken,
			TokenType:   "Bearer",
		}, ses.Info.Id)

		if err != nil {
			Alert(mod.pageNav.Pages, "auth:login:alert:err", err.Error())
			return
		}

		state.Set(mod.appState, state.BrochatUserInfo, brochatUser)

		mod.pageNav.NavigateTo(HOME_MENU_PAGE)
	})

	loginForm.AddButton("Back", func() {
		mod.pageNav.NavigateTo(WELCOME_PAGE)
	})

	grid.AddItem(loginForm, 1, 1, 1, 1, 0, 0, true)

	mod.pageNav.Register(LOGIN_PAGE, grid, true, false, func() {
		loginForm.SetFocus(0)

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
	REGISTRATION_MODAL_INFO      = "auth:register:alert:info"
	REGISTRATION_MODAL_ERR       = "auth:register:alert:err"
	REGISTRATION_SUCCESS_MESSAGE = "Registration Successful. A verification email has been sent to the email address provided."
)

func (mod *AuthModule) setupRegistrationPage() {
	grid := tview.NewGrid()
	grid.SetBackgroundColor(DefaultBackgroundColor)

	grid.SetRows(4, 0, 6)
	grid.SetColumns(0, 70, 0)

	registrationForm := tview.NewForm()
	registrationForm.SetBorder(true).SetTitle(" BroChat - Register ").SetTitleAlign(tview.AlignCenter)
	registrationForm.SetBackgroundColor(AccentBackgroundColor)
	registrationForm.SetFieldBackgroundColor(AccentColorTwoColorCode)
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

			request := &idam.UserRegistrationRequest{
				Email:    email,
				Password: password,
				Username: username,
			}

			_, err := mod.userAuthClient.Register("brochat", request)

			if err != nil {
				errMessage := err.Error()

				idamErr, ok := err.(*idam.ErrorResponse)

				if ok {
					switch idamErr.Code {
					case idam.RequestValidationFailure:
						errMessage = "Registration Failed - Request Validation Error"
						alertErrors(mod.pageNav.Pages, REGISTRATION_MODAL_ERR, errMessage, idamErr.Details)
						return
					case idam.InvalidCredentials:
						errMessage = "Registration Failed - Invalid Credentials"
					case idam.UnhandledError:
						errMessage = "Registration Failed - An Unexpected Error Occurred"
					}
				}

				Alert(mod.pageNav.Pages, REGISTRATION_MODAL_ERR, errMessage)
				return
			}

			AlertWithDoneFunc(mod.pageNav.Pages, REGISTRATION_MODAL_INFO, REGISTRATION_SUCCESS_MESSAGE, func(buttonIndex int, buttonLabel string) {
				mod.pageNav.Pages.HidePage(REGISTRATION_MODAL_INFO).RemovePage(REGISTRATION_MODAL_INFO)
				mod.pageNav.NavigateTo(WELCOME_PAGE)
			})
		}).
		AddButton("Back", func() { mod.pageNav.NavigateTo(WELCOME_PAGE) })

	grid.AddItem(registrationForm, 1, 1, 1, 1, 0, 0, true)

	mod.pageNav.Register(REGISTER_PAGE, grid, true, false, func() {
		registrationForm.SetFocus(0)

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

func alertErrors(pages *tview.Pages, id, errMessage string, messages []string) {
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
