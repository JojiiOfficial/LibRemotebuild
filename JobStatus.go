package libremotebuild

// JobState a state of a job
type JobState uint8

// ...
const (
	JobWaiting JobState = iota
	JobCancelled
	JobFailed
	JobRunning
	JobDone
	JobPaused
)

func (js JobState) String() string {
	switch js {
	case JobWaiting:
		return "Waiting"
	case JobCancelled:
		return "Cancelled"
	case JobFailed:
		return "Failed"
	case JobRunning:
		return "Running"
	case JobDone:
		return "Done"
	case JobPaused:
		return "Paused"
	}

	return "<invaild>"
}
