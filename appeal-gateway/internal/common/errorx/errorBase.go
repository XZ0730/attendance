package errorx

import "encoding/json"

const defaultCode = 1001

type CodeError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CodeErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewCodeError(status int, message string) error {
	ce := &CodeErrorResponse{Status: status, Message: message}
	res, _ := json.Marshal(ce)

	return &CodeError{Status: status, Message: string(res)}
}

func NewDefaultError(message string) error {
	return NewCodeError(defaultCode, message)
}

func (e *CodeError) Error() string {
	return e.Message
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Status:  e.Status,
		Message: e.Message,
	}
}
