package rabbit

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/smtp"

	"github.com/streadway/amqp"
)

// Email details structure
type CartOpen struct {
	Email string
	Name  string
}

// Listener function to consume messages from RabbitMQ
func ListenForPriceUpdates() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"price_updates", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			var user CartOpen
			err := json.Unmarshal(d.Body, &user)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Send price change email
			if err := SendPriceChangeEmail(user); err != nil {
				log.Printf("Failed to send email to %s: %v", user.Email, err)
			} else {
				log.Printf("Email sent to %s\n", user.Email)
			}
		}
	}()

	// Keep the listener running
	select {}
}

// SendPriceChangeEmail sends an email about price changes
func SendPriceChangeEmail(user CartOpen) error {
	// Parse the price change template
	tmpl, err := template.ParseFiles("templates/price_update.tmpl")
	if err != nil {
		return err
	}

	// Create a buffer to hold the rendered template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, user); err != nil {
		return err
	}

	from := "joacomateocortez@gmail.com"
	password := "nxdw bukh weno nutr"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Email details
	subject := "Precio Actualizado"
	message := []byte("Subject: " + subject + "\r\n\r\n" + body.String())
	to := []string{user.Email}

	// Send the email
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
}
