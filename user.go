package todo

import (
	"time"
)

type RefreshSession struct {
	ID        string    `db:"id"`
	UserUUID  string    `db:"user_uuid"`
	TokenHash string    `db:"token_hash"`
	ClientIP  string    `db:"client_ip"`
	TokenID   string    `db:"token_id"`
	ExpiresAt time.Time `db:"expires_at"`
	IsUsed    bool      `db:"is_used"`
	UserEmail string    `db:"user_email"`
}
