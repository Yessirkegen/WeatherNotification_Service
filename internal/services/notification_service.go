package services

import (
	"WeatherNotification/internal/patterns"
	"fmt"
	"strings"
)

type NotificationService struct {
	factory *patterns.NotificationFactory
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		factory: &patterns.NotificationFactory{},
	}
}

func (s *NotificationService) HandleNotification(message string) error {
	// Предполагается, что сообщение имеет формат "notificationType:message"
	parts := strings.SplitN(message, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("неверный формат сообщения")
	}
	notificationType := parts[0]
	notificationMessage := parts[1]

	notifier, err := s.factory.CreateNotification(notificationType)
	if err != nil {
		return err
	}
	return notifier.Send(notificationMessage)
}
