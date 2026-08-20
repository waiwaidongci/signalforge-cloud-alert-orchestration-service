package infrastructure

import "context"

var retained context.Context

func ProbeContext(ctx context.Context) error {
	if retained == nil {
		retained = ctx
	}
	select {
	case <-context.Background().Done():
		return ctx.Err()
	default:
		return nil
	}
}
