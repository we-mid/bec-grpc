package mq_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mq "bec-grpc/mq"
)

var (
	host       = "localhost"
	port       = 50051
	_, b, _, _ = runtime.Caller(0)
	filePath   = filepath.Join(filepath.Dir(b), "WX20230824-165323@2x.png")

	s    *grpc.Server
	conn *grpc.ClientConn
	c    mq.MQClient
)

// How can I do test setup using the testing package in Go
// https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestBasic(t *testing.T) {
	key := "bec-grpc.mq.test.basic"
	ch, _ := mq.Conn().Channel()
	ch.QueuePurge(key, false)
	ch.Close()
	// publish
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Publish(ctx, &mq.PublishRequest{Name: key, Body: key})
	if err != nil {
		t.Fatalf("c.Publish(...) = %v, %v", r, err)
	}
	want := true
	if want != r.GetOk() {
		t.Fatalf("r.GetOk() = %v, want %v", r.GetOk(), want)
	}
	// consume
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
	defer cancel1()
	r1, err := c.Consume(ctx1, &mq.ConsumeRequest{Name: key})
	if err != nil {
		t.Fatalf("c.Consume(...) = %v, %v", r1, err)
	}
	want1 := key
	if want1 != r1.GetBody() {
		t.Fatalf("r1.GetBody() = %v, want %v", r1.GetBody(), want1)
	}
}

func setup() {
	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s = grpc.NewServer()
	mq.RegisterMQServer(s, &mq.Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c = mq.NewMQClient(conn)
}

func shutdown() {
	if conn != nil {
		conn.Close()
	}
	if s != nil {
		s.Stop()
	}
	// var forever chan (int)
	// <-forever
}
