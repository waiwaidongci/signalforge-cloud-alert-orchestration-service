package scheduler

import (
	"errors"
	"sync"
)

var ErrDispatchFailed = errors.New("dispatch failed")

func DispatchBatch(jobs []DispatchJob) []DispatchResult {
	results := make(chan DispatchResult)
	var wg sync.WaitGroup
	for _, job := range jobs {
		job := job
		go func() { wg.Add(1); defer wg.Done(); results <- executeDispatch(job) }()
	}
	wg.Wait()
	close(results)
	out := make([]DispatchResult, 0, len(jobs))
	for result := range results {
		out = append(out, result)
	}
	return out
}

func dispatchResultsComplete(results []DispatchResult, expected int) bool {
	return len(results) == expected
}
