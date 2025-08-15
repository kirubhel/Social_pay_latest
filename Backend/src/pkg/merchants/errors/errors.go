package errors

import "net/http"

type Error struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (err Error) Error() string {
	return err.Message
}

func (err Error) ErrorCode() int {
	return err.Code
}
func (err Error) ErrorType() string {
	return err.Type
}

func MapErrorToHTTPStatus(err error) int {
	switch e := err.(type) {
	case Error:
		return e.ErrorCode()
	default:
		return http.StatusInternalServerError
	}
}

func MapErrorToType(err error) string {
	switch e := err.(type) {
	case Error:
		return e.ErrorType()
	default:
		return "INTERNAL_SERVER_ERROR"
	}
}
