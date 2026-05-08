package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 中文：ElasticsearchConfig 定义当前包使用的数据结构或接口。
// English: ElasticsearchConfig defines a data structure or interface used by this package.
type ElasticsearchConfig struct {
	// 中文：Endpoint 保存当前结构中的配置或数据值。
	// English: Endpoint stores a configuration or data value for this struct.
	Endpoint string
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string
	// 中文：Password 保存当前结构中的配置或数据值。
	// English: Password stores a configuration or data value for this struct.
	Password string
	// 中文：Index 保存当前结构中的配置或数据值。
	// English: Index stores a configuration or data value for this struct.
	Index string
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout time.Duration
}

// 中文：Elasticsearch 定义当前包使用的数据结构或接口。
// English: Elasticsearch defines a data structure or interface used by this package.
type Elasticsearch struct {
	// 中文：endpoint 保存当前结构中的配置或数据值。
	// English: endpoint stores a configuration or data value for this struct.
	endpoint string
	// 中文：username 保存当前结构中的配置或数据值。
	// English: username stores a configuration or data value for this struct.
	username string
	// 中文：password 保存当前结构中的配置或数据值。
	// English: password stores a configuration or data value for this struct.
	password string
	// 中文：index 保存当前结构中的配置或数据值。
	// English: index stores a configuration or data value for this struct.
	index string
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *http.Client
}

// 中文：NewElasticsearch 创建并返回对应组件实例。
// English: NewElasticsearch creates and returns the corresponding component instance.
func NewElasticsearch(cfg ElasticsearchConfig) (*Elasticsearch, error) {
	endpoint := strings.TrimRight(strings.TrimSpace(cfg.Endpoint), "/")
	if endpoint == "" {
		return nil, fmt.Errorf("elasticsearch endpoint is required")
	}
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return nil, fmt.Errorf("parse elasticsearch endpoint: %w", err)
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Elasticsearch{
		endpoint: endpoint,
		username: cfg.Username,
		password: cfg.Password,
		index:    cfg.Index,
		client:   &http.Client{Timeout: timeout},
	}, nil
}

// 中文：IndexDocument 执行当前包中的对应流程。
// English: IndexDocument executes the corresponding workflow in this package.
func (e *Elasticsearch) IndexDocument(ctx context.Context, index, id string, document any) error {
	index = e.resolveIndex(index)
	if index == "" {
		return fmt.Errorf("elasticsearch index is required")
	}

	method := http.MethodPost
	parts := []string{index, "_doc"}
	if id != "" {
		method = http.MethodPut
		parts = append(parts, id)
	}
	resp, err := e.doJSON(ctx, method, parts, document)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return readError(resp)
	}
	return nil
}

// 中文：Index 执行当前包中的对应流程。
// English: Index executes the corresponding workflow in this package.
func (e *Elasticsearch) Index(ctx context.Context, index string, id string, doc any) error {
	return e.IndexDocument(ctx, index, id, doc)
}

