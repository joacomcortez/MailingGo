package consume

import (
	"encoding/json"
	"fmt"
	"log"
	"mailinggo/rabbit/rschema"
	"mailinggo/rest"

	"github.com/streadway/amqp"
)

func ConsumePriceChangeNotifications() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare exchange
	err = ch.ExchangeDeclare(
		"price_updates", // name
		"direct",        // type
		false,           // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %s", err)
	}

	// Declare and bind the queue
	queue, err := ch.QueueDeclare(
		"price_updates_joaco", // queue name
		false,                 // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	err = ch.QueueBind(
		queue.Name,      // queue name
		"price_updates", // routing key
		"price_updates", // exchange name
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %s", err)
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer tag
		true,       // auto-acknowledge
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to start consuming: %s", err)
	}


	// Process messages in a goroutine
	go func() {
		log.Printf("inside the func")
		for msg := range msgs {
			log.Println("Message received")
			fmt.Printf("\n Raw message body: %s\n", string(msg.Body))
			var notification rschema.PriceChangeNotification
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				continue
			}
			fmt.Printf("Received price change notification: ArticleId: %s, New Price: %.2f\n", notification.ArticleId, notification.Price)
			if err := rest.PostPriceUpdate(notification); err != nil {
				log.Printf("Failed to handle price change notification: %s", err)
			}
		}
		
	}()

	

	fmt.Println("Waiting for messages. To exit press CTRL+C")
	select {} // Block the function to keep it running
}
