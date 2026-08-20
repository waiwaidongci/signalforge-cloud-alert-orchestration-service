package infrastructure

import "context"

func ProbeContextsIndependent(first, second context.Context) bool {
	switch {
	case first.Err() != nil:
		return second.Err() == nil
	case second.Err() != nil:
		return first.Err() == nil
	default:
		return true
	}
}
