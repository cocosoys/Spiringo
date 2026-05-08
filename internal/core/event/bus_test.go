package event

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spiringo/spiringo/internal/pkg/mq"
)

// 中文：TestBus_SubscribeAndPublish 验证相关行为符合预期。
// English: TestBus_SubscribeAndPublish verifies the related behavior.
func TestBus_SubscribeAndPublish(t *testing.T) {
	bus := NewBus(4)
	var received atomic.Int32

	bus.Subscribe("test.topic", func(ctx context.Context, e *Event) error {
		received.Add(1)
		return nil
	})

	bus.Publish(context.Background(), NewEvent("test.topic", "payload"))
	time.Sleep(50 * time.Millisecond)

	if received.Load() != 1 {
		t.Errorf("expected 1 event, got %d", received.Load())
	}
}

// 中文：TestBus_WildcardSubscribe 验证相关行为符合预期。
// English: TestBus_WildcardSubscribe verifies the related behavior.
func TestBus_WildcardSubscribe(t *testing.T) {
	bus := NewBus(4)
	var received atomic.Int32

	bus.Subscribe("user.*", func(ctx context.Context, e *Event) error {
		received.Add(1)
		return nil
	})

	bus.Publish(context.Background(), NewEvent("user.created", nil))
	bus.Publish(context.Background(), NewEvent("user.deleted", nil))
	bus.Publish(context.Background(), NewEvent("order.created", nil))
	time.Sleep(50 * time.Millisecond)

	if received.Load() != 2 {
		t.Errorf("expected 2 events from user.*, got %d", received.Load())
	}
}

// 中文：TestBus_PublishAsync 验证相关行为符合预期。
// English: TestBus_PublishAsync verifies the related behavior.
func TestBus_PublishAsync(t *testing.T) {
	bus := NewBus(4)
	ctx := context.Background()
	bus.Start(ctx)
	defer bus.Stop()

	var received atomic.Int32

	bus.Subscribe("async.topic", func(ctx context.Context, e *Event) error {
		received.Add(1)
		return nil
	})

	bus.PublishAsync(ctx, NewEvent("async.topic", nil))
	time.Sleep(100 * time.Millisecond)

	if received.Load() != 1 {
		t.Errorf("expected 1 async event, got %d", received.Load())
	}
}

// 中文：TestBus_MultipleSubscribers 验证相关行为符合预期。
// English: TestBus_MultipleSubscribers verifies the related behavior.
func TestBus_MultipleSubscribers(t *testing.T) {
	bus := NewBus(4)
	var received atomic.Int32

	handler := func(ctx context.Context, e *Event) error {
		received.Add(1)
		return nil
	}

	bus.Subscribe("multi.topic", handler)
	bus.Subscribe("multi.topic", handler)

	bus.Publish(context.Background(), NewEvent("multi.topic", nil))
	time.Sleep(50 * time.Millisecond)

	if received.Load() != 2 {
		t.Errorf("expected 2 invocations, got %d", received.Load())
	}
}

// 中文：TestBus_UnsubscribeRemovesOnlyMatchingHandler 验证相关行为符合预期。
// English: TestBus_UnsubscribeRemovesOnlyMatchingHandler verifies the related behavior.
func TestBus_UnsubscribeRemovesOnlyMatchingHandler(t *testing.T) {
	bus := NewBus(1)
	var first atomic.Int32
	var second atomic.Int32

	handlerOne := func(ctx context.Context, e *Event) error {
		first.Add(1)
		return nil
	}
	handlerTwo := func(ctx context.Context, e *Event) error {
		second.Add(1)
		return nil
	}

	bus.Subscribe("unsub.topic", handlerOne)
	bus.Subscribe("unsub.topic", handlerTwo)
	bus.Unsubscribe("unsub.topic", handlerOne)

	if count := bus.SubscriptionCount("unsub.topic"); count != 1 {
		t.Fatalf("subscription count = %d, want 1", count)
	}
	if err := bus.Publish(context.Background(), NewEvent("unsub.topic", nil)); err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}
	if first.Load() != 0 {
		t.Fatalf("first handler calls = %d, want 0", first.Load())
	}
	if second.Load() != 1 {
		t.Fatalf("second handler calls = %d, want 1", second.Load())
	}
}

