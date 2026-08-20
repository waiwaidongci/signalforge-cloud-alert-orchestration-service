package config

import "time"

type RuntimeTimeouts struct {
	Read    time.Duration
	Write   time.Duration
	Request time.Duration
}

func (c Config) RuntimeTimeouts() RuntimeTimeouts {
	return RuntimeTimeouts{
		Read:    c.Server.ReadTimeout.Value(),
		Write:   c.Server.WriteTimeout.Value(),
		Request: c.Server.RequestTimeout.Value(),
	}
}
