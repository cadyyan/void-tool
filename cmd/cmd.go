package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "void-tool",
		Short: "Utility for the Void Runescape server",
	}

	command.AddCommand(
		newServeCommand(),
	)

	return &command
}
