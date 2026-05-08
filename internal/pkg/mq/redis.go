package mq

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

// 中文：RedisStreamMQ 定义当前包使用的数据结构或接口。
// English: RedisStreamMQ defines a data structure or interface used by this package.
// RedisStreamMQ 基于Redis Streams的消息队列实现
type RedisStreamMQ struct {
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client redis.Cmdable
	// 中文：prefix 保存当前结构中的配置或数据值。
	// English: prefix stores a configuration or data value for this struct.
	prefix string
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
	// 中文：subs 保存当前结构中的配置或数据值。
	// English: subs stores a configuration or data value for this struct.
	subs map[string]context.CancelFunc
}

// 中文：NewRedisStreamMQ 创建并返回对应组件实例。
// English: NewRedisStreamMQ creates and returns the corresponding component instance.
// NewRedisStreamMQ 创建Redis Streams消息队列
func NewRedisStreamMQ(client redis.Cmdable, prefix string) *RedisStreamMQ {
	if prefix == "" {
		prefix = "spiringo:mq"
	}
	return &RedisStreamMQ{
		client: client,
		prefix: prefix,
		subs:   make(map[string]context.CancelFunc),
	}
}

// 中文：streamKey 执行当前包中的对应流程。
// English: streamKey executes the corresponding workflow in this package.
// streamKey 获取stream key
func (m *RedisStreamMQ) streamKey(topic string) string {
	return fmt.Sprintf("%s:%s", m.prefix, topic)
}

// 中文：groupKey 执行当前包中的对应流程。
// English: groupKey executes the corresponding workflow in this package.
// groupKey 获取consumer group名称
func (m *RedisStreamMQ) groupKey(topic string) string {
	return fmt.Sprintf("%s:cg", topic)
}

// 中文：Publish 执行当前包中的对应流程。
// English: Publish executes the corresponding workflow in this package.
// Publish 发布消息
func (m *RedisStreamMQ) Publish(ctx context.Context, msg *Message) error {
	if msg == nil {
		return fmt.Errorf("redis stream message is required")
	}
	if msg.Topic == "" {
		return fmt.Errorf("redis stream topic is required")
	}
	if m.client == nil {
		return fmt.Errorf("redis stream client is required")
	}
	key := m.streamKey(msg.Topic)
	values := map[string]interface{}{
		"body": string(msg.Value),
		"key":  msg.Key,
	}
	for k, v := range msg.Headers {
		values["header:"+k] = v
	}

	_, err := m.client.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		Values: values,
		ID:     "*",
	}).Result()
	return err
}

// 中文：PublishBatch 执行当前包中的对应流程。
// English: PublishBatch executes the corresponding workflow in this package.
// PublishBatch 批量发布消息
func (m *RedisStreamMQ) PublishBatch(ctx context.Context, msgs []*Message) error {
	for _, msg := range msgs {
		if err := m.Publish(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

// 中文：Subscribe 执行当前包中的对应流程。
// English: Subscribe executes the corresponding workflow in this package.
// Subscribe 订阅消息
func (m *RedisStreamMQ) Subscribe(ctx context.Context, topic string, handler func(msg *Message) error) error {
	if topic == "" {
		return fmt.Errorf("redis stream topic is required")
	}
	if handler == nil {
		return fmt.Errorf("redis stream handler is required")
	}
	if m.client == nil {
		return fmt.Errorf("redis stream client is required")
	}
	key := m.streamKey(topic)
	group := m.groupKey(topic)
	consumer := fmt.Sprintf("consumer-%s", topic)

	// 创建consumer group（如果不存在）
	_ = m.client.XGroupCreateMkStream(ctx, key, group, "0").Err()

	// 启动消费goroutine
	subCtx, cancel := context.WithCancel(ctx)

	m.mu.Lock()
	m.subs[topic] = cancel
	m.mu.Unlock()

	go func() {
		defer cancel()
		for {
			select {
			case <-subCtx.Done():
				return
			default:
			}

			streams, err := m.client.XReadGroup(subCtx, &redis.XReadGroupArgs{
				Group:    group,
				Consumer: consumer,
				Streams:  []string{key, ">"},
				Count:    10,
				Block:    0,
			}).Result()
			if err != nil {
				if subCtx.Err() != nil {
					return
				}
				continue
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					msg := &Message{
						Topic:   topic,
						Key:     getString(message.Values, "key"),
						Value:   []byte(getString(message.Values, "body")),
						Headers: extractHeaders(message.Values),
					}

					if err := handler(msg); err != nil {
						continue
					}

					// ACK消息
					_ = m.client.XAck(subCtx, key, group, message.ID).Err()
				}
			}
		}
	}()

	return nil
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
// Close 关闭消息队列
func (m *RedisStreamMQ) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, cancel := range m.subs {
		cancel()
	}
	m.subs = make(map[string]context.CancelFunc)
	return nil
}

// 中文：getString 执行当前包中的对应流程。
// English: getString executes the corresponding workflow in this package.
// getString 从map中安全获取字符串值
func getString(values map[string]interface{}, key string) string {
	if v, ok := values[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// 中文：extractHeaders 执行当前包中的对应流程。
// English: extractHeaders executes the corresponding workflow in this package.
// extractHeaders 提取headers
func extractHeaders(values map[string]interface{}) map[string]string {
	headers := make(map[string]string)
	for k, v := range values {
		if len(k) > 7 && k[:7] == "header:" {
			headers[k[7:]] = fmt.Sprintf("%v", v)
		}
	}
	return headers
}
