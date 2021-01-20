package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func bogoSort(arr []int) ([]int, error) {
	start := time.Now()
	for {
		elapsedTime := time.Since(start)
		if sort.IntsAreSorted(arr) {
			return arr, nil
		}
		if elapsedTime > time.Second*2 {
			return nil, fmt.Errorf("Max execution time is reached")
		}
		arr = shuffle(arr)
	}
}
func shuffle(arr []int) []int {
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}

func TestReal(t *testing.T) {
	taskCount := 100
	tasks := make([]Task, 0, taskCount)
	var runElapsedTasksCount int32
	var errCount int32
	workersCount := 31
	maxErrorsCount := 5
	for i := 0; i < taskCount; i++ {
		tasks = append(tasks, func() error {
			if errCount >= int32(maxErrorsCount) {
				atomic.AddInt32(&runElapsedTasksCount, 1)
			}
			arr := rand.Perm(rand.Intn(20))
			arr, err := bogoSort(arr)
			if err != nil {
				atomic.AddInt32(&errCount, 1)
			}
			return err
		})
	}

	Run(tasks, workersCount, maxErrorsCount)
	require.LessOrEqual(t, runElapsedTasksCount, int32(workersCount), "extra tasks were started")
}
