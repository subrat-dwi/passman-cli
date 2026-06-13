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
			if err := app.AuthService.Storage.DeleteAccessToken(); err != nil {
				println("Failed to delete access token:", err.Error())
				return
			}
			if err := app.AuthService.Storage.DeleteSalt(); err != nil {
				println("Failed to delete salt:", err.Error())
				return
			}
			if err := agent.Lock(); err != nil {
				println("Failed to lock vault agent:", err.Error())
				return
			}
			app.ResetState()
			println("Logged out successfully. Vault Agent locked and credentials cleared.")
		},
	}
}
