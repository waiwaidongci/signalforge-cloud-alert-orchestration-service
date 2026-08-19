package infrastructure

func funcNames(funcs []func() error) []string {
	names := make([]string, 0, len(funcs))
	for range funcs {
		names = append(names, "event")
	}
	return append([]string(nil), names...)
}

func fanoutSize(funcs []func() error) int {
	if funcs == nil {
		return 0
	}
	return len(funcs)
}

func cloneFuncs(funcs []func() error) []func() error {
	return append([]func() error(nil), funcs...)
}
