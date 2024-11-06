package patterns

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type Notification interface {
	Send(message string) error
}

type EmailNotification struct{}

func (n *EmailNotification) Send(message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "your_email@example.com")
	m.SetHeader("To", "recipient@example.com") // Здесь вы должны указать email получателя
	m.SetHeader("Subject", "Прогноз погоды")
	m.SetBody("text/plain", message)

	d := gomail.NewDialer("smtp.example.com", 587, "your_email@example.com", "your_password")

	err := d.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("ошибка отправки email: %v", err)
	}
	return nil
}

type SMSNotification struct{}

func (n *SMSNotification) Send(message string) error {
	// Реализация отправки SMS
	fmt.Println("Отправка SMS:", message)
	return nil
}

type NotificationFactory struct{}

func (f *NotificationFactory) CreateNotification(notificationType string) (Notification, error) {
	switch notificationType {
	case "email":
		return &EmailNotification{}, nil
	case "sms":
		return &SMSNotification{}, nil
	default:
		return nil, fmt.Errorf("неизвестный тип уведомления: %s", notificationType)
	}
}
