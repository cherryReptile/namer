package customErrors

import "github.com/pkg/errors"

type Error struct {
	Message    string
	Err        error
	StatusCode int
}

func New(message string, err error, code int) error {
	return Error{
		Message:    message,
		Err:        err,
		StatusCode: code,
	}
}

func Wrap(err error, message string) error {
	ce, ok := err.(Error)
	if !ok {
		return Error{
			Message: err.Error(),
			Err:     errors.Wrap(err, message),
		}
	}

	ce.Err = errors.Wrap(ce.Err, message)

	return ce
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) Cause() error {
	return errors.Cause(e.Err)
}
