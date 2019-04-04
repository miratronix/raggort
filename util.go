package raggort

import (
	"errors"
	"github.com/miratronix/lingo"
	"strings"
)

// Router defines an interface for a rapport router
type Router interface {
	Handle(request *Request) *Response
}

// Validatable defines a interface that contains a Validate() method
type Validatable interface {
	Validate() *HTTPResponse
}

// ToRequest converts a byte slice to a request object
func ToRequest(msg []byte, encoding lingo.Encoding) (*Request, error) {
	req := newEmptyRequest()

	err := encoding.Decode(msg, req)
	if err != nil {
		return nil, err
	}

	if req.ID == "" {
		return nil, errors.New("failed to get request ID")
	}

	// Trim trailing slashes in the URL
	req.HTTPRequest.URL = strings.TrimRight(strings.TrimSpace(req.HTTPRequest.URL), "/")

	// Add a leading slash to the URL if there isn't one
	if !strings.HasPrefix(req.HTTPRequest.URL, "/") {
		req.HTTPRequest.URL = "/" + req.HTTPRequest.URL
	}

	// Uppercase the method
	req.HTTPRequest.Method = strings.TrimSpace(strings.ToUpper(req.HTTPRequest.Method))

	return req, nil
}

// ToResponse converts a byte slice to a response object
func ToResponse(msg []byte, encoding lingo.Encoding) (*Response, error) {
	res := newEmptyResponse()

	err := encoding.Decode(msg, res)
	if err != nil {
		return nil, err
	}

	if res.ID == "" {
		return nil, errors.New("failed to get response ID")
	}

	return res, nil
}
