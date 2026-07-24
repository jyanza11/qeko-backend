package response

type Pagination struct {
	Page         int  `json:"page"`
	PageSize     int  `json:"pageSize"`
	Total        int  `json:"total"`
	TotalPages   int  `json:"totalPages"`
	NextPage     int  `json:"nextPage"`
	PreviousPage int  `json:"previousPage"`
	HasNext      bool `json:"hasNext"`
	HasPrevious  bool `json:"hasPrevious"`
}

type PaginationResponse[T any] struct {
	Status     string     `json:"status"`
	Message    string     `json:"message"`
	Errors     []any      `json:"errors"`
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}
