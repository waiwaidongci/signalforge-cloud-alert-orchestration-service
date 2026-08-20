package persistence

import "github.com/acme/signalforge/internal/shared/apperr"

func classifyTimelineReadError(err error) error {
	if err == nil {
		return nil
	}
	return apperr.Internal("TIMELINE_READ_FAILED", "读取时间线失败")
}
