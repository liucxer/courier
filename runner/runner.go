package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
)

func NewRunner(ctx context.Context, name string, numWorkers int) *Runner {
	if numWorkers < 1 {
		numWorkers = 1
	}
	runner := &Runner{
		Name:       name,
		ctx:        ctx,
		numWorkers: numWorkers,
	}
	return runner
}

type Runner struct {
	Name       string
	ctx        context.Context
	numWorkers int
	jobs       []*Job
}

func (r *Runner) Add(task Task, desc ...string) {
	r.jobs = append(r.jobs, &Job{
		Task: task,
		Name: fmt.Sprintf("[%s] %s", r.Name, strings.Join(desc, ";")),
	})
}

func (r *Runner) Commit() (err error) {
	total := int64(len(r.jobs))
	if total == 0 {
		return nil
	}

	errLocker := sync.RWMutex{}
	setErr := func(e error) {
		errLocker.Lock()
		defer errLocker.Unlock()
		if err == nil {
			err = e
		}
	}

	queue := make(chan *Job, total)
	defer close(queue)

	for i := int64(0); i < total; i++ {
		queue <- r.jobs[i]
	}
	r.jobs = nil

	numberWorks := int64(r.numWorkers)
	if numberWorks > total {
		numberWorks = total
	}

	goroutines := sync.WaitGroup{}

	for i := int64(0); i < numberWorks; i++ {
		goroutines.Add(1)
		go func() {
			defer goroutines.Done()
			for job := range queue {
				atomic.AddInt64(&total, -1)
				e := dispatch(r.ctx, job)
				if e != nil {
					setErr(e)
					break
				}
				// when task down
				if atomic.LoadInt64(&total) == 0 {
					break
				}
			}
		}()
	}

	goroutines.Wait()
	return
}

func dispatch(ctx context.Context, job *Job) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	select {
	case <-interrupt:
		return &Error{
			Name: job.Name,
			Type: ErrTypeInterrupt,
		}
	case <-ctx.Done():
		return &Error{
			Name: job.Name,
			Type: ErrTypeTimeout,
		}
	default:
		return job.Task(ctx)
	}
}

type Task func(ctx context.Context) error

type Job struct {
	Name string
	Task Task
}
