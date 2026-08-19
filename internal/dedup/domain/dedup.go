package domain

import "time"

type Window struct {
	Duration time.Duration
}

func (w Window) Valid() bool {
	return w.Duration > 0
}
