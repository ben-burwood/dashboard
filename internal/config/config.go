package config

import (
	"dashboard/internal/config/service"
	"dashboard/internal/config/tag"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TwiN/deepmerge"
	"github.com/google/yamlfmt/pkg/yaml"
)

const (
	// DefaultConfigurationFilePath is the default path that will be used to search for the configuration file
	// if a custom path isn't configured through the DASHBOARD_CONFIG_PATH environment variable
	DefaultConfigurationFilePath = "config/config.yml"

	DefaultTitle = "Dashboard"
)

var (
	// ErrConfigFileNotFound is an error returned when a configuration file could not be found
	ErrConfigFileNotFound = errors.New("configuration file not found")

	// ErrNoServicesInConfig is an error returned when a configuration file or directory has no services configured
	ErrNoServicesInConfig = errors.New("configuration should contain at least one Service")
)

type Config struct {
	Title           string             `yaml:"title,omitempty"`
	BackgroundColor string             `yaml:"background-color,omitempty"`
	Tags            tag.Tags           `yaml:"tags,omitempty"`
	Services        []*service.Service `yaml:"services,omitempty"`
}

// LoadConfig loads the configuration from the specified path. If the path is a directory, all .yml and .yaml files
func LoadConfig(configPath string) (*Config, error) {
	var fileInfo os.FileInfo
	var usedConfigPath string
	for _, configurationPath := range []string{configPath, DefaultConfigurationFilePath} {
		if len(configurationPath) == 0 {
			continue
		}
		var err error
		fileInfo, err = os.Stat(configurationPath)
		if err != nil {
			continue
		}
		usedConfigPath = configurationPath
		break
	}
	if len(usedConfigPath) == 0 {
		return nil, ErrConfigFileNotFound
	}

	var configBytes []byte
	var config *Config
	if fileInfo.IsDir() {
		err := walkConfigDir(configPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error walking path %s: %w", path, err)
			}
			if strings.Contains(path, "..") {
				log.Println("[config.LoadConfiguration] Ignoring configuration from %s", path)
				return nil
			}
			log.Println("[config.LoadConfiguration] Reading configuration from %s", path)
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading configuration from file %s: %w", path, err)
			}
			configBytes, err = deepmerge.YAML(configBytes, data)
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("error reading configuration from directory %s: %w", usedConfigPath, err)
		}
	} else {
		log.Println("[config.LoadConfiguration] Reading configuration from configFile=%s", usedConfigPath)
		if data, err := os.ReadFile(usedConfigPath); err != nil {
			return nil, fmt.Errorf("error reading configuration from directory %s: %w", usedConfigPath, err)
		} else {
			configBytes = data
		}
	}
	if len(configBytes) == 0 {
		return nil, ErrConfigFileNotFound
	}

	config, err := parseAndValidateConfigBytes(configBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return config, nil
}

// walkConfigDir is a wrapper for filepath.WalkDir that strips directories and non-config files
func walkConfigDir(path string, fn fs.WalkDirFunc) error {
	if len(path) == 0 {
		// If the user didn't provide a directory, we'll just use the default config file, so we can return nil now.
		return nil
	}
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d == nil || d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".yml" && ext != ".yaml" {
			return nil
		}
		return fn(path, d, err)
	})
}

// parseAndValidateConfigBytes parses a Dashboard Configuration file into a Config struct and validates its parameters
func parseAndValidateConfigBytes(yamlBytes []byte) (config *Config, err error) {
	// Replace $$ with __DASHBOARD_LITERAL_DOLLAR_SIGN__ to prevent os.ExpandEnv from treating "$$" as if it was an
	// environment variable. This allows Dashboard to support literal "$" in the configuration file.
	yamlBytes = []byte(strings.ReplaceAll(string(yamlBytes), "$$", "__DASHBOARD_LITERAL_DOLLAR_SIGN__"))
	// Expand environment variables
	yamlBytes = []byte(os.ExpandEnv(string(yamlBytes)))
	// Replace __DASHBOARD_LITERAL_DOLLAR_SIGN__ with "$" to restore the literal "$" in the configuration file
	yamlBytes = []byte(strings.ReplaceAll(string(yamlBytes), "__DASHBOARD_LITERAL_DOLLAR_SIGN__", "$"))

	// Parse configuration file
	if err = yaml.Unmarshal(yamlBytes, &config); err != nil {
		return
	}

	// Check if the configuration file at least has Services configured
	if config == nil || (len(config.Services) == 0) {
		err = ErrNoServicesInConfig
	} else {
		if err := validateServicesConfig(config); err != nil {
			return nil, err
		}
	}
	if err := validateTagsConfig(config); err != nil {
		return nil, err
	}

	return
}

// validateServicesConfig validates all services in the configuration
func validateServicesConfig(config *Config) error {
	for _, service := range config.Services {
		if err := service.ValidateAndSetDefaults(); err != nil {
			return fmt.Errorf("invalid service configuration for %s: %w", service.Title, err)
		}
	}
	return nil
}

// validateTagsConfig validates all tags in the configuration
func validateTagsConfig(config *Config) error {
	if err := config.Tags.EnsureUnique(); err != nil {
		return fmt.Errorf("invalid tag configuration: %w", err)
	}
	return nil
}
