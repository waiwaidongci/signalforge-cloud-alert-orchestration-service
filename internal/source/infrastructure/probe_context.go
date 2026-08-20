package infrastructure

import "context"

func ProbeContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func ProbeUsable(ctx context.Context) bool { return ProbeContext(ctx) == nil }
