package controller

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jesper-nord/recurringly-backend/util"
	"log"
	"net/http"
	"strings"
)

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

		log.Printf("authenticated user %s", userId.String())
		ctx := context.WithValue(r.Context(), "user", userId.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
