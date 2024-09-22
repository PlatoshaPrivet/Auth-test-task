package refresh

import (
	"crypto/rand"
	"encoding/base64"
	"test_Auth/internal/dbconn"

	"golang.org/x/crypto/bcrypt"
)

func UpdateRefresh(id string, userRefresh string) (string, string, error) {
	db := dbconn.Open()
	defer db.Close()

	var storedRefresh, email string
	err := db.QueryRow(`SELECT email, refresh_token FROM refresh_tokens WHERE user_id = $1`, id).Scan(&email, &storedRefresh)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedRefresh), []byte(userRefresh))
	if err != nil {
		return "", "", err
	}

	byteToken := make([]byte, 32)
	_, err = rand.Read(byteToken)
	if err != nil {
		return "", "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(byteToken)
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	//store data about client
	_, err = db.Exec(`UPDATE refresh_tokens SET refresh_token = $1 WHERE user_id = $2`, hashedToken, id)
	if err != nil {
		return "", "", err
	}

	return refreshToken, email, nil
}
