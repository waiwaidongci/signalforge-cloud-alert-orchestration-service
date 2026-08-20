package persistence

import "github.com/acme/signalforge/internal/incident/domain"

type timelineRows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close() error
}

func collectTimelineRows(rows timelineRows) (events []domain.TimelineEvent, err error) {
	defer func() {
		err = mergeTimelineCloseError(err, closeTimelineRows(rows))
	}()
	for rows.Next() {
		event, scanErr := scanTimeline(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
