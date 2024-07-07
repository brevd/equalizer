package auth

import (
	"time"

	"github.com/brevd/equalizer/internal"
)

type BlackListToken struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

func AddToBlacklist(token string) error {
	_, err := internal.DB.Exec("INSERT INTO blacklist (token) VALUES (?)", token)
	return err
}

func IsBlacklisted(token string) (bool, error) {
	var count int
	cutoff := time.Now().Add(-(TokenExpirationHours * time.Hour))
	err := internal.DB.QueryRow("SELECT COUNT(*) FROM blacklist WHERE token = ? AND created_at >= ?", token, cutoff).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
