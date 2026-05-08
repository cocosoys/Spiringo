package event

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"
)

// 中文：ErrNilEvent、ErrInvalidTopic、ErrNilHandler 声明当前包使用的变量。
// English: ErrNilEvent、ErrInvalidTopic、ErrNilHandler declares variables used by this package.
var (
	ErrNilEvent     = errors.New("event is required")
	ErrInvalidTopic = errors.New("event topic is required")
	ErrNilHandler   = errors.New("event handler is required")
)

// 中文：Event 定义当前包使用的数据结构或接口。
// English: Event defines a data structure or interface used by this package.
// Event is the message passed through the in-process event bus.
type Event struct {
	// 中文：Topic 保存当前结构中的配置或数据值。
	// English: Topic stores a configuration or data value for this struct.
	Topic string `json:"topic"`
	// 中文：Payload 保存当前结构中的配置或数据值。
	// English: Payload stores a configuration or data value for this struct.
	Payload any `json:"payload"`
	// 中文：Timestamp 保存当前结构中的配置或数据值。
	// English: Timestamp stores a configuration or data value for this struct.
	Timestamp time.Time `json:"timestamp"`
	// 中文：Source 保存当前结构中的配置或数据值。
	// English: Source stores a configuration or data value for this struct.
	Source string `json:"source"`
	// 中文：Metadata 保存当前结构中的配置或数据值。
	// English: Metadata stores a configuration or data value for this struct.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// 中文：NewEvent 创建并返回对应组件实例。
// English: NewEvent creates and returns the corresponding component instance.
// NewEvent creates an event with a timestamp and metadata map.
func NewEvent(topic string, payload any) *Event {
	return &Event{
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}
}

// 中文：WithSource 执行当前包中的对应流程。
// English: WithSource executes the corresponding workflow in this package.
// WithSource sets the module or subsystem that produced the event.
func (e *Event) WithSource(source string) *Event {
	e.Source = source
	return e
}

// 中文：WithMetadata 执行当前包中的对应流程。
// English: WithMetadata executes the corresponding workflow in this package.
// WithMetadata adds a metadata entry.
func (e *Event) WithMetadata(key, value string) *Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

// 中文：Handler 定义当前包使用的数据结构或接口。
// English: Handler defines a data structure or interface used by this package.
// Handler handles one event.
type Handler func(ctx context.Context, event *Event) error

// 中文：ErrorHandler 定义当前包使用的数据结构或接口。
// English: ErrorHandler defines a data structure or interface used by this package.
// ErrorHandler observes subscriber failures without blocking other handlers.
type ErrorHandler func(ctx context.Context, event *Event, err error)

// 中文：Subscription 定义当前包使用的数据结构或接口。
// English: Subscription defines a data structure or interface used by this package.
// Subscription describes one event subscription.
type Subscription struct {
	// 中文：Topic 保存当前结构中的配置或数据值。
	// English: Topic stores a configuration or data value for this struct.
	Topic string
	// 中文：Handler 保存当前结构中的配置或数据值。
	// English: Handler stores a configuration or data value for this struct.
	Handler Handler
	// 中文：Module 保存当前结构中的配置或数据值。
	// English: Module stores a configuration or data value for this struct.
	Module string
}

// 中文：EventBus 定义当前包使用的数据结构或接口。
// English: EventBus defines a data structure or interface used by this package.
// EventBus is the public event bus contract used by modules.
type EventBus interface {
	// 中文：Publish 声明该接口需要实现的行为。
	// English: Publish declares behavior required by this interface.
	Publish(ctx context.Context, event *Event) error
	// 中文：PublishAsync 声明该接口需要实现的行为。
	// English: PublishAsync declares behavior required by this interface.
	PublishAsync(ctx context.Context, event *Event) error
	// 中文：Subscribe 声明该接口需要实现的行为。
	// English: Subscribe declares behavior required by this interface.
	Subscribe(topic string, handler Handler) error
	// 中文：Unsubscribe 声明该接口需要实现的行为。
	// English: Unsubscribe declares behavior required by this interface.
	Unsubscribe(topic string, handler Handler) error
}

// 中文：Bus 定义当前包使用的数据结构或接口。
// English: Bus defines a data structure or interface used by this package.
// Bus is an in-process event bus with exact, wildcard, and async delivery.
type Bus struct {
	// 中文：subscriptions 保存当前结构中的配置或数据值。
	// English: subscriptions stores a configuration or data value for this struct.
	subscriptions map[string][]Handler
	// 中文：asyncWorkers 保存当前结构中的配置或数据值。
	// English: asyncWorkers stores a configuration or data value for this struct.
	asyncWorkers int
	// 中文：asyncQueue 保存当前结构中的配置或数据值。
	// English: asyncQueue stores a configuration or data value for this struct.
	asyncQueue chan *asyncEvent
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
	// 中文：cancel 保存当前结构中的配置或数据值。
	// English: cancel stores a configuration or data value for this struct.
	cancel context.CancelFunc
	// 中文：wg 保存当前结构中的配置或数据值。
	// English: wg stores a configuration or data value for this struct.
	wg sync.WaitGroup
	// 中文：errorHandler 保存当前结构中的配置或数据值。
	// English: errorHandler stores a configuration or data value for this struct.
	errorHandler ErrorHandler
}

// 中文：asyncEvent 定义当前包使用的数据结构或接口。
// English: asyncEvent defines a data structure or interface used by this package.
type asyncEvent struct {
	// 中文：ctx 保存当前结构中的配置或数据值。
	// English: ctx stores a configuration or data value for this struct.
	ctx context.Context
	// 中文：event 保存当前结构中的配置或数据值。
	// English: event stores a configuration or data value for this struct.
	event *Event
}

// 中文：_ 声明当前包使用的变量。
// English: _ declares variables used by this package.
var _ EventBus = (*Bus)(nil)

// 中文：NewBus 创建并返回对应组件实例。
// English: NewBus creates and returns the corresponding component instance.
// NewBus creates an event bus.
func NewBus(asyncWorkers int) *Bus {
	if asyncWorkers <= 0 {
		asyncWorkers = 4
	}
	return &Bus{
		subscriptions: make(map[string][]Handler),
		asyncWorkers:  asyncWorkers,
		asyncQueue:    make(chan *asyncEvent, 1024),
	}
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
// Start launches asynchronous workers.
func (b *Bus) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	b.cancel = cancel

	for i := 0; i < b.asyncWorkers; i++ {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case ae := <-b.asyncQueue:
					b.dispatch(ae.ctx, ae.event)
				}
			}
		}()
	}
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
// Stop stops asynchronous workers.
func (b *Bus) Stop() {
	if b.cancel != nil {
		b.cancel()
	}
	b.wg.Wait()
}

