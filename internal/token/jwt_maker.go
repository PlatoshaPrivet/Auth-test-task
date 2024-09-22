package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"test_Auth/internal/dbconn"
	"test_Auth/internal/services"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey}
}

func (maker *JWTMaker) CreateJWT(id string, email string, ip string, duration time.Duration) (string, *UserClaims, error) {
	//JWT token
	claims, err := NewUserClaims(id, email, ip, duration)
	if err != nil {
		return "", nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error signing token: %w", err)
	}

	return tokenStr, claims, nil
}

func CreateRefresh(id string, email string) (string, error) {
	db := dbconn.Open()
	defer db.Close()

	refreshTokenBytes := make([]byte, 32)
	_, err := rand.Read(refreshTokenBytes)
	if err != nil {
		return "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	//store data about client
	_, err = db.Exec(`INSERT INTO refresh_tokens (user_id, email, refresh_token) VALUES ($1, $2, $3)`, id, email, hashedRefreshToken)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (maker *JWTMaker) VerifyToken(tokenStr string, ip string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// verify the signing method
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	//checking ip if it has been changed
	if claims.IP != ip {
		services.SendWarning(claims.Email)
	}

	return claims, nil
}
