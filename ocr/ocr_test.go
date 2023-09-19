package ocr

import (
	context "context"
	"fmt"
	"log"
	"net"
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
)

func TestOCR(t *testing.T) {
	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterOCRServer(s, &Server{})
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
	c := NewOCRClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Read(ctx, &ReadRequest{FilePath: filePath})
	if err != nil {
		t.Fatalf("c.Read(...) = %v, %v", r, err)
	}
	want := "82*25*44\n82*25*44\n82*25*44\n71*36*25\n45*29*41"
	if want != r.GetText() {
		t.Fatalf("r.GetText() = %q, want %q", r.GetText(), want)
	}
	s.Stop()
}
