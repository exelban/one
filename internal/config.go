package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Docker struct {
	File          string   `yaml:"file,omitempty"`
	Image         string   `yaml:"image,omitempty"`
	Push          bool     `yaml:"push,omitempty"`
	Platforms     []string `yaml:"platforms,omitempty"`
	Args          []string `yaml:"args,omitempty"`
	ForceRecreate bool     `yaml:"force,omitempty"`
}

type SSH struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port,omitempty"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password,omitempty"`
	PrivateKey string `yaml:"privateKey,omitempty"`
	Passphrase string `yaml:"-"`

	SwarmMode bool `yaml:"swarmMode,omitempty"`
}

type Config struct {
	Name    string `yaml:"name"`
	Context string `yaml:"context,omitempty"`

	Docker *Docker `yaml:"docker,omitempty"`
	SSH    *SSH    `yaml:"ssh,omitempty"`
}

func LoadConfig() *Config {
	config := &Config{}
	if _, err := os.Stat(".one"); err != nil {
		return config
	}
	b, err := os.ReadFile(".one")
	if err != nil {
		return config
	}
	if err := yaml.Unmarshal(b, config); err != nil {
		return config
	}
	return config
}

func (c *Config) Save() error {
	b, err := yaml.Marshal(".one")
	if err != nil {
		return fmt.Errorf("error marshall configuration: %w", err)
	}
	if err := os.WriteFile(".one", b, 0644); err != nil {
		return fmt.Errorf("error write configuration file: %w", err)
	}
	return nil
}

func (c *Config) Print() string {
	str := fmt.Sprintf("Name: %s\n", c.Name)
	if c.Docker != nil {
		str += fmt.Sprint("Docker: \n")
		if c.Docker.File != "" {
			str += fmt.Sprintf("  File: %s\n", c.Docker.File)
		}
		if c.Docker.Image != "" {
			str += fmt.Sprintf("  Image: %s\n", c.Docker.Image)
		}
		if c.Docker.Push {
			str += fmt.Sprintf("  Push: %t\n", c.Docker.Push)
		}
		if len(c.Docker.Platforms) > 0 {
			str += fmt.Sprint("  Platforms: \n")
			for _, platform := range c.Docker.Platforms {
				str += fmt.Sprintf("    %s\n", platform)
			}
		}
		if len(c.Docker.Args) > 0 {
			str += fmt.Sprint("  Args:")
			for _, arg := range c.Docker.Args {
				str += fmt.Sprintf("\n    %s", arg)
			}
		}
	}
	if c.SSH != nil {
		str += fmt.Sprintf("SSH: %s", c.SSH.GetHost())
	}
	return str
}

func (s *SSH) GetHost() string {
	host := s.Host
	if !strings.Contains(host, ":") {
		if s.Port == 0 {
			s.Port = 22
		}
		host = fmt.Sprintf("%s:%d", host, s.Port)
	}
	return host
}
