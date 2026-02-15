package passwords

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/app"
	createui "github.com/subrat-dwi/passman-cli/internal/ui/create"
)

func NewCreateCmd(app *app.App) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new password",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(createui.NewCreatePasswordModel(app))
			_, err := p.Run()
			return err
		},
	}

	return createCmd
}
