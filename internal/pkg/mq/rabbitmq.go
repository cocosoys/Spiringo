package mq

import (
	"context"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 中文：RabbitMQConfig 定义当前包使用的数据结构或接口。
// English: RabbitMQConfig defines a data structure or interface used by this package.
type RabbitMQConfig struct {
	// 中文：URL 保存当前结构中的配置或数据值。
	// English: URL stores a configuration or data value for this struct.
	URL string
	// 中文：Exchange 保存当前结构中的配置或数据值。
	// English: Exchange stores a configuration or data value for this struct.
	Exchange string
	// 中文：QueuePrefix 保存当前结构中的配置或数据值。
	// English: QueuePrefix stores a configuration or data value for this struct.
	QueuePrefix string
}

// 中文：RabbitMQ 定义当前包使用的数据结构或接口。
// English: RabbitMQ defines a data structure or interface used by this package.
type RabbitMQ struct {
	// 中文：conn 保存当前结构中的配置或数据值。
	// English: conn stores a configuration or data value for this struct.
	conn *amqp.Connection
	// 中文：ch 保存当前结构中的配置或数据值。
	// English: ch stores a configuration or data value for this struct.
	ch *amqp.Channel
	// 中文：exchange 保存当前结构中的配置或数据值。
	// English: exchange stores a configuration or data value for this struct.
	exchange string
	// 中文：queuePrefix 保存当前结构中的配置或数据值。
	// English: queuePrefix stores a configuration or data value for this struct.
	queuePrefix string
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
	// 中文：closed 保存当前结构中的配置或数据值。
	// English: closed stores a configuration or data value for this struct.
	closed bool
}

// 中文：NewRabbitMQ 创建并返回对应组件实例。
// English: NewRabbitMQ creates and returns the corresponding component instance.
func NewRabbitMQ(cfg RabbitMQConfig) (*RabbitMQ, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("rabbitmq url is required")
	}
	if cfg.Exchange == "" {
		cfg.Exchange = "spiringo.events"
	}
	if cfg.QueuePrefix == "" {
		cfg.QueuePrefix = "spiringo"
	}

	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}
	if err := ch.ExchangeDeclare(cfg.Exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq declare exchange: %w", err)
	}

	return &RabbitMQ{conn: conn, ch: ch, exchange: cfg.Exchange, queuePrefix: cfg.QueuePrefix}, nil
}

// 中文：Publish 执行当前包中的对应流程。
// English: Publish executes the corresponding workflow in this package.
func (m *RabbitMQ) Publish(ctx context.Context, msg *Message) error {
	if msg == nil {
		return fmt.Errorf("rabbitmq message is required")
	}
	headers := amqp.Table{}
	for key, value := range msg.Headers {
		headers[key] = value
	}
	return m.ch.PublishWithContext(ctx, m.exchange, msg.Topic, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/octet-stream",
		MessageId:    msg.Key,
		Headers:      headers,
		Body:         msg.Value,
	})
}

// 中文：PublishBatch 执行当前包中的对应流程。
// English: PublishBatch executes the corresponding workflow in this package.
func (m *RabbitMQ) PublishBatch(ctx context.Context, msgs []*Message) error {
	for _, msg := range msgs {
		if err := m.Publish(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

// 中文：Subscribe 执行当前包中的对应流程。
// English: Subscribe executes the corresponding workflow in this package.
func (m *RabbitMQ) Subscribe(ctx context.Context, topic string, handler func(msg *Message) error) error {
	if topic == "" {
		return fmt.Errorf("rabbitmq topic is required")
	}
	if handler == nil {
		return fmt.Errorf("rabbitmq handler is required")
	}

	queueName := fmt.Sprintf("%s.%s", m.queuePrefix, topic)
	if _, err := m.ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return fmt.Errorf("rabbitmq declare queue: %w", err)
	}
	if err := m.ch.QueueBind(queueName, topic, m.exchange, false, nil); err != nil {
		return fmt.Errorf("rabbitmq bind queue: %w", err)
	}

	deliveries, err := m.ch.ConsumeWithContext(ctx, queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("rabbitmq consume: %w", err)
	}

	go func() {
		for delivery := range deliveries {
			msg := &Message{
				Topic:   topic,
				Key:     delivery.MessageId,
				Value:   delivery.Body,
				Headers: amqpHeaders(delivery.Headers),
			}
			if err := handler(msg); err != nil {
				_ = delivery.Nack(false, true)
				continue
			}
			_ = delivery.Ack(false)
		}
	}()
	return nil
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (m *RabbitMQ) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil
	}
	m.closed = true
	var err error
	if m.ch != nil {
		err = m.ch.Close()
	}
	if m.conn != nil {
		if closeErr := m.conn.Close(); err == nil {
			err = closeErr
		}
	}
	return err
}

// 中文：amqpHeaders 执行当前包中的对应流程。
// English: amqpHeaders executes the corresponding workflow in this package.
func amqpHeaders(values amqp.Table) map[string]string {
	headers := make(map[string]string, len(values))
	for key, value := range values {
		headers[key] = fmt.Sprintf("%v", value)
	}
	return headers
}
