package service

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"urlshortener-crud-consumer/models/queues"
	"urlshortener-crud-consumer/utilities"
	amqp "github.com/rabbitmq/amqp091-go"
)
type IQueueService interface{
	Send(model queues.MailQueueModel)
}
type QueueService struct {
}
func NewQueueService() *QueueService {
	return &QueueService{}
}

func (qs *QueueService) Send(model queues.MailQueueModel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body,_ := json.Marshal(&model)

	err = ch.PublishWithContext(ctx,
		utilities.CommonMailExhangeName, // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}