package project

import (
	"fmt"
	"strings"

	terminal "github.com/leodahal4/dev-kit/utils"
	"github.com/spf13/cobra"
)

func NewProjectCommand() *cobra.Command {
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Create new project",
		Long:  "Initialize any project with the basic configuration",
		RunE:  InitProject,
	}

	return projectCmd
}

// InitProject handles the initialization of a new project
func InitProject(cmd *cobra.Command, args []string) error {
	// Create a terminal provider
	prov := terminal.New("linux") // Replace "linux" with a method to detect the user's OS if needed

	// Prompt for project name
	fmt.Print("Enter project name: ")
	projectName, err := prov.ReadInput()
	if err != nil {
		return err
	}
	projectName = strings.TrimSpace(projectName)

	if projectName == "" {
		fmt.Println("Project name is required.")
		return nil
	}

	// Prompt for project description
	fmt.Print("Enter project description: ")
	projectDescription, err := prov.ReadInput()
	if err != nil {
		return err
	}
	projectDescription = strings.TrimSpace(projectDescription)

	// Prompt for microservice architecture
	fmt.Print("Is this a microservice architecture? (yes/no): ")
	microserviceInput, err := prov.ReadInput()
	if err != nil {
		return err
	}
	isMicroservice := strings.ToLower(strings.TrimSpace(microserviceInput)) == "yes"

	// Display the collected information
	fmt.Printf("Initializing project:\n")
	fmt.Printf("Name: %s\n", projectName)
	fmt.Printf("Description: %s\n", projectDescription)
	fmt.Printf("Microservice Architecture: %t\n", isMicroservice)

	// Here you can add logic to create the project structure

	return nil
}
