package exception

import "jhr.com/apirelay/model"

type CustomError struct {
	Result model.Result
}

func (ce CustomError) Error() string {
	return ce.Result.Message
}
