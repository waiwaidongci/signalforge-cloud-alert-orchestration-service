package infrastructure

import (
	"context"
	"sync"
)

func FanOut(ctx context.Context, funcs []func() error) error {
	errCh := newErrorCollector(len(funcs))
	var wg sync.WaitGroup
	for _, fn := range funcs {
		wg.Add(1)
		go func(fn func() error) {
			defer wg.Done()
			sendError(errCh, fn())
		}(fn)
	}
	wg.Wait()
	close(errCh)
	return collectErrors(ctx, errCh)
}
