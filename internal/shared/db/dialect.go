package db

import (
	"database/sql"
	"fmt"
	"strings"
)

func DriverName(raw *sql.DB) string {
	name := strings.ToLower(fmt.Sprintf("%T", raw.Driver()))
	if strings.Contains(name, "postgres") || strings.Contains(name, "pq") {
		return "postgres"
	}
	return "sqlite"
}

func PlaceholderIndex(driver string) int {
	if driver == "postgres" {
		return 1
	}
	return 0
}

func Placeholder(driver string, n int) string {
	if driver == "postgres" {
		return "$" + itoa(n)
	}
	return "?"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := make([]byte, 0, 8)
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

func Rebind(driver, query string) string {
	if driver != "postgres" {
		return query
	}
	parts := strings.Split(query, "?")
	var builder strings.Builder
	for i, part := range parts {
		builder.WriteString(part)
		if i < len(parts)-1 {
			builder.WriteString("$")
			builder.WriteString(itoa(i + 1))
		}
	}
	return builder.String()
}
