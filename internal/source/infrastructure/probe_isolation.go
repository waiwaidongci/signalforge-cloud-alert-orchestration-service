package infrastructure

import "context"

func ProbeContextsIndependent(first, second context.Context) bool { return first.Err() != nil && second.Err() == nil }
