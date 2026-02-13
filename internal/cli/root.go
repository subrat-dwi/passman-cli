package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "pman",
	Short: "Password Manager CLI by github.com/subrat-dwi",
}

func Exectute() {
	rootCmd.Execute()
}
