/*
Copyright © 2025 Chidinma Onuora <chidinmaonuora1@gmail.com>
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "env-lint",
	Short: "Validate your .env file against a JSON schema",
	Long: `env-lint is a CLI tool that helps you validate environment variables 
defined in a .env file against a user-defined JSON schema.

It supports checking for required keys and verifying types like string, number, and boolean.

Example usage:
  env-lint validate -s schema.json
`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {}
