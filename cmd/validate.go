package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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

		var schema map[string]validator.SchemaRule
		if err := json.Unmarshal(schemaData, &schema); err != nil {
			fmt.Printf("%s Invalid JSON schema: %v\n", fail("ERROR"), err)
			os.Exit(1)
		}
		fmt.Println(success("🚀 schema file loaded successfully"))

		// Validate
		fmt.Println(debug("\n🔍 Validating environment variables...\n"))

		validateRes := validator.ValidateEnv(envMap, schema)

		for key, value := range validateRes.Errors {
			fmt.Printf("%-14s %-25s %s\n", fail("ERROR"), key, value)
		}
		for key, value := range validateRes.Warnings {
			fmt.Printf("%-14s %-25s %s\n", warn("WARN"), key, value)
		}

		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		if !validateRes.Passed {
			fmt.Println(fail("❌ Validation failed. Please fix the errors above."))
			os.Exit(1)
		} else {
			fmt.Println(success("✅ All checks passed. Your .env config looks great!"))
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringVarP(&envFile, "env", "e", ".env", "Path to the .env file")
	validateCmd.Flags().StringVarP(&schemaFile, "schema", "s", "schema.json", "Path to the schema file (JSON)")
}
