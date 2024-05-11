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
	push := internal.BoolFlag(&args, "--push", "-p")

	var image string
	if cfg.Image != "" {
		image = cfg.Image
	}
	if len(args) > 0 {
		image = args[0]
	}

	if image == "" {
		return errors.New("cannot build without image")
	}

	if strings.Contains(image, ":") {
		if cfg.Build == nil {
			cfg.Build = &internal.Build{}
		}
		cfg.Image = image
	} else {
		if cfg.Image == "" {
			return errors.New("[ERROR]: provided tag without image name")
		}
		cfg.Image = fmt.Sprintf("%s:%s", strings.Split(cfg.Image, ":")[0], image)
	}

	c := cfg.Build
	if c == nil {
		c = &internal.Build{}
	}

	if cfg.Image == "" {
		return errors.New("cannot build without image name")
	}
	if !cfg.Build.Push && push {
		cfg.Build.Push = push
	}

	fmt.Printf("Building image %s\n", cfg.Image)

	command := exec.Command("docker", buildDockerBuild(cfg.Image, c)...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		return err
	}

	return nil
}

func buildDockerBuild(image string, c *internal.Build) (arr []string) {
	if c.Platform != "" {
		arr = append(arr, "buildx")
		arr = append(arr, "build")
		arr = append(arr, fmt.Sprintf("--platform=%s", c.Platform))
	} else {
		arr = append(arr, "build")
	}

	if c.Push {
		arr = append(arr, "--push")
	}

	if c.Args != "" {
		for _, arg := range strings.Split(c.Args, ",") {
			arg = strings.TrimSpace(arg)
			keyValue := strings.Split(arg, "=")
			if len(keyValue) == 2 && strings.Contains(keyValue[1], "env") {
				env := strings.TrimPrefix(keyValue[1], "env:")
				env = strings.TrimSpace(env)
				if value := os.Getenv(env); value != "" {
					keyValue[1] = value
				}
			}
			arr = append(arr, fmt.Sprintf("--build-arg=%s", strings.Join(keyValue, "=")))
		}
	}

	if c.File != "" {
		arr = append(arr, c.File)
	} else {
		arr = append(arr, ".")
	}

	if image != "" {
		arr = append(arr, fmt.Sprintf("--tag=%s", image))
	}

	return
}
