package domain

import "fmt"

type SourcePolicy struct {
	required map[string]struct{}
}

func NewSourcePolicy() *SourcePolicy {
	return &SourcePolicy{required: map[string]struct{}{}}
}

func (p *SourcePolicy) AddRule(label string) error {
	if p == nil {
		return fmt.Errorf("source policy is unavailable")
	}
	if label == "" {
		return fmt.Errorf("required label is empty")
	}
	if p.required == nil {
		p.required = map[string]struct{}{}
	}
	p.required[label] = struct{}{}
	return nil
}

func (p *SourcePolicy) Allows(labels map[string]string) bool {
	if p == nil {
		return true
	}
	for label := range p.required {
		if labels[label] == "" {
			return false
		}
	}
	return true
}
