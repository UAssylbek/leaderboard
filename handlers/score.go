package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"github.com/UAssylbek/leaderboard/models"
	"github.com/redis/go-redis/v9"
)

func SubmitScore(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("Username")
		if username == "" {
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}

		var score models.Score
		if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		score.Username = username // Устанавливаем username из токена

		// Получаем ID пользователя из базы
		var userID int
		err := db.QueryRow(
			"SELECT id FROM users WHERE username = $1",
			score.Username,
		).Scan(&userID)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		score.UserID = userID

		// Сохраняем очки в PostgreSQL
		_, err = db.Exec(
			"INSERT INTO scores (user_id, game, score) VALUES ($1, $2, $3)",
			score.UserID, score.Game, score.Score,
		)
		if err != nil {
			log.Println("Ошибка сохранения очков:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Обновляем лидерборд в Redis
		ctx := context.Background()
		redisKey := "leaderboard:global"
		err = redisClient.ZAdd(ctx, redisKey, redis.Z{
			Score:  float64(score.Score),
			Member: score.Username,
		}).Err()
		if err != nil {
			log.Println("Ошибка обновления лидерборда в Redis:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Очки успешно отправлены"})
	})
}