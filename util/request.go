package util

type CommonResponse[T any] struct {
	Body       T
	Total      *int64  `header:"X-Total-Count" doc:"Total number of items"`
	Offset     *uint64 `header:"X-Offset" doc:"Offset for pagination"`
	Limit      *uint64 `header:"X-Limit" doc:"Limit for pagination"`
	Page       *int64  `header:"X-Page" doc:"Current page number (1-indexed)"`
	TotalPages *int64  `header:"X-Total-Pages" doc:"Total number of pages"`
}

func MakeCommonResponse[T any](body T) *CommonResponse[T] {
	return &CommonResponse[T]{Body: body}
}

func MakePaginatedResponse[T any](body T, total int64, offset, limit uint64) *CommonResponse[T] {
	page := int64(1)
	if limit > 0 {
		page = int64(offset/limit) + 1
	}
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}
	return &CommonResponse[T]{
		Body:       body,
		Total:      &total,
		Offset:     &offset,
		Limit:      &limit,
		Page:       &page,
		TotalPages: &totalPages,
	}
}
