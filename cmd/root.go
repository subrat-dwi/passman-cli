package cmd

import (
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/cmd/auth"
	"github.com/subrat-dwi/passman-cli/cmd/passwords"
	"github.com/subrat-dwi/passman-cli/internal/app"
)

func NewRootCmd(app *app.App) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "pman",
		Short: "Password Manager CLI by github.com/subrat-dwi",
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(auth.NewAuthCmd(app))
	rootCmd.AddCommand(passwords.NewListCmd(app))

	return rootCmd
}
