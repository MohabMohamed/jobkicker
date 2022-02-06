package jobkicker

import (
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

type Scheduler struct {
	jobQueue *JobQueue
	logger   *log.Logger
}

func NewScheduler(loggerOutput *io.Writer, loggerFormatter *log.Formatter) *Scheduler {
	scheduler := &Scheduler{
		jobQueue: &JobQueue{
			PendingJobs: make(map[string]*Job),
			DoneJobs:    make(map[string]*time.Time),
		},
		logger: log.New()}

	if loggerFormatter != nil {
		scheduler.logger.SetFormatter(*loggerFormatter)
	}
	if loggerOutput != nil {
		scheduler.logger.SetOutput(*loggerOutput)
	}
	return scheduler
}
