package cmd

import (
	"fmt"
	"github.com/exelban/one/internal"
	"io"
	"os"
)

func RestartCMD(cfg *internal.Config, args []string) error {
	if cfg.Docker == nil {
		return fmt.Errorf("docker configuration is not provided")
	}

	copyFile := internal.BoolFlag(&args, "--copy", "-c")
	if copyFile {
		if err := internal.CopyDockerCompose(cfg); err != nil {
			return fmt.Errorf("failed to copy docker-compose file: %w", err)
		}
	}

	force := false
	for _, arg := range args {
		if arg == "-f" || arg == "--force" {
			force = true
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

	cmdName := "docker"
	cmdArgs := []string{"pull", cfg.Docker.Image}
	outPipe, errPipe, cancel, wait, err := internal.Execute(cfg, cmdName, cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	go io.Copy(os.Stdout, outPipe)
	go io.Copy(os.Stderr, errPipe)
	_ = wait()

	cmdName = "docker-compose"
	cmdArgs = []string{"up", "-d"}
	if force || cfg.Docker.ForceRecreate {
		cmdArgs = append(cmdArgs, "--force-recreate")
	}
	if cfg.Name != "" {
		cmdArgs = append(cmdArgs, cfg.Name)
	}

	_, errPipe, cancel, wait, err = internal.Execute(cfg, cmdName, cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	go io.Copy(os.Stdout, outPipe)
	go io.Copy(os.Stderr, errPipe)
	_ = wait()

	return nil
}
