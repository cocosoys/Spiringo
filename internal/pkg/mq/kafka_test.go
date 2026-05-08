package mq

import "testing"

// 中文：TestNewKafkaMQRequiresBrokers 验证相关行为符合预期。
// English: TestNewKafkaMQRequiresBrokers verifies the related behavior.
func TestNewKafkaMQRequiresBrokers(t *testing.T) {
	if _, err := NewKafkaMQ(KafkaConfig{}); err == nil {
		t.Fatal("expected missing brokers to fail")
	}
}

// 中文：TestKafkaMQTopicPrefix 验证相关行为符合预期。
// English: TestKafkaMQTopicPrefix verifies the related behavior.
func TestKafkaMQTopicPrefix(t *testing.T) {
	mq, err := NewKafkaMQ(KafkaConfig{
		Brokers:     []string{"127.0.0.1:9092"},
		TopicPrefix: "spiringo.",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer mq.Close()

	if got := mq.topic("payment.created"); got != "spiringo.payment.created" {
		t.Fatalf("topic = %q, want prefixed topic", got)
	}
}
