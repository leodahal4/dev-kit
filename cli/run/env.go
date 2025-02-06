package run

import (
	"github.com/leodahal4/dev-kit/cli/utils"
	"github.com/leodahal4/dev-kit/config"
	"github.com/sirupsen/logrus"

	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func NewEnvRun() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "env",
		Short: "Run environment",
		Long:  "Run any environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.ParseAndSaveCommand(cmd, args)
		},
		RunE: InitEnvRun,
	}

	runCmd.Flags().IntP("id", "i", 0, "id to run env")

	return runCmd
}

// InitProject handles the initialization of a new project
func InitEnvRun(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()
	envPath := cfg.Projects[0].Environments[1].Path
	_ = RunENV(envPath)
	return nil
}

func RunENV(path string) error {
	cmdToRun := exec.Command("go", "run", "main.go")
	cmdToRun.Dir = path

	// Stream output directly to stdout and stderr
	cmdToRun.Stdout = os.Stdout // Stream standard output
	cmdToRun.Stderr = os.Stderr // Stream standard error

	err := cmdToRun.Start() // Start the command
	if err != nil {
		logrus.Errorf("ERR %v", err)
		return err
	}

	err = cmdToRun.Wait() // Wait for the command to finish
	if err != nil {
		logrus.Errorf("ERR %v", err)
		return err
	}
	return nil
}
