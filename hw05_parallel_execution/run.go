package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	wg.Add(n)

	taskChan := make(chan Task, len(tasks))
	errCount := 0

	// produce
	go addTasks(tasks, taskChan)

	// consume
	mu := sync.Mutex{}

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for task := range taskChan {
				err := task()

				mu.Lock()

				if errCount >= m {
					mu.Unlock()

					return
				}

				if err != nil {
					errCount++
				}

				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func addTasks(tasks []Task, taskChan chan Task) {
	for _, t := range tasks {
		taskChan <- t
	}

	close(taskChan)
}