// 中文：BulkIndex 执行当前包中的对应流程。
// English: BulkIndex executes the corresponding workflow in this package.
func (e *Elasticsearch) BulkIndex(ctx context.Context, index string, docs []any) error {
	index = e.resolveIndex(index)
	if index == "" {
		return fmt.Errorf("elasticsearch index is required")
	}
	if len(docs) == 0 {
		return nil
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	for i, doc := range docs {
		if doc == nil {
			return fmt.Errorf("bulk document %d is nil", i)
		}
		id, source := bulkDocumentSource(doc)
		if source == nil {
			return fmt.Errorf("bulk document %d source is nil", i)
		}
		action := map[string]any{"_index": index}
		if id != "" {
			action["_id"] = id
		}
		if err := encoder.Encode(map[string]any{"index": action}); err != nil {
			return fmt.Errorf("encode bulk action %d: %w", i, err)
		}
		if err := encoder.Encode(source); err != nil {
			return fmt.Errorf("encode bulk document %d: %w", i, err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.url("_bulk"), &buf)
	if err != nil {
		return fmt.Errorf("create elasticsearch bulk request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	if e.username != "" || e.password != "" {
		req.SetBasicAuth(e.username, e.password)
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return readError(resp)
	}

	var result esBulkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode elasticsearch bulk response: %w", err)
	}
	if result.Errors {
		return result.error()
	}
	return nil
}

// 中文：DeleteDocument 执行当前包中的对应流程。
// English: DeleteDocument executes the corresponding workflow in this package.
func (e *Elasticsearch) DeleteDocument(ctx context.Context, index, id string) error {
	index = e.resolveIndex(index)
	if index == "" || id == "" {
		return fmt.Errorf("elasticsearch index and id are required")
	}

	resp, err := e.doJSON(ctx, http.MethodDelete, []string{index, "_doc", id}, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return readError(resp)
	}
	return nil
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (e *Elasticsearch) Delete(ctx context.Context, index string, id string) error {
	return e.DeleteDocument(ctx, index, id)
}

// 中文：Search 执行当前包中的对应流程。
// English: Search executes the corresponding workflow in this package.
func (e *Elasticsearch) Search(ctx context.Context, query *Query) (*Result, error) {
	if query == nil {
		query = &Query{}
	}
	index := e.resolveIndex(query.Index)
	if index == "" {
		return nil, fmt.Errorf("elasticsearch index is required")
	}

	body := buildSearchBody(query)
	resp, err := e.doJSON(ctx, http.MethodPost, []string{index, "_search"}, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, readError(resp)
	}

	var raw esSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode elasticsearch search response: %w", err)
	}
	return raw.toResult(), nil
}

// 中文：Aggregation 执行当前包中的对应流程。
// English: Aggregation executes the corresponding workflow in this package.
func (e *Elasticsearch) Aggregation(ctx context.Context, query *Query, aggFields []string) (map[string]any, error) {
	if query == nil {
		query = &Query{}
	}
	query.Size = 0
	if query.Aggregations == nil {
		query.Aggregations = make(map[string]Aggregation, len(aggFields))
	}
	for _, field := range aggFields {
		if field == "" {
			continue
		}
		query.Aggregations[field] = Aggregation{Type: "terms", Field: field, Size: 10}
	}
	result, err := e.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	return result.Aggregations, nil
}

// 中文：resolveIndex 执行当前包中的对应流程。
// English: resolveIndex executes the corresponding workflow in this package.
func (e *Elasticsearch) resolveIndex(index string) string {
	if index != "" {
		return index
	}
	return e.index
}

// 中文：doJSON 执行当前包中的对应流程。
// English: doJSON executes the corresponding workflow in this package.
func (e *Elasticsearch) doJSON(ctx context.Context, method string, parts []string, body any) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal elasticsearch request: %w", err)
		}
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, e.url(parts...), reader)
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if e.username != "" || e.password != "" {
		req.SetBasicAuth(e.username, e.password)
	}
	return e.client.Do(req)
}

// 中文：url 执行当前包中的对应流程。
// English: url executes the corresponding workflow in this package.
func (e *Elasticsearch) url(parts ...string) string {
	u := e.endpoint
	for _, part := range parts {
		u += "/" + pathEscape(part)
	}
	return u
}

// 中文：pathEscape 执行当前包中的对应流程。
// English: pathEscape executes the corresponding workflow in this package.
func pathEscape(value string) string {
	return url.PathEscape(value)
}

// 中文：buildSearchBody 执行当前包中的对应流程。
// English: buildSearchBody executes the corresponding workflow in this package.
func buildSearchBody(query *Query) map[string]any {
	size := query.Size
	if size <= 0 {
		size = 10
	}
	body := map[string]any{
		"from":  maxInt(query.From, 0),
		"size":  size,
		"query": buildQuery(query),
	}
	if len(query.Sort) > 0 {
		sort := make([]map[string]map[string]string, 0, len(query.Sort))
		for _, item := range query.Sort {
			if item.Field == "" {
				continue
			}
			order := strings.ToLower(item.Order)
			if order != "asc" {
				order = "desc"
			}
			sort = append(sort, map[string]map[string]string{
				item.Field: map[string]string{"order": order},
			})
		}
		if len(sort) > 0 {
			body["sort"] = sort
		}
	}
	if len(query.Aggregations) > 0 {
		aggs := make(map[string]any, len(query.Aggregations))
		for name, agg := range query.Aggregations {
			if agg.Field == "" {
				continue
			}
			aggType := strings.ToLower(agg.Type)
			if aggType == "" {
				aggType = "terms"
			}
			switch aggType {
			case "terms":
				size := agg.Size
				if size <= 0 {
					size = 10
				}
				aggs[name] = map[string]any{"terms": map[string]any{"field": agg.Field, "size": size}}
			case "avg", "sum", "min", "max":
				aggs[name] = map[string]any{aggType: map[string]any{"field": agg.Field}}
			}
		}
		if len(aggs) > 0 {
			body["aggs"] = aggs
		}
	}
	return body
}

// 中文：buildQuery 执行当前包中的对应流程。
// English: buildQuery executes the corresponding workflow in this package.
func buildQuery(query *Query) map[string]any {
	var must []any
	if query.Keyword != "" {
		fields := query.Fields
		if len(fields) == 0 {
			fields = []string{"*"}
		}
		must = append(must, map[string]any{
			"multi_match": map[string]any{
				"query":  query.Keyword,
				"fields": fields,
			},
		})
	}

	filters := make([]any, 0, len(query.Filters))
	for key, value := range query.Filters {
		filters = append(filters, map[string]any{"term": map[string]any{key: value}})
	}
	if len(must) == 0 && len(filters) == 0 {
		return map[string]any{"match_all": map[string]any{}}
	}
	boolQuery := map[string]any{}
	if len(must) > 0 {
		boolQuery["must"] = must
	}
	if len(filters) > 0 {
		boolQuery["filter"] = filters
	}
	return map[string]any{"bool": boolQuery}
}

