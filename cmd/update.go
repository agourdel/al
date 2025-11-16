package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update al CLI binaries",
	Long:  `Update existing al CLI binaries in /usr/local/bin with the current version.`,
	RunE:  runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Get the current executable path
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	binDir := "/usr/local/bin"

	// List of binaries to update (except algo which is a shell script)
	binaries := []string{"al", "alinit", "alnote", "allink"}

	fmt.Println("Updating al CLI...")

	for _, binary := range binaries {
		destPath := filepath.Join(binDir, binary)
		tmpPath := destPath + ".new"
		
		// Copy to temporary file first
		if err := copyFile(executable, tmpPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", binary, err)
		}

		// Make it executable
		if err := os.Chmod(tmpPath, 0755); err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("failed to make %s executable: %w", binary, err)
		}

		// Atomic rename (this works even if the file is running)
		if err := os.Rename(tmpPath, destPath); err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("failed to update %s: %w", binary, err)
		}

		fmt.Printf("✓ Updated %s\n", binary)
	}

	// Update algo as symlink/copy
	algoPath := filepath.Join(binDir, "algo")
	algoTmpPath := algoPath + ".new"
	
	os.Remove(algoPath) // Remove old version
	if err := os.Symlink(filepath.Join(binDir, "al"), algoTmpPath); err != nil {
		// If symlink fails, try copying
		if err := copyFile(executable, algoTmpPath); err != nil {
			return fmt.Errorf("failed to update algo: %w", err)
		}
		if err := os.Chmod(algoTmpPath, 0755); err != nil {
			os.Remove(algoTmpPath)
			return fmt.Errorf("failed to make algo executable: %w", err)
		}
	}
	
	if err := os.Rename(algoTmpPath, algoPath); err != nil {
		os.Remove(algoTmpPath)
		return fmt.Errorf("failed to update algo: %w", err)
	}
	
	fmt.Printf("✓ Updated algo\n")

	fmt.Println("✓ Update complete!")

	return nil
}
