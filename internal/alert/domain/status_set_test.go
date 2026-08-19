package domain

import "testing"

func TestActiveAlertStatusesIncludesAcknowledged(t *testing.T) {
	statuses := ActiveAlertStatuses()
	for _, status := range statuses {
		if status == StatusAcknowledged {
			return
		}
	}
	t.Fatal("acknowledged should be included in active statuses")
}
