package state

import "time"

type UserSession struct {
	Auth UserInfo
	Info UserAuth
}

type UserInfo struct {
	Username string
	Id       string
}

type UserAuth struct {
	AccessToken     string
	TokenExpiration time.Duration
}
