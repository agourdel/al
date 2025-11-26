package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/alex/al/storage"
	"github.com/alex/al/utils"
	"github.com/spf13/cobra"
)

type Link struct {
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

var (
	linkTarget       string
	linkURL          string
	linkKeywords     string
	linkAddKeywords  string
	linkResetKeywords string
	linkCopy         bool
)

var linkCmd = &cobra.Command{
	Use:   "link [action]",
	Short: "Manage links for projects",
	Long:  `Manage links for projects. Actions: list, add, get, edit, remove`,
}

var linkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all links",
	RunE:  runLinkList,
}

var linkAddCmd = &cobra.Command{
	Use:   "add [#name]",
	Short: "Add a new link",
	Args:  cobra.ExactArgs(1),
	RunE:  runLinkAdd,
}

var linkGetCmd = &cobra.Command{
	Use:   "get [#name/#keyword]",
	Short: "Get a link",
	Args:  cobra.ExactArgs(1),
	RunE:  runLinkGet,
}

var linkEditCmd = &cobra.Command{
	Use:   "edit [#name/#keyword]",
	Short: "Edit a link",
	Args:  cobra.ExactArgs(1),
	RunE:  runLinkEdit,
}

var linkRemoveCmd = &cobra.Command{
	Use:   "remove [#name/#keyword]",
	Short: "Remove a link",
	Args:  cobra.ExactArgs(1),
	RunE:  runLinkRemove,
}

func init() {
	// Add flags
	linkCmd.PersistentFlags().StringVarP(&linkTarget, "target", "t", "", "Target project")
	
	linkAddCmd.Flags().StringVarP(&linkURL, "url", "u", "", "Link URL (required)")
	linkAddCmd.Flags().StringVarP(&linkKeywords, "keywords", "k", "", "Keywords separated by |")
	linkAddCmd.MarkFlagRequired("url")

	linkGetCmd.Flags().BoolVarP(&linkCopy, "copy", "c", false, "Copy to clipboard")

	linkEditCmd.Flags().StringVarP(&linkURL, "url", "u", "", "New URL")
	linkEditCmd.Flags().StringVarP(&linkAddKeywords, "add_keyword", "a", "", "Add keywords")
	linkEditCmd.Flags().StringVarP(&linkResetKeywords, "reset_keyword", "r", "", "Reset keywords")
	

	// Add subcommands
	linkCmd.AddCommand(linkListCmd)
	linkCmd.AddCommand(linkAddCmd)
	linkCmd.AddCommand(linkGetCmd)
	linkCmd.AddCommand(linkEditCmd)
	linkCmd.AddCommand(linkRemoveCmd)
}

