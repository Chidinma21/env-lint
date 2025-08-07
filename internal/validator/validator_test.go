package validator

import "testing"

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func TestValidateEnv(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		schema    map[string]SchemaRule
		wantPass  bool
		wantErrs  map[string]string
		wantWarns map[string]string
	}{
		{
			name: "All keys valid",
			env: map[string]string{
				"PORT":     "3000",
				"DEBUG":    "true",
				"APP_NAME": "MyApp",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
				},
				"DEBUG": {
					Type:     "boolean",
					Required: false,
				},
				"APP_NAME": {
					Type:     "string",
					Required: true,
				},
			},
			wantPass:  true,
			wantErrs:  map[string]string{},
			wantWarns: map[string]string{},
		},
		{
			name: "Missing required key",
			env: map[string]string{
				"DEBUG":    "true",
				"APP_NAME": "MyApp",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
				},
				"DEBUG": {
					Type:     "boolean",
					Required: false,
				},
				"APP_NAME": {
					Type:     "string",
					Required: true,
				},
			},
			wantPass: false,
			wantErrs: map[string]string{
				"PORT": "Missing required key",
			},
			wantWarns: map[string]string{},
		},
		{
			name: "Missing optional key",
			env: map[string]string{
				"PORT":     "3000",
				"APP_NAME": "MyApp",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
				},
				"DEBUG": {
					Type:     "boolean",
					Required: false,
				},
				"APP_NAME": {
					Type:     "string",
					Required: true,
				},
			},
			wantPass: true,
			wantErrs: map[string]string{},
			wantWarns: map[string]string{
				"DEBUG": "Missing optional key (ok)",
			},
		},
		{
			name: "Wrong type",
			env: map[string]string{
				"PORT":     "Hello World",
				"DEBUG":    "truth",
				"APP_NAME": "MyApp",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
				},
				"DEBUG": {
					Type:     "boolean",
					Required: false,
				},
				"APP_NAME": {
					Type:     "string",
					Required: true,
				},
			},
			wantPass: false,
			wantErrs: map[string]string{
				"PORT":  "Expected number but got: Hello World",
				"DEBUG": "Expected boolean but got: truth",
			},
			wantWarns: map[string]string{},
		},
		{
			name: "Use default with correct type",
			env: map[string]string{
				"PORT":     "3000",
				"APP_NAME": "MyApp",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
				},
				"DEBUG": {
					Type:     "boolean",
					Required: false,
					Default:  true,
				},
				"APP_NAME": {
					Type:     "string",
					Required: true,
				},
			},
			wantPass: true,
			wantErrs: map[string]string{},
			wantWarns: map[string]string{
				"DEBUG": "Missing optional key — using default",
			},
		},
		{
			name: "Use default with wrong type",
			env: map[string]string{
				"PORT":     "3000",
				"APP_NAME": "MyApp",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
				},
				"DEBUG": {
					Type:     "boolean",
					Required: false,
					Default:  "truth",
				},
				"APP_NAME": {
					Type:     "string",
					Required: true,
				},
			},
			wantPass: false,
			wantErrs: map[string]string{
				"DEBUG": "Expected boolean but got: truth",
			},
			wantWarns: map[string]string{
				"DEBUG": "Missing optional key — using default",
			},
		},
		{
			name: "Allowed keys - invalid",
			env: map[string]string{
				"PORT":     "3000",
				"DEBUG":    "true",
				"APP_NAME": "MyApps",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
					Allowed:  []interface{}{3001, 3002},
				},
				"DEBUG": {
					Type:     "boolean",
					Required: false,
				},
				"APP_NAME": {
					Type:     "string",
					Required: true,
					Allowed:  []interface{}{"MyApp", "NewApp"},
				},
			},
			wantPass: false,
			wantErrs: map[string]string{
				"APP_NAME": "Value 'MyApps' is not allowed. Expected one of: [MyApp NewApp]",
				"PORT":     "Value '3000' is not allowed. Expected one of: [3001 3002]",
			},
			wantWarns: map[string]string{},
		},
		{
			name: "Allowed keys - valid",
			env: map[string]string{
				"PORT":     "3000",
				"DEBUG":    "true",
				"APP_NAME": "MyApp",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
					Allowed:  []interface{}{3000, 3002},
				},
				"DEBUG": {
					Type:     "boolean",
					Required: false,
				},
				"APP_NAME": {
					Type:     "string",
					Required: true,
					Allowed:  []interface{}{"MyApp", "NewApp"},
				},
			},
			wantPass:  true,
			wantErrs:  map[string]string{},
			wantWarns: map[string]string{},
		},
		{
			name: "Valid pattern",
			env: map[string]string{
				"EMAIL": "hello@email.com",
			},
			schema: map[string]SchemaRule{
				"EMAIL": {
					Type:     "string",
					Required: true,
					Pattern:  `^[\w.-]+@[\w.-]+\.\w+$`,
				},
			},
			wantPass:  true,
			wantErrs:  map[string]string{},
			wantWarns: map[string]string{},
		},
		{
			name: "Pattern mismatch",
			env: map[string]string{
				"EMAIL": "not-an-email",
			},
			schema: map[string]SchemaRule{
				"EMAIL": {
					Type:     "string",
					Required: true,
					Pattern:  `^[\w.-]+@[\w.-]+\.\w+$`,
				},
			},
			wantPass: false,
			wantErrs: map[string]string{
				"EMAIL": "Value does not match pattern: ^[\\w.-]+@[\\w.-]+\\.\\w+$",
			},
			wantWarns: map[string]string{},
		},
		{
			name: "Invalid pattern",
			env: map[string]string{
				"EMAIL": "not-an-email",
			},
			schema: map[string]SchemaRule{
				"EMAIL": {
					Type:     "string",
					Required: true,
					Pattern:  `*invalid[`,
				},
			},
			wantPass: true,
			wantErrs: map[string]string{},
			wantWarns: map[string]string{
				"EMAIL": "Invalid regex pattern: *invalid[",
			},
		},
		{
			name: "Valid length",
			env: map[string]string{
				"PHONE": "08020718222",
			},
			schema: map[string]SchemaRule{
				"PHONE": {
					Type:     "string",
					Required: true,
					Length:   IntPtr(11),
				},
			},
			wantPass:  true,
			wantErrs:  map[string]string{},
			wantWarns: map[string]string{},
		},
		{
			name: "Valid length",
			env: map[string]string{
				"PHONE": "08020718222",
			},
			schema: map[string]SchemaRule{
				"PHONE": {
					Type:     "string",
					Required: true,
					Length:   IntPtr(12),
				},
			},
			wantPass: false,
			wantErrs: map[string]string{
				"PHONE": "Expected string of length [12] but got: 08020718222",
			},
			wantWarns: map[string]string{},
		},
		{
			name: "Valid Min-Max Range",
			env: map[string]string{
				"PORT": "3000",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
					Min:      Float64Ptr(3000),
					Max:      Float64Ptr(9999),
				},
			},
			wantPass:  true,
			wantErrs:  map[string]string{},
			wantWarns: map[string]string{},
		},
		{
			name: "Exceed Max Range",
			env: map[string]string{
				"PORT": "10000",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
					Min:      Float64Ptr(3000),
					Max:      Float64Ptr(9999),
				},
			},
			wantPass: false,
			wantErrs: map[string]string{
				"PORT": "Expected number <= 9999.00 but got: 10000.00",
			},
			wantWarns: map[string]string{},
		},
		{
			name: "Exceed Max Range",
			env: map[string]string{
				"PORT": "1999",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:     "number",
					Required: true,
					Min:      Float64Ptr(3000),
					Max:      Float64Ptr(9999),
				},
			},
			wantPass: false,
			wantErrs: map[string]string{
				"PORT": "Expected number >= 3000.00 but got: 1999.00",
			},
			wantWarns: map[string]string{},
		},
		{
			name: "Exceed Max Range",
			env: map[string]string{
				"PORT": "1999",
			},
			schema: map[string]SchemaRule{
				"PORT": {
					Type:        "number",
					Required:    true,
					Min:         Float64Ptr(3000),
					CustomError: "Port error occurred",
				},
			},
			wantPass: false,
			wantErrs: map[string]string{
				"PORT": "Port error occurred",
			},
			wantWarns: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateEnv(tt.env, tt.schema)

			if got.Passed != tt.wantPass {
				t.Errorf("Expected pass = %v, got %v", tt.wantPass, got.Passed)
			}

			for k, v := range tt.wantErrs {
				if got.Errors[k] != v {
					t.Errorf("Expected error on %s: %s, got: %s", k, v, got.Errors[k])
				}
			}

			for k, v := range tt.wantWarns {
				if got.Warnings[k] != v {
					t.Errorf("Expected warning on %s: %s, got: %s", k, v, got.Warnings[k])
				}
			}
		})
	}
}
