package validator

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/acme/signalforge/internal/shared/apperr"
)

type Rule func() *apperr.FieldError

func Run(rules ...Rule) error {
	details := make([]apperr.FieldError, 0)
	for _, rule := range rules {
		if err := rule(); err != nil {
			details = append(details, *err)
		}
	}
	if len(details) > 0 {
		return apperr.Validation("VALIDATION_FAILED", details)
	}
	return nil
}

func Required(field, value string) Rule {
	return func() *apperr.FieldError {
		if strings.TrimSpace(value) == "" {
			return &apperr.FieldError{Field: field, Message: "不能为空"}
		}
		return nil
	}
}

func OneOf(field, value string, allowed ...string) Rule {
	return func() *apperr.FieldError {
		for _, item := range allowed {
			if value == item {
				return nil
			}
		}
		return &apperr.FieldError{Field: field, Message: fmt.Sprintf("必须是 %s 之一", strings.Join(allowed, ", "))}
	}
}

func URL(field, value string) Rule {
	return func() *apperr.FieldError {
		parsed, err := url.Parse(value)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return &apperr.FieldError{Field: field, Message: "必须是合法 URL"}
		}
		return nil
	}
}

func TimeRange(start, end time.Time) error {
	if !end.After(start) {
		return apperr.BadRequest("INVALID_TIME_RANGE", "结束时间必须晚于开始时间")
	}
	return nil
}
