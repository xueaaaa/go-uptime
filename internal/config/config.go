package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/xueaaaa/go-uptime/internal/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	SSLDefault    bool          `yaml:"ssl_default"`
	ServerAddress string        `yaml:"server_address"`
	TickInterval  time.Duration `yaml:"tick_interval"`
	Workers       int           `yaml:"workers"`
	TimeoutAfter  time.Duration `yaml:"timeout_after"`
	Database      `yaml:"database"`
}

type Database struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       string `yaml:"db"`
	MinConns int32  `yaml:"min_conns"`
	MaxConns int32  `yaml:"max_conns"`
}

func New() (Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return Config{}, err
	}

	configPath := filepath.Join(filepath.Dir(exePath), "config.yaml")
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	config := Config{
		SSLDefault:    false,
		ServerAddress: "0.0.0.0:8080",
		TickInterval:  10 * time.Second,
		Workers:       50,
		TimeoutAfter:  10 * time.Second,
		Database: Database{
			MinConns: 5,
			MaxConns: 25,
		},
	}

	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	if config.User == "" ||
		config.Password == "" ||
		config.Host == "" ||
		config.Port == "" ||
		config.DB == "" {
		return Config{}, errors.MissingConfigFields
	}

	return config, nil
}
