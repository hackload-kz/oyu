package entity

import "time"

type AuthUser struct {
	ID            int64
	Email         string
	FirstName     string
	Surname       string
	IsActive      bool
	LastLoggedIn  time.Time
	PasswordHash  string  // hex (SHA-256) из password_hash
	PasswordPlain *string // если в БД хранится открытый пароль
}
