package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/UAssylbek/leaderboard/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

// SetJWTSecret задает секретный ключ для JWT
func SetJWTSecret(key []byte) {
	jwtKey = key
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		err = db.QueryRow(
			"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
			user.Username, string(hashedPassword),
		).Scan(&user.ID)
		if err != nil {
			http.Error(w, "Пользователь уже существует или ошибка", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь зарегистрирован"})
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}

		var storedPassword string
		var id int
		err := db.QueryRow(
			"SELECT id, password FROM users WHERE username = $1",
			user.Username,
		).Scan(&id, &storedPassword)
		if err != nil {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}
}
