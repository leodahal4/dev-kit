package models

import (
	"log"
	"os"

	"github.com/leodahal4/dev-kit/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProjectConfig struct {
	gorm.Model

	ID             string              `json:"id" gorm:"primaryKey"`
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	IsValid        bool                `json:"-"`
	IsMicroservice bool                `json:"is_microservice"`
	Environments   []EnvironmentConfig `json:"environments" gorm:"foreignKey:ProjectID"`
	GlobalConfigID uint                `json:"global_config_id"`
}

type EnvironmentConfig struct {
	gorm.Model

	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Path        string `json:"path"`
	ProjectID   string `json:"project_id"`
}

type GlobalConfig struct {
	gorm.Model

	// DEBUG is a boolean value that determines whether the application is in debug mode.
	DEBUG bool `json:"debug" default:"false" required:"false"`

	// PPROF_ENABLED is a boolean value that determines whether the pprof server is enabled.
	PPROF_ENABLED bool `json:"PPROF_ENABLED" default:"false" required:"false"`

	// PPROF_PORT is the address and port for the pprof server.
	PPROF_ADD_AND_PORT string `json:"PPROF_PORT" default:"localhost:6060" required:"false"`

	// LOG_FORMAT is the format of the logs.
	LOG_FORMAT string `json:"LOG_FORMAT" default:"text" required:"false"`

	// KUBECONFIG is the path to the kubeconfig file.
	// NOTE: THIS IS ONLY USED IF API DOES NOT PROVIDE KUBECONFIG
	KUBECONFIG string `json:"KUBECONFIG" required:"false"`
	SQLITEDB   string `json:"db_path" default:".dev-kit/devkit.sqlite3"`

	CHECKED_TOOLS bool            `json:"checked_tools" yaml:"checked_tools" required:"true"`
	Projects      []ProjectConfig `json:"projects" gorm:"foreignKey:GlobalConfigID"`
	CURRENT_CMD   string          `json:"_"`
}

type ConfigRepository struct {
	db *gorm.DB
}

type RepoImpl interface {
	CreateGlobalConfig(*GlobalConfig) error
	ReadGlobalConfig() (*GlobalConfig, error)
	UpdateGlobalConfig(config *GlobalConfig) error
	DeleteGlobalConfig() error
	Migrate() error
}

func NewConfigRepository(cfg *config.GlobalConfig) RepoImpl {
	// Create the directory if it doesn't exist
	dbPath := cfg.HOME_FOLDER + "/" + cfg.SQLITEDB
	logrus.Info(dbPath)
	// Ensure the directory exists
	if err := os.MkdirAll(cfg.HOME_FOLDER, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	r := &ConfigRepository{db: db}
	r.Migrate()
	return r
}

func (r *ConfigRepository) Migrate() error {
	err := r.db.AutoMigrate(&GlobalConfig{})
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
		return err
	}
	return nil
}

// CreateGlobalConfig inserts a new GlobalConfig into the database
func (r *ConfigRepository) CreateGlobalConfig(config *GlobalConfig) error {
	return r.db.Model(&GlobalConfig{}).Create(config).Error
}

// ReadGlobalConfig retrieves the GlobalConfig from the database
func (r *ConfigRepository) ReadGlobalConfig() (*GlobalConfig, error) {
	var config GlobalConfig
	err := r.db.Model(&GlobalConfig{}).Find(&config).Error
	return &config, err
}

// UpdateGlobalConfig updates an existing GlobalConfig in the database
func (r *ConfigRepository) UpdateGlobalConfig(config *GlobalConfig) error {
	return r.db.Model(&GlobalConfig{}).Updates(config).Error
}

// DeleteGlobalConfig deletes a GlobalConfig from the database
func (r *ConfigRepository) DeleteGlobalConfig() error {
	err := r.db.Raw("DELETE FROM global_config").Error
	return err
}
