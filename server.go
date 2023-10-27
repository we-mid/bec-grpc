/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"

	// hello "google.golang.org/grpc/examples/helloworld/helloworld"
	hello "bec-grpc/helloworld"
	mq "bec-grpc/mq"
	ocr "bec-grpc/ocr"
)

var (
	// to specify host: 127.0.0.1 is safer in production deployment
	host = flag.String("host", "127.0.0.1", "The server host")
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	hello.RegisterGreeterServer(s, &hello.Server{})
	ocr.RegisterOCRServer(s, &ocr.Server{})
	mq.RegisterMQServer(s, &mq.Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
