package libremotebuild

// ClearCcache clear ccache on server
func (librb LibRB) ClearCcache() (string, error) {
	resp, err := librb.NewRequest(EPCcacheClear, nil).WithAuthFromConfig().WithMethod(POST).Do(nil)
	return resp.Message, err
}

// QueryCcache get ccache stats
func (librb LibRB) QueryCcache() (StringResponse, error) {
	var resp StringResponse
	_, err := librb.NewRequest(EPCcacheStats, nil).WithAuthFromConfig().WithMethod(GET).Do(&resp)
	return resp, err
}
