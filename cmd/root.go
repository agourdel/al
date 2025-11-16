package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "al",
	Short: "Al - CLI for managing client projects",
	Long:  `Al is a CLI tool to help manage client projects with notes, links, and shortcuts.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(goCmd)
	rootCmd.AddCommand(noteCmd)
	rootCmd.AddCommand(linkCmd)
}
