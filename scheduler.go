package jobkicker

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/google/uuid"
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
			DoneJobs:    make(map[string]time.Time),
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

func (scheduler *Scheduler) runJob(job *Job, jobId string) {
	scheduler.jobQueue.Lock()
	if doneTime, ok := scheduler.jobQueue.DoneJobs[jobId]; ok {
		scheduler.logger.Errorf("Job with id [%s] already executed at %v", jobId, doneTime)
	}

	// theoretically shouldn't happen but handeled just in case
	if _, ok := scheduler.jobQueue.PendingJobs[jobId]; !ok {
		scheduler.logger.Errorf("Job with id [%s] isn't scheduled", jobId)
	}
	scheduler.jobQueue.Unlock()
	fn := reflect.ValueOf(job.Fn)
	params := make([]reflect.Value, len(job.Args))
	for idx, param := range job.Args {
		params[idx] = reflect.ValueOf(param)
	}
	for {
		select {
		case <-job.Timer.GetWaiter():
			fn.Call(params)
			scheduler.jobQueue.Lock()
			if job.JobType == Once {
				delete(scheduler.jobQueue.PendingJobs, jobId)
			}
			executionTime := time.Now()
			scheduler.jobQueue.DoneJobs[jobId] = executionTime
			scheduler.jobQueue.Unlock()
			scheduler.logger.Infof("job with id [%s] executed in %v", jobId, executionTime)
			if job.JobType != Periodically {
				return
			}
		case <-job.cxt.Done():
			scheduler.jobQueue.Lock()
			delete(scheduler.jobQueue.PendingJobs, jobId)
			defer scheduler.jobQueue.Unlock()
			scheduler.logger.Infof("job with id [%s] cancelled successfully in %v", jobId, time.Now())
			return

		}
	}
}

func (scheduler *Scheduler) CancelJob(jobId string) error {
	scheduler.jobQueue.Lock()
	defer scheduler.jobQueue.Unlock()
	if doneTime, ok := scheduler.jobQueue.DoneJobs[jobId]; ok {

		err := fmt.Errorf(
			"Job with id [%s] can't be cancelled because it's already executed at %v",
			jobId, doneTime)

		scheduler.logger.Error(err.Error())
		return err
	}

	// theoretically shouldn't happen but handeled just in case
	if _, ok := scheduler.jobQueue.PendingJobs[jobId]; !ok {
		err := fmt.Errorf("Job with id [%s] isn't scheduled", jobId)
		scheduler.logger.Error(err.Error())
		return err
	}
	scheduler.jobQueue.PendingJobs[jobId].cancelFunc()
	return nil

}

func (scheduler *Scheduler) KickOnceAfter(delay time.Time, fn interface{}, args ...interface{}) (jobID string) {
	jobID = uuid.New().String()
	var timer *KickerTimer
	timer.InitiateNew(time.Duration(delay.Nanosecond()))
	context, cancelFunc := context.WithCancel(context.Background())
	job := &Job{
		JobType:    Once,
		Fn:         fn,
		Args:       args,
		Timer:      timer,
		cxt:        context,
		cancelFunc: cancelFunc,
	}
	scheduler.jobQueue.Lock()
	scheduler.jobQueue.PendingJobs[jobID] = job
	scheduler.jobQueue.Unlock()

	go scheduler.runJob(job, jobID)

	return jobID
}
