package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// A client for making http calls to the IDAM service's user account serving endpoints
// This client should be used for user facing calls to IDAM.
type UserAuthClient struct {
	httpClient *http.Client
	baseUrl    string
}

const (
	UserRegistrationAccountUrlSuffix  = "/api/idam/user-account/applications/:appId/register"
	UserLoginUrlSuffix                = "/api/idam/user-account/applications/:appId/login"
	UserAccountVerifyAccountUrlSuffix = "/api/idam/user-account/applications/:appId/verify-account"
	UserLogoutUrlSuffix               = "/api/idam/user-account/logout"
	InitiateUserPasswordResetUrl      = "/api/idam/user-account/applications/:appId/initiate-password-reset"
	ExecuteUserPasswordResetUrl       = "/api/idam/user-account/applications/:appId/execute-password-reset"
)

// Function to create a new IdamAuthService
func NewUserAuthClient(httpClient *http.Client, baseUrl string) *UserAuthClient {
	return &UserAuthClient{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

// Register method to call the user account registration endpoint
func (client *UserAuthClient) Register(appId string, request *UserRegistrationRequest) (*UserRegistrationResponse, error) {
	urlSuffix := strings.Replace(UserRegistrationAccountUrlSuffix, ":appId", appId, 1)

	// Parse the base URL
	base, err := url.Parse(client.baseUrl)

	if err != nil {
		return nil, err
	}

	// Parse the suffix as a URL
	suffix, err := url.Parse(urlSuffix)

	if err != nil {
		return nil, err
	}

	// Resolve to correct URL
	resolvedURL := base.ResolveReference(suffix)

	requestBodyBytes, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, resolvedURL.String(), bytes.NewReader(requestBodyBytes))
	if err != nil {
		return nil, err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Call the IDAM service /api/idam/user-account/applications/:appId/register endpoint and validate the token via http
	response, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		// UnMarhsal the response body into an ErrorResponse object
		var errorResponse ErrorResponse

		err = json.NewDecoder(response.Body).Decode(&errorResponse)

		if err != nil {
			return nil, fmt.Errorf("error decoding response body from idam service - %v", err)
		}

		return nil, &errorResponse
	}

	var usrRegResponse UserRegistrationResponse

	err = json.NewDecoder(response.Body).Decode(&usrRegResponse)

	if err != nil {
		return nil, errors.New("error decoding register user response body from idam service")
	}

	return &usrRegResponse, nil
}

// Login method to call the user account login endpoint
func (client *UserAuthClient) Login(appId string, request *UserLoginRequest) (*UserLoginResponse, error) {
	urlSuffix := strings.Replace(UserLoginUrlSuffix, ":appId", appId, 1)

	// Parse the base URL
	base, err := url.Parse(client.baseUrl)

	if err != nil {
		return nil, err
	}

	// Parse the suffix as a URL
	suffix, err := url.Parse(urlSuffix)

	if err != nil {
		return nil, err
	}

	// Resolve to correct URL
	resolvedURL := base.ResolveReference(suffix)

	requestBodyBytes, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, resolvedURL.String(), bytes.NewReader(requestBodyBytes))
	if err != nil {
		return nil, err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Call the IDAM service /api/idam/user-account/applications/:appId/login endpoint and validate the token via http
	response, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// UnMarhsal the response body into an ErrorResponse object
		var errorResponse ErrorResponse

		err = json.NewDecoder(response.Body).Decode(&errorResponse)

		if err != nil {
			return nil, fmt.Errorf("error decoding response body from idam service - %v", err)
		}

		return nil, &errorResponse
	}

	var loginResponse UserLoginResponse

	err = json.NewDecoder(response.Body).Decode(&loginResponse)

	if err != nil {
		return nil, errors.New("error decoding login user response body from idam service")
	}

	return &loginResponse, nil
}

