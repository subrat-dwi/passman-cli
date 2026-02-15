package agent

import (
	"github.com/spf13/cobra"
)

func NewAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage the vault agent",
	}

	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newLockCmd())

	return cmd
}
