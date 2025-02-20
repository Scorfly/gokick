package gokick

type Response[T any] struct {
	Result T
}

type EmptyResponse struct{}
