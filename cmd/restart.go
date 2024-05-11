package cmd

import (
	"fmt"
	"github.com/exelban/one/internal"
	"io"
	"os"
	"time"
)

func RestartCMD(cfg *internal.Config, args []string) error {
	if cfg.Name == "" && cfg.Image == "" {
		return fmt.Errorf("service name or image is not provided")
	}

	if cfg.Image != "" {
		cmdName := "docker"
		cmdArgs := []string{"pull", cfg.Image}
		outPipe, errPipe, cancel, wait, err := internal.Execute(cfg, cmdName, cmdArgs)
		if err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
		defer cancel()
		go io.Copy(os.Stdout, outPipe)
		go io.Copy(os.Stderr, errPipe)
		_ = wait()
	}

	copyFile := internal.BoolFlag(&args, "--copy", "-c")
	force := internal.BoolFlag(&args, "--force", "-f")

	if copyFile {
		if err := internal.CopyDockerCompose(cfg); err != nil {
			return fmt.Errorf("failed to copy docker-compose file: %w", err)
		}
	}

	if cfg.SSH != nil && cfg.SSH.SwarmMode {
		cmdName := "docker"
		cmdArgs := []string{"service", "update"}
		if force {
			cmdArgs = append(cmdArgs, "--force")
		}
		if cfg.Name == "" {
			return fmt.Errorf("service name is not provided")
		}
		cmdArgs = append(cmdArgs, cfg.Name)

		outPipe, errPipe, cancel, wait, err := internal.Execute(cfg, cmdName, cmdArgs)
		if err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
		defer cancel()
		go io.Copy(os.Stdout, outPipe)
		go io.Copy(os.Stderr, errPipe)
		_ = wait()
		return nil
	}

	if force {
		cmdArgs := []string{"up", "-d", "--force-recreate"}
		if cfg.Name != "" {
			cmdArgs = append(cmdArgs, cfg.Name)
		}
		outPipe, errPipe, cancel, wait, err := internal.Execute(cfg, "docker-compose", cmdArgs)
		if err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
		defer cancel()
		go io.Copy(os.Stdout, outPipe)
		go io.Copy(os.Stderr, errPipe)
		_ = wait()
		return nil
	}

	if cfg.Name == "" {
		return fmt.Errorf("service name is not provided, blue-green deployment is not possible")
	}

	fmt.Println("Starting temporary container...")

	cmdArgs := []string{"--project-name=copy", "up", "-d", "--wait", cfg.Name}
	_, _, cancel, wait, err := internal.Execute(cfg, "docker-compose", cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	_ = wait()

	time.Sleep(1 * time.Second)
	fmt.Println("Temporary container started, going to start new container...")

	cmdArgs = []string{"up", "-d", "--force-recreate", "--wait", cfg.Name}
	_, _, cancel, wait, err = internal.Execute(cfg, "docker-compose", cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	_ = wait()

	time.Sleep(1 * time.Second)
	fmt.Println("New container started, going to stop temporary container...")

	cmdArgs = []string{"--project-name=copy", "stop", cfg.Name}
	_, _, cancel, wait, err = internal.Execute(cfg, "docker-compose", cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	_ = wait()

	fmt.Println("Graceful restart completed")

	return nil
}
