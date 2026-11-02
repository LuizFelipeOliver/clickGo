package clickup

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configDir  = ".config/clickGo"
	configFile = "config.yaml"
)

type Config struct {
	Token    string `yaml:"token"`
	UserID   int    `yaml:"user_id,omitempty"`
	Username string `yaml:"username,omitempty"`
	TeamID   string `yaml:"team_id,omitempty"`
	TeamName string `yaml:"team_name,omitempty"`
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	path, err := ConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func (c *Config) HasToken() bool {
	return c.Token != ""
}
