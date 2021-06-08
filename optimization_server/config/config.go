package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	Local = "local"
	S3    = "s3"
)

type Config struct {
	Application struct {
		Port     string `yaml:"port" env:"PORT" env-default:"8080"`
		LogPath  string `yaml:"logPath" env:"LOG_PATH" env-default:"log/server.log"`
		LogLevel string `yaml:"logLevel" env:"LOG_LEVEL" env-default:"debug"`
	} `yaml:"application"`
	Script struct {
		Dir     string        `yaml:"dir" env:"SCRIPT_DIR" env-default:"script"`
		Prefix  string        `yaml:"prefix" env:"SCRIPT_PREFIX" env-default:"tmp"`
		Path    string        `yaml:"path" env:"SCRIPT_PATH" env-default:"main.py"`
		Timeout time.Duration `yaml:"timeout" env:"SCRIPT_TIMEOUT" env-default:"5000ms"`
	} `yaml:"script"`
	Storage struct {
		Type   string `yaml:"type" env:"STORAGE_TYPE" env-default:"local"`
		Region string `yaml:"region" env:"STORAGE_REGION"`
		// Bucket points either at an S3 bucket or to a local storage folder
		Bucket string `yaml:"bucket" env:"STORAGE_BUCKET"`
	} `yaml:"storage"`
}

// ReadConfig first reads the config file at provided path, then overwrites its values with environment variables of fallbacks to default values.
func ReadConfig(path string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
