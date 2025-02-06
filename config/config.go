package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var globalConfig *GlobalConfig
var defaultConfigPath string

// Default configuration values
var defaultConfig = GlobalConfig{
	DEBUG:              false,
	PPROF_ENABLED:      false,
	PPROF_ADD_AND_PORT: "localhost:6060",
	LOG_FORMAT:         "text",
	KUBECONFIG:         "",
	CHECKED_TOOLS: false,

	Projects: []ProjectConfig{},
}

type ProjectConfig struct {
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	IsValid        bool                `json:"-"`
	IsMicroservice bool                `json:"is_microservice"`
	Environments   []EnvironmentConfig `json:"environments"`
}

type EnvironmentConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Path        string `json:"path"`
}

type GlobalConfig struct {
	// DEBUG is a boolean value that determines whether the application is in debug mode.
	DEBUG bool `json:"debug" default:"false" required:"false"`

	// PPROF_ENABLED is a boolean value that determines whether the pprof server is enabled.
	PPROF_ENABLED bool `json:"PPROF_ENABLED" default:"false" required:"false"`

	// PPROF_PORT is the address and port for the pprof server.
	PPROF_ADD_AND_PORT string `json:"PPROF_PORT" default:"localhost:6060" required:"false"`

	// LOG_FORMAT is the format of the logs.
	LOG_FORMAT string `json:"LOG_FORMAT" default:"text" required:"false"`

	// KUBECONFIG is the path to the kubeconfig file.
	// NOTE: THIS IS ONLY USED IF API DOESNOT PROVIDE KUBECONFIG
	KUBECONFIG string `json:"KUBECONFIG" required:"false"`

	CHECKED_TOOLS bool `json:"checked_tools" yaml:"checked_tools" required:"true"`

	Projects []ProjectConfig `json:"projects"`
}

// LoadConfig loads the configuration from a file or environment variables
func LoadConfig(configPath string) (*GlobalConfig, error) {
	if configPath != "" {
		// Check if the configuration file exists
		if _, err := os.Stat(configPath); err == nil {
			// Read the configuration file
			data, err := os.ReadFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("error reading config file: %w", err)
			}

			// Unmarshal the JSON data into the globalConfig
			if err = yaml.Unmarshal(data, &globalConfig); err != nil {
				return nil, fmt.Errorf("error unmarshalling config file: %w", err)
			}

			// Validate and set defaults
			if err = validateAndSetDefaults(globalConfig); err != nil {
				return nil, err
			}

			logrus.Info("Successfully loaded configuration from file.")
			logrus.Infof("%+v", globalConfig)
			return globalConfig, validateConfig(globalConfig)
		} else if os.IsNotExist(err) {
			logrus.Warn("Configuration file does not exist.")
		} else {
			return nil, fmt.Errorf("error checking config file: %w", err)
		}
	}

	// If no config file is provided or found, create default config
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	// Create .dev-kit directory
	devKitDir := filepath.Join(home, ".dev-kit")
	if err := os.MkdirAll(devKitDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating .dev-kit directory: %w", err)
	}

	// Create default config file
	defaultConfigPath = filepath.Join(devKitDir, "config.yaml")
	if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
		return CreateDefaultConfig()
	}

	data, err := os.ReadFile(defaultConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the JSON data into the globalConfig
	if err = yaml.Unmarshal(data, &globalConfig); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	// Validate and set defaults

	// Validate and set defaults
	if err = validateAndSetDefaults(globalConfig); err != nil {
		return nil, err
	}

	logrus.Info("Successfully loaded default configuration.")
	return globalConfig, validateConfig(globalConfig)
}

func CreateDefaultConfig() (*GlobalConfig, error) {
	// If no config file is provided or found, create default config
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	// Create .dev-kit directory
	devKitDir := filepath.Join(home, ".dev-kit")
	if err := os.MkdirAll(devKitDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating .dev-kit directory: %w", err)
	}

	// Create default config file
	defaultConfigPath = filepath.Join(devKitDir, "config.yaml")

	if err = validateAndSetDefaults(&defaultConfig); err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return nil, fmt.Errorf("error marshalling default config: %w", err)
	}
	if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
		if err := os.WriteFile(defaultConfigPath, data, 0644); err != nil {
			return nil, fmt.Errorf("error writing default config file: %w", err)
		}
	}

	logrus.Info("Successfully loaded default configuration.")
	return globalConfig, nil
}

func validateAndSetDefaults(cfg *GlobalConfig) error {
	val := reflect.ValueOf(cfg).Elem()
	typ := val.Type()

	var errors []string

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		required := fieldType.Tag.Get("required")
		defaultVal := fieldType.Tag.Get("default")

		if strings.ToLower(required) == "true" {
			if field.Kind() == reflect.String && field.String() == "" {
				errors = append(errors, fmt.Sprintf("%s is required but not set", fieldType.Name))
				continue
			}
		}

		if defaultVal != "" {
			switch field.Kind() {
			case reflect.String:
				if field.String() == "" {
					field.SetString(defaultVal)
				}
			case reflect.Bool:
				if field.Bool() && defaultVal == "true" {
					field.SetBool(true)
				}
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, ", "))
	}

	return nil
}

func validateConfig(check *GlobalConfig) error {
	if !check.CHECKED_TOOLS {
    logrus.Errorf("tools are not checked, start with \"devkit check\" command, so that this tool can confirm all needed tools")
	}
	return nil
}
