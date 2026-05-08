package types

// 中文：PaginationRequest 定义当前包使用的数据结构或接口。
// English: PaginationRequest defines a data structure or interface used by this package.
// PaginationRequest 分页请求参数
type PaginationRequest struct {
	// 中文：Page 保存当前结构中的配置或数据值。
	// English: Page stores a configuration or data value for this struct.
	Page int `form:"page" binding:"omitempty,min=1" json:"page"`
	// 中文：PageSize 保存当前结构中的配置或数据值。
	// English: PageSize stores a configuration or data value for this struct.
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100" json:"page_size"`
}

// 中文：GetPage 执行当前包中的对应流程。
// English: GetPage executes the corresponding workflow in this package.
// GetPage 获取页码（默认1）
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// 中文：GetPageSize 执行当前包中的对应流程。
// English: GetPageSize executes the corresponding workflow in this package.
// GetPageSize 获取每页数量（默认20）
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// 中文：GetOffset 执行当前包中的对应流程。
// English: GetOffset executes the corresponding workflow in this package.
// GetOffset 计算偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// 中文：SortRequest 定义当前包使用的数据结构或接口。
// English: SortRequest defines a data structure or interface used by this package.
// SortRequest 排序请求参数
type SortRequest struct {
	// 中文：SortBy 保存当前结构中的配置或数据值。
	// English: SortBy stores a configuration or data value for this struct.
	SortBy string `form:"sort_by" json:"sort_by"`
	// 中文：SortOrder 保存当前结构中的配置或数据值。
	// English: SortOrder stores a configuration or data value for this struct.
	SortOrder string `form:"sort_order" json:"sort_order"` // asc | desc
}

// 中文：GetSortOrder 执行当前包中的对应流程。
// English: GetSortOrder executes the corresponding workflow in this package.
// GetSortOrder 获取排序方向（默认asc）
func (s *SortRequest) GetSortOrder() string {
	if s.SortOrder == "desc" {
		return "desc"
	}
	return "asc"
}
