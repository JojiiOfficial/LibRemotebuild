package libremotebuild

// LibRB data required in all requests
type LibRB struct {
	Config *RequestConfig
}

// NewLibDM create new libDM "class"
func NewLibDM(config *RequestConfig) *LibRB {
	return &LibRB{
		Config: config,
	}
}

// Request do a request using libdm
func (libdm LibRB) Request(ep Endpoint, payload, response interface{}, authorized bool) (*RestRequestResponse, error) {
	req := libdm.NewRequest(ep, payload)
	if authorized {
		req.WithAuthFromConfig()
	}
	resp, err := req.Do(response)

	if err != nil || resp.Status == ResponseError {
		return nil, NewErrorFromResponse(resp, err)
	}

	return resp, nil
}
