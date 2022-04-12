package errors

import "fmt"

type Error struct {
	err string
}

func New(message string, args ...interface{}) error {
	return Error{err: "pp: " + fmt.Sprintf(message, args...)}
}

func NewEncodeError(t interface{}) error {
	return Error{err: "pp_encode_error: " + fmt.Sprintf("Unable to encode value %+v", t)}
}

func (e Error) Error() string {
	return e.err
}
