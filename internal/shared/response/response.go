package response

type Response[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Errors  any    `json:"errors"`
	Data    T      `json:"data"`
}
