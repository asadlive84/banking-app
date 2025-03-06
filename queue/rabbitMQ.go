package queue

import (
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"

)

func ConnectToRabbitMQ(rabbitmqURL string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(rabbitmqURL)
		if err == nil {
			log.Println("Successfully connected to RabbitMQ")
			return conn, nil
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after multiple attempts: %v", err)
}