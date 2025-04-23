package service

import (
	"log"

	"github.com/sirupsen/logrus"
)

type MockEmailSender struct{}

func NewMockEmailSender() EmailSender {
	return &MockEmailSender{}
}

func (s *MockEmailSender) Send(to, subject, body string) error {
	log.Printf("[MOCK EMAIL] To: %s\nSubject: %s\nBody: %s\n", to, subject, body)
	return nil
}

func (s *MockEmailSender) SendIPChangeWarning(email, oldIP, newIP string) error {
	logrus.WithFields(logrus.Fields{
		"email": email,
		"oldIP": oldIP,
		"newIP": newIP,
	}).Warning("Будет отправлено электронное письмо с предупреждением об изменении IP-адреса.")

	return nil
}
