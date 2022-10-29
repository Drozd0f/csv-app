package errors

import "fmt"

type ErrorWithMessage struct {
	Err error
	Msg string
}

func (ewm *ErrorWithMessage) Unwrap() error {
	return ewm.Err
}

func (ewm *ErrorWithMessage) Error() string {
	return fmt.Sprintf("%s: %s", ewm.Err.Error(), ewm.Msg)
}

func NewErrorWithMessage(err error, msg string) ErrorWithMessage {
	return ErrorWithMessage{
		Err: err,
		Msg: msg,
	}
}
