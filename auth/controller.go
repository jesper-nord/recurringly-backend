package auth

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Controller struct {
	Service Service
}

func (c Controller) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, tokens, err := c.Service.Login(request.Username, request.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Printf("logged in: user '%s'", user.ID.String())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{Tokens: tokens})
}

func (c Controller) Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterUserRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, tokens, err := c.Service.RegisterUser(request.Username, request.Password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	log.Printf("registered: user '%s' with username '%s'", user.ID.String(), user.Username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{Tokens: tokens})
}

func (c Controller) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var request RefreshTokenRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := c.Service.RefreshAccessToken(request.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{Tokens: tokens})
}
