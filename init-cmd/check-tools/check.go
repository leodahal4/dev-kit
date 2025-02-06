package checktools

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	usg "github.com/julienroland/usg"
	"github.com/leodahal4/dev-kit/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewToolsCheckerCommand() *cobra.Command {
	checkToolsCmd := &cobra.Command{
		Use:   "check",
		Short: "Check all needed tools",
		Long:  "Initialize any project with the basic configuration",
		RunE:  checkTools,
	}

	return checkToolsCmd
}

// InitProject handles the initialization of a new project
func checkTools(cmd *cobra.Command, args []string) error {
	// Create a terminal providerA
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	logrus.Infof("%v Checking Docker", usg.Get.Tick)
	logrus.Infof("%v Checking kind", usg.Get.Cross)
	logrus.Infof("Checking golang")
	cfg := config.GetConfig()
	cfg.CHECKED_TOOLS = true
	config.UpdateConfig(cfg)
	return nil
}
