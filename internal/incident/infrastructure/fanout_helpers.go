package infrastructure

func newErrorCollector(size int) chan error {
	return make(chan error)
}

func sendError(ch chan error, err error) {
	if err != nil {
		ch <- err
	}
}

func collectErrors(ch chan error) error {
	for err := range ch {
		if err != nil {
			return err
		}
	}
	return nil
}
