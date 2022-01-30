package customerr

import "encoding/json"

var NotFoundErr *CustomError = NewCustomError(nil, "not found")
var NotAcceptable *CustomError = NewCustomError(nil, "not acceptable (maybe the username is not unique)")

type CustomError struct {
	Err     error  `json:"-"`
	Message string `json:"message,omitempty"`
}

func NewCustomError(err error, message string) *CustomError {
	return &CustomError{Err: err, Message: message}
}

func (e *CustomError) Error() string {
	return e.Message
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

func (e *CustomError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}
