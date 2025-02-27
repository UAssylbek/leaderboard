package db

import (
	"database/sql"
	"log"
)

func SetupPostgres(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS scores (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			game VARCHAR(50) NOT NULL,
			score INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}
	log.Println("Таблицы созданы успешно")
	return nil
}