package test

import (
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"

	"test_Auth/internal/handler"
)

const (
	host = "localhost:8082"
)

func TestAuth_HappyReg(t *testing.T) {
	testCases := []struct {
		name  string
		guid  string
		email string
		error string
	}{
		{
			name:  "Valid request",
			guid:  gofakeit.UUID(),
			email: gofakeit.Email(),
		},
		{
			name:  "Invalid GUID",
			guid:  "123",
			email: gofakeit.Email(),
			error: "field uuid is not valid UUID",
		},
		{
			name:  "Invalid email",
			guid:  gofakeit.UUID(),
			email: "notemail",
		},
		{
			name:  "Empty request",
			guid:  gofakeit.UUID(),
			email: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/register").
				WithJSON(handler.AuthUserReq{
					GUID:  tc.guid,
					Email: tc.email,
				}).
				Expect().Status(200).
				JSON().Object().
				NotContainsKey("email").
				ContainsKey("access_token").
				ContainsKey("refresh_token").
				ContainsKey("guid").
				ContainsKey("access_token_expires_at")

			if tc.error != "" {
				resp.NotContainsKey("email")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}
		})
	}
}
