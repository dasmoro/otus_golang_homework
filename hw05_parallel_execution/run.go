package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	// Place your code here
	var errCount int
	var err error
	tasksCh := make(chan struct{}, n)
	tasksWg := sync.WaitGroup{}
	tasksMtx := sync.Mutex{}
	for _, task := range tasks {
		if errCount >= m {
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
			if errCount >= m {
				return
			}
			err := task()
			if err != nil {
				tasksMtx.Lock()
				errCount++
				tasksMtx.Unlock()
			}
		}()
	}
	tasksWg.Wait()
	return err
}
