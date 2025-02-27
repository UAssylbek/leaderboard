package handlers

import (
	"context"
	"encoding/json"
	"net/http"
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