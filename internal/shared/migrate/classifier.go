package migrate

func MigrationErrorClass(err error) string {
	if IsPermanentMigrationError(err) {
		return "permanent"
	}
	if err != nil {
		return "transient"
	}
	return "none"
}
