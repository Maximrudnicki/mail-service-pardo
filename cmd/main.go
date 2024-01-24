package main

import (
	"encoding/json"
	"log"

	conf "mail_service/cmd/config"
	m "mail_service/cmd/mail"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailData struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Name    string   `json:"username"`
}

func main() {
	conf, err := conf.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	qGreeting, err := ch.QueueDeclare(
		"greeting", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgsGreeting, err := ch.Consume(
		qGreeting.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgsGreeting {
			log.Printf("Received a message: %s", d.Body)
			var emailData EmailData
			err := json.Unmarshal(d.Body, &emailData)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v", err)
				continue
			}

			m.SendHTML(
				conf.Username, 
				conf.Password, 
				conf.From, 
				emailData.To, 
				emailData.Name, 
				emailData.Subject, 
				"./templates/greeting.html",
			)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
