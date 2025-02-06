package project

import (
	"fmt"

	"github.com/leodahal4/dev-kit/cli/utils"
	"github.com/leodahal4/dev-kit/config"

	"github.com/spf13/cobra"
)

func NewProjectCommand() *cobra.Command {
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Create new project",
		Long:  "Initialize any project with the basic configuration",
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.ParseAndSaveCommand(cmd, args)
		},
		RunE: InitProject,
	}

	return projectCmd
}

// InitProject handles the initialization of a new project
func InitProject(cmd *cobra.Command, args []string) error {
	// Create a terminal provider
	projectName := utils.AskInput("Name: ", "Sample Project")
	projectDescription := utils.AskInput("Description: ", " ")
	isMicroserviceUser := utils.AskInput("Is this a microservice architecture ? (yes/no): ", "no")
	isMicroservice := isMicroserviceUser == "yes" || isMicroserviceUser == "y"

	cfg := config.GetConfig()

	// Validate for duplicate project
	if err := cfg.ValidateProject(projectName); err != nil {
		return err
	}

	cfg.Projects = append(cfg.Projects, config.ProjectConfig{
		ID:             fmt.Sprintf("%d", cfg.GetProjectNewId()),
		Name:           projectName,
		Description:    projectDescription,
		IsMicroservice: isMicroservice,
	})
	cfg.UpdateConfig()
	return nil
}