// VerifyAccount method to call the user account verify account endpoint
func (client *UserAuthClient) VerifyAccount(appId string, request *UserAccountVerificationRequest) error {
	urlSuffix := strings.Replace(UserAccountVerifyAccountUrlSuffix, ":appId", appId, 1)

	// Parse the base URL
	base, err := url.Parse(client.baseUrl)

	if err != nil {
		return err
	}

	// Parse the suffix as a URL
	suffix, err := url.Parse(urlSuffix)

	if err != nil {
		return err
	}

	// Resolve to correct URL
	resolvedURL := base.ResolveReference(suffix)

	requestBodyBytes, err := json.Marshal(request)

	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPut, resolvedURL.String(), bytes.NewReader(requestBodyBytes))
	if err != nil {
		return err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Call the IDAM service /api/idam/user-account/applications/:appId/verify-account endpoint and validate the token via http
	response, err := client.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		// UnMarhsal the response body into an ErrorResponse object
		var errorResponse ErrorResponse

		err = json.NewDecoder(response.Body).Decode(&errorResponse)

		if err != nil {
			return fmt.Errorf("error decoding response body from idam service - %v", err)
		}

		return &errorResponse
	}

	return nil
}

// Logout method to call the user account logout endpoint
// If the authToken doesnt start with "Bearer " then it will be prepeneded and added to the Authorization header of the request
func (client *UserAuthClient) Logout(authToken string) error {

	// Parse the base URL
	base, err := url.Parse(client.baseUrl)

	if err != nil {
		return err
	}

	// Parse the suffix as a URL
	suffix, err := url.Parse(UserLogoutUrlSuffix)

	if err != nil {
		return err
	}

	// Resolve to correct URL
	resolvedURL := base.ResolveReference(suffix)

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, resolvedURL.String(), nil)
	if err != nil {
		return err
	}

	// If the authToken doesnt start with "Bearer " then add it
	if !strings.HasPrefix(authToken, "Bearer ") {
		authToken = "Bearer " + authToken
	}

	// Set the Authorization header to the token
	req.Header.Set("Authorization", authToken)

	// Call the IDAM service /api/idam/user-account/logout endpoint and validate the token via http
	response, err := client.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// UnMarhsal the response body into an ErrorResponse object
		var errorResponse ErrorResponse

		err = json.NewDecoder(response.Body).Decode(&errorResponse)

		if err != nil {
			return fmt.Errorf("error decoding response body from idam service - %v", err)
		}

		return &errorResponse
	}

	return nil
}

// InitiatePasswordReset method to call the user account initiate password reset endpoint
func (client *UserAuthClient) InitiatePasswordReset(appId string, request *UserPasswordResetInitiationRequest) error {
	urlSuffix := strings.Replace(InitiateUserPasswordResetUrl, ":appId", appId, 1)

	// Parse the base URL
	base, err := url.Parse(client.baseUrl)

	if err != nil {
		return err
	}

	// Parse the suffix as a URL
	suffix, err := url.Parse(urlSuffix)

	if err != nil {
		return err
	}

	// Resolve to correct URL
	resolvedURL := base.ResolveReference(suffix)

	requestBodyBytes, err := json.Marshal(request)

	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, resolvedURL.String(), bytes.NewReader(requestBodyBytes))
	if err != nil {
		return err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Call the IDAM service /api/idam/user-account/initiate-password-reset endpoint and validate the token via http
	response, err := client.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// UnMarhsal the response body into an ErrorResponse object
		var errorResponse ErrorResponse

		err = json.NewDecoder(response.Body).Decode(&errorResponse)

		if err != nil {
			return fmt.Errorf("error decoding response body from idam service - %v", err)
		}

		return &errorResponse
	}

	return nil
}

// ExecutePasswordReset method to call the user account execute password reset endpoint
func (client *UserAuthClient) ExecutePasswordReset(appId string, request *UserPasswordResetExecutionRequest) error {
	urlSuffix := strings.Replace(ExecuteUserPasswordResetUrl, ":appId", appId, 1)

	// Parse the base URL
	base, err := url.Parse(client.baseUrl)

	if err != nil {
		return err
	}

	// Parse the suffix as a URL
	suffix, err := url.Parse(urlSuffix)

	if err != nil {
		return err
	}

	// Resolve to correct URL
	resolvedURL := base.ResolveReference(suffix)

	requestBodyBytes, err := json.Marshal(request)

	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPut, resolvedURL.String(), bytes.NewReader(requestBodyBytes))
	if err != nil {
		return err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Call the IDAM service /api/idam/user-account/execute-password-reset endpoint and validate the token via http
	response, err := client.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		// UnMarhsal the response body into an ErrorResponse object
		var errorResponse ErrorResponse

		err = json.NewDecoder(response.Body).Decode(&errorResponse)

		if err != nil {
			return fmt.Errorf("error decoding response body from idam service - %v", err)
		}

		return &errorResponse
	}

	return nil
}
