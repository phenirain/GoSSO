package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env              string     `yaml:"env" env-default:"local" json:"env"`
	ConnectionString string     `yaml:"connection_string" json:"connection-string" env-required:"true"`
	Secret           string     `yaml:"secret" json:"secret" env-required:"true"`
	GRPC             GRPCConfig `yaml:"grpc" json:"grpc" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" json:"port" env-default:"8000"`
	Timeout time.Duration `yaml:"timeout" json:"timeout" env-default:"5m"`
}

func MustLoadConfig() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("GRPC_AUTH_CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("GRPC_AUTH_CONFIG_PATH file does not exist")
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Failed to parse config with: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var configPath string

	// --config="path/to/config.yaml"
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()
	if configPath == "" {
		configPath = os.Getenv("GRPC_AUTH_CONFIG_PATH")
	}
	return configPath
}
