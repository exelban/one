package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Build struct {
	File          string   `yaml:"file,omitempty"`
	Push          bool     `yaml:"push,omitempty"`
	Platforms     string   `yaml:"platforms,omitempty"`
	Args          string   `yaml:"args,omitempty"`
	ForceRecreate bool     `yaml:"force,omitempty"`
	Context       []string `yaml:"context,omitempty"`
}

type SSH struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port,omitempty"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password,omitempty"`
	PrivateKey string `yaml:"privateKey,omitempty"`
	Passphrase string `yaml:"-"`
}

type Config struct {
	Name    string `yaml:"name"`
	Image   string `yaml:"image,omitempty"`
	Context string `yaml:"context,omitempty"`

	Build *Build `yaml:"build,omitempty"`
	SSH   *SSH   `yaml:"ssh,omitempty"`
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
	if c.Image != "" {
		str += fmt.Sprintf("Image: %s\n", c.Image)
	}
	if c.Context != "" {
		str += fmt.Sprintf("Context: %s\n", c.Context)
	}
	if c.Build != nil {
		str += fmt.Sprint("Docker: \n")
		if c.Build.File != "" {
			str += fmt.Sprintf("  File: %s\n", c.Build.File)
		}
		if c.Build.Push {
			str += fmt.Sprintf("  Push: %t\n", c.Build.Push)
		}
		if c.Build.Platforms != "" {
			str += fmt.Sprintf("  Platforms: %s\n", c.Build.Platforms)
		}
		if c.Build.Args != "" {
			str += fmt.Sprintf("  Args: %s\n", c.Build.Args)
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
