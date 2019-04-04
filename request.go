package raggort

import (
	"github.com/miratronix/lingo"
	"strconv"
	"strings"
	"sync"
)

var requestIDLock = &sync.Mutex{}
var requestID = 0

// Request defines a rapport request
type Request struct {
	ID          string       `json:"_rq"`
	HTTPRequest *HTTPRequest `json:"_b"`
}

// Method gets the request method
func (r *Request) Method() string {
	return r.HTTPRequest.Method
}

// URL gets the request URL
func (r *Request) URL() string {
	return r.HTTPRequest.URL
}

// Route gets a route representation of the request
func (r *Request) Route() string {
	return r.Method() + " " + r.URL()
}

// RouteWithoutPrefix gets the route string without the prefix
func (r *Request) RouteWithoutPrefix() string {

	// No trimmed URL set up yet, get the prefix to force the trim
	if r.HTTPRequest.trimmedURL == "" {
		r.GetPrefix()
	}

	return r.Method() + " " + r.HTTPRequest.trimmedURL
}

// GetPrefix gets the first part of the route and removes it from the underlying route
func (r *Request) GetPrefix() string {

	// Remove the leading slash
	url := strings.TrimPrefix(r.HTTPRequest.URL, "/")

	// Split it into two parts
	parts := strings.SplitN(url, "/", 2)

	// Set the URL to the second part
	r.HTTPRequest.trimmedURL = "/" + parts[1]

	// Return the first part
	return "/" + parts[0]
}

// Body gets the body of the request, parsing it into a validatable format
func (r *Request) Body(destination Validatable) *HTTPResponse {

	// If the body is not a map, we can't decode it
	cast, ok := r.HTTPRequest.Body.(map[string]interface{})
	if !ok {
		return NewBadRequestError("Please supply an object body")
	}

	// Decode the request body
	err := lingo.Map.Decode(cast, destination)
	if err != nil {
		return NewBadRequestError("Failed to decode request")
	}

	// Validate and return
	return destination.Validate()
}

// Int converts the request body to an integer
func (r *Request) Int() (int, *HTTPResponse) {
	cast, ok := r.HTTPRequest.Body.(int)
	if !ok {
		return 0, NewBadRequestError("Please supply an integer")
	}
	return cast, nil
}

// Bool converts the request body to a boolean
func (r *Request) Bool() (bool, *HTTPResponse) {
	cast, ok := r.HTTPRequest.Body.(bool)
	if !ok {
		return false, NewBadRequestError("Please supply a boolean")
	}
	return cast, nil
}

// String converts the request body to a string
func (r *Request) String() (string, *HTTPResponse) {
	cast, ok := r.HTTPRequest.Body.(string)
	if !ok {
		return "", NewBadRequestError("Please supply a string")
	}
	return cast, nil
}

// Float converts the request body to a float value
func (r *Request) Float() (float64, *HTTPResponse) {

	// Try it as a float 64
	castFloat64, ok := r.HTTPRequest.Body.(float64)
	if ok {
		return castFloat64, nil
	}

	// Try it as a float 32
	castFloat32, ok := r.HTTPRequest.Body.(float32)
	if ok {
		return float64(castFloat32), nil
	}

	// Try it as an int
	castInt, ok := r.HTTPRequest.Body.(int)
	if ok {
		return float64(castInt), nil
	}

	return 0, NewBadRequestError("Please supply a float")
}

// newEmptyRequest constructs an empty request, which is useful when we need a request to serialize bytes into
func newEmptyRequest() *Request {
	return &Request{
		HTTPRequest: &HTTPRequest{},
	}
}

// newRequest constructs a new request
func newRequest(method string, url string, body interface{}) *Request {
	return &Request{
		ID: createRequestID(),
		HTTPRequest: &HTTPRequest{
			Method: method,
			URL:    url,
			Body:   body,
		},
	}
}

// createRequestID creates an ID for a request
func createRequestID() string {
	requestIDLock.Lock()
	defer requestIDLock.Unlock()

	// Don't let request IDs get too large
	if requestID >= maxRequestID {
		requestID = 0
	}

	requestID++
	return strconv.Itoa(requestID)
}
