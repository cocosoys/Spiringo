package pagination

import "testing"

// 中文：TestNormalizeAndOffset 验证相关行为符合预期。
// English: TestNormalizeAndOffset verifies the related behavior.
func TestNormalizeAndOffset(t *testing.T) {
	q := Normalize(Query{Page: -1, PageSize: 999, SortOrder: "DESC"})
	if q.Page != 1 || q.PageSize != 100 || q.SortOrder != "desc" {
		t.Fatalf("unexpected normalized query: %+v", q)
	}
	if got := Offset(Query{Page: 3, PageSize: 10}); got != 20 {
		t.Fatalf("offset = %d, want 20", got)
	}
}

// 中文：TestOrderClauseUsesAllowedColumns 验证相关行为符合预期。
// English: TestOrderClauseUsesAllowedColumns verifies the related behavior.
func TestOrderClauseUsesAllowedColumns(t *testing.T) {
	order := OrderClause(Query{SortBy: "created_at", SortOrder: "desc"}, map[string]string{
		"created_at": "created_at",
	}, "id")
	if order != "created_at desc" {
		t.Fatalf("order = %q", order)
	}

	fallback := OrderClause(Query{SortBy: "unsafe;drop", SortOrder: "desc"}, map[string]string{}, "id")
	if fallback != "id desc" {
		t.Fatalf("fallback order = %q", fallback)
	}
}

// 中文：TestNewResult 验证相关行为符合预期。
// English: TestNewResult verifies the related behavior.
func TestNewResult(t *testing.T) {
	result := NewResult([]int{1, 2}, 5, Query{Page: 2, PageSize: 2})
	if result.TotalPages != 3 || !result.HasNext || !result.HasPrev {
		t.Fatalf("unexpected result: %+v", result)
	}
}
