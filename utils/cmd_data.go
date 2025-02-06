package utils

import (
	"github.com/leodahal4/dev-kit/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ParseAndSaveCommand(cmd *cobra.Command, _ []string) {
	err := config.ValidateConfig(config.GetConfig())
	if err != nil {
		logrus.Fatal(err)
	}
	switch cmd.Use{
	case CMD_INIT:
		cfg := config.GetConfig()
		cfg.CURRENT_CMD = CMD_INIT
	}
}


const (
	CMD_INIT = "init"
)
