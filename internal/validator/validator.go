package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SchemaRule struct {
	Type        string        `json:"type" yaml:"type"`
	Required    bool          `json:"required,omitempty" yaml:"required,omitempty"`
	Default     interface{}   `json:"default,omitempty" yaml:"default,omitempty"`
	Allowed     []interface{} `json:"allowed,omitempty" yaml:"allowed,omitempty"`
	Pattern     string        `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Length      *int          `json:"length,omitempty" yaml:"length,omitempty"`
	MaxLength   *int          `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength   *int          `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Min         *float64      `json:"min,omitempty" yaml:"min,omitempty"`
	Max         *float64      `json:"max,omitempty" yaml:"max,omitempty"`
	CustomError string        `json:"customError,omitempty" yaml:"customError,omitempty"`
}

type ValidationResult struct {
	Passed    bool
	Errors    map[string]string
	Warnings  map[string]string
	ExtraKeys []string
}

func ValidateEnv(envMap map[string]string, schema map[string]SchemaRule, failFast, strictMode bool) ValidationResult {
	errors := make(map[string]string)
	warnings := make(map[string]string)
	extraKeys := []string{}
	passed := true

	fail := func(key, msg string) ValidationResult {
		errors[key] = msg
		return ValidationResult{
			Passed:   false,
			Errors:   errors,
			Warnings: warnings,
		}
	}

	for key, rule := range schema {
		value, ok := envMap[key]

		// Missing required key
		if !ok {
			if rule.Required {
				if failFast {
					return fail(key, "Missing required key")
				}
				errors[key] = "Missing required key"
				continue
			}
			// Optional key handling
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

		// Allowed values check
		if len(rule.Allowed) > 0 {
			valueStr := fmt.Sprintf("%v", value)
			valid := false
			for _, allowed := range rule.Allowed {
				if valueStr == fmt.Sprintf("%v", allowed) {
					valid = true
					break
				}
			}
			if !valid {
				msg := fmt.Sprintf("Value '%s' is not allowed. Expected one of: %v", value, rule.Allowed)
				if rule.CustomError != "" {
					msg = rule.CustomError
				}
				if failFast {
					return fail(key, msg)
				}
				errors[key] = msg
			}
		}

		// Type checks
		switch rule.Type {
		case "string":
			if rule.Pattern != "" {
				if matched, err := regexp.MatchString(rule.Pattern, value); err != nil {
					warnings[key] = fmt.Sprintf("Invalid regex pattern: %s", rule.Pattern)
				} else if !matched {
					msg := fmt.Sprintf("Value does not match pattern: %s", rule.Pattern)
					if rule.CustomError != "" {
						msg = rule.CustomError
					}
					if failFast {
						return fail(key, msg)
					}
					errors[key] = msg
				}
			}
			if rule.Length != nil && len(value) != *rule.Length {
				msg := fmt.Sprintf("Expected string of length [%v] but got: %s", *rule.Length, value)
				if rule.CustomError != "" {
					msg = rule.CustomError
				}
				if failFast {
					return fail(key, msg)
				}
				errors[key] = msg
			}
			if rule.MaxLength != nil && len(value) > *rule.MaxLength {
				msg := fmt.Sprintf("Expected max length [%v] but got: %s", *rule.MaxLength, value)
				if rule.CustomError != "" {
					msg = rule.CustomError
				}
				if failFast {
					return fail(key, msg)
				}
				errors[key] = msg
			}
			if rule.MinLength != nil && len(value) < *rule.MinLength {
				msg := fmt.Sprintf("Expected min length [%v] but got: %s", *rule.MinLength, value)
				if rule.CustomError != "" {
					msg = rule.CustomError
				}
				if failFast {
					return fail(key, msg)
				}
				errors[key] = msg
			}

		case "number":
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				msg := fmt.Sprintf("Expected number but got: %s", value)
				if rule.CustomError != "" {
					msg = rule.CustomError
				}
				if failFast {
					return fail(key, msg)
				}
				errors[key] = msg
				continue
			}
			if rule.Min != nil && num < *rule.Min {
				msg := fmt.Sprintf("Expected number >= %.2f but got: %.2f", *rule.Min, num)
				if rule.CustomError != "" {
					msg = rule.CustomError
				}
				if failFast {
					return fail(key, msg)
				}
				errors[key] = msg
			}
			if rule.Max != nil && num > *rule.Max {
				msg := fmt.Sprintf("Expected number <= %.2f but got: %.2f", *rule.Max, num)
				if rule.CustomError != "" {
					msg = rule.CustomError
				}
				if failFast {
					return fail(key, msg)
				}
				errors[key] = msg
			}

		case "boolean":
			lower := strings.ToLower(value)
			if lower != "true" && lower != "false" {
				msg := fmt.Sprintf("Expected boolean but got: %s", value)
				if rule.CustomError != "" {
					msg = rule.CustomError
				}
				if failFast {
					return fail(key, msg)
				}
				errors[key] = msg
			}

		default:
			warnings[key] = fmt.Sprintf("Unknown type '%s' — skipping check", rule.Type)
		}
	}

	if strictMode {
		for key := range envMap {
			if _, exists := schema[key]; !exists {
				extraKeys = append(extraKeys, key)
			}
		}

		if len(extraKeys) > 0 {
			passed = false
		}
	}

	passed = len(errors) == 0 && passed

	return ValidationResult{
		Passed:    passed,
		Errors:    errors,
		Warnings:  warnings,
		ExtraKeys: extraKeys,
	}
}
