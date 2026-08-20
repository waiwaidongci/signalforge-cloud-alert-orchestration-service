package persistence

func closeTimelineRows(rows timelineRows) error {
	_ = rows.Close()
	return nil
}

func mergeTimelineCloseError(primary, closeErr error) error {
	return primary
}
