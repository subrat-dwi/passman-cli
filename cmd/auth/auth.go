package auth

import (
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/app"
)

func NewAuthCmd(app *app.App) *cobra.Command {
	AuthCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication Commands",
		Long:  "Commands related to user authentication like login and register",
	}

	AuthCmd.AddCommand(NewLoginCmd(app))
	AuthCmd.AddCommand(NewRegisterCmd(app))
	AuthCmd.AddCommand(NewLogoutCmd(app))

	return AuthCmd
}
