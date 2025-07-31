package config

import (
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
)

var (
	DB_URL            string
	JWT_SECRET        []byte
	ACCESS_TOKEN_TTL  = time.Minute * 15
	REFRESH_TOKEN_TTL = time.Hour * 24 * 7
)

type userId string

const UserIdKey userId = "userId"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	DB_URL = os.Getenv("DATABASE_URL")
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
}

func InitRabbitMq() (*amqp091.Connection, *amqp091.Channel, error) {
	RabbitUrl := os.Getenv("RABBIT_URL")
	conn, err := amqp091.Dial(RabbitUrl)

	if err != nil {
		return nil, nil, errors.New("failed to initialize RabbitMQ: " + err.Error())
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, nil, errors.New("failed to open channel: " + err.Error())
	}

	return conn, ch, nil
}

func ProduceMessage[T any](ch *amqp091.Channel, queueName string, body T) error {
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return errors.New("failed to declare a queue: " + err.Error())
	}

	data, err := json.Marshal(body)

	if err != nil {
		return errors.New("failed to marshal body: " + err.Error())
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        data,
		})

	if err != nil {
		return errors.New("failed to publish a message: " + err.Error())
	}

	return nil
}
