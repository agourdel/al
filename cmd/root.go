package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "al",
	Short: "Al - CLI for managing client projects",
	Long:  `Al is a CLI tool to help manage client projects with notes, links, and shortcuts.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Project commands
	initCmd.GroupID = "project"
	goCmd.GroupID = "project"
	noteCmd.GroupID = "project"
	linkCmd.GroupID = "project"
	
	// Setup commands
	installCmd.GroupID = "setup"
	updateCmd.GroupID = "setup"
	
	// Add command groups
	rootCmd.AddGroup(&cobra.Group{
		ID:    "project",
		Title: "Project Commands:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "setup",
		Title: "Setup Commands:",
	})
	
	// Add commands (order matters within groups)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(goCmd)
	rootCmd.AddCommand(noteCmd)
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(updateCmd)
}
