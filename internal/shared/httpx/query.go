package httpx

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/acme/signalforge/internal/shared/apperr"
)

func QueryInt(r *http.Request, key string) (int, bool, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return 0, false, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false, apperr.BadRequest("INVALID_QUERY", key+" 必须是整数")
	}
	return value, true, nil
}

func QueryBool(r *http.Request, key string) (bool, bool, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return false, false, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, false, apperr.BadRequest("INVALID_QUERY", key+" 必须是布尔值")
	}
	return value, true, nil
}

func RequirePathValue(r *http.Request, key string) (string, error) {
	value := r.PathValue(key)
	if strings.TrimSpace(value) == "" {
		return "", apperr.BadRequest("MISSING_PATH_PARAMETER", key+" 不能为空")
	}
	return value, nil
}

func DecodeJSON(r *http.Request, target any) error {
	return DecodeWithLimit(r, target)
}
