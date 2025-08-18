package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/chidinma21/env-lint/internal/validator"
	"github.com/chidinma21/env-lint/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var envPath string
var format string

func generateSchemaFromMap(envMap map[string]string) map[string]validator.SchemaRule {
	schema := make(map[string]validator.SchemaRule)

	for key, value := range envMap {
		schema[key] = validator.SchemaRule{
			Type:     utils.GuessType(value),
			Required: false,
			Default:  value,
		}
	}

	return schema
}

var generateSchemaCmd = &cobra.Command{
	Use:   "generate-schema",
	Short: "Generate schema from an existing .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		envMap, err := godotenv.Read(envPath)
		if err != nil {
			return fmt.Errorf("error reading env file: %v", err)
		}

		schema := generateSchemaFromMap(envMap)

		switch format {
		case "json":
			out, _ := json.MarshalIndent(schema, "", "   ")
			fmt.Println(string(out))
		case "yaml", "yml":
			out, _ := yaml.Marshal(schema)
			fmt.Println(string(out))
		default:
			return fmt.Errorf("unsupported format: %s", format)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateSchemaCmd)

	generateSchemaCmd.Flags().StringVarP(&envPath, "env", "e", ".env", "Path to the .env file for schema generation")
	generateSchemaCmd.Flags().StringVarP(&format, "format", "f", "json", "Output format: json, yaml, yml")
}
