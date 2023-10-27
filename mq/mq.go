// https://www.rabbitmq.com/tutorials/tutorial-one-go.html
// https://github.com/rabbitmq/rabbitmq-tutorials/blob/main/go/send.go
// https://github.com/rabbitmq/rabbitmq-tutorials/blob/main/go/receive.go
package mq

import (
	"context"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	UnimplementedMQServer
}

const (
	exchangeName      = ""
	delayExchangeName = "exchange.delay.1"
)

var (
	url  = ""
	conn *amqp.Connection
)

func FailedOnError(err error, message string) {
	if err != nil {
		log.Fatal(message, "Error: ", err.Error())
	}
}

func init() {
	var err error
	url = os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/" // default
	}
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
	name := in.GetName()
	body := in.GetBody()
	delay := in.GetDelay()
	// fmt.Printf("Publish: name=%q, body=%q, delay=%v\n", name, body, delay)

	var ch *amqp.Channel
	if ch, err = conn.Channel(); err != nil {
		return
	}
	if err = queueDeclare(ch, name); err != nil {
		return
	}
	publishing := amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(body),
		DeliveryMode: amqp.Persistent, // 消息持久化
	}
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

func (s *Server) Consume(ctx context.Context, in *ConsumeRequest) (r *ConsumeReply, err error) {
	r = &ConsumeReply{Ok: false, Body: "", Tag: 0}
	name := in.GetName()
	autoAck := in.GetAutoAck()
	// fmt.Printf("Consume: name=%q, autoAck=%v\n", name, autoAck)

	var ch *amqp.Channel
	if ch, err = conn.Channel(); err != nil {
		return
	}
	if err = queueDeclare(ch, name); err != nil {
		return
	}
	msgs, err := ch.Consume(
		name,    // queue
		"",      // consumer
		autoAck, // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return
	}
	select {
	case d := <-msgs:
		ch.Close()
		body := string(d.Body[:]) // bytes to string
		log.Printf(" [x] Received a message: %s\n", body)
		r = &ConsumeReply{Ok: true, Body: body, Tag: d.DeliveryTag}
		return
	// case <-time.After(time.Second):
	case <-time.After(500 * time.Millisecond):
		ch.Close()
		// return nil, errors.New("Consume timeout")
		// won't throw error but return Ok:false
		return
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

func (s *Server) Ack(ctx context.Context, in *AckRequest) (r *AckReply, err error) {
	r = &AckReply{Ok: false}
	// typ, tag := in.GetType(), in.GetTag()
	// todo
	return
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
