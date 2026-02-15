package agent

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/agent"
)

func newLockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lock",
		Short: "Lock the agent (forget encryption key)",
		Run: func(cmd *cobra.Command, args []string) {
			err := agent.Lock()
			if err != nil {
				fmt.Printf("Failed to lock agent: %v\n", err)
				return
			}
			fmt.Println("Vault Agent locked. Key wiped from memory.")
		},
	}
}
