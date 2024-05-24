package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App       `yaml:"app"`
		HTTP      `yaml:"http"`
		Log       `yaml:"logger"`
		CommsRepo `yaml:"comms_repo"`
		PostsRepo `yaml:"posts_repo"`
		PG        `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name     string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version  string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		InMemory bool   `env-required:"true" yaml:"in_memory" env:"IN_MEMORY_STORAGE"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log
	Log struct {
		Level string `env-required:"true" yaml:"level"   env:"LOG_LEVEL"`
	}

	CommsRepo struct {
		Limit  int `env-default:"100" yaml:"limit" env:"COMMSREPO_LIMIT"`
		Offset int `env-default:"0" yaml:"offset" env:"COMMSREPO_OFFSET"`
	}

	PostsRepo struct {
		Limit  int `env-default:"100" yaml:"limit" env:"POSTSREPO_LIMIT"`
		Offset int `env-default:"0" yaml:"offset" env:"POSTSREPO_OFFSET"`
	}

	// PG
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
