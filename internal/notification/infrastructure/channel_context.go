package infrastructure

import "context"

func channelContextError(ctx context.Context) error { return ctx.Err() }
