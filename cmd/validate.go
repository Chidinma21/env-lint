package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/chidinma21/env-lint/internal/validator"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	success = color.New(color.FgGreen).SprintFunc()
	warn    = color.New(color.FgYellow).SprintFunc()
	fail    = color.New(color.FgHiRed).SprintFunc()
	debug   = color.New(color.FgCyan).SprintFunc()
)

var envFile string
var schemaFile string
var suppressWarnings bool
var strictMode bool
var failFast bool

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate your .env file against a defined JSON schema",
	Long: `env-lint validate checks your .env configuration against a JSON schema.

You can specify which keys are required and what type of value (string, number, boolean) each should have.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load .env file
		envMap, err := godotenv.Read(envFile)
		if err != nil {
			fmt.Printf("%s Failed to read .env file: %v\n", fail("❌"), err)
			os.Exit(1)
		}
		fmt.Println(success("🚀 .env file loaded successfully"))

		// Load schema
		schemaData, err := os.ReadFile(schemaFile)
		if err != nil {
			fmt.Printf("%s Failed to read schema file: %v\n", fail("❌"), err)
			os.Exit(1)
		}

		ext := filepath.Ext(schemaFile)
		var schema map[string]validator.SchemaRule
		switch ext {
		case ".json":
			if err := json.Unmarshal(schemaData, &schema); err != nil {
				fmt.Printf("%s Invalid JSON schema: %v\n", fail("❌"), err)
				os.Exit(1)
			}
		case ".yaml", ".yml":
			if err := yaml.Unmarshal(schemaData, &schema); err != nil {
				fmt.Printf("%s Invalid YAML schema: %v\n", fail("❌"), err)
				os.Exit(1)
			}
		default:
			fmt.Printf("%s Unsupported schema format: %s\n", fail("❌"), ext)
			os.Exit(1)
		}

		fmt.Println(success("🚀 schema file loaded successfully"))

		// Validate
		fmt.Println(debug("\n🔍 Validating environment variables..."))

		validateRes := validator.ValidateEnv(envMap, schema, failFast, strictMode)

		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		if !validateRes.Passed {
			for key, value := range validateRes.Errors {
				fmt.Printf("%-14s %-25s %s\n", fail("ERROR"), key, value)
			}
			if !suppressWarnings {
				printValidationWarnings(validateRes.Warnings)
			}

			if strictMode {
				fmt.Println("\n")
				fmt.Println(fail("❌ Strict Mode: Extra keys found in .env not in schema: "))
				for _, k := range validateRes.ExtraKeys {
					fmt.Printf("   - %s\n", k)
				}
			}

			fmt.Println("\n")
			fmt.Println(fail("❌ Validation failed. Please fix the errors above."))

			os.Exit(1)
		} else {
			fmt.Println(success("✅ All checks passed. Your .env config looks great!"))
			if !suppressWarnings {
				printValidationWarnings(validateRes.Warnings)
			}
			fmt.Println("\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringVarP(&envFile, "env", "e", ".env", "Path to the .env file")
	validateCmd.Flags().StringVarP(&schemaFile, "schema", "s", "schema.json", "Path to the schema file (JSON)")
	validateCmd.Flags().BoolVarP(&suppressWarnings, "suppress-warnings", "w", false, "Suppress warning messages in output")
	validateCmd.Flags().BoolVarP(&strictMode, "strict-mode", "t", false, "Fail if extra keys exist in .env that are not in schema")
	validateCmd.Flags().BoolVarP(&failFast, "fail-fast", "f", false, "Stop validation after the first error")
}

func printValidationWarnings(warnings map[string]string) {
	if len(warnings) == 0 {
		return
	}

	for key, value := range warnings {
		fmt.Printf("%-14s %-25s %s\n", warn("WARN"), key, value)
	}
}
