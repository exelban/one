package cmd

import (
	"errors"
	"fmt"
	"github.com/exelban/one/internal"
	"os"
	"os/exec"
	"strings"
)

func BuildCMD(cfg *internal.Config, args []string) error {
	var image string
	if cfg.Docker != nil && cfg.Docker.Image != "" {
		image = cfg.Docker.Image
	}
	if len(args) > 0 {
		image = args[0]
	}

	if image == "" {
		return errors.New("cannot build without image")
	}

	if strings.Contains(image, ":") {
		if cfg.Docker == nil {
			cfg.Docker = &internal.Docker{}
		}
		cfg.Docker.Image = image
	} else {
		if cfg.Docker == nil || cfg.Docker.Image == "" {
			return errors.New("[ERROR]: provided tag without image name")
		}
		cfg.Docker.Image = fmt.Sprintf("%s:%s", strings.Split(cfg.Docker.Image, ":")[0], image)
	}

	if cfg.Docker == nil || cfg.Docker.Image == "" {
		return errors.New("cannot build without image name")
	}

	command := exec.Command("docker", buildDockerBuild(cfg.Docker)...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		return err
	}

	return nil
}

func buildDockerBuild(c *internal.Docker) (arr []string) {
	if len(c.Platforms) > 0 {
		arr = append(arr, "buildx")
	} else {
		arr = append(arr, "build")
	}

	if c.Push {
		arr = append(arr, "--push")
	}
	if len(c.Platforms) > 0 {
		str := "--platform="
		for i, platform := range c.Platforms {
			if i > 0 {
				str += ","
			}
			str += platform
		}
		arr = append(arr, str)
	}
	if len(c.Args) > 0 {
		for _, arg := range c.Args {
			arr = append(arr, fmt.Sprintf("--build-arg=%s", arg))
		}
	}
	if c.Image != "" {
		arr = append(arr, fmt.Sprintf("--tag=%s", c.Image))
	}
	if c.File != "" {
		arr = append(arr, c.File)
	} else {
		arr = append(arr, ".")
	}
	return
}
