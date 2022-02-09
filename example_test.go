package jobkicker_test

import (
	"fmt"
	"time"

	"github.com/MohabMohamed/jobkicker"
)

func ExampleJobKicker_KickOnceAfter() {
	task := func() {
		println("jobkicker is awesome")
	}
	jk := jobkicker.NewScheduler(nil, nil)
	// time.Date(year ,month ,day ,hour ,min ,sec,nsec,loc)
	// every field equals zeo expect seconds equals 3
	delay := time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC)
	jk.KickOnceAfter(delay, task)
	time.Sleep(4 * time.Second)
	// Output:
	// jobkicker is awesome
	// INFO[0003] job with id [5e1b8baa-5133-483b-9188-8179ecc8aea4] executed in 2022-02-08 22:43:08.701094256 +0200 EET m=+3.002681837
}

func ExampleJobKicker_KickOnceAt() {
	task := func() {
		println("jobkicker is awesome")
	}
	jk := jobkicker.NewScheduler(nil, nil)
	// time.Date(year ,month ,day ,hour ,min ,sec,nsec,loc)
	// KickOnceAt takes the date to run the task in
	// in this example I run it with date equals now + 3 seconds
	fmt.Println(time.Now())
	runAt := time.Now().Add(3 * time.Second)
	jk.KickOnceAt(runAt, task)
	time.Sleep(4 * time.Second)
	// Output:
	// 2022-02-09 12:39:26.645499827 +0200 EET m=+0.000422605
	// jobkicker is awesome
	// INFO[0003] job with id [0477e575-22a1-48a2-851f-017b5aeb9ea4] executed in 2022-02-09 12:39:29.646510482 +0200 EET m=+3.001433319
}

func ExampleJobKicker_KickPeriodicallyEvery() {
	task := func() {
		println("jobkicker is awesome")
	}
	jk := jobkicker.NewScheduler(nil, nil)
	// time.Date(year ,month ,day ,hour ,min ,sec,nsec,loc)
	// every field equals zeo expect seconds equals 3
	// as KickPeriodicallyEvery takes the delay between every execution
	delay := time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC)
	jk.KickPeriodicallyEvery(delay, task)
	time.Sleep(10 * time.Second)
	// Output:
	// jobkicker is awesome
	// INFO[0003] job with id [d8ae6a99-1d2b-434f-8a9c-db77a4e3e844] executed in 2022-02-09 12:43:54.758287996 +0200 EET m=+3.001531441
	// jobkicker is awesome
	// INFO[0009] job with id [d8ae6a99-1d2b-434f-8a9c-db77a4e3e844] executed in 2022-02-09 12:43:57.760335343 +0200 EET m=+6.003578840
	// jobkicker is awesome
	// INFO[0011] job with id [d8ae6a99-1d2b-434f-8a9c-db77a4e3e844] executed in 2022-02-09 12:44:00.758331522 +0200 EET m=+9.001575009
}

func ExampleJobKicker_CancelJob() {
	task := func() {
		println("jobkicker is awesome")
	}
	jk := jobkicker.NewScheduler(nil, nil)
	// time.Date(year ,month ,day ,hour ,min ,sec,nsec,loc)
	// every field equals zeo expect seconds equals 3
	// as KickPeriodicallyEvery takes the delay between every execution
	delay := time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC)
	jobId := jk.KickPeriodicallyEvery(delay, task)
	// sleep for 4 seconds to let it run once before cancelling
	time.Sleep(4 * time.Second)

	err := jk.CancelJob(jobId)
	if err != nil {
		fmt.Printf("error while cancelling a job: %s", err.Error())
	}
	time.Sleep(10 * time.Second)
	// Output:
	// jobkicker is awesome
	// INFO[0003] job with id [e7e0d6ab-6aa5-4b52-9f60-3c74e5585b15] executed in 2022-02-09 12:49:14.102550679 +0200 EET m=+3.001345112
	// INFO[0016] job with id [e7e0d6ab-6aa5-4b52-9f60-3c74e5585b15] cancelled successfully in 2022-02-09 12:49:15.101968589 +0200 EET m=+4.000763046
}
