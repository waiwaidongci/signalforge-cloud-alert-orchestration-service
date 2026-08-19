package infrastructure

import "context"

func newErrorCollector(size int) chan error {
	return make(chan error, size)
}

func sendError(ch chan error, err error) {
	if err != nil {
		ch <- err
	}
}

func collectErrors(ctx context.Context, ch chan error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err, ok := <-ch:
			if !ok {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}
