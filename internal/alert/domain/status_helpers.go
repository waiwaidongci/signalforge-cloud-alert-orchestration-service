package domain

func NextStatus(action string) Status {
	if action == "acknowledge" {
		return StatusAcknowledged
	}
	if action == "close" {
		return StatusResolved
	}
	return StatusFiring
}

func IsActive(status Status) bool {
	return status == StatusFiring || status == StatusAcknowledged || status == StatusSuppressed
}

func StatusLabel(status Status) string {
	switch status {
	case StatusFiring:
		return "告警中"
	case StatusAcknowledged:
		return "已确认"
	case StatusSuppressed:
		return "已抑制"
	case StatusResolved:
		return "已恢复"
	default:
		return "未知"
	}
}
