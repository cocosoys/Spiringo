package mq

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// 中文：KafkaConfig 定义当前包使用的数据结构或接口。
// English: KafkaConfig defines a data structure or interface used by this package.
type KafkaConfig struct {
	// 中文：Brokers 保存当前结构中的配置或数据值。
	// English: Brokers stores a configuration or data value for this struct.
	Brokers []string
	// 中文：ClientID 保存当前结构中的配置或数据值。
	// English: ClientID stores a configuration or data value for this struct.
	ClientID string
	// 中文：GroupID 保存当前结构中的配置或数据值。
	// English: GroupID stores a configuration or data value for this struct.
	GroupID string
	// 中文：TopicPrefix 保存当前结构中的配置或数据值。
	// English: TopicPrefix stores a configuration or data value for this struct.
	TopicPrefix string
}

// 中文：KafkaMQ 定义当前包使用的数据结构或接口。
// English: KafkaMQ defines a data structure or interface used by this package.
type KafkaMQ struct {
	// 中文：writer 保存当前结构中的配置或数据值。
	// English: writer stores a configuration or data value for this struct.
	writer *kafka.Writer
	// 中文：brokers 保存当前结构中的配置或数据值。
	// English: brokers stores a configuration or data value for this struct.
	brokers []string
	// 中文：clientID 保存当前结构中的配置或数据值。
	// English: clientID stores a configuration or data value for this struct.
	clientID string
	// 中文：groupID 保存当前结构中的配置或数据值。
	// English: groupID stores a configuration or data value for this struct.
	groupID string
	// 中文：topicPrefix 保存当前结构中的配置或数据值。
	// English: topicPrefix stores a configuration or data value for this struct.
	topicPrefix string
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
	// 中文：readers 保存当前结构中的配置或数据值。
	// English: readers stores a configuration or data value for this struct.
	readers []*kafka.Reader
}

// 中文：NewKafkaMQ 创建并返回对应组件实例。
// English: NewKafkaMQ creates and returns the corresponding component instance.
func NewKafkaMQ(cfg KafkaConfig) (*KafkaMQ, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers are required")
	}
	if cfg.GroupID == "" {
		cfg.GroupID = "spiringo"
	}
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
	}
	return &KafkaMQ{
		writer:      writer,
		brokers:     append([]string(nil), cfg.Brokers...),
		clientID:    cfg.ClientID,
		groupID:     cfg.GroupID,
		topicPrefix: strings.TrimSpace(cfg.TopicPrefix),
	}, nil
}

// 中文：Publish 执行当前包中的对应流程。
// English: Publish executes the corresponding workflow in this package.
func (m *KafkaMQ) Publish(ctx context.Context, msg *Message) error {
	if msg == nil {
		return fmt.Errorf("kafka message is required")
	}
	return m.writer.WriteMessages(ctx, m.kafkaMessage(msg))
}

// 中文：PublishBatch 执行当前包中的对应流程。
// English: PublishBatch executes the corresponding workflow in this package.
func (m *KafkaMQ) PublishBatch(ctx context.Context, msgs []*Message) error {
	kafkaMsgs := make([]kafka.Message, 0, len(msgs))
	for _, msg := range msgs {
		if msg == nil {
			continue
		}
		kafkaMsgs = append(kafkaMsgs, m.kafkaMessage(msg))
	}
	if len(kafkaMsgs) == 0 {
		return nil
	}
	return m.writer.WriteMessages(ctx, kafkaMsgs...)
}

// 中文：Subscribe 执行当前包中的对应流程。
// English: Subscribe executes the corresponding workflow in this package.
func (m *KafkaMQ) Subscribe(ctx context.Context, topic string, handler func(msg *Message) error) error {
	if topic == "" {
		return fmt.Errorf("kafka topic is required")
	}
	if handler == nil {
		return fmt.Errorf("kafka handler is required")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: m.brokers,
		GroupID: m.groupID,
		Topic:   m.topic(topic),
	})
	m.mu.Lock()
	m.readers = append(m.readers, reader)
	m.mu.Unlock()

	go func() {
		defer reader.Close()
		for {
			record, err := reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				time.Sleep(200 * time.Millisecond)
				continue
			}
			msg := &Message{
				Topic:   topic,
				Key:     string(record.Key),
				Value:   record.Value,
				Headers: kafkaHeaders(record.Headers),
			}
			if err := handler(msg); err != nil {
				continue
			}
			_ = reader.CommitMessages(ctx, record)
		}
	}()
	return nil
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (m *KafkaMQ) Close() error {
	m.mu.Lock()
	readers := append([]*kafka.Reader(nil), m.readers...)
	m.readers = nil
	m.mu.Unlock()

	var err error
	for _, reader := range readers {
		if readerErr := reader.Close(); err == nil {
			err = readerErr
		}
	}
	if m.writer != nil {
		if writerErr := m.writer.Close(); err == nil {
			err = writerErr
		}
	}
	return err
}

// 中文：kafkaMessage 执行当前包中的对应流程。
// English: kafkaMessage executes the corresponding workflow in this package.
func (m *KafkaMQ) kafkaMessage(msg *Message) kafka.Message {
	headers := make([]kafka.Header, 0, len(msg.Headers))
	for key, value := range msg.Headers {
		headers = append(headers, kafka.Header{Key: key, Value: []byte(value)})
	}
	return kafka.Message{
		Topic:   m.topic(msg.Topic),
		Key:     []byte(msg.Key),
		Value:   msg.Value,
		Headers: headers,
	}
}

// 中文：topic 执行当前包中的对应流程。
// English: topic executes the corresponding workflow in this package.
func (m *KafkaMQ) topic(topic string) string {
	if m.topicPrefix == "" {
		return topic
	}
	return m.topicPrefix + topic
}

// 中文：kafkaHeaders 执行当前包中的对应流程。
// English: kafkaHeaders executes the corresponding workflow in this package.
func kafkaHeaders(values []kafka.Header) map[string]string {
	headers := make(map[string]string, len(values))
	for _, header := range values {
		headers[header.Key] = string(header.Value)
	}
	return headers
}
