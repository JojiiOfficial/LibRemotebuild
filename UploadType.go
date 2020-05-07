package libremotebuild

// UploadType type of upload destination
type UploadType uint8

// ...
const (
	NoUploadType UploadType = iota
	DataManagerUploadType
)
