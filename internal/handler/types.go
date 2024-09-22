package handler

import "time"

type AuthUserReq struct {
	GUID  string `json:"guid"`
	Email string `json:"email"`
}
type refreshUserReq struct {
	GUID         string `json:"guid"`
	RefreshToken string `json:"refresh_token"`
}

type AuthUserRes struct {
	GUID                 string    `json:"guid"`
	AccessToken          string    `json:"access_token"`
	RefreshToken         string    `json:"refresh_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
