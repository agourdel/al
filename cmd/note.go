package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/alex/al/storage"
	"github.com/alex/al/utils"
	"github.com/spf13/cobra"
)

type Note struct {
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Encrypted bool      `json:"encrypted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	noteTarget    string
	noteEncrypted bool
	noteBody      string
	noteCopy      bool
)

var noteCmd = &cobra.Command{
	Use:   "note [action] [name]",
	Short: "Manage notes for projects",
	Long:  `Manage notes for projects. Actions: list, add, get, edit, remove`,
}

var noteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	RunE:  runNoteList,
}

var noteAddCmd = &cobra.Command{
	Use:   "add [#name]",
	Short: "Add a new note",
	Args:  cobra.ExactArgs(1),
	RunE:  runNoteAdd,
}

var noteGetCmd = &cobra.Command{
	Use:   "get [#name]",
	Short: "Get a note",
	Args:  cobra.ExactArgs(1),
	RunE:  runNoteGet,
}

var noteEditCmd = &cobra.Command{
	Use:   "edit [#name]",
	Short: "Edit a note",
	Args:  cobra.ExactArgs(1),
	RunE:  runNoteEdit,
}

var noteRemoveCmd = &cobra.Command{
	Use:   "remove [#name]",
	Short: "Remove a note",
	Args:  cobra.ExactArgs(1),
	RunE:  runNoteRemove,
}

func init() {
	// Add flags
	noteCmd.PersistentFlags().StringVarP(&noteTarget, "target", "t", "", "Target project")
	noteAddCmd.Flags().BoolVarP(&noteEncrypted, "chiffre", "c", false, "Encrypt the note")
	noteAddCmd.Flags().StringVarP(&noteBody, "body", "b", "", "Note body (no editor)")
	noteEditCmd.Flags().StringVarP(&noteBody, "body", "b", "", "Note body (no editor)")
	noteGetCmd.Flags().BoolVar(&noteCopy, "cp", false, "Copy to clipboard")

	// Add subcommands
	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteAddCmd)
	noteCmd.AddCommand(noteGetCmd)
	noteCmd.AddCommand(noteEditCmd)
	noteCmd.AddCommand(noteRemoveCmd)
}

func getProjectPath() (string, error) {
	if noteTarget != "" {
		_, project, err := storage.FindProjectByShortcut(noteTarget)
		if err != nil {
			return "", fmt.Errorf("project '%s' not found", noteTarget)
		}
		return project.Path, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return cwd, nil
}

func getNotesDir(projectPath string) string {
	return filepath.Join(storage.GetLocalDir(projectPath), "notes")
}

func getNoteFilePath(projectPath, noteName string) string {
	noteName = strings.TrimPrefix(noteName, "#")
	return filepath.Join(getNotesDir(projectPath), noteName+".json")
}

func loadNote(projectPath, noteName string) (*Note, error) {
	notePath := getNoteFilePath(projectPath, noteName)
	data, err := os.ReadFile(notePath)
	if err != nil {
		return nil, err
	}

	var note Note
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func saveNote(projectPath string, note *Note) error {
	notePath := getNoteFilePath(projectPath, note.Name)
	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(notePath, data, 0600)
}

func listNotes(projectPath string) ([]Note, error) {
	notesDir := getNotesDir(projectPath)
	
	entries, err := os.ReadDir(notesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Note{}, nil
		}
		return nil, err
	}

	var notes []Note
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		noteName := strings.TrimSuffix(entry.Name(), ".json")
		note, err := loadNote(projectPath, noteName)
		if err != nil {
			continue
		}
		notes = append(notes, *note)
	}

	return notes, nil
}

func findSimilarNotes(projectPath, noteName string, maxDistance int) ([]string, error) {
	notes, err := listNotes(projectPath)
	if err != nil {
		return nil, err
	}

	var noteNames []string
	for _, note := range notes {
		noteNames = append(noteNames, note.Name)
	}

	return utils.FindSimilarStrings(noteName, noteNames, maxDistance), nil
}

func runNoteList(cmd *cobra.Command, args []string) error {
	projectPath, err := getProjectPath()
	if err != nil {
		return err
	}

	notes, err := listNotes(projectPath)
	if err != nil {
		return err
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	config, _ := storage.LoadConfig()
	previewLength := config.PreviewLength

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Name\tDate\tPreview")
	fmt.Fprintln(w, "----\t----\t-------")

	for _, note := range notes {
		date := note.UpdatedAt.Format("2006-01-02")
		preview := "**chiffrée**"
		if !note.Encrypted {
			preview = utils.TruncateString(note.Content, previewLength)
			preview = strings.ReplaceAll(preview, "\n", " ")
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", note.Name, date, preview)
	}

	w.Flush()
	return nil
}

func runNoteAdd(cmd *cobra.Command, args []string) error {
	projectPath, err := getProjectPath()
	if err != nil {
		return err
	}

	noteName := strings.TrimPrefix(args[0], "#")

	// Check if note already exists
	if _, err := loadNote(projectPath, noteName); err == nil {
		return fmt.Errorf("note '%s' already exists", noteName)
	}

	var content string
	var password string

	if noteEncrypted {
		password, err = utils.ReadPassword("Enter encryption password: ")
		if err != nil {
			return err
		}
		confirmPassword, err := utils.ReadPassword("Confirm password: ")
		if err != nil {
			return err
		}
		if password != confirmPassword {
			return fmt.Errorf("passwords do not match")
		}
	}

	if noteBody != "" {
		content = noteBody
	} else {
		// Create temporary file for editing
		tmpFile := filepath.Join(os.TempDir(), "al_note_"+noteName+".txt")
		if err := os.WriteFile(tmpFile, []byte(""), 0600); err != nil {
			return err
		}
		defer os.Remove(tmpFile)

		if err := utils.OpenEditor(tmpFile); err != nil {
			return err
		}

		data, err := os.ReadFile(tmpFile)
		if err != nil {
			return err
		}
		content = string(data)
	}

	// Encrypt if needed
	if noteEncrypted {
		encrypted, err := utils.Encrypt(content, password)
		if err != nil {
			return err
		}
		content = encrypted
	}

	note := &Note{
		Name:      noteName,
		Content:   content,
		Encrypted: noteEncrypted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := saveNote(projectPath, note); err != nil {
		return err
	}

	fmt.Printf("✓ Note '%s' created\n", noteName)
	return nil
}

func runNoteGet(cmd *cobra.Command, args []string) error {
	projectPath, err := getProjectPath()
	if err != nil {
		return err
	}

	noteName := strings.TrimPrefix(args[0], "#")

	note, err := loadNote(projectPath, noteName)
	if err != nil {
		// Try to find similar notes
		similar, _ := findSimilarNotes(projectPath, noteName, 3)
		if len(similar) > 0 {
			fmt.Printf("Note '%s' not found. Did you mean:\n", noteName)
			for _, s := range similar {
				fmt.Printf("  - %s\n", s)
			}
			return nil
		}
		return fmt.Errorf("note '%s' not found", noteName)
	}

	content := note.Content
	if note.Encrypted {
		password, err := utils.ReadPassword("Enter decryption password: ")
		if err != nil {
			return err
		}

		decrypted, err := utils.Decrypt(content, password)
		if err != nil {
			return err
		}
		content = decrypted
	}

	if noteCopy {
		if err := utils.CopyToClipboard(content); err != nil {
			return err
		}
		fmt.Println("✓ Note copied to clipboard")
	} else {
		fmt.Println(content)
	}

	return nil
}

func runNoteEdit(cmd *cobra.Command, args []string) error {
	projectPath, err := getProjectPath()
	if err != nil {
		return err
	}

	noteName := strings.TrimPrefix(args[0], "#")

	note, err := loadNote(projectPath, noteName)
	if err != nil {
		// Try to find similar notes
		similar, _ := findSimilarNotes(projectPath, noteName, 3)
		if len(similar) > 0 {
			fmt.Printf("Note '%s' not found. Did you mean:\n", noteName)
			for _, s := range similar {
				fmt.Printf("  - %s\n", s)
			}
			return nil
		}
		return fmt.Errorf("note '%s' not found", noteName)
	}

	var password string
	content := note.Content

	if note.Encrypted {
		password, err = utils.ReadPassword("Enter decryption password: ")
		if err != nil {
			return err
		}

		decrypted, err := utils.Decrypt(content, password)
		if err != nil {
			return err
		}
		content = decrypted
	}

	if noteBody != "" {
		content = noteBody
	} else {
		// Create temporary file for editing
		tmpFile := filepath.Join(os.TempDir(), "al_note_"+noteName+".txt")
		if err := os.WriteFile(tmpFile, []byte(content), 0600); err != nil {
			return err
		}
		defer os.Remove(tmpFile)

		if err := utils.OpenEditor(tmpFile); err != nil {
			return err
		}

		data, err := os.ReadFile(tmpFile)
		if err != nil {
			return err
		}
		content = string(data)
	}

	// Re-encrypt if needed
	if note.Encrypted {
		encrypted, err := utils.Encrypt(content, password)
		if err != nil {
			return err
		}
		content = encrypted
	}

	note.Content = content
	note.UpdatedAt = time.Now()

	if err := saveNote(projectPath, note); err != nil {
		return err
	}

	fmt.Printf("✓ Note '%s' updated\n", noteName)
	return nil
}

func runNoteRemove(cmd *cobra.Command, args []string) error {
	projectPath, err := getProjectPath()
	if err != nil {
		return err
	}

	noteName := strings.TrimPrefix(args[0], "#")

	if _, err := loadNote(projectPath, noteName); err != nil {
		// Try to find similar notes
		similar, _ := findSimilarNotes(projectPath, noteName, 3)
		if len(similar) > 0 {
			fmt.Printf("Note '%s' not found. Did you mean:\n", noteName)
			for _, s := range similar {
				fmt.Printf("  - %s\n", s)
			}
			return nil
		}
		return fmt.Errorf("note '%s' not found", noteName)
	}

	if !utils.AskConfirmation(fmt.Sprintf("Are you sure you want to delete note '%s'?", noteName)) {
		fmt.Println("Cancelled.")
		return nil
	}

	notePath := getNoteFilePath(projectPath, noteName)
	if err := os.Remove(notePath); err != nil {
		return err
	}

	fmt.Printf("✓ Note '%s' deleted\n", noteName)
	return nil
}
