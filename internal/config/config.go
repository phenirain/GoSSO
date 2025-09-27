package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env              string     `yaml:"env" env-default:"local" json:"env"`
	ConnectionString string     `yaml:"connection_string" json:"connection-string" env-required:"true"`
	AllowedOrigins   []string   `mapstructure:"allowed_origins"`
	Secret           []byte     `yaml:"secret" json:"secret" env-required:"true"`
	HTTP             HTTPConfig `yaml:"http" json:"http" env-required:"true"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port" json:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" json:"timeout" env-default:"5m"`
}

func LoadConfig() (*Config, error) {

	viper.SetConfigFile("config/config.yaml")

	var cfg Config
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling config file")
	}
	return &cfg, nil
}
