package search

import "context"

// 中文：Query 定义当前包使用的数据结构或接口。
// English: Query defines a data structure or interface used by this package.
type Query struct {
	// 中文：Index 保存当前结构中的配置或数据值。
	// English: Index stores a configuration or data value for this struct.
	Index string
	// 中文：Keyword 保存当前结构中的配置或数据值。
	// English: Keyword stores a configuration or data value for this struct.
	Keyword string
	// 中文：Fields 保存当前结构中的配置或数据值。
	// English: Fields stores a configuration or data value for this struct.
	Fields []string
	// 中文：Filters 保存当前结构中的配置或数据值。
	// English: Filters stores a configuration or data value for this struct.
	Filters map[string]any
	// 中文：From 保存当前结构中的配置或数据值。
	// English: From stores a configuration or data value for this struct.
	From int
	// 中文：Size 保存当前结构中的配置或数据值。
	// English: Size stores a configuration or data value for this struct.
	Size int
	// 中文：Sort 保存当前结构中的配置或数据值。
	// English: Sort stores a configuration or data value for this struct.
	Sort []Sort
	// 中文：Aggregations 保存当前结构中的配置或数据值。
	// English: Aggregations stores a configuration or data value for this struct.
	Aggregations map[string]Aggregation
}

// 中文：Sort 定义当前包使用的数据结构或接口。
// English: Sort defines a data structure or interface used by this package.
type Sort struct {
	// 中文：Field 保存当前结构中的配置或数据值。
	// English: Field stores a configuration or data value for this struct.
	Field string
	// 中文：Order 保存当前结构中的配置或数据值。
	// English: Order stores a configuration or data value for this struct.
	Order string
}

// 中文：SortField 定义当前包使用的数据结构或接口。
// English: SortField defines a data structure or interface used by this package.
type SortField = Sort

// 中文：SearchQuery 定义当前包使用的数据结构或接口。
// English: SearchQuery defines a data structure or interface used by this package.
type SearchQuery = Query

// 中文：SearchResult 定义当前包使用的数据结构或接口。
// English: SearchResult defines a data structure or interface used by this package.
type SearchResult = Result

// 中文：Aggregation 定义当前包使用的数据结构或接口。
// English: Aggregation defines a data structure or interface used by this package.
type Aggregation struct {
	// 中文：Type 保存当前结构中的配置或数据值。
	// English: Type stores a configuration or data value for this struct.
	Type string
	// 中文：Field 保存当前结构中的配置或数据值。
	// English: Field stores a configuration or data value for this struct.
	Field string
	// 中文：Size 保存当前结构中的配置或数据值。
	// English: Size stores a configuration or data value for this struct.
	Size int
}

// 中文：Hit 定义当前包使用的数据结构或接口。
// English: Hit defines a data structure or interface used by this package.
type Hit struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"id"`
	// 中文：Score 保存当前结构中的配置或数据值。
	// English: Score stores a configuration or data value for this struct.
	Score float64 `json:"score"`
	// 中文：Source 保存当前结构中的配置或数据值。
	// English: Source stores a configuration or data value for this struct.
	Source map[string]any `json:"source"`
}

// 中文：BulkDocument 定义当前包使用的数据结构或接口。
// English: BulkDocument defines a data structure or interface used by this package.
type BulkDocument struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string
	// 中文：Source 保存当前结构中的配置或数据值。
	// English: Source stores a configuration or data value for this struct.
	Source any
}

// 中文：DocumentIDer 定义当前包使用的数据结构或接口。
// English: DocumentIDer defines a data structure or interface used by this package.
type DocumentIDer interface {
	// 中文：DocumentID 声明该接口需要实现的行为。
	// English: DocumentID declares behavior required by this interface.
	DocumentID() string
}

// 中文：Result 定义当前包使用的数据结构或接口。
// English: Result defines a data structure or interface used by this package.
type Result struct {
	// 中文：Total 保存当前结构中的配置或数据值。
	// English: Total stores a configuration or data value for this struct.
	Total int64 `json:"total"`
	// 中文：Hits 保存当前结构中的配置或数据值。
	// English: Hits stores a configuration or data value for this struct.
	Hits []Hit `json:"hits"`
	// 中文：Aggregations 保存当前结构中的配置或数据值。
	// English: Aggregations stores a configuration or data value for this struct.
	Aggregations map[string]any `json:"aggregations,omitempty"`
}

// 中文：Engine 定义当前包使用的数据结构或接口。
// English: Engine defines a data structure or interface used by this package.
type Engine interface {
	// 中文：IndexDocument 声明该接口需要实现的行为。
	// English: IndexDocument declares behavior required by this interface.
	IndexDocument(ctx context.Context, index, id string, document any) error
	// 中文：BulkIndex 声明该接口需要实现的行为。
	// English: BulkIndex declares behavior required by this interface.
	BulkIndex(ctx context.Context, index string, docs []any) error
	// 中文：DeleteDocument 声明该接口需要实现的行为。
	// English: DeleteDocument declares behavior required by this interface.
	DeleteDocument(ctx context.Context, index, id string) error
	// 中文：Search 声明该接口需要实现的行为。
	// English: Search declares behavior required by this interface.
	Search(ctx context.Context, query *Query) (*Result, error)
	// 中文：Aggregation 声明该接口需要实现的行为。
	// English: Aggregation declares behavior required by this interface.
	Aggregation(ctx context.Context, query *Query, aggFields []string) (map[string]any, error)
}

// 中文：Search 定义当前包使用的数据结构或接口。
// English: Search defines a data structure or interface used by this package.
type Search interface {
	// 中文：Index 声明该接口需要实现的行为。
	// English: Index declares behavior required by this interface.
	Index(ctx context.Context, index string, id string, doc any) error
	// 中文：BulkIndex 声明该接口需要实现的行为。
	// English: BulkIndex declares behavior required by this interface.
	BulkIndex(ctx context.Context, index string, docs []any) error
	// 中文：Delete 声明该接口需要实现的行为。
	// English: Delete declares behavior required by this interface.
	Delete(ctx context.Context, index string, id string) error
	// 中文：Search 声明该接口需要实现的行为。
	// English: Search declares behavior required by this interface.
	Search(ctx context.Context, query *SearchQuery) (*SearchResult, error)
	// 中文：Aggregation 声明该接口需要实现的行为。
	// English: Aggregation declares behavior required by this interface.
	Aggregation(ctx context.Context, query *SearchQuery, aggFields []string) (map[string]any, error)
}
