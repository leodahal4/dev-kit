package init_cmd

import (
	checktools "github.com/leodahal4/dev-kit/cli/init-cmd/check-tools"
	"github.com/leodahal4/dev-kit/cli/init-cmd/environment"
	"github.com/leodahal4/dev-kit/cli/init-cmd/project"
	"github.com/leodahal4/dev-kit/config"
	"github.com/spf13/cobra"
)

var RootHelp = `This command can be used for initializing this command itself, and also for
varuous resources which can be managed from this cli.`

var example = `
	devkit init // for initializing config file
	devkit init project // for initializing new project
`

func NewInitCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:                   "init",
		Short:                 "Initialize Command",
		Long:                  RootHelp,
		GroupID:               "init",
		Example:               example,
		RunE:                  initCmd,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		SilenceErrors:         true,
	}

	initCmd.AddCommand(project.NewProjectCommand())
	initCmd.AddCommand(checktools.NewToolsCheckerCommand())
	initCmd.AddCommand(environment.NewEnvCommand())
	return initCmd
}

func initCmd(cmd *cobra.Command, _ []string) error {
	_, err := config.CreateDefaultConfig()
	if err != nil {
		return err
	}
	return cmd.Help()
}
