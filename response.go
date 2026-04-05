package gokick

type Response[T any] struct {
	Result T
}

// Pagination carries cursor-based pagination metadata from Kick API v2 list endpoints.
type Pagination struct {
	NextCursor string `json:"next_cursor"`
}

// PaginatedResponse is like Response but includes pagination (e.g. GET /public/v2/categories).
type PaginatedResponse[T any] struct {
	Result     T          `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type EmptyResponse struct{}
