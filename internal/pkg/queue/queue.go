package queue

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 中文：ErrTaskRequired、ErrTaskName、ErrHandlerRequired、... 声明当前包使用的变量。
// English: ErrTaskRequired、ErrTaskName、ErrHandlerRequired、... declares variables used by this package.
var (
	ErrTaskRequired    = errors.New("queue task is required")
	ErrTaskName        = errors.New("queue task name is required")
	ErrHandlerRequired = errors.New("queue handler is required")
	ErrHandlerNotFound = errors.New("queue handler not found")
	ErrClosed          = errors.New("queue is closed")
)

// 中文：Task 定义当前包使用的数据结构或接口。
// English: Task defines a data structure or interface used by this package.
type Task struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"id"`
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name"`
	// 中文：Payload 保存当前结构中的配置或数据值。
	// English: Payload stores a configuration or data value for this struct.
	Payload any `json:"payload,omitempty"`
	// 中文：RunAt 保存当前结构中的配置或数据值。
	// English: RunAt stores a configuration or data value for this struct.
	RunAt time.Time `json:"run_at,omitempty"`
	// 中文：MaxRetries 保存当前结构中的配置或数据值。
	// English: MaxRetries stores a configuration or data value for this struct.
	MaxRetries int `json:"max_retries,omitempty"`
	// 中文：Attempts 保存当前结构中的配置或数据值。
	// English: Attempts stores a configuration or data value for this struct.
	Attempts int `json:"attempts,omitempty"`
	// 中文：Metadata 保存当前结构中的配置或数据值。
	// English: Metadata stores a configuration or data value for this struct.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// 中文：Handler 定义当前包使用的数据结构或接口。
// English: Handler defines a data structure or interface used by this package.
type Handler func(ctx context.Context, task *Task) error

// 中文：ErrorHandler 定义当前包使用的数据结构或接口。
// English: ErrorHandler defines a data structure or interface used by this package.
type ErrorHandler func(ctx context.Context, task *Task, err error)

// 中文：Queue 定义当前包使用的数据结构或接口。
// English: Queue defines a data structure or interface used by this package.
type Queue interface {
	// 中文：Register 声明该接口需要实现的行为。
	// English: Register declares behavior required by this interface.
	Register(name string, handler Handler) error
	// 中文：Enqueue 声明该接口需要实现的行为。
	// English: Enqueue declares behavior required by this interface.
	Enqueue(ctx context.Context, task *Task) error
	// 中文：EnqueueAfter 声明该接口需要实现的行为。
	// English: EnqueueAfter declares behavior required by this interface.
	EnqueueAfter(ctx context.Context, name string, payload any, delay time.Duration) (*Task, error)
	// 中文：Start 声明该接口需要实现的行为。
	// English: Start declares behavior required by this interface.
	Start(ctx context.Context)
	// 中文：Close 声明该接口需要实现的行为。
	// English: Close declares behavior required by this interface.
	Close() error
	// 中文：Stats 声明该接口需要实现的行为。
	// English: Stats declares behavior required by this interface.
	Stats() Stats
}

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
type Config struct {
	// 中文：Workers 保存当前结构中的配置或数据值。
	// English: Workers stores a configuration or data value for this struct.
	Workers int
	// 中文：Buffer 保存当前结构中的配置或数据值。
	// English: Buffer stores a configuration or data value for this struct.
	Buffer int
	// 中文：MaxRetries 保存当前结构中的配置或数据值。
	// English: MaxRetries stores a configuration or data value for this struct.
	MaxRetries int
	// 中文：RetryDelay 保存当前结构中的配置或数据值。
	// English: RetryDelay stores a configuration or data value for this struct.
	RetryDelay time.Duration
}

// 中文：Stats 定义当前包使用的数据结构或接口。
// English: Stats defines a data structure or interface used by this package.
type Stats struct {
	// 中文：Submitted 保存当前结构中的配置或数据值。
	// English: Submitted stores a configuration or data value for this struct.
	Submitted uint64
	// 中文：Succeeded 保存当前结构中的配置或数据值。
	// English: Succeeded stores a configuration or data value for this struct.
	Succeeded uint64
	// 中文：Failed 保存当前结构中的配置或数据值。
	// English: Failed stores a configuration or data value for this struct.
	Failed uint64
	// 中文：Retried 保存当前结构中的配置或数据值。
	// English: Retried stores a configuration or data value for this struct.
	Retried uint64
	// 中文：Dropped 保存当前结构中的配置或数据值。
	// English: Dropped stores a configuration or data value for this struct.
	Dropped uint64
}

