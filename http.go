package raggort

import "net/http"

const (
	get          = "GET"
	post         = "POST"
	put          = "PUT"
	patch        = "PATCH"
	del          = "DELETE"
	maxRequestID = 10000
)

// HTTPRequest defines a rapport-http request
type HTTPRequest struct {
	Method     string      `json:"_m"`
	URL        string      `json:"_u"`
	Body       interface{} `json:"_b"`
	trimmedURL string      `json:"-"`
}

// NewGet constructs a rapport get request
func NewGet(url string) *Request {
	return newRequest(get, url, nil)
}

// NewPost constructs a rapport post request
func NewPost(url string, body interface{}) *Request {
	return newRequest(post, url, body)
}

// NewPut constructs a rapport put request
func NewPut(url string, body interface{}) *Request {
	return newRequest(put, url, body)
}

// NewPatch constructs a rapport patch request
func NewPatch(url string, body interface{}) *Request {
	return newRequest(patch, url, body)
}

// NewDelete constructs a rapport delete request
func NewDelete(url string) *Request {
	return newRequest(del, url, nil)
}

// HTTPResponse defines a http response body
type HTTPResponse struct {
	RawStatus int         `json:"_s"`
	RawBody   interface{} `json:"_b"`
}

// NewHTTPResponse creates a new HTTP error response
func NewHTTPResponse() *HTTPResponse {
	return &HTTPResponse{
		RawStatus: http.StatusOK,
		RawBody:   map[string]interface{}{},
	}
}

// Status sets the status for the response in a fluent manner
func (h *HTTPResponse) Status(status int) *HTTPResponse {
	h.RawStatus = status
	return h
}

// Body sets the response body using a fluent API
func (h *HTTPResponse) Body(body interface{}) *HTTPResponse {
	h.RawBody = body
	return h
}

// IsOk determines if the status code of the response is ok (2xx)
func (h *HTTPResponse) IsOk() bool {
	return h.RawStatus >= 200 && h.RawStatus <= 299
}

// Response converts this HTTP response to a full rapport response
func (h *HTTPResponse) Response(request *Request) *Response {

	// Not a 2xx status code, return an error response
	if !h.IsOk() {
		return &Response{
			ID:    request.ID,
			Error: h,
		}
	}

	return &Response{
		ID:   request.ID,
		Body: h,
	}
}
