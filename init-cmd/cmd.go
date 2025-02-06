package init_cmd

import (
	"github.com/leodahal4/dev-kit/config"
	"github.com/leodahal4/dev-kit/init-cmd/project"
	checktools "github.com/leodahal4/dev-kit/init-cmd/check-tools"
	"github.com/spf13/cobra"
)

var RootHelp = `This command can be used for initializing this command itself, and also for
varuous resources which can be managed from this cli.`

var example = `
	devkit init // for initializing config file
	devkit init project // for initializing new project
`

func NewInitCommand() *cobra.Command {
	initCmd := &cobra.Command {
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
	return initCmd
}

func initCmd(cmd *cobra.Command, _ []string) error {
	_, err := config.CreateDefaultConfig()
	if err != nil {
		return err
	}
	return cmd.Help()
}
