package libremotebuild

import (
	"net/http"
)

// ResponseStatus the status of response
type ResponseStatus uint8

const (
	// ResponseError if there was an error
	ResponseError ResponseStatus = 0
	// ResponseSuccess if the response is successful
	ResponseSuccess ResponseStatus = 1
)

const (
	// HeaderStatus headername for status in response
	HeaderStatus string = "X-Response-Status"

	// HeaderStatusMessage headername for status in response
	HeaderStatusMessage string = "X-Response-Message"

	// HeaderContentType contenttype of response
	HeaderContentType string = "Content-Type"

	// HeaderRequest request content
	HeaderRequest string = "Request"

	// HeaderContentLength request content length
	HeaderContentLength string = "ContentLength"
)

// LoginResponse response for login
type LoginResponse struct {
	Token     string `json:"token"`
	Namespace string `json:"ns"`
}

// RestRequestResponse the response of a rest call
type RestRequestResponse struct {
	HTTPCode int
	Status   ResponseStatus
	Message  string
	Headers  *http.Header
}

// StringResponse response containing only one string
type StringResponse struct {
	String string `json:"content"`
}

// StringSliceResponse response containing only one string slice
type StringSliceResponse struct {
	Slice []string `json:"slice"`
}
