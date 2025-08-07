package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SchemaRule struct {
	Type      string        `json:"type"`
	Required  bool          `json:"required"`
	Default   interface{}   `json:"default"`
	Allowed   []interface{} `json:"allowed"`
	Pattern   string        `json:"pattern"`
	Length    *int          `json:"length"`
	MaxLength *int          `json:"maxLength"`
	MinLength *int          `json:"minLength"`
	Min       *float64      `json:"min"`
	Max       *float64      `json:"max"`
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
			matched, err := regexp.MatchString(rule.Pattern, value)
			if err != nil {
				warnings[key] = fmt.Sprintf("Invalid regex pattern: %s", rule.Pattern)
			} else if !matched {
				errors[key] = fmt.Sprintf("Value does not match pattern: %s", rule.Pattern)
			}

			if rule.Length != nil && len(value) != *rule.Length {
				errors[key] = fmt.Sprintf("Expected string of length [%v] but got: %s", *rule.Length, value)
			}

			if rule.MaxLength != nil && len(value) > *rule.MaxLength {
				errors[key] = fmt.Sprintf("Expected max length [%v] but got: %s", *rule.MaxLength, value)
			}

			if rule.MinLength != nil && len(value) < *rule.MinLength {
				errors[key] = fmt.Sprintf("Expected min length [%v] but got: %s", *rule.MinLength, value)
			}
		case "number":
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				errors[key] = fmt.Sprintf("Expected number but got: %s", value)
			}

			if rule.Min != nil && num < *rule.Min {
				errors[key] = fmt.Sprintf("Expected number >= %.2f but got: %.2f", *rule.Min, num)
			}

			if rule.Max != nil && num > *rule.Max {
				errors[key] = fmt.Sprintf("Expected number <= %.2f but got: %.2f", *rule.Max, num)
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
