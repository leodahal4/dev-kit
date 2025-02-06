package run

import (
	"github.com/leodahal4/dev-kit/utils"
	"github.com/spf13/cobra"
)

func NewRun() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run env or project",
		Long:  "Run any environment or projects",
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.ParseAndSaveCommand(cmd, args)
		},
	}

	runCmd.AddCommand(NewEnvRun())
	runCmd.AddCommand(NewProjectRun())

	return runCmd
}
