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
	r1, _ := c.Consume(ctx1, &mq.ConsumeRequest{Name: key})
	if r1 != nil {
		t.Fatalf("r1 = %v, want nil", r1)
	}
	time.Sleep(3 * time.Second)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	r2, _ := c.Consume(ctx2, &mq.ConsumeRequest{Name: key})
	if r2 == nil {
		t.Fatalf("r2 = %v, want not nil", r2)
	}
}
