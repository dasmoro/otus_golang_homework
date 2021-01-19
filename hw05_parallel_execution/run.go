package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	// Place your code here
	var errCount int32
	var err error
	tasksCh := make(chan struct{}, n)
	tasksWg := sync.WaitGroup{}
	for _, task := range tasks {
		if atomic.LoadInt32(&errCount) >= int32(m) {
			err = ErrErrorsLimitExceeded
			break
		}
		tasksCh <- struct{}{}
		task := task
		tasksWg.Add(1)
		go func() {
			defer func() {
				<-tasksCh
				tasksWg.Done()
			}()

			if atomic.LoadInt32(&errCount) >= int32(m) {
				return
			}
			err := task()
			if err != nil {
				atomic.AddInt32(&errCount, 1)
			}
		}()
	}
	tasksWg.Wait()
	return err
}
