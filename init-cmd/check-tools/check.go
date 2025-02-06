package checktools

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
 usg "github.com/julienroland/usg"
 "github.com/charmbracelet/bubbles/spinner"
 	"github.com/charmbracelet/lipgloss"
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
  return nil
}
