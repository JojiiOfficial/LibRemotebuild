package libremotebuild

// JobType type of job
type JobType uint8

// ...
const (
	JobNoBuild JobType = iota
	JobAUR
)

func (jt JobType) String() string {
	switch jt {
	case JobNoBuild:
		return "NoJob"
	case JobAUR:
		return "buildAUR"
	}

	return ""
}

// ParseJobType parse a jobtype from string
func ParseJobType(inp string) JobType {
	switch inp {
	case "NoJob":
		return JobNoBuild
	case "buildAUR":
		return JobAUR
	}

	return JobNoBuild
}
