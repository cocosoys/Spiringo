package event

import (
	"context"
	"encoding/json"

	"github.com/spiringo/spiringo/internal/pkg/mq"
)

// 中文：MQBridgeConfig 定义当前包使用的数据结构或接口。
// English: MQBridgeConfig defines a data structure or interface used by this package.
type MQBridgeConfig struct {
	// 中文：IncludeTopics 保存当前结构中的配置或数据值。
	// English: IncludeTopics stores a configuration or data value for this struct.
	IncludeTopics []string
	// 中文：ExcludeTopics 保存当前结构中的配置或数据值。
	// English: ExcludeTopics stores a configuration or data value for this struct.
	ExcludeTopics []string
	// 中文：Source 保存当前结构中的配置或数据值。
	// English: Source stores a configuration or data value for this struct.
	Source string
}

// 中文：MQBridge 定义当前包使用的数据结构或接口。
// English: MQBridge defines a data structure or interface used by this package.
type MQBridge struct {
	// 中文：bus 保存当前结构中的配置或数据值。
	// English: bus stores a configuration or data value for this struct.
	bus *Bus
	// 中文：queue 保存当前结构中的配置或数据值。
	// English: queue stores a configuration or data value for this struct.
	queue mq.MQ
	// 中文：include 保存当前结构中的配置或数据值。
	// English: include stores a configuration or data value for this struct.
	include map[string]struct{}
	// 中文：exclude 保存当前结构中的配置或数据值。
	// English: exclude stores a configuration or data value for this struct.
	exclude map[string]struct{}
	// 中文：source 保存当前结构中的配置或数据值。
	// English: source stores a configuration or data value for this struct.
	source string
	// 中文：attached 保存当前结构中的配置或数据值。
	// English: attached stores a configuration or data value for this struct.
	attached bool
}

// 中文：NewMQBridge 创建并返回对应组件实例。
// English: NewMQBridge creates and returns the corresponding component instance.
func NewMQBridge(bus *Bus, queue mq.MQ, cfg MQBridgeConfig) *MQBridge {
	return &MQBridge{
		bus:     bus,
		queue:   queue,
		include: topicSet(cfg.IncludeTopics),
		exclude: topicSet(cfg.ExcludeTopics),
		source:  cfg.Source,
	}
}

// 中文：Register 执行当前包中的对应流程。
// English: Register executes the corresponding workflow in this package.
func (b *MQBridge) Register() {
	if b == nil || b.attached || b.bus == nil {
		return
	}
	if err := b.bus.Subscribe("*", b.Handle); err != nil {
		return
	}
	b.attached = true
}

// 中文：Handle 执行当前包中的对应流程。
// English: Handle executes the corresponding workflow in this package.
func (b *MQBridge) Handle(ctx context.Context, e *Event) error {
	if b == nil || b.queue == nil || e == nil || !b.shouldPublish(e.Topic) {
		return nil
	}

	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	headers := make(map[string]string, len(e.Metadata)+2)
	for key, value := range e.Metadata {
		headers[key] = value
	}
	if e.Source != "" {
		headers["event_source"] = e.Source
	} else if b.source != "" {
		headers["event_source"] = b.source
	}
	headers["event_topic"] = e.Topic

	return b.queue.Publish(ctx, &mq.Message{
		Topic:   e.Topic,
		Key:     headers["key"],
		Value:   payload,
		Headers: headers,
	})
}

// 中文：shouldPublish 执行当前包中的对应流程。
// English: shouldPublish executes the corresponding workflow in this package.
func (b *MQBridge) shouldPublish(topic string) bool {
	if _, ok := b.exclude[topic]; ok {
		return false
	}
	if len(b.include) == 0 {
		return true
	}
	_, ok := b.include[topic]
	return ok
}

// 中文：topicSet 执行当前包中的对应流程。
// English: topicSet executes the corresponding workflow in this package.
func topicSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value != "" {
			result[value] = struct{}{}
		}
	}
	return result
}
