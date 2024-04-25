package errors

import "fmt"

type error struct {
	code    int64
	message string
}

func (e *error) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.message)
}

func (e *error) Msg() string {
	return e.message
}

func (e *error) Code() int64 {
	return e.code
}

func (e *error) RefineError(err ...interface{}) *error {
	return new(e.Code(), e.message+", "+fmt.Sprint(err...))
}

func new(code int64, msg string) *error {
	return &error{
		code:    code,
		message: msg,
	}
}
