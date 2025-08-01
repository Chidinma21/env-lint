package validator

import (
	"fmt"
	"strconv"
	"strings"
)

type SchemaRule struct {
	Type     string        `json:"type"`
	Required bool          `json:"required"`
	Default  interface{}   `json:"default"`
	Allowed  []interface{} `json:"allowed"`
}

type ValidationResult struct {
	Passed   bool
	Errors   map[string]string
	Warnings map[string]string
}

func ValidateEnv(envMap map[string]string, schema map[string]SchemaRule) ValidationResult {
	errors := make(map[string]string)
	warnings := make(map[string]string)

	for key, rule := range schema {
		value, ok := envMap[key]
		if !ok {
			if rule.Required {
				errors[key] = "Missing required key"
				continue
			} else {
				if rule.Default != nil {
					defaultStr := fmt.Sprintf("%v", rule.Default)
					envMap[key] = defaultStr
					value = defaultStr
					warnings[key] = "Missing optional key — using default"
				} else {
					warnings[key] = "Missing optional key (ok)"
					continue
				}
			}
		}

		if len(rule.Allowed) > 0 {
			isValid := false
			valueStr := fmt.Sprintf("%v", value)
			for _, allowed := range rule.Allowed {
				allowedStr := fmt.Sprintf("%v", allowed)
				if valueStr == allowedStr {
					isValid = true
					break
				}
			}
			if !isValid {
				errors[key] = fmt.Sprintf("Value '%s' is not allowed. Expected one of: %v", value, rule.Allowed)
			}
		}

		switch rule.Type {
		case "string":
			// string is always valid
		case "number":
			if _, err := strconv.Atoi(value); err != nil {
				errors[key] = fmt.Sprintf("Expected number but got: %s", value)
			}
		case "boolean":
			lower := strings.ToLower(value)
			if lower != "true" && lower != "false" {
				errors[key] = fmt.Sprintf("Expected boolean but got: %s", value)
			}
		default:
			warnings[key] = fmt.Sprintf("Unknown type '%s' — skipping check", rule.Type)
		}
	}

	return ValidationResult{
		Passed:   len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}
}
