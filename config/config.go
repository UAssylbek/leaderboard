package config

import "os"

// Config хранит конфигурационные параметры приложения
type Config struct {
	PostgresURL string
	RedisAddr   string
	JWTSecret   string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() Config {
	return Config{
		PostgresURL: getEnv("POSTGRES_URL", "postgres://user:password@localhost:5432/leaderboard?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "my_secret_key"),
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}