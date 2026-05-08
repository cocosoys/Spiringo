package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/spiringo/spiringo/internal/modules/notification/model"
	"github.com/spiringo/spiringo/internal/modules/notification/repository"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
type Config struct {
	// 中文：Events 保存当前结构中的配置或数据值。
	// English: Events stores a configuration or data value for this struct.
	Events []string `yaml:"events" mapstructure:"events"`
	// 中文：Webhook 保存当前结构中的配置或数据值。
	// English: Webhook stores a configuration or data value for this struct.
	Webhook WebhookConfig `yaml:"webhook" mapstructure:"webhook"`
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email EmailConfig `yaml:"email" mapstructure:"email"`
	// 中文：Inbox 保存当前结构中的配置或数据值。
	// English: Inbox stores a configuration or data value for this struct.
	Inbox InboxConfig `yaml:"inbox" mapstructure:"inbox"`
}

// 中文：WebhookConfig 定义当前包使用的数据结构或接口。
// English: WebhookConfig defines a data structure or interface used by this package.
type WebhookConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：URLs 保存当前结构中的配置或数据值。
	// English: URLs stores a configuration or data value for this struct.
	URLs []string `yaml:"urls" mapstructure:"urls"`
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout string `yaml:"timeout" mapstructure:"timeout"`
	// 中文：Headers 保存当前结构中的配置或数据值。
	// English: Headers stores a configuration or data value for this struct.
	Headers map[string]string `yaml:"headers" mapstructure:"headers"`
}

// 中文：EmailConfig 定义当前包使用的数据结构或接口。
// English: EmailConfig defines a data structure or interface used by this package.
type EmailConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Host 保存当前结构中的配置或数据值。
	// English: Host stores a configuration or data value for this struct.
	Host string `yaml:"host" mapstructure:"host"`
	// 中文：Port 保存当前结构中的配置或数据值。
	// English: Port stores a configuration or data value for this struct.
	Port int `yaml:"port" mapstructure:"port"`
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `yaml:"username" mapstructure:"username"`
	// 中文：Password 保存当前结构中的配置或数据值。
	// English: Password stores a configuration or data value for this struct.
	Password string `yaml:"password" mapstructure:"password"`
	// 中文：From 保存当前结构中的配置或数据值。
	// English: From stores a configuration or data value for this struct.
	From string `yaml:"from" mapstructure:"from"`
	// 中文：To 保存当前结构中的配置或数据值。
	// English: To stores a configuration or data value for this struct.
	To []string `yaml:"to" mapstructure:"to"`
}

// 中文：InboxConfig 定义当前包使用的数据结构或接口。
// English: InboxConfig defines a data structure or interface used by this package.
type InboxConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
}

// 中文：Message 定义当前包使用的数据结构或接口。
// English: Message defines a data structure or interface used by this package.
type Message struct {
	// 中文：Event 保存当前结构中的配置或数据值。
	// English: Event stores a configuration or data value for this struct.
	Event string `json:"event"`
	// 中文：Severity 保存当前结构中的配置或数据值。
	// English: Severity stores a configuration or data value for this struct.
	Severity string `json:"severity"`
	// 中文：Subject 保存当前结构中的配置或数据值。
	// English: Subject stores a configuration or data value for this struct.
	Subject string `json:"subject"`
	// 中文：Content 保存当前结构中的配置或数据值。
	// English: Content stores a configuration or data value for this struct.
	Content string `json:"content"`
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `json:"tenant_id,omitempty"`
	// 中文：RecipientID 保存当前结构中的配置或数据值。
	// English: RecipientID stores a configuration or data value for this struct.
	RecipientID string `json:"recipient_id,omitempty"`
	// 中文：Payload 保存当前结构中的配置或数据值。
	// English: Payload stores a configuration or data value for this struct.
	Payload map[string]any `json:"payload,omitempty"`
	// 中文：Timestamp 保存当前结构中的配置或数据值。
	// English: Timestamp stores a configuration or data value for this struct.
	Timestamp time.Time `json:"timestamp"`
}

// 中文：InboxFilter 定义当前包使用的数据结构或接口。
// English: InboxFilter defines a data structure or interface used by this package.
type InboxFilter struct {
	// 中文：Page 保存当前结构中的配置或数据值。
	// English: Page stores a configuration or data value for this struct.
	Page int
	// 中文：PageSize 保存当前结构中的配置或数据值。
	// English: PageSize stores a configuration or data value for this struct.
	PageSize int
	// 中文：Event 保存当前结构中的配置或数据值。
	// English: Event stores a configuration or data value for this struct.
	Event string
	// 中文：RecipientID 保存当前结构中的配置或数据值。
	// English: RecipientID stores a configuration or data value for this struct.
	RecipientID string
	// 中文：UnreadOnly 保存当前结构中的配置或数据值。
	// English: UnreadOnly stores a configuration or data value for this struct.
	UnreadOnly bool
}

