package helloworld

import (
	context "context"
	"fmt"
	"log"
	"net"
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
	name       = "bec-grpc"
)

func TestOCR(t *testing.T) {
	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterGreeterServer(s, &Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &HelloRequest{Name: name})
	if err != nil {
		t.Fatalf("c.SayHello(...) = %v, %v", r, err)
	}
	want := "Hello " + name
	if want != r.GetMessage() {
		t.Fatalf("r.GetMessage() = %q, want %q", r.GetMessage(), want)
	}
	s.Stop()
}
