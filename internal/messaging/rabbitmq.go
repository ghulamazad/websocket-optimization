package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

var rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"

// PublishMessage publishes a message to a RabbitMQ priority queue.
func PublishMessage(priority int, message string) error {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"priority_queue",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{"x-max-priority": int32(10)},
	)
	if err != nil {
		log.Println("Error to declare a queue:", err)
		return err
	}

	err = ch.Publish(
		"", q.Name, false, false, amqp.Publishing{
			Priority: uint8(priority),
			Body:     []byte(message),
		},
	)
	if err != nil {
		log.Println("Failed to publish a message:", err)
		return err
	}

	log.Printf("Published message: %s with priority %d", message, priority)
	return nil
}
