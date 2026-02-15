package auth

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"github.com/subrat-dwi/passman-cli/internal/ui/register"
)

func NewRegisterCmd(app *app.App) *cobra.Command {
	registerCmd := &cobra.Command{
		Use:     "register",
		Short:   "Register a new account",
		Long:    "Register a new account with your email and master password",
		Aliases: []string{"signup", "join"},
		RunE: func(cmd *cobra.Command, args []string) error {

			p := tea.NewProgram(register.NewRegisterModel(app))
			result, err := p.Run()
			if err != nil {
				return err
			}

			if m, ok := result.(register.RegisterModel); ok && m.Success() {
				fmt.Println("Registration Successful")
			}
			return nil
		},
	}
	return registerCmd
}
