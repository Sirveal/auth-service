package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	todo "helloapp"
	"helloapp/pkg/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	signingKey      = "qrkjk#4#%35FSFJlja#4353KSFjH"
	accessTokenTTL  = 12 * time.Minute
	refreshTokenTTL = 24 * time.Hour * 7
)

type AuthService struct {
	repo        repository.Authorization
	emailSender EmailSender
}

type tokenClaims struct {
	jwt.StandardClaims
	UserUUID string `json:"uuid"`
	ClientIP string `json:"ip"`
	TokenID  string `json:"jti"`
}

func NewAuthService(repo repository.Authorization, emailSender EmailSender) *AuthService {
	return &AuthService{
		repo:        repo,
		emailSender: emailSender,
	}
}

func (s *AuthService) GenerateTokenPair(userUUID, clientIP string) (*TokenPair, error) {
	tokenID := generateTokenID()

	accessToken, err := s.createAccessToken(userUUID, clientIP, tokenID)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания access token: %w", err)
	}

	refreshToken, refreshHash, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации refresh token: %w", err)
	}

	session := todo.RefreshSession{
		UserUUID:  userUUID,
		TokenHash: refreshHash,
		ClientIP:  clientIP,
		TokenID:   tokenID,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}

	if err := s.repo.SaveRefreshSession(&session); err != nil {
		return nil, fmt.Errorf("ошибка сохранения refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshTokens(encodedToken, clientIP string) (*TokenPair, error) {
	sessions, err := s.repo.GetRefreshSessions(encodedToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения сессий: %w", err)
	}

	var validSession *todo.RefreshSession
	for _, session := range sessions {
		if err := bcrypt.CompareHashAndPassword([]byte(session.TokenHash), []byte(encodedToken)); err == nil {
			validSession = &session
			break
		}
	}

	if validSession == nil {
		logrus.Debugf("не найдена сессия для токена: %s", encodedToken[:10])
		return nil, errors.New("invalid refresh token")
	}

	if validSession.ClientIP != clientIP {
		if err := s.emailSender.SendIPChangeWarning(validSession.UserEmail, validSession.ClientIP, clientIP); err != nil {
			logrus.Errorf("ошибка отправки уведомления о смене IP: %v", err)
		}
		return nil, errors.New("обнаружена смена IP-адреса")
	}

	if err := s.repo.InvalidateRefreshToken(validSession.ID); err != nil {
		return nil, fmt.Errorf("couldn't invalidate token: %w", err)
	}

	return s.GenerateTokenPair(validSession.UserUUID, clientIP)
}

func (s *AuthService) generateRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}

	refreshToken := base64.URLEncoding.EncodeToString(b)

	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return refreshToken, string(hash), nil
}

func (s *AuthService) createAccessToken(userUUID, clientIP, tokenID string) (string, error) {
	claims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        tokenID,
		},
		UserUUID: userUUID,
		ClientIP: clientIP,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(signingKey))
}

func generateTokenID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
