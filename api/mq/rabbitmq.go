package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	UNIQUE_REQUESTS_COUNT = "unique_requests_count"
)

type RabbitMQ struct {
	url     string
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQ(user, password, host string) *RabbitMQ {
	return &RabbitMQ{
		url: fmt.Sprintf("amqp://%s:%s@%s/", user, password, host),
	}
}

func (r *RabbitMQ) SetupConnection() error {
	if err := r.Connect(); err != nil {
		return err
	}

	if err := r.OpenChannel(); err != nil {
		return err
	}

	if err := r.QueueDeclare(); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) Connect() error {
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return err
	}
	r.conn = conn
	return nil
}

func (r *RabbitMQ) Close() error {
	return r.conn.Close()
}

func (r *RabbitMQ) OpenChannel() error {
	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}
	r.channel = channel
	return nil
}

func (r *RabbitMQ) QueueDeclare() error {
	queue, err := r.channel.QueueDeclare(
		UNIQUE_REQUESTS_COUNT,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.queue = queue
	return nil
}

func (r *RabbitMQ) Write(msg []byte) (int, error) {
	err := r.channel.Publish(
		"",
		UNIQUE_REQUESTS_COUNT,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)
	return 0, err
}
