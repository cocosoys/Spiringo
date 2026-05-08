package mq

import "testing"

// 中文：TestNewRabbitMQRequiresURL 验证相关行为符合预期。
// English: TestNewRabbitMQRequiresURL verifies the related behavior.
func TestNewRabbitMQRequiresURL(t *testing.T) {
	if _, err := NewRabbitMQ(RabbitMQConfig{}); err == nil {
		t.Fatal("expected missing url to fail")
	}
}

// 中文：TestAMQPHeadersConvertsValues 验证相关行为符合预期。
// English: TestAMQPHeadersConvertsValues verifies the related behavior.
func TestAMQPHeadersConvertsValues(t *testing.T) {
	headers := amqpHeaders(map[string]interface{}{
		"tenant": "acme",
		"retry":  2,
	})

	if headers["tenant"] != "acme" || headers["retry"] != "2" {
		t.Fatalf("unexpected headers: %#v", headers)
	}
}
