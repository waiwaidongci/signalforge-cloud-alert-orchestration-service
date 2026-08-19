package domain

const StatusAcknowledged Status = "acknowledged"

func (s Status) Valid() bool {
	switch s {
	case StatusFiring, StatusResolved, StatusSuppressed, StatusAcknowledged:
		return true
	default:
		return false
	}
}

func AllStatuses() []string {
	return []string{string(StatusFiring), string(StatusResolved), string(StatusSuppressed), string(StatusAcknowledged)}
}
