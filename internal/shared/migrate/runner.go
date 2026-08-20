package migrate

import (
	"context"
	"fmt"
)

type Driver interface{ Apply(context.Context) error }

func applyMigration(ctx context.Context, driver Driver) error { return driver.Apply(ctx) }

func RunMigrations(ctx context.Context, driver Driver) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := applyMigration(ctx, driver)
	if err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
