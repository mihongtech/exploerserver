package resp

import "net/http"

type RequestErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *RequestErr) Error() string {
	return err.Message
}

var (
	InternalServerErr = &RequestErr{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
	MethodNotAllowedErr = &RequestErr{
		Code:    http.StatusMethodNotAllowed,
		Message: http.StatusText(http.StatusMethodNotAllowed),
	}
	BadRequestErr = &RequestErr{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
	}
	NotFoundErr = &RequestErr{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}
)
