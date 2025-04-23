package repository

import (
	todo "helloapp"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	SaveRefreshSession(session *todo.RefreshSession) error
	InvalidateRefreshToken(id string) error
	GetRefreshSessions(tokenHash string) ([]todo.RefreshSession, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
