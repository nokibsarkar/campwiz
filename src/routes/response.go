package routes

type ResponseSingle[T any] struct {
	Data T `json:"data"`
}
type ResponseList[T any] struct {
	Data []T `json:"data"`
}
type ResponseError struct {
	Detail string `json:"detail"`
}
