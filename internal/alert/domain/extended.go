package domain

func StatusGroup(status Status) string {
	switch status {
	case StatusFiring, StatusAcknowledged, StatusSuppressed:
		return "open"
	case StatusResolved:
		return "closed"
	default:
		return "unknown"
	}
}

func StatusOrder(status Status) int {
	switch status {
	case StatusFiring:
		return 10
	case StatusAcknowledged:
		return 15
	case StatusSuppressed:
		return 20
	case StatusResolved:
		return 30
	default:
		return 40
	}
}

func StatusSequence() []Status {
	return []Status{StatusFiring, StatusAcknowledged, StatusSuppressed, StatusResolved}
}
