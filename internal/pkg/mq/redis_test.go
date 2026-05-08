package mq

import (
	"context"
	"testing"
)

// 中文：TestRedisStreamMQRejectsInvalidPublish 验证相关行为符合预期。
// English: TestRedisStreamMQRejectsInvalidPublish verifies the related behavior.
func TestRedisStreamMQRejectsInvalidPublish(t *testing.T) {
	mq := NewRedisStreamMQ(nil, "")

	if err := mq.Publish(context.Background(), nil); err == nil {
		t.Fatal("expected nil message to fail")
	}
	if err := mq.Publish(context.Background(), &Message{}); err == nil {
		t.Fatal("expected empty topic to fail")
	}
	if err := mq.Publish(context.Background(), &Message{Topic: "events"}); err == nil {
		t.Fatal("expected nil redis client to fail")
	}
}

// 中文：TestRedisStreamMQRejectsInvalidSubscribe 验证相关行为符合预期。
// English: TestRedisStreamMQRejectsInvalidSubscribe verifies the related behavior.
func TestRedisStreamMQRejectsInvalidSubscribe(t *testing.T) {
	mq := NewRedisStreamMQ(nil, "")

	if err := mq.Subscribe(context.Background(), "", func(*Message) error { return nil }); err == nil {
		t.Fatal("expected empty topic to fail")
	}
	if err := mq.Subscribe(context.Background(), "events", nil); err == nil {
		t.Fatal("expected nil handler to fail")
	}
	if err := mq.Subscribe(context.Background(), "events", func(*Message) error { return nil }); err == nil {
		t.Fatal("expected nil redis client to fail")
	}
}
