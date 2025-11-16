package cmd

import (
	"fmt"


	"github.com/alex/al/storage"
	"github.com/alex/al/utils"
	"github.com/spf13/cobra"
)

var goCmd = &cobra.Command{
	Use:   "go [shortcut]",
	Short: "Copy project path to clipboard",
	Long: `Copy the project directory path to clipboard using its name or any of its shortcuts.
	
Example: al go myproject
Then: cd <paste>`,
	Args: cobra.ExactArgs(1),
	RunE: runGo,
}

func runGo(cmd *cobra.Command, args []string) error {
	shortcut := args[0]

	// Find project
	name, project, err := storage.FindProjectByShortcut(shortcut)
	if err != nil {
		// Try to find similar projects
		projects, loadErr := storage.LoadProjects()
		if loadErr != nil {
			return fmt.Errorf("project '%s' not found", shortcut)
		}

		// Collect all possible shortcuts
		var allShortcuts []string
		for projName, proj := range projects {
			allShortcuts = append(allShortcuts, projName)
			allShortcuts = append(allShortcuts, proj.Shortcuts...)
		}

		similar := utils.FindSimilarStrings(shortcut, allShortcuts, 3)
		if len(similar) > 0 {
			fmt.Printf("Project '%s' not found. Did you mean:\n", shortcut)
			for _, s := range similar {
				fmt.Printf("  - %s\n", s)
			}
			return fmt.Errorf("project not found")
		}

		return fmt.Errorf("project '%s' not found", shortcut)
	}

	// Copy path to clipboard
	if err := utils.CopyToClipboard(project.Path); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	fmt.Printf("✓ Copied path to clipboard: %s\n", project.Path)
	fmt.Printf("  Project: %s\n", name)

	return nil
}
