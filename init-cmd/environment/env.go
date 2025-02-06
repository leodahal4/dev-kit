package environment

import (
	"fmt"
	"os"

	"github.com/leodahal4/dev-kit/config"
	"github.com/leodahal4/dev-kit/utils"

	"github.com/spf13/cobra"
)

func NewEnvCommand() *cobra.Command {
	EnvCmd := &cobra.Command{
		Use:   "env",
		Short: "Create new Env",
		Long:  "Initialize any Env with the basic configuration",
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.ParseAndSaveCommand(cmd, args)
		},
		RunE: InitEnv,
	}

	return EnvCmd
}

// InitEnv handles the initialization of a new Env
func InitEnv(cmd *cobra.Command, args []string) error {
	// Create a terminal provider
	projectID := utils.AskInput("Project ID: ", "1")
	EnvName := utils.AskInput("Name: ", "Sample Env")
	EnvDescription := utils.AskInput("Description: ", " ")
	EnvPath := utils.AskInput("Path: ", "/path/to/env")

	// Check if the user entered '.' and resolve it to the current working directory
	if EnvPath == "." {
		var err error
		EnvPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %v", err)
		}
	}

	cfg := config.GetConfig()

	// Validate for duplicate Env and path
	if err := cfg.ValidateEnv(projectID, EnvName, EnvPath); err != nil {
		return err
	}

	// Find the project and append the new environment
	for i, project := range cfg.Projects {
		if project.ID == projectID {
			cfg.Projects[i].Environments = append(cfg.Projects[i].Environments, config.EnvironmentConfig{
				Name:        EnvName,
				Description: EnvDescription,
				Path:        EnvPath,
			})
			break
		}
	}

	cfg.UpdateConfig()
	return nil
}
