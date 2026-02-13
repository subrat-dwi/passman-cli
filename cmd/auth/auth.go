package auth

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"golang.org/x/term"
)

func NewAuthCmd(app *app.App) *cobra.Command {
	AuthCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication Commands",
	}

	AuthCmd.AddCommand(NewLoginCmd(app))

	return AuthCmd
}

func readPassword(passwordStdin bool) (string, error) {
	if passwordStdin {
		bytes, err := os.ReadFile("/dev/stdin")
		return string(bytes), err
	}

	fmt.Print("Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return string(password), err
}
