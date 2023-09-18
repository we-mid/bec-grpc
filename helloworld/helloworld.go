package helloworld

import (
	context "context"
	"log"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	log.Printf("hello.SayHello name: %v", in.GetName())
	return &HelloReply{Message: "Hello " + in.GetName()}, nil
}
