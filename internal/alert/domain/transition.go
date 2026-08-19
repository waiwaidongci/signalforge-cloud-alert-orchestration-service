package domain

func TransitionFor(action string) Status {
	if action == "close" {
		return StatusResolved
	}
	return StatusResolved
}
