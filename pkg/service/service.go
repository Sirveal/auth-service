package service

import (
	"helloapp/pkg/repository"
)

type Authorization interface {
	GenerateTokenPair(userUUID, clientIP string) (*TokenPair, error)
	RefreshTokens(refreshToken, clientIP string) (*TokenPair, error)
}

type EmailSender interface {
	Send(to, subject, body string) error
	SendIPChangeWarning(userEmail string, oldIP string, newIP string) error
}

type Service struct {
	Authorization
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewService(repos *repository.Repository, emailSender EmailSender) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, emailSender),
	}
}
