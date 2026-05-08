package mq

import "context"

// 中文：Message 定义当前包使用的数据结构或接口。
// English: Message defines a data structure or interface used by this package.
// Message 消息
type Message struct {
	// 中文：Topic 保存当前结构中的配置或数据值。
	// English: Topic stores a configuration or data value for this struct.
	Topic string `json:"topic"`
	// 中文：Key 保存当前结构中的配置或数据值。
	// English: Key stores a configuration or data value for this struct.
	Key string `json:"key"`
	// 中文：Value 保存当前结构中的配置或数据值。
	// English: Value stores a configuration or data value for this struct.
	Value []byte `json:"value"`
	// 中文：Headers 保存当前结构中的配置或数据值。
	// English: Headers stores a configuration or data value for this struct.
	Headers map[string]string `json:"headers,omitempty"`
}

// 中文：MQ 定义当前包使用的数据结构或接口。
// English: MQ defines a data structure or interface used by this package.
// MQ 消息队列接口
type MQ interface {
	// 中文：Publish 声明该接口需要实现的行为。
	// English: Publish declares behavior required by this interface.
	Publish(ctx context.Context, msg *Message) error
	// 中文：PublishBatch 声明该接口需要实现的行为。
	// English: PublishBatch declares behavior required by this interface.
	PublishBatch(ctx context.Context, msgs []*Message) error
	// 中文：Subscribe 声明该接口需要实现的行为。
	// English: Subscribe declares behavior required by this interface.
	Subscribe(ctx context.Context, topic string, handler func(msg *Message) error) error
	// 中文：Close 声明该接口需要实现的行为。
	// English: Close declares behavior required by this interface.
	Close() error
}
