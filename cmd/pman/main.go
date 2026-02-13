package main

import (
	"github.com/subrat-dwi/passman-cli/cmd"
	"github.com/subrat-dwi/passman-cli/internal/app"
)

func main() {
	app := app.New()

	rootCmd := cmd.NewRootCmd(app)
	rootCmd.Execute()
}
