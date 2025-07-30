package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

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

type SchemaRule struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

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

		var schema map[string]SchemaRule
		if err := json.Unmarshal(schemaData, &schema); err != nil {
			fmt.Printf("%s Invalid JSON schema: %v\n", fail("ERROR"), err)
			os.Exit(1)
		}
		fmt.Println(success("🚀 schema file loaded successfully"))

		// Validate
		fmt.Println(debug("\n🔍 Validating environment variables...\n"))
		failed := false

		for key, rule := range schema {
			value, exists := envMap[key]

			if !exists {
				if rule.Required {
					fmt.Printf("%-14s %-25s %s\n", fail("ERROR"), key, "Missing required key")
					failed = true
				} else {
					fmt.Printf("%-14s %-25s %s\n", warn("WARN"), key, "Optional key missing (ok)")
				}
				continue
			}

			switch rule.Type {
			case "string":
				// always valid
			case "number":
				if _, err := strconv.Atoi(value); err != nil {
					fmt.Printf("%-14s %-25s %s\n", fail("ERROR"), key, fmt.Sprintf("Expected number but got: %s", value))
					failed = true
				}
			case "boolean":
				val := strings.ToLower(value)
				if val != "true" && val != "false" {
					fmt.Printf("%-14s %-25s %s\n", fail("ERROR"), key, fmt.Sprintf("Expected boolean but got: %s", value))
					failed = true
				}
			default:
				fmt.Printf("%-14s %-25s %s\n", warn("WARN"), key, fmt.Sprintf("Unknown type '%s' — skipping check", rule.Type))
			}
		}

		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		if failed {
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
