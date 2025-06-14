package dtos

type InterceptorResponse struct {
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
	Data         any    `json:"data,omitempty"`
}