// 中文：InboxItem 定义当前包使用的数据结构或接口。
// English: InboxItem defines a data structure or interface used by this package.
type InboxItem struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"id"`
	// 中文：Event 保存当前结构中的配置或数据值。
	// English: Event stores a configuration or data value for this struct.
	Event string `json:"event"`
	// 中文：Severity 保存当前结构中的配置或数据值。
	// English: Severity stores a configuration or data value for this struct.
	Severity string `json:"severity"`
	// 中文：RecipientID 保存当前结构中的配置或数据值。
	// English: RecipientID stores a configuration or data value for this struct.
	RecipientID string `json:"recipient_id,omitempty"`
	// 中文：Subject 保存当前结构中的配置或数据值。
	// English: Subject stores a configuration or data value for this struct.
	Subject string `json:"subject"`
	// 中文：Content 保存当前结构中的配置或数据值。
	// English: Content stores a configuration or data value for this struct.
	Content string `json:"content,omitempty"`
	// 中文：Payload 保存当前结构中的配置或数据值。
	// English: Payload stores a configuration or data value for this struct.
	Payload map[string]any `json:"payload,omitempty"`
	// 中文：ReadAt 保存当前结构中的配置或数据值。
	// English: ReadAt stores a configuration or data value for this struct.
	ReadAt *time.Time `json:"read_at,omitempty"`
	// 中文：CreatedAt 保存当前结构中的配置或数据值。
	// English: CreatedAt stores a configuration or data value for this struct.
	CreatedAt time.Time `json:"created_at"`
}

// 中文：Service 定义当前包使用的数据结构或接口。
// English: Service defines a data structure or interface used by this package.
type Service struct {
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：repo 保存当前结构中的配置或数据值。
	// English: repo stores a configuration or data value for this struct.
	repo *repository.NotificationRepository
	// 中文：httpClient 保存当前结构中的配置或数据值。
	// English: httpClient stores a configuration or data value for this struct.
	httpClient *http.Client
	// 中文：enabledEvents 保存当前结构中的配置或数据值。
	// English: enabledEvents stores a configuration or data value for this struct.
	enabledEvents map[string]struct{}
}

// 中文：New 创建并返回对应组件实例。
// English: New creates and returns the corresponding component instance.
func New(config Config, repos ...*repository.NotificationRepository) *Service {
	timeout := 5 * time.Second
	if config.Webhook.Timeout != "" {
		if parsed, err := time.ParseDuration(config.Webhook.Timeout); err == nil && parsed > 0 {
			timeout = parsed
		}
	}

	events := config.Events
	if len(events) == 0 {
		events = []string{"payment.failed", "tenant.suspended"}
	}
	enabledEvents := make(map[string]struct{}, len(events))
	for _, eventName := range events {
		if eventName != "" {
			enabledEvents[eventName] = struct{}{}
		}
	}

	var repo *repository.NotificationRepository
	if len(repos) > 0 {
		repo = repos[0]
	}

	return &Service{
		config:        config,
		repo:          repo,
		httpClient:    &http.Client{Timeout: timeout},
		enabledEvents: enabledEvents,
	}
}

// 中文：Notify 执行当前包中的对应流程。
// English: Notify executes the corresponding workflow in this package.
func (s *Service) Notify(ctx context.Context, msg Message) error {
	if !s.acceptsEvent(msg.Event) {
		return nil
	}
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	var errs []error
	if s.config.Inbox.Enabled {
		if err := s.saveInbox(ctx, msg); err != nil {
			errs = append(errs, err)
		}
	}
	if s.config.Webhook.Enabled {
		if err := s.sendWebhook(ctx, msg); err != nil {
			errs = append(errs, err)
		}
	}
	if s.config.Email.Enabled {
		if err := s.sendEmail(msg); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// 中文：ListInbox 执行当前包中的对应流程。
// English: ListInbox executes the corresponding workflow in this package.
func (s *Service) ListInbox(ctx context.Context, filter InboxFilter) ([]*InboxItem, int64, error) {
	if s.repo == nil {
		return nil, 0, fmt.Errorf("notification inbox store is not configured")
	}
	items, total, err := s.repo.List(ctx, repository.NotificationFilter{
		Page:        filter.Page,
		PageSize:    filter.PageSize,
		Event:       filter.Event,
		RecipientID: filter.RecipientID,
		UnreadOnly:  filter.UnreadOnly,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]*InboxItem, 0, len(items))
	for _, item := range items {
		out = append(out, inboxItemFromModel(item))
	}
	return out, total, nil
}

// 中文：MarkRead 执行当前包中的对应流程。
// English: MarkRead executes the corresponding workflow in this package.
func (s *Service) MarkRead(ctx context.Context, id, recipientID string) error {
	if s.repo == nil {
		return fmt.Errorf("notification inbox store is not configured")
	}
	return s.repo.MarkRead(ctx, id, recipientID)
}

// 中文：acceptsEvent 执行当前包中的对应流程。
// English: acceptsEvent executes the corresponding workflow in this package.
func (s *Service) acceptsEvent(eventName string) bool {
	if len(s.enabledEvents) == 0 {
		return true
	}
	_, ok := s.enabledEvents[eventName]
	return ok
}

// 中文：sendWebhook 执行当前包中的对应流程。
// English: sendWebhook executes the corresponding workflow in this package.
func (s *Service) sendWebhook(ctx context.Context, msg Message) error {
	if len(s.config.Webhook.URLs) == 0 {
		return fmt.Errorf("notification webhook enabled but no urls configured")
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal notification webhook payload: %w", err)
	}

	var errs []error
	for _, target := range s.config.Webhook.URLs {
		if target == "" {
			continue
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
		if err != nil {
			errs = append(errs, fmt.Errorf("create webhook request %s: %w", target, err))
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		for key, value := range s.config.Webhook.Headers {
			req.Header.Set(key, value)
		}

		resp, err := s.httpClient.Do(req)
		if err != nil {
			errs = append(errs, fmt.Errorf("send webhook %s: %w", target, err))
			continue
		}
		_ = resp.Body.Close()
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			errs = append(errs, fmt.Errorf("send webhook %s: status %d", target, resp.StatusCode))
		}
	}
	return errors.Join(errs...)
}

// 中文：saveInbox 执行当前包中的对应流程。
// English: saveInbox executes the corresponding workflow in this package.
func (s *Service) saveInbox(ctx context.Context, msg Message) error {
	if s.repo == nil {
		return fmt.Errorf("notification inbox enabled but store is not configured")
	}
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("marshal notification inbox payload: %w", err)
	}
	item := &model.Notification{
		Event:       strings.TrimSpace(msg.Event),
		Severity:    firstNonEmpty(strings.TrimSpace(msg.Severity), "info"),
		RecipientID: strings.TrimSpace(msg.RecipientID),
		Subject:     firstNonEmpty(strings.TrimSpace(msg.Subject), strings.TrimSpace(msg.Event)),
		Content:     strings.TrimSpace(msg.Content),
		Payload:     string(payload),
	}
	item.TenantID = strings.TrimSpace(msg.TenantID)
	return s.repo.Create(ctx, item)
}

// 中文：inboxItemFromModel 执行当前包中的对应流程。
// English: inboxItemFromModel executes the corresponding workflow in this package.
func inboxItemFromModel(item *model.Notification) *InboxItem {
	out := &InboxItem{
		ID:          item.ID,
		Event:       item.Event,
		Severity:    item.Severity,
		RecipientID: item.RecipientID,
		Subject:     item.Subject,
		Content:     item.Content,
		ReadAt:      item.ReadAt,
		CreatedAt:   item.CreatedAt,
	}
	if item.Payload != "" && item.Payload != "null" {
		var payload map[string]any
		if err := json.Unmarshal([]byte(item.Payload), &payload); err == nil {
			out.Payload = payload
		}
	}
	return out
}

// 中文：sendEmail 执行当前包中的对应流程。
// English: sendEmail executes the corresponding workflow in this package.
func (s *Service) sendEmail(msg Message) error {
	cfg := s.config.Email
	if cfg.Host == "" || cfg.Port == 0 || cfg.From == "" || len(cfg.To) == 0 {
		return fmt.Errorf("notification email enabled but smtp host, port, from, or recipients are missing")
	}

	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	var auth smtp.Auth
	if cfg.Username != "" || cfg.Password != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	subject := sanitizeHeader(firstNonEmpty(msg.Subject, msg.Event))
	body := strings.TrimSpace(msg.Content)
	if body == "" {
		bodyBytes, _ := json.MarshalIndent(msg, "", "  ")
		body = string(bodyBytes)
	}

	raw := strings.Join([]string{
		"From: " + cfg.From,
		"To: " + strings.Join(cfg.To, ","),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	if err := smtp.SendMail(addr, auth, cfg.From, cfg.To, []byte(raw)); err != nil {
		return fmt.Errorf("send email notification: %w", err)
	}
	return nil
}

// 中文：sanitizeHeader 执行当前包中的对应流程。
// English: sanitizeHeader executes the corresponding workflow in this package.
func sanitizeHeader(value string) string {
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	return strings.TrimSpace(value)
}

// 中文：firstNonEmpty 执行当前包中的对应流程。
// English: firstNonEmpty executes the corresponding workflow in this package.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
