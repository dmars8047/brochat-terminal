package state

import "time"

type UserSession struct {
	Auth UserAuth
	Info UserInfo
}

type UserInfo struct {
	Username string
	Id       string
}

type UserAuth struct {
	AccessToken     string
	TokenExpiration time.Time
}
