package routes

type ResponseSingle[T any] struct {
	Data T `json:"data"`
}
type ResponseList[T any] struct {
	Data          []T    `json:"data"`
	ContinueToken string `json:"next"`
	PreviousToken string `json:"prev"`
}
type ResponseError struct {
	Detail string `json:"detail"`
}
