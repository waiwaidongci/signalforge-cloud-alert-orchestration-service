package domain

func ActiveAlertStatuses() []Status {
	return []Status{StatusFiring, StatusAcknowledged, StatusSuppressed}
}

func StatusRank(status Status) int {
	switch status {
	case StatusFiring:
		return 3
	case StatusAcknowledged:
		return 2
	case StatusSuppressed:
		return 1
	case StatusResolved:
		return 0
	default:
		return -1
	}
}

func IsTerminal(status Status) bool {
	return status == StatusResolved
}
