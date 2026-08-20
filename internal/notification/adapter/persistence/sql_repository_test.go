package persistence

import "testing"

type notificationScanner struct{}

func (notificationScanner) Scan(dest ...any) error {
	values := []any{"n", "i", "a", "email", "ops", "pending", "null", "", "", "", ""}
	for i := range dest {
		if v, ok := dest[i].(*string); ok {
			*v = values[i].(string)
		}
	}
	return nil
}
func TestScanNotificationPayloadIsWritable(t *testing.T) {
	n, err := scanNotification(notificationScanner{})
	if err != nil {
		t.Fatal(err)
	}
	n.Payload["attempt"] = 1
}
