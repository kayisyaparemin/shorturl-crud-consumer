package consumer

import (
	"log"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	request "urlshortener-crud-consumer/models/requests"
	"urlshortener-crud-consumer/core/extensions"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Consume(services *extensions.Services) {
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan struct{})
	go func() {
		for d := range msgs {
			request :=&request.QueueModel{}
			log.Printf("Received a message: %s", d.Body)
			err := json.Unmarshal(d.Body,request)
			if err != nil {
				log.Panic(err)
			}
			go services.UrlShorthener.StartProcess(request)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
