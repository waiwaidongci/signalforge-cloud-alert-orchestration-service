package infrastructure

func FanoutWindow(funcs []func() error) int {
	return len(funcs)
}

func FanoutState(funcs []func() error) string {
	if len(funcs) == 0 {
		return "idle"
	}
	return "busy"
}

func FanoutLabel(funcs []func() error) string {
	if FanoutState(funcs) == "busy" {
		return "processing"
	}
	return "ready"
}
