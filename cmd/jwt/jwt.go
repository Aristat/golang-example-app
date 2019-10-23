package jwt

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:           "jwt",
	Short:         "Tools for generate JWT",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	Cmd.AddCommand(tokenCmd)
}
