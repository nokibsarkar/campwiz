package routes

type ResponseSingle[T any] struct {
	Data T `json:"data"`
}
type ResponseList[T any] struct {
	Data          []T    `json:"data"`
	ContinueToken string `json:"continueToken"`
	PreviousToken string `json:"previousToken"`
}
type ResponseError struct {
	Detail string `json:"detail"`
}
