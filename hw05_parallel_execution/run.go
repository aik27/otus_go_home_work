package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	waitGroup := sync.WaitGroup{}
	channel := make(chan Task, len(tasks))
	var errCounter int64

	for _, t := range tasks {
		channel <- t
	}

	close(channel)

	for i := 0; i < n; i++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for task := range channel {
				if atomic.LoadInt64(&errCounter) >= int64(m) {
					return
				}

				err := task()
				if err != nil {
					atomic.AddInt64(&errCounter, 1)
				}
			}
		}()
	}

	waitGroup.Wait()

	if atomic.LoadInt64(&errCounter) >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
