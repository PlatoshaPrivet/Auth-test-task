package handler

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"test_Auth/internal/dbconn"
	"test_Auth/internal/refresh"
	"test_Auth/internal/token"
	"time"

	"github.com/google/uuid"
)

type handler struct {
	TokenMaker *token.JWTMaker
}

func NewHandler(secretKey string) *handler {
	return &handler{
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *handler) AuthUser(w http.ResponseWriter, r *http.Request) {
	var u AuthUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(u.GUID)
	if err != nil {
		http.Error(w, "incorrect GUID", http.StatusBadRequest)
		return
	}
	_, err = mail.ParseAddress(u.Email)
	if err != nil {
		http.Error(w, "incorrect email", http.StatusBadRequest)
		return
	}

	ip := r.RemoteAddr

	// create JWT
	accessToken, accessClaims, err := h.TokenMaker.CreateJWT(u.GUID, u.Email, ip, 15*time.Minute)
	if err != nil {
		http.Error(w, "error on creating JWT token", http.StatusInternalServerError)
		return
	}

	//create refresh
	refreshToken, err := token.CreateRefresh(u.GUID, u.Email)
	if err != nil {
		http.Error(w, "error on creating Refresh token", http.StatusInternalServerError)
		return
	}

	res := AuthUserRes{
		GUID:                 u.GUID,
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) HdlRefresh(w http.ResponseWriter, r *http.Request) {
	db := dbconn.Open()
	defer db.Close()

	var u refreshUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	ip := r.RemoteAddr

	refreshToken, email, err := refresh.UpdateRefresh(u.GUID, u.RefreshToken)
	if err != nil {
		http.Error(w, "error on update token", http.StatusInternalServerError)
		return
	}

	accessToken, accessClaims, err := h.TokenMaker.CreateJWT(u.GUID, email, ip, 15*time.Minute)
	if err != nil {
		http.Error(w, "error on creating token", http.StatusInternalServerError)
		return
	}

	res := AuthUserRes{
		GUID:                 u.GUID,
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
