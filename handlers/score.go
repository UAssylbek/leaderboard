package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/UAssylbek/leaderboard/models"
	"github.com/redis/go-redis/v9"
)

func SubmitScore(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("Username")
		if username == "" {
			SendError(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}

		var score models.Score
		if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
			SendError(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		score.Username = username

		var userID int
		err := db.QueryRow(
			"SELECT id FROM users WHERE username = $1",
			score.Username,
		).Scan(&userID)
		if err != nil {
			SendError(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		score.UserID = userID

		_, err = db.Exec(
			"INSERT INTO scores (user_id, game, score) VALUES ($1, $2, $3)",
			score.UserID, score.Game, score.Score,
		)
		if err != nil {
			log.Println("Ошибка сохранения очков:", err)
			SendError(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		redisKey := "leaderboard:global"
		err = redisClient.ZAdd(ctx, redisKey, redis.Z{
			Score:  float64(score.Score),
			Member: score.Username,
		}).Err()
		if err != nil {
			log.Println("Ошибка обновления лидерборда в Redis:", err)
			SendError(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		updateMsg := fmt.Sprintf("%s updated score to %d", score.Username, score.Score)
		err = redisClient.Publish(ctx, "leaderboard_updates", updateMsg).Err()
		if err != nil {
			log.Println("Ошибка публикации в Redis Pub/Sub:", err)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Очки успешно обновлены"})
	})
}
