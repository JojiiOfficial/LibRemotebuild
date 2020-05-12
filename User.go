package libremotebuild

import "strings"

// Login login into the server
func (librb LibRB) Login(username, password string) (*LoginResponse, error) {
	var response LoginResponse

	// Do http request
	resp, err := librb.NewRequest(EPLogin, CredentialsRequest{
		Password:  password,
		Username:  strings.ToLower(username),
		MachineID: librb.Config.MachineID,
	}).Do(&response)

	// Return new error on ... error
	if err != nil || resp.Status == ResponseError {
		return nil, NewErrorFromResponse(resp, err)
	}

	return &response, nil
}

// Register create a new account. Return true on success
func (librb LibRB) Register(username, password string) (*RestRequestResponse, error) {
	// Do http request
	resp, err := librb.NewRequest(EPRegister, CredentialsRequest{
		Username: strings.ToLower(username),
		Password: password,
	}).Do(nil)

	if err != nil || resp.Status == ResponseError {
		return resp, NewErrorFromResponse(resp, err)
	}

	return resp, nil
}

// Ping pings a server the REST way to
// ensure it is reachable
func (librb LibRB) Ping() (*StringResponse, error) {
	var response StringResponse

	// Do ping request
	req := librb.NewRequest(EPPing, PingRequest{Payload: "ping"})
	if librb.Config.SessionToken != "" {
		req.WithAuthFromConfig()
	}
	_, err := req.Do(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
