package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	Local = "local"
	S3    = "s3"
)

type Config struct {
	Application struct {
		Port     string `yaml:"port" env:"PORT" env-default:"8090"`
		LogPath  string `yaml:"logPath" env:"LOG_PATH" env-default:"log/server.log"`
		LogLevel string `yaml:"logLevel" env:"LOG_LEVEL" env-default:"debug"`
	} `yaml:"application"`
	Storage struct {
		Type   string `yaml:"type" env:"STORAGE_TYPE" env-default:"local"`
		Region string `yaml:"region" env:"STORAGE_REGION"`
		// Bucket points either at an S3 bucket or to a local storage folder
		Bucket string `yaml:"bucket" env:"STORAGE_BUCKET"`
	} `yaml:"storage"`
	OptSrv struct {
		Endpoint string `yaml:"endpoint" env:"OPT_SRV_ENDPOINT" env-default:"127.0.0.1"`
		Port     string `yaml:"port" env:"OPT_SRV_PORT" env-default:"8090"`
	} `yaml:"opt_srv"`
}

// ReadConfig first reads the config file at provided path, then overwrites its values with environment variables of fallbacks to default values.
func ReadConfig(path string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
