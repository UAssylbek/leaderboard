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
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем токен (заглушка, позже добавим проверку JWT)

		var score models.Score
		if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}

		// Получаем ID пользователя из базы по username
		var userID int
		err := db.QueryRow(
			"SELECT id FROM users WHERE username = $1",
			score.Username, // Добавим поле Username в models.Score
		).Scan(&userID)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		score.UserID = userID // Заполняем UserID для записи в БД

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

		// Обновляем лидерборд в Redis (глобальный лидерборд)
		ctx := context.Background()
		redisKey := "leaderboard:global"
		err = redisClient.ZAdd(ctx, redisKey, redis.Z{
			Score:  float64(score.Score),
			Member: score.Username, // Используем Username как идентификатор в Redis
		}).Err()
		if err != nil {
			log.Println("Ошибка обновления лидерборда в Redis:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Очки успешно отправлены"})
	}
}