// 中文：Subscribe 执行当前包中的对应流程。
// English: Subscribe executes the corresponding workflow in this package.
// Subscribe adds a handler for a topic. Topic "*" receives all events and
// "module.*" receives prefixed events such as "module.created".
func (b *Bus) Subscribe(topic string, handler Handler) error {
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return ErrInvalidTopic
	}
	if handler == nil {
		return ErrNilHandler
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscriptions[topic] = append(b.subscriptions[topic], handler)
	return nil
}

// 中文：SetErrorHandler 执行当前包中的对应流程。
// English: SetErrorHandler executes the corresponding workflow in this package.
// SetErrorHandler registers a process-wide observer for handler errors.
func (b *Bus) SetErrorHandler(handler ErrorHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.errorHandler = handler
}

// 中文：Unsubscribe 执行当前包中的对应流程。
// English: Unsubscribe executes the corresponding workflow in this package.
// Unsubscribe removes one matching handler for a topic.
func (b *Bus) Unsubscribe(topic string, handler Handler) error {
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return ErrInvalidTopic
	}
	if handler == nil {
		return ErrNilHandler
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.subscriptions[topic]
	for i, h := range handlers {
		if sameHandler(h, handler) {
			b.subscriptions[topic] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	if len(b.subscriptions[topic]) == 0 {
		delete(b.subscriptions, topic)
	}
	return nil
}

// 中文：Publish 执行当前包中的对应流程。
// English: Publish executes the corresponding workflow in this package.
// Publish synchronously dispatches an event.
func (b *Bus) Publish(ctx context.Context, event *Event) error {
	if event == nil {
		return ErrNilEvent
	}
	b.dispatch(ctx, event)
	return nil
}

// 中文：PublishAsync 执行当前包中的对应流程。
// English: PublishAsync executes the corresponding workflow in this package.
// PublishAsync queues an event for async dispatch, falling back to synchronous
// delivery when the queue is full.
func (b *Bus) PublishAsync(ctx context.Context, event *Event) error {
	if event == nil {
		return ErrNilEvent
	}
	select {
	case b.asyncQueue <- &asyncEvent{ctx: ctx, event: event}:
		return nil
	default:
		b.dispatch(ctx, event)
		return nil
	}
}

// 中文：dispatch 执行当前包中的对应流程。
// English: dispatch executes the corresponding workflow in this package.
func (b *Bus) dispatch(ctx context.Context, event *Event) {
	if event == nil {
		return
	}

	b.mu.RLock()
	handlers := make([]Handler, 0)
	errorHandler := b.errorHandler
	if topicHandlers, ok := b.subscriptions[event.Topic]; ok {
		handlers = append(handlers, topicHandlers...)
	}
	for topic, topicHandlers := range b.subscriptions {
		if topic == event.Topic {
			continue
		}
		if topic == "*" {
			handlers = append(handlers, topicHandlers...)
			continue
		}
		if strings.HasSuffix(topic, ".*") {
			prefix := topic[:len(topic)-2]
			if strings.HasPrefix(event.Topic, prefix+".") {
				handlers = append(handlers, topicHandlers...)
			}
		}
	}
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil && errorHandler != nil {
			errorHandler(ctx, event, err)
		}
	}
}

// 中文：SubscriptionCount 执行当前包中的对应流程。
// English: SubscriptionCount executes the corresponding workflow in this package.
// SubscriptionCount returns the number of handlers registered for a topic.
func (b *Bus) SubscriptionCount(topic string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subscriptions[topic])
}

// 中文：sameHandler 执行当前包中的对应流程。
// English: sameHandler executes the corresponding workflow in this package.
func sameHandler(a, b Handler) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}
