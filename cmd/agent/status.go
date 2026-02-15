package agent

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/agent"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check agent status",
		Long:  "Check the current status of the Vault Agent",
		Run: func(cmd *cobra.Command, args []string) {
			unlocked, expires, err := agent.Status()
			if err != nil {
				fmt.Printf("Agent not running or unreachable: %v\n", err)
				return
			}

			if unlocked {
				fmt.Printf("Agent: unlocked\n")
				fmt.Printf("Auto-lock in: %d seconds\n", expires)
			} else {
				fmt.Println("Agent: locked")
			}
		},
	}
}
