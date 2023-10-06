// https://www.rabbitmq.com/tutorials/tutorial-one-go.html
// https://github.com/rabbitmq/rabbitmq-tutorials/blob/main/go/send.go
// https://github.com/rabbitmq/rabbitmq-tutorials/blob/main/go/receive.go
package mq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	UnimplementedMQServer
}

const (
	url = "amqp://guest:guest@localhost:5672/"
)

func (s *Server) Publish(ctx context.Context, in *PublishRequest) (*PublishReply, error) {
	q, ch, conn, err := getQueue(in.GetName())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := in.GetBody()
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return nil, err
	}
	log.Printf(" [x] Sent a message: %s\n", body)
	return &PublishReply{Ok: true}, nil
}

func (s *Server) Consume(ctx context.Context, in *ConsumeRequest) (*ConsumeReply, error) {
	q, ch, conn, err := getQueue(in.GetName())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer ch.Close()

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
		return nil, err
	}

	d := <-msgs
	body := string(d.Body[:]) // bytes to string
	log.Printf(" [x] Received a message: %s\n", body)
	return &ConsumeReply{Body: body}, nil
	// var forever chan struct{}
	// go func() {
	// 	for d := range msgs {
	// 		log.Printf("Received a message: %s", d.Body)
	// 	}
	// }()
	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// <-forever
}

func getQueue(name string) (*amqp.Queue, *amqp.Channel, *amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}

	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, nil, err
	}
	return &q, ch, conn, nil
}
