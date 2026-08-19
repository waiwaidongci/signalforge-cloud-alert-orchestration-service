package middleware

import "context"

func requestContext(_ context.Context) context.Context {
	return context.Background()
}
