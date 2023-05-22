package errorx

const defaultCode = 1001

type CodeError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CodeErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Status: code, Message: msg}
}

func NewDefaultError(msg string) error {
	return NewCodeError(defaultCode, msg)
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
