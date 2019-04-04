package raggort

import (
	"net/http"
)

// Response defines a rapport response structure
type Response struct {
	ID    string        `json:"_rs"`
	Body  *HTTPResponse `json:"_b"`
	Error *HTTPResponse `json:"_e"`
}

// IsError determines if the response is an error
func (h *Response) IsError() bool {
	return h.Body == nil || h.Body.RawStatus == 0
}

// newEmptyResponse constructs an empty response
func newEmptyResponse() *Response {
	return &Response{
		Body:  &HTTPResponse{},
		Error: &HTTPResponse{},
	}
}

// newTimeoutResponse creates a new HTTP timeout response
func newTimeoutResponse(id string) *Response {
	return &Response{
		ID: id,
		Error: &HTTPResponse{
			RawStatus: http.StatusGatewayTimeout,
			RawBody:   map[string]interface{}{},
		},
	}
}
