package run

import (
	"sync"

	"github.com/leodahal4/dev-kit/cli/utils"
	"github.com/leodahal4/dev-kit/config"

	"github.com/spf13/cobra"
)

func NewProjectRun() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "project",
		Short: "Run environment",
		Long:  "Run any environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.ParseAndSaveCommand(cmd, args)
		},
		RunE: InitProjectRun,
	}

	runCmd.Flags().IntP("id", "i", 0, "id to run env")

	return runCmd
}

// InitProject handles the initialization of a new project
func InitProjectRun(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()
	envs := cfg.Projects[0].Environments
	var wg sync.WaitGroup
	for _, env := range envs {
		wg.Add(1)
		go func() {
			err := RunENV(env.Path)
			if err != nil {
				wg.Done()
			}
		}()
	}
	wg.Wait()
	return nil
}
