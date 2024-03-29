package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dmars8047/brolib/chat"
	"github.com/dmars8047/broterm/internal/config"
	"github.com/dmars8047/broterm/internal/state"
	"github.com/dmars8047/broterm/internal/ui"
	"github.com/dmars8047/idamlib/idam"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

func main() {

	const applicationVersion = "v0.1.11"

	// Look for the update command line argument
	if len(os.Args) > 1 {
		if os.Args[1] == "update" {
			fmt.Println("Updating Broterm...")

			var downloadUrl string

			// Get the os (windows, linux, darwin)
			if runtime.GOOS == "windows" {
				downloadUrl = "https://dev.brochat.app/binaries/Windows/broterm.exe"
			} else if runtime.GOOS == "linux" {
				downloadUrl = "https://dev.brochat.app/binaries/Linux/broterm"
			} else if runtime.GOOS == "darwin" {
				downloadUrl = "https://dev.brochat.app/binaries/OSX/broterm"
			} else {
				fmt.Println("Unsupported operating system")
				return
			}

			// Replace the current executable with the latest version
			sourceLocation := os.Args[0]

			err := downloadFile(downloadUrl, sourceLocation)

			if err != nil {
				log.Fatalf("Error updating Broterm: %v", err)
			}

			fmt.Printf("\nThe latest Broterm executable has been downloaded to %s.\n", sourceLocation)
			return
		}

		if os.Args[1] == "version" {
			fmt.Printf("Broterm Version - %s\n", applicationVersion)
			return
		}

		helpCommands := map[string]struct{}{
			"help":   {},
			"?":      {},
			"/?":     {},
			"/help":  {},
			"-?":     {},
			"--?":    {},
			"-help":  {},
			"--help": {},
		}

		if _, ok := helpCommands[os.Args[1]]; ok {
			fmt.Println("Broterm is a terminal based chat application that allows users to chat with friends and create chat rooms.\n" +
				"Usage: broterm [update|version|help]\n" +
				"\nupdate - updates the Broterm application\n" +
				"version - displays the version of the Broterm application\n" +
				"help - displays this help message\n" +
				"\nFor more information, visit https://dev.brochat.app")
			return
		}
	}

	// Configure logging
	config, file, err := provisionConfigFile()

	if err != nil {
		log.Fatalf("Broterm Version - %s\n\nFatal error: log files could not be configured - %v", applicationVersion, err)
	}

	defer file.Close()

	if config.LoggingEnabled {
		log.Printf("Broterm Version - %s\n\nBroterm logging is enabled. Writing logs to %s\n", applicationVersion, file.Name())
		log.SetOutput(file)
	} else {
		log.Printf("Broterm Version - %s\n\nBroterm logging is disabled.\n", applicationVersion)
		// supress logging
		log.SetOutput(&NullWriter{})
	}

	// Setup the http client
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	const hostAddr = "dev.marshall-labs.com"

	// Setup dependencies
	userAuthClient := idam.NewUserAuthClient(httpClient, "https://"+hostAddr)

	brochatClient := chat.NewBroChatClient(httpClient, "https://"+hostAddr)

	// Configure the application
	app := tview.NewApplication()

	// Setup the application context
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	appContext := state.NewApplicationContext(context, config.Theme)

	// Setup the page navigator
	nav := ui.NewNavigator(appContext)

	dialer := &websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	feedClient := state.NewFeedClient(dialer, hostAddr, brochatClient, appContext)

	// Setup the welcome page
	welcomePage := ui.NewWelcomePage(applicationVersion)
	welcomePage.Setup(app, appContext, nav)

	// Setup the app settings page
	appSettingsPage := ui.NewAppSettingsPage(config.LoggingEnabled)
	appSettingsPage.Setup(app, appContext, nav)

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

	// Set the background color of the navs pages
	theme := appContext.GetTheme()
	nav.Pages.SetBackgroundColor(theme.BackgroundColor)
	theme.ApplyGlobals()

	// Start the application.
	err = app.SetRoot(nav.Pages, true).Run()

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}

// IO.Writer that does nothing
type NullWriter struct{}

// Write does nothing
func (NullWriter) Write(p []byte) (int, error) {
	return len(p), nil

}

// provisionConfigFile creates a config directory in the user's home directory if it does not already exist
// It will read the config.json file in the config directory and return a ConfigSettings struct with the values from the file
// It also creates a log file in the user's home directory and returns a file handle to it
// If the log file already exists, it will be opened and the file handle will be returned
// If the log file does not exist, it will be created and the file handle will be returned
// If the log file cannot be created or opened, an error will be returned
// This function also administers the log directory by deleting the oldest log file if the number of log files exceeds MAX_NUM_LOG_FILES.
func provisionConfigFile() (*config.ConfigSettings, *os.File, error) {
	const maxNumLogFiles = 10

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}

	// Create the path to the config directory
	configDir := filepath.Join(homeDir, config.DEFAULT_CONFIG_DIRECTORY_NAME)

	// Create the config directory if it does not exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, os.ModePerm)

		if err != nil {
			return nil, nil, err
		}
	}

	// Check if there are too many log dirEntries and delete the oldest one
	dirEntries, err := os.ReadDir(configDir)

	if err != nil {
		return nil, nil, err
	}

	// Get the config.json file and read it into a new ConfigSettings struct
	configFilePath := filepath.Join(configDir, config.CONFIG_FILE_NAME)

	configSettings := config.NewConfigSettings()

	// If the file does not exist then create it, otherwise read it into the ConfigSettings struct
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		bytesToSave, err := json.Marshal(configSettings)

		if err != nil {
			return nil, nil, err
		}

		err = os.WriteFile(configFilePath, bytesToSave, 0644)

		if err != nil {
			return nil, nil, err
		}
	} else {
		configBytes, err := os.ReadFile(configFilePath)

		if err != nil {
			return nil, nil, err
		}

		err = json.Unmarshal(configBytes, configSettings)

		if err != nil {
			return nil, nil, err
		}
	}

	if !configSettings.LoggingEnabled {
		return configSettings, nil, nil
	}

	if len(dirEntries) >= maxNumLogFiles {
		oldestFile, err := dirEntries[0].Info()

		if err != nil {
			return nil, nil, err
		}

		for _, dirEntry := range dirEntries {
			file, err := dirEntry.Info()

			if err != nil {
				return nil, nil, err
			}

			if file.ModTime().Before(oldestFile.ModTime()) {
				oldestFile = file
			}
		}
		err = os.Remove(fmt.Sprintf("%s/%s", config.DEFAULT_CONFIG_DIRECTORY_NAME, oldestFile.Name()))

		if err != nil {
			return nil, nil, err
		}
	}

	logfileName := fmt.Sprintf("broterm_%s.log", time.Now().Format("2006_01_02"))

	logFilePath := filepath.Join(configDir, logfileName)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	return configSettings, file, err
}

func downloadFile(url, destination string) error {
	// Download the latest version of the application
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("Error downloading the latest version of Broterm: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error downloading the latest version of Broterm: %v", resp.Status)
	}

	// Create the file, overwriting if it already exists
	out, err := os.Create(destination)

	if err != nil {
		log.Fatalf("Error creating the file: %v", err)
	}

	defer out.Close()

	// Write the body to the file
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		log.Fatalf("Error writing the file: %v", err)
	}

	return nil
}
