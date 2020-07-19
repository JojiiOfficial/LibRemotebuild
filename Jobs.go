package libremotebuild

import (
	"fmt"
	"sort"
	"time"
)

// AddJob a job
func (librb LibRB) AddJob(jobType JobType, uploadType UploadType, args map[string]string, disableCcache bool) (*AddJobResponse, error) {
	var response AddJobResponse

	// Do http request
	resp, err := librb.NewRequest(EPJobAdd, AddJobRequest{
		Type:          jobType,
		UploadType:    uploadType,
		Args:          args,
		DisableCcache: disableCcache,
	}).WithAuthFromConfig().
		WithMethod(PUT).
		Do(&response)

	// Return new error on ... error
	if err != nil || resp.Status == ResponseError {
		return nil, NewErrorFromResponse(resp, err)
	}

	return &response, nil
}

// ListJobs list running jobs
func (librb LibRB) ListJobs() (*ListJobsResponse, error) {
	var response ListJobsResponse

	// Do http request
	resp, err := librb.NewRequest(EPJobs, nil).WithAuthFromConfig().
		WithMethod(GET).
		Do(&response)

	// Return new error on ... error
	if err != nil || resp.Status == ResponseError {
		return nil, NewErrorFromResponse(resp, err)
	}

	// Sort jobs
	sort.Sort(SortByJob(response.Jobs))

	return &response, nil
}

// SetJobState pauses a running or queued job
func (librb LibRB) SetJobState(jobID uint, state JobState) error {
	switch state {
	case JobPaused, JobRunning:
	default:
		return fmt.Errorf("Invalid state to set job to")
	}

	endpoint := EPJobPause
	if state == JobRunning {
		endpoint = EPJobResume
	}

	// Do http request
	resp, err := librb.NewRequest(endpoint, JobRequest{
		JobID: jobID,
	}).WithAuthFromConfig().
		WithMethod(PUT).
		Do(nil)

	// Return new error on ... error
	if err != nil || resp.Status == ResponseError {
		return NewErrorFromResponse(resp, err)
	}

	return nil
}

// CancelJob cancel a running or queued job
func (librb LibRB) CancelJob(jobID uint) error {
	// Do http request
	resp, err := librb.NewRequest(EPJobCancel, JobRequest{
		JobID: jobID,
	}).WithAuthFromConfig().
		WithMethod(POST).
		Do(nil)

	// Return new error on ... error
	if err != nil || resp.Status == ResponseError {
		return NewErrorFromResponse(resp, err)
	}

	return nil
}

// Logs for a job
func (librb LibRB) Logs(jobID uint, since time.Time) (*RestRequestResponse, error) {
	// Do http request
	resp, err := librb.NewRequest(EPJobLogs, JobLogsRequest{
		Since: since,
		JobID: jobID,
	}).WithAuthFromConfig().
		WithNoBodyClose().
		WithMethod(GET).
		Do(nil)

	// Return new error on ... error
	if err != nil || resp.Status == ResponseError {
		return nil, NewErrorFromResponse(resp, err)
	}

	return resp, nil
}
