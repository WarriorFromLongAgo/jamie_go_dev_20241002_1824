package types

type PageReq struct {
	PageNum  uint64 `json:"pageNum"`
	PageSize uint64 `json:"pageSize"`
}

type PageResp struct {
	PageNum   uint64 `json:"pageNum,string"`
	PageSize  uint64 `json:"pageSize,string"`
	TotalPage uint64 `json:"totalPage,string"`
}

type GenericPageReq[T any] struct {
	PageReq
	Req T `json:"req"`
}

type GenericPageResp[T any] struct {
	PageResp
	List []T `json:"list"`
}

type SortType string

const (
	SortTypeAsc  SortType = "asc"
	SortTypeDesc SortType = "desc"
)
