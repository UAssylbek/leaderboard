package handlers

import (
	"net/http"
	"github.com/redis/go-redis/v9"
)

// GetLeaderboard - заглушка для получения глобального лидерборда
func GetLeaderboard(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Лидерборд пока не реализован"))
	}
}

// GetRank - заглушка для получения ранга пользователя
func GetRank(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ранг пока не реализован"))
	}
}