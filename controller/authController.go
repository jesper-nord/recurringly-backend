package controller

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"recurringly-backend/dto"
	"recurringly-backend/entity"
	"recurringly-backend/util"
	"strings"
	"time"
)

type AuthController struct {
	Database *gorm.DB
}

func (c AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user entity.User
	err = c.Database.Where("email = ?", strings.ToLower(request.Email)).Take(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenPair, err := generateTokenPair(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.LoginResponse{Tokens: tokenPair})
}

func (c AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterUserRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.Database.Where("email = ?", strings.ToLower(request.Email)).Take(&entity.User{}).Error
	if err == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := entity.User{
		Email:    strings.ToLower(request.Email),
		Password: string(password),
	}
	err = c.Database.Create(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenPair, err := generateTokenPair(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.LoginResponse{Tokens: tokenPair})
}

func (c AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var request dto.RefreshTokenRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := jwt.MapClaims{}
	token, err := util.ParseJwt(request.RefreshToken, claims)
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user entity.User
	err = c.Database.Take(&user, userId).Error
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenPair, err := generateTokenPair(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.RefreshTokenResponse{Tokens: tokenPair})
}

func generateTokenPair(user entity.User) (dto.TokenPair, error) {
	secret := []byte(os.Getenv("JWT_SIGNING_SECRET"))

	accessClaims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
		"iat":   time.Now().Unix(),
		"jti":   uuid.New().String(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := at.SignedString(secret)
	if err != nil {
		return dto.TokenPair{}, err
	}

	refreshClaims := jwt.MapClaims{
		"sub": user.ID,
		// should practically never expire
		"exp": time.Now().Add(time.Hour * 24 * 365 * 10).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := rt.SignedString(secret)
	if err != nil {
		return dto.TokenPair{}, err
	}

	return dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
