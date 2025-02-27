package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // Не возвращаем пароль в JSON
}

type Score struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"` // Добавляем поле для передачи в JSON
	Game      string `json:"game"`
	Score     int    `json:"score"`
	CreatedAt string `json:"created_at"`
}