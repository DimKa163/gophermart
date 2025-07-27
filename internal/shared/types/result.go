package types

type AppCode string

const (
	Created   AppCode = "created"
	Updated   AppCode = "updated"
	NotFound  AppCode = "not_found"
	Duplicate AppCode = "duplicate"
	NoChange  AppCode = "no_change"
	Problem   AppCode = "problem"
)

type AppResult[T any] struct {
	Code    AppCode
	Payload T
	Error   error
}
