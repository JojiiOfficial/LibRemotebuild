package libremotebuild

// AURBuild build an AUR package
type AURBuild struct {
	LibRB
	args          map[string]string
	UploadType    UploadType
	DisableCcache bool
}

// NewAURBuild build an AUR package
func (Librb LibRB) NewAURBuild(packageName string) *AURBuild {
	return &AURBuild{
		LibRB: Librb,
		args: map[string]string{
			AURPackage: packageName,
		},
	}
}

// WithoutCcache disables ccache
func (aurBuild *AURBuild) WithoutCcache() *AURBuild {
	aurBuild.DisableCcache = true
	return aurBuild
}

// WithDmanager use dmnager for uplaod
func (aurBuild *AURBuild) WithDmanager(username, token, host string) {
	aurBuild.UploadType = DataManagerUploadType
	aurBuild.args[DMToken] = token
	aurBuild.args[DMUser] = username
	aurBuild.args[DMHost] = host
}

// CreateJob build AUR package
func (aurBuild *AURBuild) CreateJob() (*AddJobResponse, error) {
	return aurBuild.LibRB.AddJob(JobAUR, aurBuild.UploadType, aurBuild.args, aurBuild.DisableCcache)
}
