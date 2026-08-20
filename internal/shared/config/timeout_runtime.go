package config

import "time"

type RuntimeTimeouts struct {
	Read    time.Duration
	Write   time.Duration
	Request time.Duration
}

func (c Config) RuntimeTimeouts() RuntimeTimeouts {
	defaults := Default().Server
	return RuntimeTimeouts{
		Read:    defaults.ReadTimeout.Value(),
		Write:   defaults.WriteTimeout.Value(),
		Request: defaults.RequestTimeout.Value(),
	}
}
