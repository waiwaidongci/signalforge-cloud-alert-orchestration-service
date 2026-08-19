package domain

const StatusAcknowledged Status = "acknowledged"

func (s Status) Valid() bool {
	switch s {
	case StatusFiring, StatusAcknowledged, StatusResolved, StatusSuppressed:
		return true
	default:
		return false
	}
}

func AllStatuses() []string {
	return []string{string(StatusFiring), string(StatusAcknowledged), string(StatusResolved), string(StatusSuppressed)}
}
