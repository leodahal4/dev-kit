package environment

import (
	"context"
	"fmt"
	"os"

	"github.com/leodahal4/dev-kit/cli/utils"
	"github.com/leodahal4/dev-kit/config"
	pb "github.com/leodahal4/dev-kit/protos"
	"google.golang.org/grpc"

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
	EnvCmd.Flags().BoolP("server", "s", false, "print DevKit version")

	return EnvCmd
}

// InitEnv handles the initialization of a new Env
func InitEnv(cmd *cobra.Command, args []string) error {
	serverFlag, _ := cmd.Flags().GetBool("server")
	if serverFlag {
		return createEnvOnServer()
	}
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

// createEnvOnServer handles the creation of an environment using the gRPC server
func createEnvOnServer() error {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewConfigServiceClient(conn)

	// Gather environment details
	projectID := utils.AskInput("Project ID: ", "1")
	EnvName := utils.AskInput("Name: ", "Sample Env")
	EnvDescription := utils.AskInput("Description: ", " ")
	EnvPath := utils.AskInput("Path: ", "/path/to/env")

	// Create the environment on the server
	env := &pb.EnvironmentConfig{
		Name:        EnvName,
		Description: EnvDescription,
		Path:        EnvPath,
	}

	// Call the gRPC method to create the environment
	_, err = client.CreateEnvironment(context.Background(), &pb.CreateEnvironmentRequest{
		ProjectId:   projectID,
		Environment: env,
	})
	if err != nil {
		return fmt.Errorf("error creating environment on server: %v", err)
	}

	fmt.Println("Environment created successfully on server.")
	return nil
}
