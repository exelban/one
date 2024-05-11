package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"strings"
)

func DockerCompose(cfg *Config) *Config {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil
	}

	fileName := ""
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "docker-compose") {
			fileName = file.Name()
			break
		}
	}
	if fileName == "" {
		return nil
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	defer f.Close()

	type service struct {
		Image         string `yaml:"image"`
		ContainerName string `yaml:"container_name"`
	}

	compose := struct {
		Services map[string]service `yaml:"services"`
	}{}
	if err := yaml.NewDecoder(f).Decode(&compose); err != nil {
		return nil
	}

	var s service
	if cfg.Name != "" {
		s = compose.Services[cfg.Name]
	} else {
		for _, v := range compose.Services {
			s = v
			break
		}
	}

	return &Config{
		Name:  s.ContainerName,
		Image: s.Image,
	}
}

func CopyDockerCompose(cfg *Config) error {
	if cfg.SSH == nil {
		return fmt.Errorf("no ssh config")
	}

	files, err := os.ReadDir(".")
	if err != nil {
		return nil
	}

	fileName := ""
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "docker-compose") {
			fileName = file.Name()
			break
		}
	}
	if fileName == "" {
		return nil
	}

	params := []string{}
	if cfg.SSH.Port != 0 && cfg.SSH.Port != 22 {
		params = append(params, "-P", fmt.Sprintf("%d", cfg.SSH.Port))
	}
	params = append(params, fileName)
	params = append(params, fmt.Sprintf("%s@%s:./", cfg.SSH.Username, cfg.SSH.Host))

	command := exec.Command("scp", params...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		return err
	}

	return nil
}
