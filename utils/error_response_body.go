package utils

type ErrorResponseBody struct {
	ClientOperation string `json:"client_operation"`
	Message         string `json:"message"`
	ContextString   string `json:"context_string"`
	RequestBody     string `json:"request_body"`
	Detail          string `json:"detail"`
}
