package jobkicker

import (
	"fmt"
	"testing"
	"time"
)

func testingTask1() {
	println("Hello from job 1")
}

func testingTaskWithParams(name string, age int) {
	fmt.Printf("hello %s, your age is %d\n", name, age)
}

func TestNewScheduler(t *testing.T) {
	jk := NewScheduler(nil, nil)
	jk.jobQueue.RLock()
	pendingJobsSize := len(jk.jobQueue.PendingJobs)
	doneJobsSize := len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 0 {
		t.Errorf("newly intialized jobkicker's pendingjobs should be empty but it's size %d", pendingJobsSize)
	}
	if doneJobsSize != 0 {
		t.Errorf("newly intialized jobkicker's donejobs should be empty but it's size %d", doneJobsSize)
	}
	t.Log("newly intialized jobkicker's jobqueue is empty passed")

	jk.jobQueue.RUnlock()
}
func TestKickOnceAfter(t *testing.T) {
	jk := NewScheduler(nil, nil)

	delay := time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC)
	jk.KickOnceAfter(delay, testingTask1)
	jk.jobQueue.RLock()
	pendingJobsSize := len(jk.jobQueue.PendingJobs)
	doneJobsSize := len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 0 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 0, doneJobsSize)
	}
	jk.jobQueue.RUnlock()
	t.Log("scheduling new job passed")

	time.Sleep(4 * time.Second)
	jk.jobQueue.RLock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 0 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 0, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	jk.jobQueue.RUnlock()
	t.Log("finishing scheduled job passed")

}

func TestKickOnceAt(t *testing.T) {
	jk := NewScheduler(nil, nil)

	runAt := time.Now().Add(3 * time.Second)
	jk.KickOnceAt(runAt, testingTask1)
	jk.jobQueue.RLock()
	pendingJobsSize := len(jk.jobQueue.PendingJobs)
	doneJobsSize := len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 0 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 0, doneJobsSize)
	}
	jk.jobQueue.RUnlock()
	t.Log("scheduling new job passed")

	time.Sleep(4 * time.Second)
	jk.jobQueue.RLock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 0 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 0, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	jk.jobQueue.RUnlock()
	t.Log("finishing scheduled job passed")

}

func TestKickPeriodicallyEvery(t *testing.T) {
	jk := NewScheduler(nil, nil)

	delay := time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC)
	jobID := jk.KickPeriodicallyEvery(delay, testingTask1)
	jk.jobQueue.RLock()
	pendingJobsSize := len(jk.jobQueue.PendingJobs)
	doneJobsSize := len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 0 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 0, doneJobsSize)
	}
	jk.jobQueue.RUnlock()
	t.Log("scheduling new job passed")

	time.Sleep(4 * time.Second)
	jk.jobQueue.RLock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	firstFinishTime := jk.jobQueue.DoneJobs[jobID]
	jk.jobQueue.RUnlock()

	time.Sleep(4 * time.Second)
	jk.jobQueue.RLock()
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
	jk.jobQueue.RUnlock()

	err := jk.CancelJob(jobID)
	if err != nil {
		t.Errorf("can't cancel job due to error: %s", err.Error())
	}
	// a little delay before locking the jobQueue
	// to allow locking it from cancel side to delete
	// it from the pending jobs map
	time.Sleep(time.Second)
	jk.jobQueue.RLock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 0 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 0, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	jk.jobQueue.RUnlock()

	t.Log("finishing scheduled job passed")

}

func TestKickWithParams(t *testing.T) {
	jk := NewScheduler(nil, nil)

	delay := time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC)
	jk.KickOnceAfter(delay, testingTaskWithParams, "Mohab", 25)
	jk.jobQueue.RLock()
	pendingJobsSize := len(jk.jobQueue.PendingJobs)
	doneJobsSize := len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 1 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 1, pendingJobsSize)
	}
	if doneJobsSize != 0 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 0, doneJobsSize)
	}
	jk.jobQueue.RUnlock()
	t.Log("scheduling new job with arguments passed")

	time.Sleep(4 * time.Second)
	jk.jobQueue.RLock()
	pendingJobsSize = len(jk.jobQueue.PendingJobs)
	doneJobsSize = len(jk.jobQueue.DoneJobs)
	if pendingJobsSize != 0 {
		t.Errorf("jobkicker's pendingjobs should have size %d but found it's size %d", 0, pendingJobsSize)
	}
	if doneJobsSize != 1 {
		t.Errorf("jobkicker's donejobs should have size %d but found it's size %d", 1, doneJobsSize)
	}
	jk.jobQueue.RUnlock()
	t.Log("finishing scheduled job passed")
}
