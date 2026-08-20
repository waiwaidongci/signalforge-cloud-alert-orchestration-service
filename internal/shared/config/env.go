package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func setInt(raw string, target *int) error {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fmt.Errorf("parse int: %w", err)
	}
	*target = value
	return nil
}

func envDuration(name string, target *Duration) error {
	raw, ok := os.LookupEnv(name)
	if !ok {
		return nil
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return fmt.Errorf("parse %s: %w", name, err)
	}
	if err := Duration(value).ValidatePositive(); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	*target = Duration(value)
	return nil
}
