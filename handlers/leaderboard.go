package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetLeaderboard(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		redisKey := "leaderboard:global"

		// Получаем топ-10 пользователей из Redis Sorted Set
		leaderboard, err := redisClient.ZRevRangeWithScores(ctx, redisKey, 0, 9).Result()
		if err != nil {
			http.Error(w, "Ошибка получения лидерборда", http.StatusInternalServerError)
			return
		}

		// Форматируем результат в JSON
		type LeaderboardEntry struct {
			Username string  `json:"username"`
			Score    float64 `json:"score"`
		}
		var result []LeaderboardEntry
		for _, entry := range leaderboard {
			result = append(result, LeaderboardEntry{
				Username: entry.Member.(string),
				Score:    entry.Score,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func GetRank(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Не указан username", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		redisKey := "leaderboard:global"

		// Получаем ранг пользователя (нумерация с 0, добавляем 1 для читаемости)
		rank, err := redisClient.ZRevRank(ctx, redisKey, username).Result()
		if err != nil {
			http.Error(w, "Пользователь не найден в лидерборде", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int64{"rank": rank + 1})
	}
}

func GetTopPlayers(db *sql.DB) http.HandlerFunc {
	return JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		if period == "" {
			period = "day" // По умолчанию за день
		}

		var since time.Time
		switch period {
		case "day":
			since = time.Now().Add(-24 * time.Hour)
		case "week":
			since = time.Now().Add(-7 * 24 * time.Hour)
		case "month":
			since = time.Now().Add(-30 * 24 * time.Hour)
		default:
			http.Error(w, "Неверный период. Используйте: day, week, month", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
			SELECT u.username, SUM(s.score) as total_score
			FROM scores s
			JOIN users u ON s.user_id = u.id
			WHERE s.created_at >= $1
			GROUP BY u.username
			ORDER BY total_score DESC
			LIMIT 5
		`, since)
		if err != nil {
			http.Error(w, "Ошибка получения топ игроков", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type TopPlayer struct {
			Username   string `json:"username"`
			TotalScore int    `json:"total_score"`
		}
		var topPlayers []TopPlayer
		for rows.Next() {
			var player TopPlayer
			if err := rows.Scan(&player.Username, &player.TotalScore); err != nil {
				http.Error(w, "Ошибка обработки данных", http.StatusInternalServerError)
				return
			}
			topPlayers = append(topPlayers, player)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(topPlayers)
	})
}
