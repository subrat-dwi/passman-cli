package cmd

import (
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/cmd/agent"
	"github.com/subrat-dwi/passman-cli/cmd/auth"
	"github.com/subrat-dwi/passman-cli/cmd/passwords"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"github.com/subrat-dwi/passman-cli/internal/service"
	"github.com/subrat-dwi/passman-cli/internal/ui/styles"
)

func NewRootCmd(app *app.App) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:     "pman",
		Short:   "Password Manager CLIent for Passman",
		Long:    "A secure CLI client for the Passman - Cloud Password Manager.\nOfficial Website: passman.subratdwivedi.dev",
		Version: service.Version(),
	}

	// Custom version template with styling
	rootCmd.SetVersionTemplate(styles.Version("{{.Version}}") + "\n")

	rootCmd.SuggestionsMinimumDistance = 2

	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(auth.NewAuthCmd(app))
	rootCmd.AddCommand(passwords.NewListCmd(app))
	rootCmd.AddCommand(passwords.NewCreateCmd(app))
	rootCmd.AddCommand(agent.NewAgentCmd())

	return rootCmd
}
