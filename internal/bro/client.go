package bro

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	GET_USER_URL_SUFFIX              = "/api/brochat/user"
	GET_USERS_URL_SUFFIX             = "/api/brochat/users"
	GET_CHANNEL_URL_SUFFIX           = "/api/brochat/channels/:channelId"
	GET_CHANNEL_MESSAGES_URL_SUFFIX  = "/api/brochat/channels/:channelId/messages"
	SEND_FRIEND_REQUEST_URL_SUFFIX   = "/api/brochat/friends/send-friend-request"
	ACCEPT_FRIEND_REQUEST_URL_SUFFIX = "/api/brochat/friends/accept-friend-request"
)

type BroChatUserClient struct {
	httpClient *http.Client
	baseUrl    string
}

func NewBroChatClient(httpClient *http.Client, baseUrl string) *BroChatUserClient {
	return &BroChatUserClient{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

type AuthInfo struct {
	// The JWT token.
	AccessToken string
	// The type of the token. Most likely "Bearer".
	TokenType string
}

// Get User
func (c *BroChatUserClient) GetUser(authInfo *AuthInfo, userId string) (*User, error) {

	base, err := url.Parse(c.baseUrl)

	if err != nil {
		return nil, err
	}

	suffix, err := url.Parse(GET_USER_URL_SUFFIX)

	if err != nil {
		return nil, err
	}

	resolvedUrl := base.ResolveReference(suffix)

	// Create a new request using http
	req, err := http.NewRequest("GET", resolvedUrl.String(), nil)

	if err != nil {
		return nil, err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", authInfo.TokenType+" "+authInfo.AccessToken)

	// Send req using http Client
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return nil, errors.New("user not found")
		}

		if res.StatusCode == 401 {
			return nil, errors.New("unauthorized")
		}

		if res.StatusCode == 403 {
			return nil, errors.New("forbidden")
		}

		return nil, errors.New("unexpected status code")
	}

	var user User

	err = json.NewDecoder(res.Body).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get Users
func (c *BroChatUserClient) GetUsers(authInfo *AuthInfo, excludeFriends, excludeSelf bool, page, pageSize int, usernameFilter string) ([]UserInfo, error) {
	base, err := url.Parse(c.baseUrl)

	if err != nil {
		return nil, err
	}

	suffix, err := url.Parse(GET_USERS_URL_SUFFIX)

	if err != nil {
		return nil, err
	}

	resolvedUrl := base.ResolveReference(suffix)

	// Add query params
	q := resolvedUrl.Query()

	q.Add("exclude-friends", strconv.FormatBool(excludeFriends))
	q.Add("exclude-self", strconv.FormatBool(excludeSelf))
	q.Add("page", strconv.Itoa(page))
	q.Add("page-size", strconv.Itoa(pageSize))

	if usernameFilter != "" {
		q.Add("username-filter", usernameFilter)
	}

	resolvedUrl.RawQuery = q.Encode()

	// Create a new request using http
	req, err := http.NewRequest("GET", resolvedUrl.String(), nil)

	if err != nil {
		return nil, err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", authInfo.TokenType+" "+authInfo.AccessToken)

	// Send req using http Client
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 401 {
			return nil, errors.New("unauthorized")
		}

		if res.StatusCode == 403 {
			return nil, errors.New("forbidden")
		}

		return nil, errors.New("unexpected status code")
	}

	var users = make([]UserInfo, 0)

	err = json.NewDecoder(res.Body).Decode(&users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// Get Channel
func (c *BroChatUserClient) GetChannelManifest(authInfo *AuthInfo, channelId string) (*ChannelManifest, error) {
	base, err := url.Parse(c.baseUrl)

	if err != nil {
		return nil, err
	}

	suffix, err := url.Parse(strings.Replace(GET_CHANNEL_URL_SUFFIX, ":channelId", channelId, 1))

	if err != nil {
		return nil, err
	}

	resolvedUrl := base.ResolveReference(suffix)

	// Create a new request using http
	req, err := http.NewRequest("GET", resolvedUrl.String(), nil)

	if err != nil {
		return nil, err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", authInfo.TokenType+" "+authInfo.AccessToken)

	// Send req using http Client
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return nil, errors.New("channel not found")
		} else if res.StatusCode == 401 {
			return nil, errors.New("unauthorized")
		} else if res.StatusCode == 403 {
			return nil, errors.New("forbidden")
		}

		return nil, errors.New("unexpected status code")
	}

	var channel ChannelManifest

	err = json.NewDecoder(res.Body).Decode(&channel)

	if err != nil {
		return nil, err
	}

	return &channel, nil
}

// Get Channel Messages
func (c *BroChatUserClient) GetChannelMessages(authInfo *AuthInfo, channelId string) ([]ChatMessage, error) {
	base, err := url.Parse(c.baseUrl)

	if err != nil {
		return nil, err
	}

	suffix, err := url.Parse(strings.Replace(GET_CHANNEL_MESSAGES_URL_SUFFIX, ":channelId", channelId, 1))

	if err != nil {
		return nil, err
	}

	resolvedUrl := base.ResolveReference(suffix)

	// Create a new request using http
	req, err := http.NewRequest("GET", resolvedUrl.String(), nil)

	if err != nil {
		return nil, err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", authInfo.TokenType+" "+authInfo.AccessToken)

	// Send req using http Client
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 401 {
			return nil, errors.New("unauthorized")
		} else if res.StatusCode == 403 {
			return nil, errors.New("forbidden")
		}

		return nil, errors.New("unexpected status code")
	}

	var channels = make([]ChatMessage, 0)

	err = json.NewDecoder(res.Body).Decode(&channels)

	if err != nil {
		return nil, err
	}

	return channels, nil
}

// Send Friend Request
func (c *BroChatUserClient) SendFriendRequest(authInfo *AuthInfo, request *SendFriendRequestRequest) error {
	base, err := url.Parse(c.baseUrl)

	if err != nil {
		return err
	}

	suffix, err := url.Parse(SEND_FRIEND_REQUEST_URL_SUFFIX)

	if err != nil {
		return err
	}

	resolvedUrl := base.ResolveReference(suffix)

	requestBodyBytes, err := json.Marshal(request)

	if err != nil {
		return err
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPut, resolvedUrl.String(), bytes.NewReader(requestBodyBytes))

	if err != nil {
		return err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", authInfo.TokenType+" "+authInfo.AccessToken)

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Send req using http Client
	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		if res.StatusCode == 401 {
			return errors.New("unauthorized")
		} else if res.StatusCode == 403 {
			return errors.New("forbidden")
		} else if res.StatusCode == 404 {
			return errors.New("user not found")
		} else if res.StatusCode == 409 {
			return errors.New("friend request already exists or users are already a friend")
		} else if res.StatusCode == 400 {
			return errors.New("bad request")
		}

		return errors.New("unexpected status code")
	}

	return nil
}

// Accept Friend Request
func (c *BroChatUserClient) AcceptFriendRequest(authInfo *AuthInfo, request *AcceptFriendRequestRequest) error {
	base, err := url.Parse(c.baseUrl)

	if err != nil {
		return err
	}

	suffix, err := url.Parse(ACCEPT_FRIEND_REQUEST_URL_SUFFIX)

	if err != nil {
		return err
	}

	resolvedUrl := base.ResolveReference(suffix)

	requestBodyBytes, err := json.Marshal(request)

	if err != nil {
		return err
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPut, resolvedUrl.String(), bytes.NewReader(requestBodyBytes))

	if err != nil {
		return err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", authInfo.TokenType+" "+authInfo.AccessToken)

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Send req using http Client
	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		if res.StatusCode == 401 {
			return errors.New("unauthorized")
		} else if res.StatusCode == 403 {
			return errors.New("forbidden")
		} else if res.StatusCode == 404 {
			return errors.New("user not found or friend request not found")
		} else if res.StatusCode == 400 {
			return errors.New("bad request")
		}

		return errors.New("unexpected status code")
	}

	return nil
}
