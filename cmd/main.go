package main

import (
	"WeatherNotification/internal/services"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	service := services.NewNotificationService()

	// Подключение к RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка открытия канала RabbitMQ: %s", err)
	}
	defer ch.Close()

	// Объявление обменника
	err = ch.ExchangeDeclare(
		"notifications_exchange", // имя обменника
		"fanout",                 // тип
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		log.Fatalf("Ошибка объявления обменника: %s", err)
	}

	// Объявление очереди
	q, err := ch.QueueDeclare(
		"",    // имя очереди (пустая строка для автоматического генерации)
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Ошибка объявления очереди: %s", err)
	}

	// Привязка очереди к обменнику
	err = ch.QueueBind(
		q.Name,                   // имя очереди
		"",                       // routing key
		"notifications_exchange", // имя обменника
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Ошибка привязки очереди к обменнику: %s", err)
	}

	// Потребление сообщений
	msgs, err := ch.Consume(
		q.Name, // имя очереди
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Ошибка регистрации потребителя: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			message := string(d.Body)
			err := service.HandleNotification(message)
			if err != nil {
				log.Println("Ошибка обработки уведомления:", err)
			}
		}
	}()

	fmt.Println("Ожидание сообщений. Нажмите CTRL+C для выхода.")
	<-forever
}
