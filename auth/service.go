package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SIGNING_SECRET"))
)

type authService struct {
	Repository Repository
}

func NewService(repository Repository) Service {
	return &authService{
		Repository: repository,
	}
}

func (a *authService) Login(username, password string) (*User, TokenPair, error) {
	user, err := a.Repository.FindUser(username)
	if user == nil {
		return nil, TokenPair{}, fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, TokenPair{}, fmt.Errorf("password does not match: %w", err)
	}
	return a.generateTokenPair(user)
}

func (a *authService) RegisterUser(username, password string) (*User, TokenPair, error) {
	username = strings.ToLower(username)
	_, err := a.Repository.FindUser(username)
	if err == nil {
		return nil, TokenPair{}, fmt.Errorf("user already exists: %w", err)
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, TokenPair{}, err
	}

	user, err := a.Repository.SaveUser(&User{
		Username: username,
		Password: string(pw),
	})
	if err != nil {
		return nil, TokenPair{}, err
	}
	return a.generateTokenPair(user)
}

func (a *authService) RefreshAccessToken(refreshToken string) (TokenPair, error) {
	claims := jwt.MapClaims{}
	token, err := ParseJwt(refreshToken, claims)
	if err != nil || !token.Valid {
		return TokenPair{}, errors.New("invalid refresh token")
	}
	sub, err := claims.GetSubject()
	if err != nil {
		return TokenPair{}, err
	}
	userId, err := strconv.Atoi(sub)
	if err != nil {
		return TokenPair{}, err
	}
	user, err := a.Repository.GetUser(UserId(userId))
	if err != nil {
		return TokenPair{}, err
	}

	_, tokens, err := a.generateTokenPair(user)
	return tokens, err
}

func (a *authService) generateTokenPair(user *User) (*User, TokenPair, error) {
	env := os.Getenv("APP_ENV")
	accessClaims := jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Minute * 15).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"aud": "recurringly-" + env,
		"iss": "recurringly",
		"jti": uuid.New().String(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := at.SignedString(jwtSecret)
	if err != nil {
		return nil, TokenPair{}, err
	}

	refreshClaims := jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		// should practically never expire
		"exp": time.Now().Add(time.Hour * 24 * 365 * 10).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"aud": "recurringly-" + env,
		"iss": "recurringly",
		"jti": uuid.New().String(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := rt.SignedString(jwtSecret)
	if err != nil {
		return nil, TokenPair{}, err
	}

	return user, TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ParseJwt(token string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
}
