// jobkicker is A Golang in-process cron task scheduler that kicks (executes)
// them once in specified time or periodically.
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

type JobKicker struct {
	jobQueue *JobQueue
	logger   *log.Logger
}

// NewScheduler generates new JobKicker (scheduler) and enables to configure the logger.
//
// by passing `loggerOutput *io.Writer` like a file or any thing implements io.Writer
// to write the logs to. if nil passed it's write to `os.Stderr`
//
// by passing `loggerFormatter *log.Formatter` it allows to format the logs before logging it
// by passing any struct implements the interface:
//  type Formatter interface {
//  	Format(*Entry) ([]byte, error)
//  }
//
// https://github.com/sirupsen/logrus used for logging so you can try to pass `&log.JSONFormatter{}`
// for example to format the logs as json.
// and if nil is passed it will use `logrus.TextFormatter` by default.
func NewScheduler(loggerOutput *io.Writer, loggerFormatter *log.Formatter) *JobKicker {
	jobKicker := &JobKicker{
		jobQueue: &JobQueue{
			PendingJobs: make(map[string]*Job),
			DoneJobs:    make(map[string]time.Time),
		},
		logger: log.New()}

	if loggerFormatter != nil {
		jobKicker.logger.SetFormatter(*loggerFormatter)
	}
	if loggerOutput != nil {
		jobKicker.logger.SetOutput(*loggerOutput)
	}
	return jobKicker
}

// runJob excutes the job and removes it after running from the pendingjobs map
// if it should only run once or if it got canceled. and it saves the last time
// job ran in the doneJobs map
//
// it takes the job that it should run with it's id.
func (jobKicker *JobKicker) runJob(job *Job, jobId string) {
	jobKicker.jobQueue.RLock()
	if doneTime, ok := jobKicker.jobQueue.DoneJobs[jobId]; ok {
		jobKicker.logger.Errorf("Job with id [%s] already executed at %v", jobId, doneTime)
	}

	// theoretically shouldn't happen but handeled just in case
	if _, ok := jobKicker.jobQueue.PendingJobs[jobId]; !ok {
		jobKicker.logger.Errorf("Job with id [%s] isn't scheduled", jobId)
	}
	jobKicker.jobQueue.RUnlock()
	fn := reflect.ValueOf(job.Fn)
	params := make([]reflect.Value, len(job.Args))
	for idx, param := range job.Args {
		params[idx] = reflect.ValueOf(param)
	}
	waiter := job.Timer.GetWaiter()
	for {
		select {
		case <-waiter:
			fn.Call(params)
			jobKicker.jobQueue.Lock()
			if job.JobType == Once {
				delete(jobKicker.jobQueue.PendingJobs, jobId)
			}
			executionTime := time.Now()
			jobKicker.jobQueue.DoneJobs[jobId] = executionTime
			jobKicker.jobQueue.Unlock()
			jobKicker.logger.Infof("job with id [%s] executed in %v", jobId, executionTime)
			if job.JobType != Periodically {
				return
			}
		case <-job.cxt.Done():
			jobKicker.jobQueue.Lock()
			delete(jobKicker.jobQueue.PendingJobs, jobId)
			defer jobKicker.jobQueue.Unlock()
			jobKicker.logger.Infof("job with id [%s] cancelled successfully in %v", jobId, time.Now())
			return

		}
	}
}

// CancelJob cancels a job with it's id by using a cancel context that the job struct holds
func (jobKicker *JobKicker) CancelJob(jobId string) error {
	jobKicker.jobQueue.RLock()
	defer jobKicker.jobQueue.RUnlock()

	jobType := Once
	if job, ok := jobKicker.jobQueue.PendingJobs[jobId]; ok {
		jobType = job.JobType
	}

	if doneTime, ok := jobKicker.jobQueue.DoneJobs[jobId]; ok && jobType == Once {

		err := fmt.Errorf(
			"Job with id [%s] can't be cancelled because it's already executed at %v",
			jobId, doneTime)

		jobKicker.logger.Error(err.Error())
		return err
	}

	// theoretically shouldn't happen but handeled just in case
	if _, ok := jobKicker.jobQueue.PendingJobs[jobId]; !ok {
		err := fmt.Errorf("Job with id [%s] isn't scheduled", jobId)
		jobKicker.logger.Error(err.Error())
		return err
	}
	jobKicker.jobQueue.PendingJobs[jobId].cancelFunc()
	return nil

}

// KickOnceAfter excutes a function once after a delay and returns the job id like:

func (jobKicker *JobKicker) KickOnceAfter(delay time.Time, fn interface{}, args ...interface{}) (jobID string) {
	jobID = uuid.New().String()
	delayDuration := delayToDuration(delay)
	timer := InitiateNewKickerTimer(delayDuration)
	context, cancelFunc := context.WithCancel(context.Background())
	job := &Job{
		JobType:    Once,
		Fn:         fn,
		Args:       args,
		Timer:      timer,
		cxt:        context,
		cancelFunc: cancelFunc,
	}
	jobKicker.jobQueue.Lock()
	jobKicker.jobQueue.PendingJobs[jobID] = job
	jobKicker.jobQueue.Unlock()

	go jobKicker.runJob(job, jobID)
	return
}

// KickOnceAt excutes a function once at a certain time and returns the job id
func (jobKicker *JobKicker) KickOnceAt(runAt time.Time, fn interface{}, args ...interface{}) (jobID string) {
	jobID = uuid.New().String()
	duration := time.Until(runAt)
	timer := InitiateNewKickerTimer(duration)
	context, cancelFunc := context.WithCancel(context.Background())
	job := &Job{
		JobType:    Once,
		Fn:         fn,
		Args:       args,
		Timer:      timer,
		cxt:        context,
		cancelFunc: cancelFunc,
	}
	jobKicker.jobQueue.Lock()
	jobKicker.jobQueue.PendingJobs[jobID] = job
	jobKicker.jobQueue.Unlock()

	go jobKicker.runJob(job, jobID)
	return
}

// KickPeriodicallyEvery excutes a function Periodically with the passed interval between every execution and returns the job id
func (jobKicker *JobKicker) KickPeriodicallyEvery(delay time.Time, fn interface{}, args ...interface{}) (jobID string) {
	jobID = uuid.New().String()
	delayDuration := delayToDuration(delay)
	ticker := InitiateNewKickerTicker(delayDuration)
	context, cancelFunc := context.WithCancel(context.Background())
	job := &Job{
		JobType:    Periodically,
		Fn:         fn,
		Args:       args,
		Timer:      ticker,
		cxt:        context,
		cancelFunc: cancelFunc,
	}
	jobKicker.jobQueue.Lock()
	jobKicker.jobQueue.PendingJobs[jobID] = job
	jobKicker.jobQueue.Unlock()

	go jobKicker.runJob(job, jobID)
	return
}
