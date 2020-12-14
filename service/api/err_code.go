/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/12/1
   Description :
-------------------------------------------------
*/

package api

var (
	OK                    = &Error{Code: 0, Message: "ok"}
	ServiceInternalError  = &Error{Code: 1, Message: "service internal error"}
	ParamError            = &Error{Code: 2, Message: "param error"}
	AuthorizationRequired = &Error{Code: 3, Message: "authorization required"}
	AuthorizationError    = &Error{Code: 4, Message: "authorization error"}
)

type Error struct {
	Code    int
	Message string
	Err     error
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}
func (e Error) WithMessage(msg string) Error {
	e.Message = msg
	return e
}
func (e Error) WithError(err error) Error {
	e.Err = err
	return e
}

func decodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch v := err.(type) {
	case *Error:
		return v.Code, v.Message
	case Error:
		return v.Code, v.Message
	}

	return ServiceInternalError.Code, ServiceInternalError.Message
}
