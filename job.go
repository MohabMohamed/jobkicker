package jobkicker

import "time"

type JobType int8

const (
	Once JobType = iota
	Periodically
)

type Job struct {
	JobType JobType
	Fn      interface{}
	Args    []interface{}
	Timer   *ITimer
}

type JobQueue struct {
	PendingJobs map[string]*Job
	DoneJobs    map[string]*time.Time //Done jobs with it's finish time
}
