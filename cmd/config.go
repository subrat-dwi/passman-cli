// cmd/config.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/config" // adjust import path
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Passman CLI configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]
		if err := config.Set(key, value); err != nil {
			return fmt.Errorf("failed to set config: %w", err)
		}
		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		value := config.Get(args[0])
		fmt.Printf("api_base_url: %s\n", value)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd, configGetCmd)
}
