package auth

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"github.com/subrat-dwi/passman-cli/internal/ui/login"
)

func NewLoginCmd(app *app.App) *cobra.Command {
	loginCmd := &cobra.Command{
		Use:     "login",
		Short:   "Login into your Passman account",
		Long:    "Login into your account using your email and master password",
		Aliases: []string{"signin"},
		RunE: func(cmd *cobra.Command, args []string) error {

			p := tea.NewProgram(login.NewLoginModel(app))
			result, err := p.Run()
			if err != nil {
				return err
			}

			if m, ok := result.(login.LoginModel); ok && m.Success() {
				fmt.Println("Login Successful")
			}
			return nil
		},
	}
	return loginCmd
}
