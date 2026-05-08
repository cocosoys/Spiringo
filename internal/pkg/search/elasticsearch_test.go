package search

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 中文：_、_ 声明当前包使用的变量。
// English: _、_ declares variables used by this package.
var (
	_ Engine = (*Elasticsearch)(nil)
	_ Search = (*Elasticsearch)(nil)
)

// 中文：TestElasticsearchSearchBuildsRequestAndParsesResponse 验证相关行为符合预期。
// English: TestElasticsearchSearchBuildsRequestAndParsesResponse verifies the related behavior.
func TestElasticsearchSearchBuildsRequestAndParsesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/articles/_search" {
			t.Fatalf("path = %s, want /articles/_search", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if body["size"].(float64) != 5 {
			t.Fatalf("size = %v, want 5", body["size"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"hits":{
				"total":{"value":1},
				"hits":[{"_id":"doc-1","_score":1.25,"_source":{"title":"Go search"}}]
			},
			"aggregations":{"status":{"buckets":[{"key":"published","doc_count":1}]}}
		}`))
	}))
	defer server.Close()

	engine, err := NewElasticsearch(ElasticsearchConfig{Endpoint: server.URL, Index: "articles"})
	if err != nil {
		t.Fatalf("NewElasticsearch returned error: %v", err)
	}
	result, err := engine.Search(context.Background(), &Query{
		Keyword: "go",
		Fields:  []string{"title"},
		Filters: map[string]any{"status": "published"},
		Size:    5,
		Aggregations: map[string]Aggregation{
			"status": {Type: "terms", Field: "status.keyword", Size: 3},
		},
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if result.Total != 1 {
		t.Fatalf("total = %d, want 1", result.Total)
	}
	if len(result.Hits) != 1 || result.Hits[0].ID != "doc-1" {
		t.Fatalf("hits = %#v, want doc-1", result.Hits)
	}
	if result.Aggregations["status"] == nil {
		t.Fatal("missing status aggregation")
	}
}

// 中文：TestElasticsearchIndexDocumentUsesDefaultIndex 验证相关行为符合预期。
// English: TestElasticsearchIndexDocumentUsesDefaultIndex verifies the related behavior.
func TestElasticsearchIndexDocumentUsesDefaultIndex(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/articles/_doc/doc-1" {
			t.Fatalf("path = %s, want /articles/_doc/doc-1", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	engine, err := NewElasticsearch(ElasticsearchConfig{Endpoint: server.URL, Index: "articles"})
	if err != nil {
		t.Fatalf("NewElasticsearch returned error: %v", err)
	}
	if err := engine.IndexDocument(context.Background(), "", "doc-1", map[string]any{"title": "Go"}); err != nil {
		t.Fatalf("IndexDocument returned error: %v", err)
	}
}

// 中文：TestElasticsearchBulkIndexWritesNDJSON 验证相关行为符合预期。
// English: TestElasticsearchBulkIndexWritesNDJSON verifies the related behavior.
func TestElasticsearchBulkIndexWritesNDJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/_bulk" {
			t.Fatalf("path = %s, want /_bulk", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "application/x-ndjson") {
			t.Fatalf("content-type = %s, want application/x-ndjson", got)
		}

		scanner := bufio.NewScanner(r.Body)
		lines := make([]map[string]any, 0, 4)
		for scanner.Scan() {
			var line map[string]any
			if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
				t.Fatalf("decode bulk line: %v", err)
			}
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			t.Fatalf("scan bulk body: %v", err)
		}
		if len(lines) != 4 {
			t.Fatalf("bulk lines = %d, want 4", len(lines))
		}
		firstAction := lines[0]["index"].(map[string]any)
		if firstAction["_index"] != "articles" || firstAction["_id"] != "doc-1" {
			t.Fatalf("first action = %#v", firstAction)
		}
		if lines[1]["title"] != "Go" {
			t.Fatalf("first source = %#v", lines[1])
		}
		secondAction := lines[2]["index"].(map[string]any)
		if secondAction["_id"] != "doc-2" {
			t.Fatalf("second action = %#v", secondAction)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":false,"items":[{"index":{"_index":"articles","_id":"doc-1","status":201}}]}`))
	}))
	defer server.Close()

	engine, err := NewElasticsearch(ElasticsearchConfig{Endpoint: server.URL, Index: "articles"})
	if err != nil {
		t.Fatalf("NewElasticsearch returned error: %v", err)
	}
	if err := engine.BulkIndex(context.Background(), "", []any{
		BulkDocument{ID: "doc-1", Source: map[string]any{"title": "Go"}},
		map[string]any{"id": "doc-2", "title": "Search"},
	}); err != nil {
		t.Fatalf("BulkIndex returned error: %v", err)
	}
}

// 中文：TestElasticsearchBulkIndexReportsItemError 验证相关行为符合预期。
// English: TestElasticsearchBulkIndexReportsItemError verifies the related behavior.
func TestElasticsearchBulkIndexReportsItemError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"errors": true,
			"items": [
				{"index":{"_index":"articles","_id":"doc-1","status":400,"error":{"type":"mapper_parsing_exception"}}}
			]
		}`))
	}))
	defer server.Close()

	engine, err := NewElasticsearch(ElasticsearchConfig{Endpoint: server.URL, Index: "articles"})
	if err != nil {
		t.Fatalf("NewElasticsearch returned error: %v", err)
	}
	if err := engine.BulkIndex(context.Background(), "", []any{BulkDocument{ID: "doc-1", Source: map[string]any{"title": "Go"}}}); err == nil {
		t.Fatal("BulkIndex returned nil error, want item error")
	}
}

// 中文：TestElasticsearchSearchReturnsStatusError 验证相关行为符合预期。
// English: TestElasticsearchSearchReturnsStatusError verifies the related behavior.
func TestElasticsearchSearchReturnsStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusBadGateway)
	}))
	defer server.Close()

	engine, err := NewElasticsearch(ElasticsearchConfig{Endpoint: server.URL, Index: "articles"})
	if err != nil {
		t.Fatalf("NewElasticsearch returned error: %v", err)
	}
	if _, err := engine.Search(context.Background(), &Query{}); err == nil {
		t.Fatal("Search returned nil error, want status error")
	}
}
