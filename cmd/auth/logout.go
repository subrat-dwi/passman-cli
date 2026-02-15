package auth

import (
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/agent"
	"github.com/subrat-dwi/passman-cli/internal/app"
)

func NewLogoutCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:     "logout",
		Short:   "Logout from Passman",
		Long:    "Logout and clear your credentials",
		Aliases: []string{"signout", "out"},
		Run: func(cmd *cobra.Command, args []string) {
			app.AuthService.Storage.DeleteAccessToken()
			app.AuthService.Storage.DeleteSalt()
			agent.Lock()
			app.ResetState()
			println("Logged out successfully. Vault Agent locked and credentials cleared.")
		},
	}
}