// 中文：readError 执行当前包中的对应流程。
// English: readError executes the corresponding workflow in this package.
func readError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	message := strings.TrimSpace(string(body))
	if message == "" {
		message = resp.Status
	}
	return fmt.Errorf("elasticsearch request failed: status=%d body=%s", resp.StatusCode, message)
}

// 中文：esSearchResponse 定义当前包使用的数据结构或接口。
// English: esSearchResponse defines a data structure or interface used by this package.
type esSearchResponse struct {
	// 中文：Hits 保存当前结构中的配置或数据值。
	// English: Hits stores a configuration or data value for this struct.
	Hits struct {
		Total json.RawMessage `json:"total"`
		Hits  []struct {
			ID     string         `json:"_id"`
			Score  float64        `json:"_score"`
			Source map[string]any `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
	// 中文：Aggregations 保存当前结构中的配置或数据值。
	// English: Aggregations stores a configuration or data value for this struct.
	Aggregations map[string]any `json:"aggregations"`
}

// 中文：toResult 执行当前包中的对应流程。
// English: toResult executes the corresponding workflow in this package.
func (r esSearchResponse) toResult() *Result {
	hits := make([]Hit, 0, len(r.Hits.Hits))
	for _, hit := range r.Hits.Hits {
		hits = append(hits, Hit{ID: hit.ID, Score: hit.Score, Source: hit.Source})
	}
	return &Result{
		Total:        parseTotal(r.Hits.Total),
		Hits:         hits,
		Aggregations: r.Aggregations,
	}
}

// 中文：parseTotal 执行当前包中的对应流程。
// English: parseTotal executes the corresponding workflow in this package.
func parseTotal(raw json.RawMessage) int64 {
	if len(raw) == 0 {
		return 0
	}
	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		return n
	}
	var obj struct {
		Value int64 `json:"value"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil {
		return obj.Value
	}
	asString := strings.Trim(string(raw), `"`)
	parsed, _ := strconv.ParseInt(asString, 10, 64)
	return parsed
}

// 中文：maxInt 执行当前包中的对应流程。
// English: maxInt executes the corresponding workflow in this package.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 中文：bulkDocumentSource 执行当前包中的对应流程。
// English: bulkDocumentSource executes the corresponding workflow in this package.
func bulkDocumentSource(doc any) (string, any) {
	switch v := doc.(type) {
	case BulkDocument:
		return v.ID, v.Source
	case *BulkDocument:
		if v == nil {
			return "", nil
		}
		return v.ID, v.Source
	case DocumentIDer:
		return v.DocumentID(), doc
	case map[string]any:
		if id, ok := stringField(v, "_id"); ok {
			return id, doc
		}
		if id, ok := stringField(v, "id"); ok {
			return id, doc
		}
	}
	return "", doc
}

// 中文：stringField 执行当前包中的对应流程。
// English: stringField executes the corresponding workflow in this package.
func stringField(values map[string]any, key string) (string, bool) {
	raw, ok := values[key]
	if !ok {
		return "", false
	}
	value, ok := raw.(string)
	return value, ok && value != ""
}

// 中文：esBulkResponse 定义当前包使用的数据结构或接口。
// English: esBulkResponse defines a data structure or interface used by this package.
type esBulkResponse struct {
	// 中文：Errors 保存当前结构中的配置或数据值。
	// English: Errors stores a configuration or data value for this struct.
	Errors bool `json:"errors"`
	// 中文：Items 保存当前结构中的配置或数据值。
	// English: Items stores a configuration or data value for this struct.
	Items []map[string]bulkItem `json:"items"`
}

// 中文：bulkItem 定义当前包使用的数据结构或接口。
// English: bulkItem defines a data structure or interface used by this package.
type bulkItem struct {
	// 中文：Index 保存当前结构中的配置或数据值。
	// English: Index stores a configuration or data value for this struct.
	Index string `json:"_index"`
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"_id"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status int `json:"status"`
	// 中文：Error 保存当前结构中的配置或数据值。
	// English: Error stores a configuration or data value for this struct.
	Error any `json:"error"`
}

// 中文：error 执行当前包中的对应流程。
// English: error executes the corresponding workflow in this package.
func (r esBulkResponse) error() error {
	for _, item := range r.Items {
		for action, result := range item {
			if result.Error != nil || result.Status >= http.StatusBadRequest {
				return fmt.Errorf("elasticsearch bulk %s failed: index=%s id=%s status=%d error=%v", action, result.Index, result.ID, result.Status, result.Error)
			}
		}
	}
	return fmt.Errorf("elasticsearch bulk request failed")
}