// 中文：MemoryQueue 定义当前包使用的数据结构或接口。
// English: MemoryQueue defines a data structure or interface used by this package.
type MemoryQueue struct {
	// 中文：cfg 保存当前结构中的配置或数据值。
	// English: cfg stores a configuration or data value for this struct.
	cfg Config
	// 中文：tasks 保存当前结构中的配置或数据值。
	// English: tasks stores a configuration or data value for this struct.
	tasks chan *Task
	// 中文：handlers 保存当前结构中的配置或数据值。
	// English: handlers stores a configuration or data value for this struct.
	handlers map[string]Handler

	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
	// 中文：startMu 保存当前结构中的配置或数据值。
	// English: startMu stores a configuration or data value for this struct.
	startMu sync.Mutex
	// 中文：wg 保存当前结构中的配置或数据值。
	// English: wg stores a configuration or data value for this struct.
	wg sync.WaitGroup
	// 中文：ctx 保存当前结构中的配置或数据值。
	// English: ctx stores a configuration or data value for this struct.
	ctx context.Context
	// 中文：cancel 保存当前结构中的配置或数据值。
	// English: cancel stores a configuration or data value for this struct.
	cancel context.CancelFunc
	// 中文：started 保存当前结构中的配置或数据值。
	// English: started stores a configuration or data value for this struct.
	started bool
	// 中文：errorHandler 保存当前结构中的配置或数据值。
	// English: errorHandler stores a configuration or data value for this struct.
	errorHandler ErrorHandler

	// 中文：nextID 保存当前结构中的配置或数据值。
	// English: nextID stores a configuration or data value for this struct.
	nextID atomic.Uint64
	// 中文：submitted 保存当前结构中的配置或数据值。
	// English: submitted stores a configuration or data value for this struct.
	submitted atomic.Uint64
	// 中文：succeeded 保存当前结构中的配置或数据值。
	// English: succeeded stores a configuration or data value for this struct.
	succeeded atomic.Uint64
	// 中文：failed 保存当前结构中的配置或数据值。
	// English: failed stores a configuration or data value for this struct.
	failed atomic.Uint64
	// 中文：retried 保存当前结构中的配置或数据值。
	// English: retried stores a configuration or data value for this struct.
	retried atomic.Uint64
	// 中文：dropped 保存当前结构中的配置或数据值。
	// English: dropped stores a configuration or data value for this struct.
	dropped atomic.Uint64
}

// 中文：NewMemoryQueue 创建并返回对应组件实例。
// English: NewMemoryQueue creates and returns the corresponding component instance.
func NewMemoryQueue(cfg Config) *MemoryQueue {
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}
	if cfg.Buffer <= 0 {
		cfg.Buffer = 1024
	}
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 0
	}
	if cfg.RetryDelay <= 0 {
		cfg.RetryDelay = time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &MemoryQueue{
		cfg:      cfg,
		tasks:    make(chan *Task, cfg.Buffer),
		handlers: make(map[string]Handler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// 中文：Register 执行当前包中的对应流程。
// English: Register executes the corresponding workflow in this package.
func (q *MemoryQueue) Register(name string, handler Handler) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrTaskName
	}
	if handler == nil {
		return ErrHandlerRequired
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[name] = handler
	return nil
}

// 中文：SetErrorHandler 执行当前包中的对应流程。
// English: SetErrorHandler executes the corresponding workflow in this package.
func (q *MemoryQueue) SetErrorHandler(handler ErrorHandler) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.errorHandler = handler
}

// 中文：Enqueue 执行当前包中的对应流程。
// English: Enqueue executes the corresponding workflow in this package.
func (q *MemoryQueue) Enqueue(ctx context.Context, task *Task) error {
	task, err := q.prepareTask(task)
	if err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	now := time.Now()
	if !task.RunAt.IsZero() && task.RunAt.After(now) {
		if err := q.schedule(ctx, task, time.Until(task.RunAt)); err != nil {
			return err
		}
		q.submitted.Add(1)
		return nil
	}

	if err := q.enqueueReady(ctx, task); err != nil {
		return err
	}
	q.submitted.Add(1)
	return nil
}

// 中文：EnqueueAfter 执行当前包中的对应流程。
// English: EnqueueAfter executes the corresponding workflow in this package.
func (q *MemoryQueue) EnqueueAfter(ctx context.Context, name string, payload any, delay time.Duration) (*Task, error) {
	task, err := q.prepareTask(&Task{
		Name:    name,
		Payload: payload,
	})
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if delay > 0 {
		task.RunAt = time.Now().Add(delay)
	}

	if !task.RunAt.IsZero() && task.RunAt.After(time.Now()) {
		if err := q.schedule(ctx, task, time.Until(task.RunAt)); err != nil {
			return nil, err
		}
		q.submitted.Add(1)
		return task, nil
	}

	if err := q.enqueueReady(ctx, task); err != nil {
		return nil, err
	}
	q.submitted.Add(1)
	return task, nil
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (q *MemoryQueue) Start(ctx context.Context) {
	q.startMu.Lock()
	defer q.startMu.Unlock()
	if q.started {
		return
	}
	if ctx != nil {
		q.ctx, q.cancel = context.WithCancel(ctx)
	}
	q.started = true

	for i := 0; i < q.cfg.Workers; i++ {
		q.wg.Add(1)
		go q.worker()
	}
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (q *MemoryQueue) Close() error {
	q.startMu.Lock()
	if q.cancel != nil {
		q.cancel()
	}
	q.startMu.Unlock()
	q.wg.Wait()
	return nil
}

// 中文：Stats 执行当前包中的对应流程。
// English: Stats executes the corresponding workflow in this package.
func (q *MemoryQueue) Stats() Stats {
	return Stats{
		Submitted: q.submitted.Load(),
		Succeeded: q.succeeded.Load(),
		Failed:    q.failed.Load(),
		Retried:   q.retried.Load(),
		Dropped:   q.dropped.Load(),
	}
}

// 中文：prepareTask 执行当前包中的对应流程。
// English: prepareTask executes the corresponding workflow in this package.
func (q *MemoryQueue) prepareTask(task *Task) (*Task, error) {
	if task == nil {
		return nil, ErrTaskRequired
	}
	name := strings.TrimSpace(task.Name)
	if name == "" {
		return nil, ErrTaskName
	}

	cp := *task
	cp.Name = name
	if cp.ID == "" {
		cp.ID = q.newID()
	}
	if len(cp.Metadata) > 0 {
		cp.Metadata = cloneStringMap(cp.Metadata)
	}
	return &cp, nil
}

// 中文：newID 执行当前包中的对应流程。
// English: newID executes the corresponding workflow in this package.
func (q *MemoryQueue) newID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), q.nextID.Add(1))
}

