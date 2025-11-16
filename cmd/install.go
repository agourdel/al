package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alex/al/storage"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install al CLI to the system",
	Long:  `Install al CLI binaries to /usr/local/bin and create the global configuration directory.`,
	RunE:  runInstall,
}

func runInstall(cmd *cobra.Command, args []string) error {
	// Ensure global directory exists
	if err := storage.EnsureGlobalDir(); err != nil {
		return fmt.Errorf("failed to create global directory: %w", err)
	}

	// Get the current executable path
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	binDir := "/usr/local/bin"

	// List of binaries to install (except algo which is a shell script)
	binaries := []string{"al", "alinit", "alnote", "allink"}

	fmt.Println("Installing al CLI...")

	for _, binary := range binaries {
		destPath := filepath.Join(binDir, binary)
		
		// Copy executable
		if err := copyFile(executable, destPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", binary, err)
		}

		// Make it executable
		if err := os.Chmod(destPath, 0755); err != nil {
			return fmt.Errorf("failed to make %s executable: %w", binary, err)
		}

		fmt.Printf("✓ Installed %s to %s\n", binary, destPath)
	}

	// Create algo as symlink to al (no longer needs shell wrapper)
	algoPath := filepath.Join(binDir, "algo")
	os.Remove(algoPath) // Remove old version if exists
	if err := os.Symlink(filepath.Join(binDir, "al"), algoPath); err != nil {
		// If symlink fails, try copying
		if err := copyFile(executable, algoPath); err != nil {
			return fmt.Errorf("failed to install algo: %w", err)
		}
		if err := os.Chmod(algoPath, 0755); err != nil {
			return fmt.Errorf("failed to make algo executable: %w", err)
		}
	}
	fmt.Printf("✓ Installed algo to %s\n", algoPath)

	globalDir, _ := storage.GetGlobalDir()
	fmt.Printf("\n✓ Global configuration directory created at %s\n", globalDir)
	fmt.Println("✓ Installation complete!")
	fmt.Println("\nUsage:")
	fmt.Println("  algo myproject    # Copies path to clipboard")
	fmt.Println("  cd <paste>        # Then paste the path")

	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0755)
}
