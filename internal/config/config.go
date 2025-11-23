package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SSHServers []SSHServer `yaml:"ssh_servers"`
	Tunnels    []Tunnel    `yaml:"tunnels"`
}

type SSHServer struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password,omitempty"`
	KeyFile  string `yaml:"key_file,omitempty"`
}

type Tunnel struct {
	ServerName string `yaml:"server_name"`
	RemoteAddr string `yaml:"remote_addr"`
	LocalAddr  string `yaml:"local_addr"`
	Mode       string `yaml:"mode"` // "local" or "remote"
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
