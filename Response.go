package libremotebuild

import (
	"net/http"
	"time"
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
	Token string `json:"token"`
}

// RestRequestResponse the response of a rest call
type RestRequestResponse struct {
	HTTPCode int
	Status   ResponseStatus
	Message  string
	Headers  *http.Header
	Response *http.Response
}

// StringResponse response containing only one string
type StringResponse struct {
	String string `json:"content"`
}

// StringSliceResponse response containing only one string slice
type StringSliceResponse struct {
	Slice []string `json:"slice"`
}

// AddJobResponse response for adding a job
type AddJobResponse struct {
	ID       uint `json:"id"`
	Position int  `json:"pos"`
}

// JobInfo info of job
type JobInfo struct {
	ID           uint          `json:"id"`
	Info         string        `json:"info"`
	Position     uint          `json:"pos"`
	BuildType    JobType       `json:"jobtype"`
	UploadType   UploadType    `json:"uploadtype"`
	Status       JobState      `json:"state"`
	RunningSince time.Time     `json:"rs,omitempty"`
	Duration     time.Duration `json:"dr"`
}

// ListJobsResponse list of queued jobs
type ListJobsResponse struct {
	Jobs []JobInfo `json:"jobs"`
}

// SortByJob sort jobs
type SortByJob []JobInfo

func (a SortByJob) Len() int           { return len(a) }
func (a SortByJob) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByJob) Less(i, j int) bool { return a[i].ID < a[j].ID }
