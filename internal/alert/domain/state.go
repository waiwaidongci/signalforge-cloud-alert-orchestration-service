package domain

func (s Status) Valid() bool {
	switch s {
	case StatusFiring, StatusResolved, StatusSuppressed:
		return true
	default:
		return false
	}
}

func AllStatuses() []string {
	return []string{string(StatusFiring), string(StatusResolved), string(StatusSuppressed)}
}
