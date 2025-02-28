package handlers

import (
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			SendError(w, "Требуется токен авторизации", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			SendError(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			SendError(w, "Неверный или просроченный токен", http.StatusUnauthorized)
			return
		}

		r.Header.Set("Username", claims.Username)
		next(w, r)
	}
}