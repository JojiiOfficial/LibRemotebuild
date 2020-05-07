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
)
