package validator

import "testing"

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
					Default: true,
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
					Default: "truth",
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
