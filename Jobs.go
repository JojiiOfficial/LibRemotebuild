package libremotebuild

// AddJob a job
func (librb LibRB) AddJob(jobType JobType, uploadType UploadType, args map[string]string) (*AddJobResponse, error) {
	var response AddJobResponse

	// Do http request
	resp, err := librb.NewRequest(EPJobAdd, AddJobRequest{
		Type:       jobType,
		UploadType: uploadType,
		Args:       args,
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

	return &response, nil
}

// CancelJob cancel a running or queued job
func (librb LibRB) CancelJob(jobID uint) error {
	// Do http request
	resp, err := librb.NewRequest(EPJobCancel, CancelJobRequest{
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
