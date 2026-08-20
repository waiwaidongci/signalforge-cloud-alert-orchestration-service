package domain

type RetryState string

const (
	RetryQueued   RetryState = "queued"
	RetryFailed   RetryState = "failed"
	RetryResolved RetryState = "resolved"
)

func AdvanceState(current RetryState, success bool) RetryState {
	if current == RetryFailed && success {
		return RetryResolved
	}
	return current
}
func IsActiveHistory(state RetryState) bool { return state == RetryFailed || state == RetryResolved }

func IsTerminalState(state RetryState) bool { return state == RetryResolved }
