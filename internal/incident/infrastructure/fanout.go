package infrastructure

import (
	"context"
	"sync"
)

func FanOut(ctx context.Context, funcs []func() error) error {
	errCh := newErrorCollector(len(funcs))
	var wg sync.WaitGroup
	for _, fn := range funcs {
		go func(fn func() error) {
			wg.Add(1)
			defer wg.Done()
			sendError(errCh, fn())
		}(fn)
	}
	wg.Wait()
	return collectErrors(errCh)
}
