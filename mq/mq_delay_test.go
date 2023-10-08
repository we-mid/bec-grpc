package mq_test

import (
	context "context"
	"testing"
	"time"

	mq "bec-grpc/mq"
)

func TestDelay(t *testing.T) {
	key := "bec-grpc.mq.test.delay"
	ch, _ := mq.Conn().Channel()
	ch.QueuePurge(key, false)
	ch.Close()
	// publish
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Publish(ctx, &mq.PublishRequest{Name: key, Body: key, Delay: 5000})
	if err != nil {
		t.Fatalf("c.Publish(...) = %v, %v", r, err)
	}
	want := true
	if want != r.GetOk() {
		t.Fatalf("r.GetOk() = %v, want %v", r.GetOk(), want)
	}
	// consume
	time.Sleep(3 * time.Second)
	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel1()
	r1, err := c.Consume(ctx1, &mq.ConsumeRequest{Name: key})
	if err != nil {
		t.Fatalf("c.Consume(...) = %v, %v", r1, err)
	}
	if r1.GetOk() != false {
		t.Fatalf("r1.GetOk() = %v, want false", r1.GetOk())
	}
	time.Sleep(3 * time.Second)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	r2, err := c.Consume(ctx2, &mq.ConsumeRequest{Name: key})
	if err != nil {
		t.Fatalf("c.Consume(...) = %v, %v", r1, err)
	}
	if r2.GetOk() != true {
		t.Fatalf("r2.GetOk() = %v, want true", r1.GetOk())
	}
}
