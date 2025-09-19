package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
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

type RedisConfig struct {
	client *redis.Client
	ctx    context.Context
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	DB_URL = os.Getenv("DATABASE_URL")
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		client: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		}),
		ctx: context.Background(),
	}
}

func (r *RedisConfig) Set(key, value string, exp time.Duration) error {
	return r.client.Set(r.ctx, key, value, exp).Err()
}

func (r *RedisConfig) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisConfig) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisConfig) Keys(pattern string) ([]string, error) {
	return r.client.Keys(r.ctx, pattern).Result()
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

type ConsumedMessage[T any] struct {
	Body          T
	CorrelationId uuid.UUID
}

func WaitForResponseQueue(ch *amqp091.Channel, queueName string, correlationId uuid.UUID) (*ResponsePythonHandler, error) {
	timeout := time.Second * 30
	consumerTag := fmt.Sprintf("consumer-tag-%s", correlationId)
	stringCorrelationId := correlationId.String()

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

	defer ch.Cancel(consumerTag, false)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case msg, ok := <-msgs:
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

				timer.Stop()

				return &parsedResponse, nil
			} else {
				msg.Nack(false, true)
			}
		case <-timer.C:
			return nil, errors.New("timed out waiting for response")
		}
	}
}
