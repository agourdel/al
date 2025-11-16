package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	GlobalDirName  = ".al_global"
	LocalDirName   = ".al_local"
	ProjectsFile   = "projects"
	ConfigFile     = "config"
)

type Project struct {
	Path      string   `json:"path"`
	Shortcuts []string `json:"shortcuts"`
}

type Config struct {
	PreviewLength int `json:"preview_length"`
}

// GetGlobalDir returns the path to the global .al_global directory
func GetGlobalDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, GlobalDirName), nil
}

// EnsureGlobalDir creates the global directory if it doesn't exist
func EnsureGlobalDir() error {
	globalDir, err := GetGlobalDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(globalDir, 0755); err != nil {
		return fmt.Errorf("failed to create global directory: %w", err)
	}

	// Create projects file if it doesn't exist
	projectsPath := filepath.Join(globalDir, ProjectsFile)
	if _, err := os.Stat(projectsPath); os.IsNotExist(err) {
		emptyProjects := make(map[string]Project)
		if err := SaveProjects(emptyProjects); err != nil {
			return err
		}
	}

	// Create config file if it doesn't exist
	configPath := filepath.Join(globalDir, ConfigFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := Config{PreviewLength: 60}
		if err := SaveConfig(defaultConfig); err != nil {
			return err
		}
	}

	return nil
}

// LoadProjects loads the projects map from the global directory
func LoadProjects() (map[string]Project, error) {
	globalDir, err := GetGlobalDir()
	if err != nil {
		return nil, err
	}

	projectsPath := filepath.Join(globalDir, ProjectsFile)
	data, err := os.ReadFile(projectsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]Project), nil
		}
		return nil, fmt.Errorf("failed to read projects file: %w", err)
	}

	var projects map[string]Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("failed to parse projects file: %w", err)
	}

	return projects, nil
}

// SaveProjects saves the projects map to the global directory
func SaveProjects(projects map[string]Project) error {
	globalDir, err := GetGlobalDir()
	if err != nil {
		return err
	}

	projectsPath := filepath.Join(globalDir, ProjectsFile)
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal projects: %w", err)
	}

	if err := os.WriteFile(projectsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write projects file: %w", err)
	}

	return nil
}

// LoadConfig loads the configuration from the global directory
func LoadConfig() (Config, error) {
	globalDir, err := GetGlobalDir()
	if err != nil {
		return Config{}, err
	}

	configPath := filepath.Join(globalDir, ConfigFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{PreviewLength: 60}, nil
		}
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves the configuration to the global directory
func SaveConfig(config Config) error {
	globalDir, err := GetGlobalDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(globalDir, ConfigFile)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// FindProjectByShortcut finds a project by any of its shortcuts
func FindProjectByShortcut(shortcut string) (string, Project, error) {
	projects, err := LoadProjects()
	if err != nil {
		return "", Project{}, err
	}

	shortcut = strings.ToLower(shortcut)

	for name, project := range projects {
		for _, s := range project.Shortcuts {
			if strings.ToLower(s) == shortcut {
				return name, project, nil
			}
		}
	}

	return "", Project{}, fmt.Errorf("project not found")
}

// GetLocalDir returns the path to the local .al_local directory in the given path
func GetLocalDir(projectPath string) string {
	return filepath.Join(projectPath, LocalDirName)
}

// EnsureLocalDir creates the local directory if it doesn't exist
func EnsureLocalDir(projectPath string) error {
	localDir := GetLocalDir(projectPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}
	return nil
}

// GetCurrentProjectName returns the name of the current project
func GetCurrentProjectName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	projects, err := LoadProjects()
	if err != nil {
		return "", err
	}

	for name, project := range projects {
		if project.Path == cwd {
			return name, nil
		}
	}

	return "", fmt.Errorf("current directory is not an al project")
}
