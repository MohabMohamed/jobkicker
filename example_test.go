package jobkicker_test

import (
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
