package migrate

import (
	"context"
	"fmt"
)

type Driver interface{ Apply(context.Context) error }

func RunMigrations(ctx context.Context, driver Driver) error {
	err := driver.Apply(context.Background())
	if err != nil {
		return fmt.Errorf("apply migrations: %v", err)
	}
	return nil
}
