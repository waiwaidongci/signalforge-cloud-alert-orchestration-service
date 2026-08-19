package httpx

import (
	"strconv"
)

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Offset  int `json:"-"`
}

type PageMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func ParsePagination(values map[string][]string, defaults ...int) Pagination {
	perPage := 20
	if len(defaults) > 0 && defaults[0] > 0 {
		perPage = defaults[0]
	}
	page := atoiDefault(first(values["page"]), 1)
	size := atoiDefault(first(values["per_page"]), perPage)
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = perPage
	}
	if size > 100 {
		size = 100
	}
	return Pagination{Page: page, PerPage: size, Offset: (page - 1) * size}
}

func NewPageMeta(page Pagination, total int) PageMeta {
	pages := (total + page.PerPage - 1) / page.PerPage
	if pages < 1 {
		pages = 1
	}
	return PageMeta{Page: page.Page, PerPage: page.PerPage, Total: total, TotalPages: pages}
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func atoiDefault(raw string, fallback int) int {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
