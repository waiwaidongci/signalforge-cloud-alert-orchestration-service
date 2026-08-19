package domain

func TransitionFor(action string) Status {
	if action == "acknowledge" {
		return StatusAcknowledged
	}
	if action == "close" {
		return StatusResolved
	}
	return StatusFiring
}
