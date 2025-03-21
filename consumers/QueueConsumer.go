package consumer
import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"urlshortener-crud-consumer/core/extensions"
	request "urlshortener-crud-consumer/models/requests"
	"urlshortener-crud-consumer/models/responses"
	"urlshortener-crud-consumer/utilities"

	amqp "github.com/rabbitmq/amqp091-go"
)
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func Consume(services *extensions.Services) {
	response := new([]responses.Settings)
	utilities.GetJSON("http://configurationapi.test-microservices/api/configurations/rabbitqueuesettings/list", response)

	rabbitSettings := responses.GetRabbitSettings(response)
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s",
	rabbitSettings.UserName,
	rabbitSettings.Password,
	rabbitSettings.Host,
	strings.ToLower(rabbitSettings.VirtualHost),
	))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("shorturl-list.exchange", "direct", true, false, false, false, nil)
	failOnError(err, "Ana Exchange tanımlanamadı")

	err = ch.ExchangeDeclare("shorturl-list.exchange.dlx", "direct", true, false, false, false, nil)
	failOnError(err, "DLX Exchange tanımlanamadı")

	_, err = ch.QueueDeclare(
		"shorturl-list.queue",
		true,
		false,
		false,
		false,
		amqp.Table{"x-dead-letter-exchange": "shorturl-list.exchange.dlx"},
	)
	failOnError(err, "Ana kuyruk oluşturulamadı")

	_, err = ch.QueueDeclare(
		"shorturl-list.dlx",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange": "shorturl-list.exchange",
			"x-message-ttl":          int32(30000),
		},
	)
	failOnError(err, "DLX kuyruk oluşturulamadı")

	err = ch.QueueBind("shorturl-list.queue", "", "shorturl-list.exchange", false, nil)
	failOnError(err, "Ana kuyruk bağlanamadı")

	err = ch.QueueBind("shorturl-list.dlx", "", "shorturl-list.exchange.dlx", false, nil)
	failOnError(err, "DLX kuyruk bağlanamadı")

	fmt.Println("RabbitMQ yapılandırması başarıyla tamamlandı.")

	msgs, err := ch.Consume(
		"shorturl-list.queue", 
		"",
		true, 
		false,
		false,
		false,
		nil,   
	)
	failOnError(err, "Kuyruktan mesaj tüketme hatası")

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			req := &request.QueueModel{}
			log.Printf("Received a message: %s", d.Body)

			err := json.Unmarshal(d.Body, req)
			if err != nil {
				log.Printf("JSON parse hatası: %v", err)
				continue
			}

			go services.UrlShorthener.StartProcess(req)
		}
	}()

	log.Printf(" [*] Mesajlar dinleniyor. Çıkmak için CTRL+C")
	<-forever
}
