package controller

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jesper-nord/recurringly-backend/util"
	"github.com/rs/cors"
	"net/http"
	"os"
	"strings"
)

var (
	clientHost   = os.Getenv("CLIENT_HOST")
	corsSettings = cors.New(cors.Options{
		AllowedOrigins:   []string{clientHost},
		AllowedMethods:   []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodHead},
		AllowCredentials: true,
	})
)

func CorsMiddleware(next http.Handler) http.Handler {
	return corsSettings.Handler(next)
}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if len(auth) == 0 || !strings.Contains(auth, "Bearer ") {
			http.Error(w, "missing authentication", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		_, err := util.ParseJwt(strings.Replace(auth, "Bearer ", "", 1), claims)
		if err != nil {
			http.Error(w, "invalid authentication", http.StatusUnauthorized)
			return
		}

		userId, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			http.Error(w, "missing userId", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", userId.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
