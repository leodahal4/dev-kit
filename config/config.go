package config

import (
	"errors"
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

const (
	defaultConfigFileName = "config.yaml"
	devKitDirName         = ".dev-kit"
)

// Default configuration values
var defaultConfig = GlobalConfig{
	DEBUG:              false,
	PPROF_ENABLED:      false,
	PPROF_ADD_AND_PORT: "localhost:6060",
	LOG_FORMAT:         "text",
	KUBECONFIG:         "",
	CHECKED_TOOLS:      false,

	Projects: []ProjectConfig{},
}

type ProjectConfig struct {
	ID             string              `json:"id"`
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

	Projects    []ProjectConfig `json:"projects"`
	CURRENT_CMD string          `json:"_"`
}

// LoadConfig loads the configuration from a file or environment variables
func LoadConfig(configPath string) (*GlobalConfig, error) {
	if configPath != "" {
		if err := loadConfigFromFile(configPath); err != nil {
			return nil, err
		}
	}

	return loadDefaultConfig()
}

func loadConfigFromFile(configPath string) error {
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			logrus.Errorf("error reading config file: %v", err)
			return err
		}

		if err = yaml.Unmarshal(data, &globalConfig); err != nil {
			logrus.Errorf("error unmarshalling config file: %v", err)
			return err
		}

		return validateAndSetDefaults(globalConfig)
	} else if os.IsNotExist(err) {
		logrus.Warn("Configuration file does not exist.")
		return nil
	} else {
		logrus.Errorf("error checking config file: %v", err)
		return err
	}
}

func loadDefaultConfig() (*GlobalConfig, error) {
	home, err := homedir.Dir()
	if err != nil {
		logrus.Errorf("error getting home directory: %v", err)
		return nil, err
	}

	devKitDir := filepath.Join(home, devKitDirName)
	if err := os.MkdirAll(devKitDir, os.ModePerm); err != nil {
		logrus.Errorf("error creating .dev-kit directory: %v", err)
		return nil, err
	}

	defaultConfigPath = filepath.Join(devKitDir, defaultConfigFileName)
	if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
		if _, err := CreateDefaultConfig(); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(defaultConfigPath)
	if err != nil {
		logrus.Errorf("error reading default config file: %v", err)
		return nil, err
	}

	if err = yaml.Unmarshal(data, &globalConfig); err != nil {
		logrus.Errorf("error unmarshalling default config file: %v", err)
		return nil, err
	}

	return globalConfig, validateAndSetDefaults(globalConfig)
}

// CreateDefaultConfig creates the default configuration file
func CreateDefaultConfig() (*GlobalConfig, error) {
	if err := validateAndSetDefaults(&defaultConfig); err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		logrus.Errorf("error marshalling default config: %v", err)
		return nil, err
	}

	home, err := homedir.Dir()
	if err != nil {
		logrus.Errorf("error getting home directory: %v", err)
		return nil, err
	}

	devKitDir := filepath.Join(home, devKitDirName)
	if err := os.MkdirAll(devKitDir, os.ModePerm); err != nil {
		logrus.Errorf("error creating .dev-kit directory: %v", err)
		return nil, err
	}

	defaultConfigPath = filepath.Join(devKitDir, defaultConfigFileName)
	if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
		if err := os.WriteFile(defaultConfigPath, data, 0644); err != nil {
			logrus.Errorf("error writing default config file: %v", err)
			return nil, err
		}
	}

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

func ValidateConfig(check *GlobalConfig) error {
	if !check.CHECKED_TOOLS {
		// logrus.Errorf("tools are not checked, start with \"devkit check\" command, so that this tool can confirm all needed tools")
		return errors.New("tools are not checked, start with \"devkit check\" command, so that this tool can confirm all needed tools")
	}
	return nil
}

func GetConfig() *GlobalConfig {
	return globalConfig
}

func (cfg *GlobalConfig) UpdateConfig() {
	data, _ := yaml.Marshal(cfg)
	os.WriteFile(defaultConfigPath, data, 0644)
}

func UpdateConfig(config *GlobalConfig) {
	data, _ := yaml.Marshal(config)
	os.WriteFile(defaultConfigPath, data, 0644)
}

// ValidateEnv checks if an environment with the same name already exists in the project
func (cfg *GlobalConfig) ValidateEnv(projectID, envName, path string) error {
	// Check if the project exists by ID
	var project *ProjectConfig
	for _, p := range cfg.Projects {
		if p.ID == projectID {
			project = &p
			break
		}
	}

	if project == nil {
		return fmt.Errorf("project with ID '%s' does not exist", projectID)
	}

	// Check for duplicate environment names
	for _, env := range project.Environments {
		if env.Name == envName {
			return fmt.Errorf("environment with name '%s' already exists in project with ID '%s'", envName, projectID)
		}
	}

	// Validate the path (you can add more specific path validation as needed)
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	return nil
}

// ValidateProject checks if a project with the same name already exists
func (cfg *GlobalConfig) ValidateProject(projectName string) error {
	for _, project := range cfg.Projects {
		if project.Name == projectName {
			return fmt.Errorf("project with name '%s' already exists", projectName)
		}
	}
	return nil
}

func (cfg *GlobalConfig) GetProjectNewId() int {
	// Generate a new project ID based on the current number of projects
	return len(cfg.Projects) + 1
}
