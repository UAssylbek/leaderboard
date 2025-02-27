package main

import (
	"log"
	"net/http"
	"time"

	"github.com/UAssylbek/leaderboard/config"
	"github.com/UAssylbek/leaderboard/db"
	"github.com/UAssylbek/leaderboard/handlers"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	log.Println("Ожидание запуска баз данных...")
	time.Sleep(5 * time.Second)

	// Инициализация PostgreSQL
	pgDB, err := db.InitPostgres(cfg.PostgresURL)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}
	defer pgDB.Close()

	// Инициализация Redis
	redisClient, err := db.InitRedis(cfg.RedisAddr, "", 0)
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}
	defer redisClient.Close()

	// Настройка таблиц
	if err := db.SetupPostgres(pgDB); err != nil {
		log.Fatal("Ошибка настройки PostgreSQL:", err)
	}

	// Передаем JWTSecret в handlers (пока глобально через пакет handlers)
	handlers.SetJWTSecret([]byte(cfg.JWTSecret))

	// Регистрация обработчиков
	http.HandleFunc("/register", handlers.Register(pgDB))
	http.HandleFunc("/login", handlers.Login(pgDB))
	http.HandleFunc("/submit-score", handlers.SubmitScore(pgDB, redisClient))
	http.HandleFunc("/leaderboard", handlers.GetLeaderboard(redisClient))
	http.HandleFunc("/rank", handlers.GetRank(redisClient))
	http.HandleFunc("/top-players", handlers.GetTopPlayers(pgDB))

	// Запуск сервера
	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
