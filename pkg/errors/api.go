// error: Custom error type
// using method:
// err := errors.New(10000, "Example error msg")
package errors

type Error interface {
	Error() string
	// Msg() string
	// Code() int64
	// RefineError(err ...interface{}) *error
}

func New(code int64, msg string) Error {
	return new(code, msg)
}
