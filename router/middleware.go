package router

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jesper-nord/recurringly-backend/auth"
	"net/http"
	"strings"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if len(authHeader) == 0 || !strings.Contains(authHeader, "Bearer ") {
			http.Error(w, "missing authentication", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		_, err := auth.ParseJwt(strings.Replace(authHeader, "Bearer ", "", 1), claims)
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
