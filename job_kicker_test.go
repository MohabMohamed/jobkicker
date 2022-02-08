package jobkicker

import (
	"testing"
	"time"
)

func testingTask1() {
	println("Hello from job 1")
}

func TestNewScheduler(t *testing.T) {
	jk := NewScheduler(nil, nil)
	jk.jobQueue.Lock()
	pendingJobsSize := len(jk.jobQueue.PendingJobs)
	doneJobsSize := len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 0 {
		t.Errorf("newly intialized jobkicker's pendingjobs should be empty but it's size %d", pendingJobsSize)
	}
	if doneJobsSize != 0 {
		t.Errorf("newly intialized jobkicker's donejobs should be empty but it's size %d", doneJobsSize)
	}
	t.Log("newly intialized jobkicker's jobqueue is empty passed")

	jk.jobQueue.Unlock()
}
func TestKickOnceAfter(t *testing.T) {
	jk := NewScheduler(nil, nil)

	delay := time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC)
	jk.KickOnceAfter(delay, testingTask1)
	jk.jobQueue.Lock()
	pendingJobsSize := len(jk.jobQueue.PendingJobs)
	doneJobsSize := len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 0 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 0, doneJobsSize)
	}
	jk.jobQueue.Unlock()
	t.Log("scheduling new job passed")

	time.Sleep(4 * time.Second)
	jk.jobQueue.Lock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 0 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 0, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	jk.jobQueue.Unlock()
	t.Log("finishing scheduled job passed")

}

func TestKickOnceAt(t *testing.T) {
	jk := NewScheduler(nil, nil)

	runAt := time.Now().Add(3 * time.Second)
	jk.KickOnceAt(runAt, testingTask1)
	jk.jobQueue.Lock()
	pendingJobsSize := len(jk.jobQueue.PendingJobs)
	doneJobsSize := len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 0 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 0, doneJobsSize)
	}
	jk.jobQueue.Unlock()
	t.Log("scheduling new job passed")

	time.Sleep(4 * time.Second)
	jk.jobQueue.Lock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 0 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 0, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	jk.jobQueue.Unlock()
	t.Log("finishing scheduled job passed")

}

func TestKickPeriodicallyEvery(t *testing.T) {
	jk := NewScheduler(nil, nil)

	delay := time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC)
	jobID := jk.KickPeriodicallyEvery(delay, testingTask1)
	jk.jobQueue.Lock()
	pendingJobsSize := len(jk.jobQueue.PendingJobs)
	doneJobsSize := len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 0 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 0, doneJobsSize)
	}
	jk.jobQueue.Unlock()
	t.Log("scheduling new job passed")

	time.Sleep(4 * time.Second)
	jk.jobQueue.Lock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	firstFinishTime := jk.jobQueue.DoneJobs[jobID]
	jk.jobQueue.Unlock()

	time.Sleep(4 * time.Second)
	jk.jobQueue.Lock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	if firstFinishTime.Equal(jk.jobQueue.DoneJobs[jobID]) {
		t.Errorf("Job's first running time and second running time are equal first:%v and second: %v", firstFinishTime, jk.jobQueue.DoneJobs[jobID])
	}
	jk.jobQueue.Unlock()

	err := jk.CancelJob(jobID)
	if err != nil {
		t.Errorf("can't cancel job due to error: %s", err.Error())
	}
	// a little delay before locking the jobQueue
	// to allow locking it from cancel side to delete
	// it from the pending jobs map
	time.Sleep(time.Second)
	jk.jobQueue.Lock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 0 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 0, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	jk.jobQueue.Unlock()

	t.Log("finishing scheduled job passed")

}
