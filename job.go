package jobkicker

import (
	"context"
	"sync"
	"time"
)

type JobType int8

// the type of the frequency the function will run with is only once or periodically
const (
	Once JobType = iota
	Periodically
)

// Job is the representation of the job to run
type Job struct {
	JobType    JobType
	Fn         interface{}
	Args       []interface{}
	Timer      ITimer
	cxt        context.Context
	cancelFunc context.CancelFunc
}

// the collections that hold the the pending jobs to be executed and the executed jobs
// jobs with the last time it ran.
//
// Note: periodic jobs stays in PendingJobs map after execution unless got canceled and Done
// jobs in this case holds the last execution time.
type JobQueue struct {
	sync.Mutex
	PendingJobs map[string]*Job
	DoneJobs    map[string]time.Time //Done jobs with it's last execution time
}
