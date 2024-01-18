package bro

import (
	"errors"
	"time"
)

type RelationshipType uint8

const (
	// This is the default relationship type. It is used when two users are not friends.
	RELATIONSHIP_TYPE_DEFAULT RelationshipType = 1 << iota
	// This relationship type is used when two users are friends.
	RELATIONSHIP_TYPE_FRIEND
	// This relationship type is applied when the user has recieved a friend request from another user.
	RELATIONSHIP_TYPE_FRIEND_REQUEST_RECIEVED
	// This relationship type is applied when the user has sent a friend request to another user.
	RELATIONSHIP_TYPE_FRIENDSHIP_REQUESTED
)

type ChannelType uint8

const (
	// A channel that is used for direct messaging between two users.
	CHANNEL_TYPE_DIRECT_MESSAGE ChannelType = iota
)

// A ChatMessage represents a text message sent in to chat channel.
type ChatMessage struct {
	// The ID of the message.
	Id string `json:"id"`
	// The ID of the channel that the message was sent in.
	ChannelId string `json:"channel_id"`
	// The ID of the user that sent the message.
	SenderUserId string `json:"sender_user_id"`
	// The content of the message.
	Content string `json:"content"`
	// The time that the message was sent.
	RecievedAtUtc time.Time `json:"recieved_at_utc"`
}

type ChatMessageRequest struct {
	// The ID of the channel that the message is being sent in.
	ChannelId string `json:"channel_id"`
	// The content of the message.
	Content string `json:"content"`
}

type UserRelationship struct {
	// The id of the user that the relationship is with.
	UserId string `json:"user_id"`
	// The type of relationship.
	Type RelationshipType `json:"type"`
	// Direct Message Channel Id
	DirectMessageChannelId string `json:"direct_message_channel_id"`
	// Username of the user the relationship is with
	Username string `json:"username"`
	// When the user was last online
	LastOnlineUtc time.Time `json:"last_online_utc"`
	// IsOnline is true if the user is online
	IsOnline bool `json:"is_online"`
}

type User struct {
	// The user's Id. This is the same as the Id in the idam service.
	Id string `json:"id"`
	// The user's username. This is the same as the username.
	Username string `json:"username"`
	// The users relationships list.
	Relationships []UserRelationship `json:"relationships"`
	// When the user was last online
	LastOnlineUtc time.Time `json:"last_online_utc"`
	// CreatedAtUtc is when the user was created
	CreatedAtUtc time.Time `json:"created_at_utc"`
}

type UserInfo struct {
	// The user's ID. This is the same as the ID in the idam service.
	ID string `json:"id"`
	// The user's username.
	Username string `json:"username"`
	// When the user was last online
	LastOnlineUtc time.Time `json:"last_online_utc"`
}

type SendFriendRequestRequest struct {
	// The ID of the user that the friend request is being sent to.
	RequestedUserId string `json:"requested_user_id"`
}

type AcceptFriendRequestRequest struct {
	// The ID of the user that sent the friend request.
	InitiatingUserId string `json:"initiating_user_id"`
}

type ChannelManifest struct {
	// The ID of the channel.
	ID string `json:"id"`
	// The type of the channel.
	Type ChannelType `json:"type"`
	// The users that are members of the channel. This is a list of user info.
	Users []UserInfo `json:"users"`
}

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrChannelNotFound = errors.New("channel not found")
)
