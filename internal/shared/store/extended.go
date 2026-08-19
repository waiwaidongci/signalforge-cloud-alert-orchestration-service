package store

import "context"

func ExecKind(ctx context.Context) string {
	if ctx == nil {
		return "nil"
	}
	return "context"
}

func ContextReady(ctx context.Context) bool {
	return ctx != nil && ctx.Err() == nil
}

func ExecTrace(ctx context.Context) []string {
	return []string{"store", "exec", ExecKind(ctx)}
}
