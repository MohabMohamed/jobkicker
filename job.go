package jobkicker

import (
	"context"
	"sync"
	"time"
)

type JobType int8

const (
	Once JobType = iota
	Periodically
)

type Job struct {
	JobType    JobType
	Fn         interface{}
	Args       []interface{}
	Timer      ITimer
	cxt        context.Context
	cancelFunc context.CancelFunc
}

type JobQueue struct {
	sync.Mutex
	PendingJobs map[string]*Job
	DoneJobs    map[string]time.Time //Done jobs with it's last execution time
}
