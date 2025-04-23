package service

import (
	todo "helloapp"
	"helloapp/pkg/repository"
	"testing"
	"time"
)

type mockRepo struct {
	sessions map[string]*todo.RefreshSession
}

func newMockRepo() repository.Authorization {
	return &mockRepo{
		sessions: make(map[string]*todo.RefreshSession),
	}
}

// Implement Authorization interface methods
func (r *mockRepo) SaveRefreshSession(session *todo.RefreshSession) error {
	r.sessions[session.TokenHash] = session
	return nil
}

func (r *mockRepo) GetRefreshSessions(tokenHash string) ([]todo.RefreshSession, error) {
	var sessions []todo.RefreshSession
	for _, s := range r.sessions {
		if !s.IsUsed && time.Now().Before(s.ExpiresAt) {
			sessions = append(sessions, *s)
		}
	}
	return sessions, nil
}

func (r *mockRepo) InvalidateRefreshToken(id string) error {
	for _, s := range r.sessions {
		if s.ID == id {
			s.IsUsed = true
			return nil
		}
	}
	return nil
}

func TestIPChangeDetection(t *testing.T) {
	repo := newMockRepo()
	emailSender := NewMockEmailSender()
	service := NewAuthService(repo, emailSender)

	userUUID := "test-uuid"
	originalIP := "192.168.1.1"

	tokens, err := service.GenerateTokenPair(userUUID, originalIP)
	if err != nil {
		t.Fatalf("не удалось сгенерировать токены: %v", err)
	}

	for _, session := range repo.(*mockRepo).sessions {
		session.UserEmail = "test@example.com"
	}

	newIP := "10.0.0.1"
	_, err = service.RefreshTokens(tokens.RefreshToken, newIP)

	if err == nil {
		t.Error("ожидалась ошибка при смене IP, получено нулевое значение")
	}

	if err.Error() != "обнаружена смена IP-адреса" {
		t.Errorf("ошибка 'обнаружена смена IP-адреса' получена: %v", err)
	}
}

func TestSuccessfulTokenRefresh(t *testing.T) {
	repo := newMockRepo()
	emailSender := NewMockEmailSender()
	service := NewAuthService(repo, emailSender)

	userUUID := "test-uuid"
	clientIP := "192.168.1.1"

	tokens, err := service.GenerateTokenPair(userUUID, clientIP)
	if err != nil {
		t.Fatalf("не удалось сгенерировать токены: %v", err)
	}
	time.Sleep(time.Second)

	newTokens, err := service.RefreshTokens(tokens.RefreshToken, clientIP)
	if err != nil {
		t.Errorf("ожидалось успешное обновление: %v", err)
	}

	if newTokens.AccessToken == tokens.AccessToken {
		t.Error("новый access token должен отличаться от старого")
	}
	if newTokens.RefreshToken == tokens.RefreshToken {
		t.Error("новый refresh token должен отличаться от старого")
	}
}
