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
	err = c.Database.Model(&entity.User{}).Where("email = ?", strings.ToLower(request.Email)).Take(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := generateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.LoginResponse{User: util.UserToApiModel(user), Token: token})
}

func (c AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterUserRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.Database.Model(&entity.User{}).Where("email = ?", strings.ToLower(request.Email)).Take(&entity.User{}).Error
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.UserToApiModel(user))
}

func generateToken(user entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
		"iat": time.Now().Unix(),
		"jti": uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
