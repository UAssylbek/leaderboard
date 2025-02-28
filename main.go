package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/UAssylbek/leaderboard/config"
	"github.com/UAssylbek/leaderboard/db"
	"github.com/UAssylbek/leaderboard/handlers"
)

func main() {
	cfg := config.LoadConfig()

	log.Println("Ожидание запуска баз данных...")
	time.Sleep(5 * time.Second)

	pgDB, err := db.InitPostgres(cfg.PostgresURL)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}
	defer pgDB.Close()

	redisClient, err := db.InitRedis(cfg.RedisAddr, "", 0)
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}
	defer redisClient.Close()

	if err := db.SetupPostgres(pgDB); err != nil {
		log.Fatal("Ошибка настройки PostgreSQL:", err)
	}

	handlers.SetJWTSecret([]byte(cfg.JWTSecret))

	// Подписчик на канал leaderboard_updates
	go func() {
		ctx := context.Background()
		pubsub := redisClient.Subscribe(ctx, "leaderboard_updates")
		defer pubsub.Close()

		log.Println("Подписчик запущен, слушает канал leaderboard_updates...")
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Println("Ошибка получения сообщения из Redis:", err)
				return
			}
			log.Printf("Обновление лидерборда: %s", msg.Payload)
		}
	}()

	// Регистрация обработчиков с JWT
	http.HandleFunc("/register", handlers.Register(pgDB))
	http.HandleFunc("/login", handlers.Login(pgDB))
	http.HandleFunc("/submit-score", handlers.SubmitScore(pgDB, redisClient))
	http.HandleFunc("/leaderboard", handlers.JWTMiddleware(handlers.GetLeaderboard(redisClient))) // Добавляем JWT
	http.HandleFunc("/rank", handlers.JWTMiddleware(handlers.GetRank(redisClient)))             // Добавляем JWT
	http.HandleFunc("/top-players", handlers.GetTopPlayers(pgDB))

	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
