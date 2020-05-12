package libremotebuild

import "strings"

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

// ParseUploadType an uploadtype string
func ParseUploadType(s string) UploadType {
	s = strings.ToLower(strings.TrimSpace(s))

	switch s {
	case strings.ToLower(DataManagerUploadType.String()):
		return DataManagerUploadType
	}

	return NoUploadType
}
