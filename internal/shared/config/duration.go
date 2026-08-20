package config

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration wraps time.Duration so YAML values such as "10s" decode naturally.
type Duration time.Duration

func (d Duration) Value() time.Duration {
	return time.Duration(d)
}

func (d Duration) ValidatePositive() error {
	if d <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}

func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("duration must be a scalar value")
	}
	parsed, err := time.ParseDuration(node.Value)
	if err != nil {
		return fmt.Errorf("parse duration %q: %w", node.Value, err)
	}
	if err := Duration(parsed).ValidatePositive(); err != nil {
		return err
	}
	*d = Duration(parsed)
	return nil
}
