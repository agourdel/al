package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alex/al/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [shortcuts...]",
	Short: "Initialize the current directory as an al project",
	Long: `Initialize the current directory as an al project with optional shortcuts.
The directory name is automatically added as a shortcut.
Additional shortcuts can be specified separated by pipes (|).

Example: al init bar|mad|tes`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	// Ensure global directory exists
	if err := storage.EnsureGlobalDir(); err != nil {
		return fmt.Errorf("failed to initialize global directory: %w", err)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get directory name
	dirName := filepath.Base(cwd)

	// Parse shortcuts
	shortcuts := []string{dirName}
	if len(args) > 0 {
		// Join all args and split by pipe
		allArgs := strings.Join(args, " ")
		additionalShortcuts := strings.Split(allArgs, "|")
		for _, s := range additionalShortcuts {
			s = strings.TrimSpace(s)
			if s != "" && s != dirName {
				shortcuts = append(shortcuts, s)
			}
		}
	}

	// Check if project already exists
	projects, err := storage.LoadProjects()
	if err != nil {
		return fmt.Errorf("failed to load projects: %w", err)
	}

	// Check if this path is already a project
	for existingName, existingProject := range projects {
		if existingProject.Path == cwd {
			return fmt.Errorf("this directory is already initialized as project '%s'", existingName)
		}
	}

	// Check if any shortcut is already used
	for _, shortcut := range shortcuts {
		for existingName, existingProject := range projects {
			for _, existingShortcut := range existingProject.Shortcuts {
				if strings.ToLower(existingShortcut) == strings.ToLower(shortcut) {
					return fmt.Errorf("shortcut '%s' is already used by project '%s'", shortcut, existingName)
				}
			}
		}
	}

	// Create local .al_local directory
	if err := storage.EnsureLocalDir(cwd); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Create subdirectories for notes and links
	localDir := storage.GetLocalDir(cwd)
	notesDir := filepath.Join(localDir, "notes")
	linksDir := filepath.Join(localDir, "links")

	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes directory: %w", err)
	}

	if err := os.MkdirAll(linksDir, 0755); err != nil {
		return fmt.Errorf("failed to create links directory: %w", err)
	}

	// Add project to global registry
	projects[dirName] = storage.Project{
		Path:      cwd,
		Shortcuts: shortcuts,
	}

	if err := storage.SaveProjects(projects); err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	fmt.Printf("✓ Initialized project '%s' in %s\n", dirName, cwd)
	fmt.Printf("✓ Shortcuts: %s\n", strings.Join(shortcuts, ", "))

	return nil
}
