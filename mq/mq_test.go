package mq

import (
	context "context"
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
)

var (
	host       = "localhost"
	port       = 50051
	_, b, _, _ = runtime.Caller(0)
	filePath   = filepath.Join(filepath.Dir(b), "WX20230824-165323@2x.png")

	s    *grpc.Server
	conn *grpc.ClientConn
	c    MQClient
	name = "my_channel"
	body = `{"a":"hello","b":"world"}`
)

// How can I do test setup using the testing package in Go
// https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestPublish(t *testing.T) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Publish(ctx, &PublishRequest{Name: name, Body: body})
	if err != nil {
		t.Fatalf("c.Publish(...) = %v, %v", r, err)
	}
	want := true
	if want != r.GetOk() {
		t.Fatalf("r.GetOk() = %v, want %v", r.GetOk(), want)
	}
}

func TestConsume(t *testing.T) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Consume(ctx, &ConsumeRequest{Name: name})
	if err != nil {
		t.Fatalf("c.Consume(...) = %v, %v", r, err)
	}
	want := body
	if want != r.GetBody() {
		t.Fatalf("r.GetBody() = %v, want %v", r.GetBody(), want)
	}
}

func setup() {
	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s = grpc.NewServer()
	RegisterMQServer(s, &Server{})
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
	c = NewMQClient(conn)
}

func shutdown() {
	if conn != nil {
		conn.Close()
	}
	if s != nil {
		s.Stop()
	}
}
