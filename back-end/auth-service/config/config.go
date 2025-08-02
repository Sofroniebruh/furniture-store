package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

type ResponsePythonHandler struct {
	StatusCode int    `json:"status_code"`
	Data       string `json:"data"`
}

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

func DeclareQueue(
	ch *amqp091.Channel,
	queueName string,
	durable bool,
	autoDelete bool,
	exclusive bool,
	noWait bool) (amqp091.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		noWait,
		nil,
	)

	if err != nil {
		return amqp091.Queue{}, errors.New("failed to declare a queue: " + err.Error())
	}

	return q, nil
}

func ProduceMessage[T any](requestQueue amqp091.Queue, replyQueue amqp091.Queue, ch *amqp091.Channel, body T, correlationId uuid.UUID) error {
	data, err := json.Marshal(body)
	stringCorrelationId := correlationId.String()

	if err != nil {
		return errors.New("failed to marshal body: " + err.Error())
	}

	log.Println("ID", stringCorrelationId)
	log.Println("Reply", replyQueue.Name)

	err = ch.Publish(
		"",
		requestQueue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType:   "application/json",
			Body:          data,
			CorrelationId: stringCorrelationId,
			ReplyTo:       replyQueue.Name,
		})

	if err != nil {
		return errors.New("failed to publish a message: " + err.Error())
	}

	return nil
}

func WaitForResponseQueue(ch *amqp091.Channel, queueName string, correlationId uuid.UUID) (*ResponsePythonHandler, error) {
	timeout := time.Second * 30
	consumerTag := fmt.Sprintf("consumer-tag-%s", correlationId)
	stringCorrelationId := correlationId.String()
	println("ID: ", stringCorrelationId)

	msgs, err := ch.Consume(
		queueName,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to consume from response queue: %w", err)
	}

	log.Println("goida")

	defer ch.Cancel(consumerTag, false)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case msg, ok := <-msgs:
			log.Println("Received a message: ", string(msg.Body), msg.CorrelationId)
			log.Println(msg.CorrelationId == stringCorrelationId)
			if !ok {
				return nil, errors.New("failed to receive a message from response queue")
			}
			if msg.CorrelationId == stringCorrelationId {
				msg.Ack(false)
				var parsedResponse ResponsePythonHandler
				err = json.Unmarshal(msg.Body, &parsedResponse)

				if err != nil {
					return nil, errors.New("failed to unmarshal response: " + err.Error())
				}

				return &parsedResponse, nil
			} else {
				msg.Nack(false, true)
			}
		case <-timer.C:
			return nil, errors.New("timed out waiting for response")
		}
	}
}