// 中文：worker 执行当前包中的对应流程。
// English: worker executes the corresponding workflow in this package.
func (q *MemoryQueue) worker() {
	defer q.wg.Done()
	for {
		select {
		case <-q.ctx.Done():
			return
		case task := <-q.tasks:
			q.handle(task)
		}
	}
}

// 中文：handle 执行当前包中的对应流程。
// English: handle executes the corresponding workflow in this package.
func (q *MemoryQueue) handle(task *Task) {
	handler, errorHandler := q.lookup(task.Name)
	if handler == nil {
		q.failed.Add(1)
		q.notifyError(context.Background(), task, fmt.Errorf("%w: %s", ErrHandlerNotFound, task.Name), errorHandler)
		return
	}

	task.Attempts++
	err := handler(q.ctx, task)
	if err == nil {
		q.succeeded.Add(1)
		return
	}

	maxRetries := q.maxRetries(task)
	if task.Attempts <= maxRetries {
		q.retried.Add(1)
		delay := q.cfg.RetryDelay * time.Duration(task.Attempts)
		scheduleErr := q.schedule(q.ctx, task, delay)
		if scheduleErr == nil {
			return
		}
		err = errors.Join(err, scheduleErr)
	}

	q.failed.Add(1)
	q.notifyError(q.ctx, task, err, errorHandler)
}

// 中文：lookup 执行当前包中的对应流程。
// English: lookup executes the corresponding workflow in this package.
func (q *MemoryQueue) lookup(name string) (Handler, ErrorHandler) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.handlers[name], q.errorHandler
}

// 中文：maxRetries 执行当前包中的对应流程。
// English: maxRetries executes the corresponding workflow in this package.
func (q *MemoryQueue) maxRetries(task *Task) int {
	if task.MaxRetries > 0 {
		return task.MaxRetries
	}
	return q.cfg.MaxRetries
}

// 中文：schedule 执行当前包中的对应流程。
// English: schedule executes the corresponding workflow in this package.
func (q *MemoryQueue) schedule(ctx context.Context, task *Task, delay time.Duration) error {
	if delay <= 0 {
		return q.enqueueReady(ctx, task)
	}

	select {
	case <-q.ctx.Done():
		q.dropped.Add(1)
		return ErrClosed
	default:
	}

	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			q.dropped.Add(1)
		case <-q.ctx.Done():
			q.dropped.Add(1)
		case <-timer.C:
			if err := q.enqueueReady(context.Background(), task); err != nil {
				q.dropped.Add(1)
				errorHandler := q.currentErrorHandler()
				q.notifyError(context.Background(), task, err, errorHandler)
			}
		}
	}()
	return nil
}

// 中文：enqueueReady 执行当前包中的对应流程。
// English: enqueueReady executes the corresponding workflow in this package.
func (q *MemoryQueue) enqueueReady(ctx context.Context, task *Task) error {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-q.ctx.Done():
		q.dropped.Add(1)
		return ErrClosed
	default:
	}

	select {
	case q.tasks <- task:
		return nil
	case <-ctx.Done():
		q.dropped.Add(1)
		return ctx.Err()
	case <-q.ctx.Done():
		q.dropped.Add(1)
		return ErrClosed
	}
}

// 中文：currentErrorHandler 执行当前包中的对应流程。
// English: currentErrorHandler executes the corresponding workflow in this package.
func (q *MemoryQueue) currentErrorHandler() ErrorHandler {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.errorHandler
}

// 中文：notifyError 执行当前包中的对应流程。
// English: notifyError executes the corresponding workflow in this package.
func (q *MemoryQueue) notifyError(ctx context.Context, task *Task, err error, handler ErrorHandler) {
	if handler != nil {
		handler(ctx, task, err)
	}
}

// 中文：cloneStringMap 执行当前包中的对应流程。
// English: cloneStringMap executes the corresponding workflow in this package.
func cloneStringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
