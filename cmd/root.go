package cmd

import (
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/cmd/agent"
	"github.com/subrat-dwi/passman-cli/cmd/auth"
	"github.com/subrat-dwi/passman-cli/cmd/passwords"
	"github.com/subrat-dwi/passman-cli/internal/app"
)

func NewRootCmd(app *app.App) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "pman",
		Short: "Password Manager CLIent for Passman",
		Long:  "A secure CLI client for the Passman - Cloud Password Management Platform.\nDeveloped by Subrat | www.subratdwivedi.dev",
	}

	rootCmd.SuggestionsMinimumDistance = 2

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(auth.NewAuthCmd(app))
	rootCmd.AddCommand(passwords.NewListCmd(app))
	rootCmd.AddCommand(passwords.NewCreateCmd(app))
	rootCmd.AddCommand(agent.NewAgentCmd())

	return rootCmd
}
