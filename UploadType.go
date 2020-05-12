package libremotebuild

// UploadType type of upload destination
type UploadType uint8

// ...
const (
	NoUploadType UploadType = iota
	DataManagerUploadType
)

func (ut UploadType) String() string {
	switch ut {
	case NoUploadType:
		return "no upload"
	case DataManagerUploadType:
		return "DataManager"
	}

	return "<invalid>"
}
