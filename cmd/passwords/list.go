package passwords

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"github.com/subrat-dwi/passman-cli/internal/ui/list"
)

func NewListCmd(app *app.App) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all your saved passwords",
		RunE: func(cmd *cobra.Command, args []string) error {
			passwords, err := app.PasswordService.List()
			if err != nil {
				return err
			}

			p := tea.NewProgram(list.NewListModel(app, passwords))
			_, err = p.Run()
			return err
		},
	}

	return listCmd
}
