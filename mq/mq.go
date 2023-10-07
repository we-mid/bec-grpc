// https://www.rabbitmq.com/tutorials/tutorial-one-go.html
// https://github.com/rabbitmq/rabbitmq-tutorials/blob/main/go/send.go
// https://github.com/rabbitmq/rabbitmq-tutorials/blob/main/go/receive.go
package mq

import (
	"context"
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	UnimplementedMQServer
}

const (
	url               = "amqp://guest:guest@localhost:5672/"
	exchangeName      = ""
	delayExchangeName = "exchange.delay.1"
)

var (
	conn *amqp.Connection
)

func FailedOnError(err error, message string) {
	if err != nil {
		log.Fatal(message, "Error: ", err.Error())
	}
}

func init() {
	var err error
	conn, err = amqp.Dial(url)
	FailedOnError(err, "failed to create connection to RabbitMQ server")

	ch, err := conn.Channel()
	FailedOnError(err, "Failed to open a channel")

	// https://github.com/search?q=%22PublishWithContext%22+%22x-delay%22+language%3AGo&ref=opensearch&type=code
	// https://github.com/a2htray/gorabbitmq/blob/1ed88cb3460a55268ccad02b06a67c809457c6c5/Delayed%20message/main.go#L49
	// https://github.com/jassue/gin-wire/blob/0025a433e5047c4a403553142dc5be1e1c7f52fc/app/compo/mq/rabbitmq/rabbitmq.go#L134
	err = ch.ExchangeDeclare(
		delayExchangeName,   // Exchange name
		"x-delayed-message", // Exchange type
		true,                // Durable
		false,               // Auto-deleted
		false,               // Internal
		false,               // No-wait
		amqp.Table{
			"x-delayed-type": "direct", // Set delayed exchange type as direct
		},
	)
	FailedOnError(err, "Failed to declare the delayed exchange")
}

func (s *Server) Publish(ctx context.Context, in *PublishRequest) (r *PublishReply, err error) {
	r = &PublishReply{Ok: false}
	var ch *amqp.Channel
	if ch, err = conn.Channel(); err != nil {
		return
	}
	name := in.GetName()
	if err = queueDeclare(ch, name); err != nil {
		return
	}
	body := in.GetBody()
	publishing := amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(body),
		DeliveryMode: amqp.Persistent, // 消息持久化
	}
	delay := in.GetDelay()
	xg := exchangeName
	if delay > 0 {
		publishing.Headers = amqp.Table{
			"x-delay": delay,
		}
		xg = delayExchangeName
		if err = ch.QueueBind(name, name, xg, false, nil); err != nil {
			return
		}
	}
	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = ch.PublishWithContext(
		ctx1,
		xg,    // exchange
		name,  // routing key
		false, // mandatory
		false, // immediate
		publishing,
	); err != nil {
		return
	}
	log.Printf(" [x] Sent a message: (delay=%v) %s\n", delay, body)
	r = &PublishReply{Ok: true}
	return
}

func (s *Server) Consume(ctx context.Context, in *ConsumeRequest) (*ConsumeReply, error) {
	var ch *amqp.Channel
	var err error
	if ch, err = conn.Channel(); err != nil {
		return nil, err
	}
	name := in.GetName()
	if err = queueDeclare(ch, name); err != nil {
		return nil, err
	}
	msgs, err := ch.Consume(
		name,  // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}
	select {
	case d := <-msgs:
		ch.Close()
		body := string(d.Body[:]) // bytes to string
		log.Printf(" [x] Received a message: %s\n", body)
		return &ConsumeReply{Body: body}, nil
	case <-time.After(time.Second):
		ch.Close()
		return nil, errors.New("Consume timeout")
	}
	// var forever chan struct{}
	// go func() {
	// 	for d := range msgs {
	// 		log.Printf("Received a message: %s", d.Body)
	// 	}
	// }()
	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// <-forever
}

func queueDeclare(ch *amqp.Channel, name string) error {
	_, err := ch.QueueDeclare(
		name, // name
		// false, // durable
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func Conn() *amqp.Connection {
	return conn
}
