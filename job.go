package jobkicker

import "time"

type JobType int8

const (
	Once JobType = iota
	periodically
)

type Job struct {
	jobType JobType
	fn      interface{}
	args    []interface{}
	timer   *time.Timer
	ticker  *time.Ticker
}

type JobQueue struct {
	pendingJobs map[string]*Job
	doneJobs    map[string]*time.Time //Done jobs with it's finish time
}
