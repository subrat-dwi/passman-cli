package auth

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/app"
)

func NewLoginCmd(app *app.App) *cobra.Command {
	var (
		email         string
		passwordStdin bool
	)

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login into your account",
		RunE: func(cmd *cobra.Command, args []string) error {
			if email == "" {
				return errors.New("email required")
			}

			password, err := readPassword(passwordStdin)
			if err != nil {
				return err
			}

			if err = app.AuthService.Login(email, password); err != nil {
				return err
			}

			fmt.Println("Login Successful")
			return nil
		},
	}

	loginCmd.Flags().StringVarP(&email, "email", "e", "", "Email for login")
	loginCmd.Flags().BoolVarP(&passwordStdin, "password-stdin", "", false, "Read password from stdin")
	loginCmd.MarkFlagRequired("email")

	return loginCmd
}
