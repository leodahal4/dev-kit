package cmd

import (
	"fmt"

	init_cmd "github.com/leodahal4/dev-kit/cli/init-cmd"
	"github.com/leodahal4/dev-kit/cli/run"
	"github.com/leodahal4/dev-kit/cli/utils"
	"github.com/leodahal4/dev-kit/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version = "0.0.1-alphav1"

var RootHelp = `DevKit is a command-line tool designed to simplify the setup of development environments. 
It provides a streamlined process for configuring tools, dependencies, and project structures 
to help developers get started quickly and efficiently. 

Features:
- Easy installation of required tools
- Configuration of project settings, supporting microservices architecture
- Management of dependencies
- Live Reloading / Hot Reloading
- Customizable templates for various project types

Docs: https://github.com/leodahal4/dev-kit.git

Usage:
Run 'devkit init' to start setting up your development environment.
Run 'devkit help' for more information on available commands and options.`

var cfgPath string

var Cmd = &cobra.Command{
	Use:                   "devkit",
	Short:                 "DevKit, prepared by Dev for Dev",
	Long:                  RootHelp,
	Run:                   RootCmdRun,
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		cobra.OnInitialize(initConfig)
		utils.ParseAndSaveCommand(cmd, args)
	},
}

func Execute() {
	Cmd.Flags().BoolP("version", "v", false, "print DevKit version")
	Cmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "base project directory eg. github.com/spf13/")
	Cmd.AddCommand(init_cmd.NewInitCommand())
	Cmd.AddCommand(run.NewRun())
	Cmd.AddGroup(&cobra.Group{
		ID:    "init",
		Title: "Init Commands",
	})

	err := Cmd.Execute()
	if err != nil {
		logrus.Errorf("Err: %s", err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	_, err := config.LoadConfig(cfgPath)
	if err != nil {
		logrus.Fatalf("%s", err.Error())
	}
}

func RootCmdRun(cmd *cobra.Command, _ []string) {
	versionFlag, _ := cmd.Flags().GetBool("version")

	if versionFlag {
		fmt.Printf("Version: %s\n", version)
		return
	}
	err := cmd.Help()
	if err != nil {
		logrus.Fatal(err)
	}
}
