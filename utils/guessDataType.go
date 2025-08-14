package utils

import (
	"strconv"
	"strings"
)

func GuessType(val string) string {
	lower := strings.ToLower(val)

	if lower == "true" || lower == "false" {
		return "boolean"
	}
	if _, err := strconv.Atoi(val); err == nil {
		return "number"
	}
	if _, err := strconv.ParseFloat(val, 64); err == nil {
		return "number"
	}
	return "string"
}
