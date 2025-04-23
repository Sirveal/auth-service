package repository

import (
	"errors"
	"fmt"
	todo "helloapp"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) SaveRefreshSession(session *todo.RefreshSession) error {
	query := `
        INSERT INTO refresh_tokens 
        (user_uuid, token_hash, client_ip, token_id, expires_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	return r.db.QueryRow(
		query,
		session.UserUUID,
		session.TokenHash,
		session.ClientIP,
		session.TokenID,
		session.ExpiresAt,
	).Scan(&session.ID)
}

func (r *AuthPostgres) GetRefreshSessions(tokenHash string) ([]todo.RefreshSession, error) {
	query := `
        SELECT rt.id, rt.user_uuid, rt.token_hash, rt.client_ip, rt.token_id, 
               rt.expires_at, rt.is_used, u.email as user_email
        FROM refresh_tokens rt
        JOIN users u ON u.uuid = rt.user_uuid
        WHERE rt.is_used = false 
        AND rt.expires_at > NOW()`
	var sessions []todo.RefreshSession
	if err := r.db.Select(&sessions, query); err != nil {
		return nil, fmt.Errorf("ошибка получения refresh сессий: %w", err)
	}

	return sessions, nil
}

func (r *AuthPostgres) InvalidateRefreshToken(id string) error {
	query := `UPDATE refresh_tokens SET is_used = true WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("не удалось сделать токен недействительным: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("refresh token не найден")
	}

	return nil
}