func getLinkProjectPath() (string, error) {
	if linkTarget != "" {
		_, project, err := storage.FindProjectByShortcut(linkTarget)
		if err != nil {
			return "", fmt.Errorf("project '%s' not found", linkTarget)
		}
		return project.Path, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return cwd, nil
}

func getLinksDir(projectPath string) string {
	return filepath.Join(storage.GetLocalDir(projectPath), "links")
}

func getLinkFilePath(projectPath, linkName string) string {
	linkName = strings.TrimPrefix(linkName, "#")
	return filepath.Join(getLinksDir(projectPath), linkName+".json")
}

func loadLink(projectPath, linkName string) (*Link, error) {
	linkPath := getLinkFilePath(projectPath, linkName)
	data, err := os.ReadFile(linkPath)
	if err != nil {
		return nil, err
	}

	var link Link
	if err := json.Unmarshal(data, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

func saveLink(projectPath string, link *Link) error {
	linkPath := getLinkFilePath(projectPath, link.Name)
	data, err := json.MarshalIndent(link, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(linkPath, data, 0600)
}

func listLinks(projectPath string) ([]Link, error) {
	linksDir := getLinksDir(projectPath)
	
	entries, err := os.ReadDir(linksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Link{}, nil
		}
		return nil, err
	}

	var links []Link
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		linkName := strings.TrimSuffix(entry.Name(), ".json")
		link, err := loadLink(projectPath, linkName)
		if err != nil {
			continue
		}
		links = append(links, *link)
	}

	return links, nil
}

func findLinkByNameOrKeyword(projectPath, identifier string) (*Link, error) {
	identifier = strings.TrimPrefix(identifier, "#")
	identifier = strings.ToLower(identifier)

	links, err := listLinks(projectPath)
	if err != nil {
		return nil, err
	}

	// First try exact match by name
	for _, link := range links {
		if strings.ToLower(link.Name) == identifier {
			return &link, nil
		}
	}

	// Then try keywords
	for _, link := range links {
		for _, keyword := range link.Keywords {
			if strings.ToLower(keyword) == identifier {
				return &link, nil
			}
		}
	}

	return nil, fmt.Errorf("link not found")
}

func findSimilarLinks(projectPath, linkIdentifier string, maxDistance int) ([]string, error) {
	links, err := listLinks(projectPath)
	if err != nil {
		return nil, err
	}

	var identifiers []string
	for _, link := range links {
		identifiers = append(identifiers, link.Name)
		identifiers = append(identifiers, link.Keywords...)
	}

	return utils.FindSimilarStrings(linkIdentifier, identifiers, maxDistance), nil
}

func runLinkList(cmd *cobra.Command, args []string) error {
	projectPath, err := getLinkProjectPath()
	if err != nil {
		return err
	}

	links, err := listLinks(projectPath)
	if err != nil {
		return err
	}

	if len(links) == 0 {
		fmt.Println("No links found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Name\tLink\tKeywords")
	fmt.Fprintln(w, "----\t----\t--------")

	for _, link := range links {
		keywords := strings.Join(link.Keywords, ", ")
		fmt.Fprintf(w, "%s\t%s\t%s\n", link.Name, link.URL, keywords)
	}

	w.Flush()
	return nil
}

func runLinkAdd(cmd *cobra.Command, args []string) error {
	projectPath, err := getLinkProjectPath()
	if err != nil {
		return err
	}

	linkName := strings.TrimPrefix(args[0], "#")

	// Check if link already exists
	if _, err := loadLink(projectPath, linkName); err == nil {
		return fmt.Errorf("link '%s' already exists", linkName)
	}

	// Parse keywords
	var keywords []string
	if linkKeywords != "" {
		parts := strings.Split(linkKeywords, "|")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				keywords = append(keywords, part)
			}
		}
	}

	link := &Link{
		Name:     linkName,
		URL:      linkURL,
		Keywords: keywords,
	}

	if err := saveLink(projectPath, link); err != nil {
		return err
	}

	fmt.Printf("✓ Link '%s' created\n", linkName)
	return nil
}

func runLinkGet(cmd *cobra.Command, args []string) error {
	projectPath, err := getLinkProjectPath()
	if err != nil {
		return err
	}

	identifier := args[0]

	link, err := findLinkByNameOrKeyword(projectPath, identifier)
	if err != nil {
		// Try to find similar links
		similar, _ := findSimilarLinks(projectPath, identifier, 3)
		if len(similar) > 0 {
			fmt.Printf("Link '%s' not found. Did you mean:\n", identifier)
			for _, s := range similar {
				fmt.Printf("  - %s\n", s)
			}
			return nil
		}
		return fmt.Errorf("link '%s' not found", identifier)
	}

	if linkCopy {
		if err := utils.CopyToClipboard(link.URL); err != nil {
			return err
		}
		
	} 
		
	fmt.Printf("URL: %s | (%s)\n", link.URL, strings.Join(link.Keywords, ", "))
	
	return nil
}

func runLinkEdit(cmd *cobra.Command, args []string) error {
	projectPath, err := getLinkProjectPath()
	if err != nil {
		return err
	}

	identifier := args[0]

	link, err := findLinkByNameOrKeyword(projectPath, identifier)
	if err != nil {
		// Try to find similar links
		similar, _ := findSimilarLinks(projectPath, identifier, 3)
		if len(similar) > 0 {
			fmt.Printf("Link '%s' not found. Did you mean:\n", identifier)
			for _, s := range similar {
				fmt.Printf("  - %s\n", s)
			}
			return nil
		}
		return fmt.Errorf("link '%s' not found", identifier)
	}

	// Update URL if provided
	if linkURL != "" {
		link.URL = linkURL
	}

	// Add keywords if provided
	if linkAddKeywords != "" {
		parts := strings.Split(linkAddKeywords, "|")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" && !contains(link.Keywords, part) {
				link.Keywords = append(link.Keywords, part)
			}
		}
	}

	// Reset keywords if provided
	if cmd.Flags().Changed("reset_keyword") || cmd.Flags().Changed("rk") {
		if linkResetKeywords == "" {
			link.Keywords = []string{}
		} else {
			parts := strings.Split(linkResetKeywords, "|")
			var newKeywords []string
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part != "" {
					newKeywords = append(newKeywords, part)
				}
			}
			link.Keywords = newKeywords
		}
	}

	if err := saveLink(projectPath, link); err != nil {
		return err
	}

	fmt.Printf("✓ Link '%s' updated\n", link.Name)
	return nil
}

func runLinkRemove(cmd *cobra.Command, args []string) error {
	projectPath, err := getLinkProjectPath()
	if err != nil {
		return err
	}

	identifier := args[0]

	link, err := findLinkByNameOrKeyword(projectPath, identifier)
	if err != nil {
		// Try to find similar links
		similar, _ := findSimilarLinks(projectPath, identifier, 3)
		if len(similar) > 0 {
			fmt.Printf("Link '%s' not found. Did you mean:\n", identifier)
			for _, s := range similar {
				fmt.Printf("  - %s\n", s)
			}
			return nil
		}
		return fmt.Errorf("link '%s' not found", identifier)
	}

	if !utils.AskConfirmation(fmt.Sprintf("Are you sure you want to delete link '%s'?", link.Name)) {
		fmt.Println("Cancelled.")
		return nil
	}

	linkPath := getLinkFilePath(projectPath, link.Name)
	if err := os.Remove(linkPath); err != nil {
		return err
	}

	fmt.Printf("✓ Link '%s' deleted\n", link.Name)
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
