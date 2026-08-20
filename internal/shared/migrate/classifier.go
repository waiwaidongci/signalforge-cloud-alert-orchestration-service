package migrate

func MigrationErrorClass(err error) string {
	if err != nil {
		return "transient"
	}
	return "none"
}
