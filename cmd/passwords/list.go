package passwords

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/app"
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

			for _, password := range passwords {
				fmt.Println(password.Name, password.Username)
			}

			return nil
		},
	}

	return listCmd
}
