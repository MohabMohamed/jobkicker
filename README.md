# jobkicker
[![Go Reference](https://pkg.go.dev/badge/github.com/MohabMohamed/jobkicker.svg)](https://pkg.go.dev/github.com/MohabMohamed/jobkicker)

jobkicker is A Golang in-process cron task scheduler that kicks (executes) them once in specified time or periodically.

## Features

- Kick (execute) jobs to run after some time once or periodically.
- Kick (execute) jobs to run at certain time.
- Ability to cancel any job with it's id.
- Logs execution and canceling jobs with the flexibility to save the logs to any thing implements `io.Writer` interface like a file or even to implement a writer to write the logs to something like elastic search and pass it to the scheduler
- flexibility to format the logs as you wish by passing the [Formatter](https://pkg.go.dev/github.com/sirupsen/logrus#Formatter) interface from [sirupsen/logrus](https://github.com/sirupsen/logrus)
- Uses language built-in time.Time to reduce design complexity by not using something like cron time format.
- Run multiply scheduled jobs concurrently at the same time.

## Third-party liberaries

- [sirupsen/logrus](https://github.com/sirupsen/logrus) for logging.
- [google/uuid](github.com/google/uuid) for generating job's id.

## Main components

### [JobKicker](https://pkg.go.dev/github.com/MohabMohamed/jobkicker#JobKicker)
The main type in the package which is the scheduler that kicks new jobs to run and cancels them and holds the `JobQueue` and the `Logger`.
### [ITimer](https://pkg.go.dev/github.com/MohabMohamed/jobkicker#ITimer)
An interface type which `KickerTimer` and `KickerTicker` implement it which they are just wrappers for `time.Timer` and `time.Ticker`.
### [Job](https://pkg.go.dev/github.com/MohabMohamed/jobkicker#Job)
Job struct that holds the function with it's arguments and the timer to execute it.

### [JobQueue](https://pkg.go.dev/github.com/MohabMohamed/jobkicker#JobQueue)
JobQueue is just holding `PendingJobs` which is just a map for the pending jobs to be excuted `map[string]*Job` and `DoneJobs` which is a map for executed jobs with it's last execution time and a read/write mutex to lock this two maps when accessed.

#### Note: periodically executed functions stays in `PendingJobs` map after execution (unless canceled) and `DoneJobs` keeps track of the time of it's last execution

## Package Apis
### `func NewScheduler(loggerOutput *io.Writer, loggerFormatter *log.Formatter) *JobKicker`

Returns a new JobKicker (scheduler) and takes:

- `loggerOutput` which any type implements the interface `io.Writer` to write the logs to, and if nil passed it will write to `os.Stderr`. the interface `io.Writer` is:

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

- `loggerFormatter` which any type implements the interface `logrus.Formatter` interface, and if nil passed it use `logrus.TextFormatter` by default, you can try to pass `&logrus.JSONFormatter{}` to format the logs as json or pass your custom formatter that implements:

```go
type Formatter interface {
	Format(*Entry) ([]byte, error)
}
```

### `func (jobKicker *JobKicker) KickOnceAfter(delay time.Time, fn interface{}, args ...interface{}) (jobID string)`
Runs a function once after a given delay, the delay is a `time.Time` type with all fields zero expect the time to runs it after, as if you want to run it after 3 hours and 30 minutes create a new time with `time.Date(year ,month ,day ,hour ,min ,sec,nsec,loc)` with all fields parameters equal zero expect `hour` = 3 and `min` = 30, and the second parameter is the function to run and the rest of the parameters are the function arguments if any.

example:

```go
import (
	"time"

	"github.com/MohabMohamed/jobkicker"
)

func main() {
	task := func() {
		println("jobkicker is awesome")
	}
	jk := jobkicker.NewScheduler(nil, nil)
	// time.Date(year ,month ,day ,hour ,min ,sec,nsec,loc)
	// every field equals zeo expect seconds equals 3
	delay := time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC)
	jk.KickOnceAfter(delay, task)
	time.Sleep(4 * time.Second)
}
```

```log
Output:

jobkicker is awesome
INFO[0003] job with id [5e1b8baa-5133-483b-9188-8179ecc8aea4] executed in 2022-02-08 22:43:08.701094256 +0200 EET m=+3.002681837
```

### `func (jobKicker *JobKicker) KickOnceAt(runAt time.Time, fn interface{}, args ...interface{}) (jobID string)`
Runs a function once at a certain time, for example if you want to run a function at `1 march 2022 13:30 am` you should create a `time.Time` with this certain time like `time.Date(2022 ,3 ,1 ,13 ,30 ,0 , 0, time.UTC)` and the function will run at that time. the second parameter is the function and the rest are the function arguments.

example:
```go
package main

import (
	"fmt"
	"time"

	"github.com/MohabMohamed/jobkicker"
)

func main() {
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
}
```

```log
Output:

2022-02-09 12:39:26.645499827 +0200 EET m=+0.000422605
jobkicker is awesome
INFO[0003] job with id [0477e575-22a1-48a2-851f-017b5aeb9ea4] executed in 2022-02-09 12:39:29.646510482 +0200 EET m=+3.001433319
```


### `func (jobKicker *JobKicker) KickPeriodicallyEvery(delay time.Time, fn interface{}, args ...interface{}) (jobID string)`

Runs the function every some specified time intervals it takes the delay like `KickOnceAfter` so if you pass `time.Time` with 3 seconds it will run the function every 3 seconds, and the second parameter is the function and the rest are the function arguments.

example:

example:

```go
package main

import (
	"time"

	"github.com/MohabMohamed/jobkicker"
)

func main() {
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
}

```

```log
Output:

jobkicker is awesome
INFO[0003] job with id [d8ae6a99-1d2b-434f-8a9c-db77a4e3e844] executed in 2022-02-09 12:43:54.758287996 +0200 EET m=+3.001531441
jobkicker is awesome
INFO[0009] job with id [d8ae6a99-1d2b-434f-8a9c-db77a4e3e844] executed in 2022-02-09 12:43:57.760335343 +0200 EET m=+6.003578840
jobkicker is awesome
INFO[0011] job with id [d8ae6a99-1d2b-434f-8a9c-db77a4e3e844] executed in 2022-02-09 12:44:00.758331522 +0200 EET m=+9.001575009
```

### `func (jobKicker *JobKicker) CancelJob(jobId string) error`
Cancels the scheduling of a job, it takes it's id and return error if the job of type run once and already ran  or if it's not already scheduled (maybe wrong id given)

example:

example:
```go
package main

import (
	"fmt"
	"time"

	"github.com/MohabMohamed/jobkicker"
)

func main() {
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
}
```

```log
Output:

2022-02-09 12:39:26.645499827 +0200 EET m=+0.000422605
jobkicker is awesome
INFO[0003] job with id [0477e575-22a1-48a2-851f-017b5aeb9ea4] executed in 2022-02-09 12:39:29.646510482 +0200 EET m=+3.001433319
```


### `func (jobKicker *JobKicker) KickPeriodicallyEvery(delay time.Time, fn interface{}, args ...interface{}) (jobID string)`

Runs the function every some specified time intervals it takes the delay like `KickOnceAfter` so if you pass `time.Time` with 3 seconds it will run the function every 3 seconds, and the second parameter is the function and the rest are the function arguments.

example:

example:

```go
package main

import (
	"fmt"
	"time"

	"github.com/MohabMohamed/jobkicker"
)

func main() {
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
}
```

```log
Output:

jobkicker is awesome
INFO[0003] job with id [e7e0d6ab-6aa5-4b52-9f60-3c74e5585b15] executed in 2022-02-09 12:49:14.102550679 +0200 EET m=+3.001345112
INFO[0016] job with id [e7e0d6ab-6aa5-4b52-9f60-3c74e5585b15] cancelled successfully in 2022-02-09 12:49:15.101968589 +0200 EET m=+4.000763046
```

## Combining the Api
you can use the apis in any combination to kick the jobs, for example if you want to kick a job every year at `1 march 2022 13:30 am` you should create a `time.Time` with this certain time like `time.Date(2022 ,3 ,1 ,13 ,30 ,0 , 0, time.UTC)` and delay with `1 year` and pass `KickPeriodicallyEvery` with arguments the delay and the function to kick and it's arguments to `KickOnceAt` the date specified like:

```go
package main

import (
	"fmt"
	"time"

	"github.com/MohabMohamed/jobkicker"
)

func main() {
	task := func(name string) {
		fmt.Printf("jobkicker is awesome. don't you agree,%s?\n", name)
	}
	jk := jobkicker.NewScheduler(nil, nil)
	date := time.Date(2022, 3, 1, 13, 30, 0, 0, time.UTC)
	delay := time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC)

	jk.KickPeriodicallyEvery(delay, jk.KickOnceAt, date, task, "Mohab")

	// block the main goroutine, could be server.listen() or any thing
	for true {
	}
}
```

Get creative using jobkicker, and keep kicking these tasks :D

## Some trade-offs while designing jobkicker

- Made the job execution (timers and context) self contained in the job to make it easier to cancel.

- Used Read/Write mutex instead of regular mutex as the only write operations made to the JobQueue when remove the job from pendeningJobs (Run once job and executed or a canceled job) and when executing job adding the last execution time to the done jobs, the rest is read operations.

- Use regular map with RWMutex instead of sync.map as I have 2 maps so with 2 sync.map both of them will lock and unlock and both of them  need to be locked at the same time so  lock and unlock a mutex and lock and unlock another one will be performance costly more than using 1 mutex.

- Using already built-in `time.Time` instead of rolling of my solution to handle delay and time as every go developer is familier with them so it would easier for the user.

## Contribution

check [contribution guide](./CONTRIBUTION.md) and the [Reference](https://pkg.go.dev/github.com/MohabMohamed/jobkicker)

## Future improvements

I'm considering maybe to add the ability to consist the tasks execution in redis as option as if the client code that using jobkicker got down and up again could reschedule the tasks that already scheduled.

Maybe adding the ability to schedule the tasks in distributed environment as if a task ran on a machine it shouldn't  run from another one.
