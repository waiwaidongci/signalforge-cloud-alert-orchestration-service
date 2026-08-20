package persistence

import (
	"errors"
	"fmt"
)

func closeTimelineRows(rows timelineRows) error {
	return rows.Close()
}

func mergeTimelineCloseError(primary, closeErr error) error {
	if primary == nil {
		return closeErr
	}
	if closeErr == nil {
		return primary
	}
	return fmt.Errorf("timeline iteration and close failed: %w", errors.Join(primary, closeErr))
}
