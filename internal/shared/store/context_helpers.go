package store

import "context"

func execContext(_ context.Context) context.Context {
	return context.Background()
}
