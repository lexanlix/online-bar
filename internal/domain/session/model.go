package session

import "time"

type Session struct {
	RefreshToken string `json:"refreshToken"`

	// Дата исчетения refresh токена
	ExpiresAt time.Time `json:"expiresAt"`
}
