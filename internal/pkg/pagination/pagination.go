package pagination

import (
	"fmt"
	"math"
	"strings"
)

// 中文：DefaultPage、DefaultPageSize、MaxPageSize 声明当前包使用的常量。
// English: DefaultPage、DefaultPageSize、MaxPageSize declares constants used by this package.
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// 中文：Query 定义当前包使用的数据结构或接口。
// English: Query defines a data structure or interface used by this package.
type Query struct {
	// 中文：Page 保存当前结构中的配置或数据值。
	// English: Page stores a configuration or data value for this struct.
	Page int
	// 中文：PageSize 保存当前结构中的配置或数据值。
	// English: PageSize stores a configuration or data value for this struct.
	PageSize int
	// 中文：SortBy 保存当前结构中的配置或数据值。
	// English: SortBy stores a configuration or data value for this struct.
	SortBy string
	// 中文：SortOrder 保存当前结构中的配置或数据值。
	// English: SortOrder stores a configuration or data value for this struct.
	SortOrder string
}

// 中文：Result 定义当前包使用的数据结构或接口。
// English: Result defines a data structure or interface used by this package.
type Result[T any] struct {
	// 中文：List 保存当前结构中的配置或数据值。
	// English: List stores a configuration or data value for this struct.
	List []T `json:"list"`
	// 中文：Total 保存当前结构中的配置或数据值。
	// English: Total stores a configuration or data value for this struct.
	Total int64 `json:"total"`
	// 中文：Page 保存当前结构中的配置或数据值。
	// English: Page stores a configuration or data value for this struct.
	Page int `json:"page"`
	// 中文：PageSize 保存当前结构中的配置或数据值。
	// English: PageSize stores a configuration or data value for this struct.
	PageSize int `json:"page_size"`
	// 中文：TotalPages 保存当前结构中的配置或数据值。
	// English: TotalPages stores a configuration or data value for this struct.
	TotalPages int `json:"total_pages"`
	// 中文：HasNext 保存当前结构中的配置或数据值。
	// English: HasNext stores a configuration or data value for this struct.
	HasNext bool `json:"has_next"`
	// 中文：HasPrev 保存当前结构中的配置或数据值。
	// English: HasPrev stores a configuration or data value for this struct.
	HasPrev bool `json:"has_prev"`
}

// 中文：Normalize 执行当前包中的对应流程。
// English: Normalize executes the corresponding workflow in this package.
func Normalize(q Query) Query {
	if q.Page <= 0 {
		q.Page = DefaultPage
	}
	if q.PageSize <= 0 {
		q.PageSize = DefaultPageSize
	}
	if q.PageSize > MaxPageSize {
		q.PageSize = MaxPageSize
	}
	q.SortOrder = strings.ToLower(strings.TrimSpace(q.SortOrder))
	if q.SortOrder != "desc" {
		q.SortOrder = "asc"
	}
	q.SortBy = strings.TrimSpace(q.SortBy)
	return q
}

// 中文：Offset 执行当前包中的对应流程。
// English: Offset executes the corresponding workflow in this package.
func Offset(q Query) int {
	q = Normalize(q)
	return (q.Page - 1) * q.PageSize
}

// 中文：Limit 执行当前包中的对应流程。
// English: Limit executes the corresponding workflow in this package.
func Limit(q Query) int {
	return Normalize(q).PageSize
}

// 中文：OrderClause 执行当前包中的对应流程。
// English: OrderClause executes the corresponding workflow in this package.
func OrderClause(q Query, allowed map[string]string, fallback string) string {
	q = Normalize(q)
	column := allowed[q.SortBy]
	if column == "" {
		column = fallback
	}
	if column == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", column, q.SortOrder)
}

// 中文：NewResult 创建并返回对应组件实例。
// English: NewResult creates and returns the corresponding component instance.
func NewResult[T any](list []T, total int64, q Query) Result[T] {
	q = Normalize(q)
	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(q.PageSize)))
	}
	return Result[T]{
		List:       list,
		Total:      total,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalPages: totalPages,
		HasNext:    q.Page < totalPages,
		HasPrev:    q.Page > 1 && totalPages > 0,
	}
}
