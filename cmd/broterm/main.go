package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/broterm/internal/ui"
	"github.com/dmars8047/idamlib/idam"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

func main() {

	// Configure logging
	file, err := provisionLogFile()

	if err != nil {
		log.Fatalf("Fatal error: log files could not be configured - %v", err)
	}

	defer file.Close()

	log.SetOutput(file)

	// Setup the http client
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Setup dependencies
	userAuthClient := idam.NewUserAuthClient(httpClient, "http://localhost:8083")

	brochatClient := chat.NewBroChatClient(httpClient, "http://localhost:8083")

	// Configure the application
	app := tview.NewApplication()

	// Setup the page navigator
	nav := ui.NewNavigator()

	// Setup the application context
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	appContext := state.NewApplicationContext(context)

	dialer := &websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	feedClient := state.NewFeedClient(dialer, "localhost:8083", brochatClient, appContext)

	// Setup the welcome page
	welcomePage := ui.NewWelcomePage()
	welcomePage.Setup(app, appContext, nav)

	// Setup the registration page
	registrationPage := ui.NewRegistrationPage(userAuthClient)
	registrationPage.Setup(app, appContext, nav)

	// Setup the login page
	loginPage := ui.NewLoginPage(userAuthClient, brochatClient, feedClient)
	loginPage.Setup(app, appContext, nav)

	// Setup the forgot password page
	forgotPasswordPage := ui.NewForgotPasswordPage(userAuthClient)
	forgotPasswordPage.Setup(app, appContext, nav)

	// Setup the chat page
	chatPage := ui.NewChatPage(brochatClient, feedClient)
	chatPage.Setup(app, appContext, nav)

	// Setup the home page
	homePage := ui.NewHomePage(userAuthClient)
	homePage.Setup(app, appContext, nav)

	// Setup the friends list page
	friendsListPage := ui.NewFriendsListPage(brochatClient, feedClient)
	friendsListPage.Setup(app, appContext, nav)

	// Setup the find a friend page
	findAFriendPage := ui.NewFindAFriendPage(brochatClient)
	findAFriendPage.Setup(app, appContext, nav)

	// Setup the accept friend request page
	acceptFriendRequestPage := ui.NewAcceptFriendRequestPage(brochatClient, feedClient)
	acceptFriendRequestPage.Setup(app, appContext, nav)

	// Setup the room list page
	roomListPage := ui.NewRoomListPage(brochatClient, feedClient)
	roomListPage.Setup(app, appContext, nav)

	// Setup the room editor page
	roomEditorPage := ui.NewRoomEditorPage(brochatClient)
	roomEditorPage.Setup(app, appContext, nav)

	// Setup the room finder page
	roomFinderPage := ui.NewRoomFinderPage(brochatClient)
	roomFinderPage.Setup(app, appContext, nav)

	// Start the application.
	err = app.SetRoot(nav.Pages, true).Run()

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}

const DEFAULT_LOG_DIRECTORY_NAME = ".broterm_logs"
const MAX_NUM_LOG_FILES = 10

// provisionLogFile creates a log file in the user's home directory and returns a file handle to it
// If the log file already exists, it will be opened and the file handle will be returned
// If the log file does not exist, it will be created and the file handle will be returned
// If the log file cannot be created or opened, an error will be returned
// This function also administers the log directory by deleting the oldest log file if the number of log files exceeds MAX_NUM_LOG_FILES.
func provisionLogFile() (*os.File, error) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Create the path to the log directory
	logDir := filepath.Join(homeDir, DEFAULT_LOG_DIRECTORY_NAME)

	// Create the log directory if it does not exist
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.Mkdir(logDir, os.ModePerm)

		if err != nil {
			return nil, err
		}
	}

	// Check if there are too many log dirEntries and delete the oldest one
	dirEntries, err := os.ReadDir(logDir)

	if err != nil {
		return nil, err
	}

	if len(dirEntries) >= MAX_NUM_LOG_FILES {
		oldestFile, err := dirEntries[0].Info()

		if err != nil {
			return nil, err
		}

		for _, dirEntry := range dirEntries {
			file, err := dirEntry.Info()

			if err != nil {
				return nil, err
			}

			if file.ModTime().Before(oldestFile.ModTime()) {
				oldestFile = file
			}
		}
		err = os.Remove(fmt.Sprintf("%s/%s", DEFAULT_LOG_DIRECTORY_NAME, oldestFile.Name()))

		if err != nil {
			return nil, err
		}
	}

	logfileName := fmt.Sprintf("broterm_%s.log", time.Now().Format("2006_01_02"))

	logFilePath := filepath.Join(logDir, logfileName)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	return file, err
}