// 中文：TestBus_RejectsInvalidPublishAndSubscription 验证相关行为符合预期。
// English: TestBus_RejectsInvalidPublishAndSubscription verifies the related behavior.
func TestBus_RejectsInvalidPublishAndSubscription(t *testing.T) {
	bus := NewBus(1)
	if err := bus.Publish(context.Background(), nil); !errors.Is(err, ErrNilEvent) {
		t.Fatalf("publish nil err = %v, want ErrNilEvent", err)
	}
	if err := bus.PublishAsync(context.Background(), nil); !errors.Is(err, ErrNilEvent) {
		t.Fatalf("publish async nil err = %v, want ErrNilEvent", err)
	}
	if err := bus.Subscribe("", func(context.Context, *Event) error { return nil }); !errors.Is(err, ErrInvalidTopic) {
		t.Fatalf("subscribe empty topic err = %v, want ErrInvalidTopic", err)
	}
	if err := bus.Subscribe("topic", nil); !errors.Is(err, ErrNilHandler) {
		t.Fatalf("subscribe nil handler err = %v, want ErrNilHandler", err)
	}
	if err := bus.Unsubscribe("topic", nil); !errors.Is(err, ErrNilHandler) {
		t.Fatalf("unsubscribe nil handler err = %v, want ErrNilHandler", err)
	}
}

// 中文：TestBus_ErrorHandlerObservesFailures 验证相关行为符合预期。
// English: TestBus_ErrorHandlerObservesFailures verifies the related behavior.
func TestBus_ErrorHandlerObservesFailures(t *testing.T) {
	bus := NewBus(1)
	var observed atomic.Int32
	var delivered atomic.Int32

	bus.SetErrorHandler(func(ctx context.Context, e *Event, err error) {
		if e.Topic == "fail.topic" && err != nil {
			observed.Add(1)
		}
	})
	bus.Subscribe("fail.topic", func(ctx context.Context, e *Event) error {
		return errors.New("boom")
	})
	bus.Subscribe("fail.topic", func(ctx context.Context, e *Event) error {
		delivered.Add(1)
		return nil
	})

	bus.Publish(context.Background(), NewEvent("fail.topic", nil))

	if observed.Load() != 1 {
		t.Fatalf("observed errors = %d, want 1", observed.Load())
	}
	if delivered.Load() != 1 {
		t.Fatalf("delivered handlers = %d, want 1", delivered.Load())
	}
}

// 中文：TestMQBridgePublishesEvents 验证相关行为符合预期。
// English: TestMQBridgePublishesEvents verifies the related behavior.
func TestMQBridgePublishesEvents(t *testing.T) {
	queue := &fakeMQ{}
	bus := NewBus(1)
	bridge := NewMQBridge(bus, queue, MQBridgeConfig{Source: "test"})
	bridge.Register()

	event := NewEvent("user.created", map[string]string{"id": "u1"}).WithMetadata("key", "u1")
	if err := bus.Publish(context.Background(), event); err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}

	if len(queue.messages) != 1 {
		t.Fatalf("message count = %d, want 1", len(queue.messages))
	}
	msg := queue.messages[0]
	if msg.Topic != "user.created" || msg.Key != "u1" {
		t.Fatalf("unexpected message: %+v", msg)
	}
	if string(msg.Value) == "" {
		t.Fatal("empty message payload")
	}
	if msg.Headers["event_source"] != "test" {
		t.Fatalf("event_source = %q, want test", msg.Headers["event_source"])
	}
}

// 中文：fakeMQ 定义当前包使用的数据结构或接口。
// English: fakeMQ defines a data structure or interface used by this package.
type fakeMQ struct {
	// 中文：messages 保存当前结构中的配置或数据值。
	// English: messages stores a configuration or data value for this struct.
	messages []*mq.Message
}

// 中文：Publish 执行当前包中的对应流程。
// English: Publish executes the corresponding workflow in this package.
func (f *fakeMQ) Publish(ctx context.Context, msg *mq.Message) error {
	f.messages = append(f.messages, msg)
	return nil
}

// 中文：PublishBatch 执行当前包中的对应流程。
// English: PublishBatch executes the corresponding workflow in this package.
func (f *fakeMQ) PublishBatch(ctx context.Context, msgs []*mq.Message) error {
	f.messages = append(f.messages, msgs...)
	return nil
}

// 中文：Subscribe 执行当前包中的对应流程。
// English: Subscribe executes the corresponding workflow in this package.
func (f *fakeMQ) Subscribe(ctx context.Context, topic string, handler func(msg *mq.Message) error) error {
	return nil
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (f *fakeMQ) Close() error { return nil }
