package cmd

import (
	"bufio"
	"fmt"
	"github.com/exelban/one/internal"
	"io"
	"os"
	"time"
)

func RestartCMD(cfg *internal.Config, args []string) error {
	prompt()

	name := cfg.Name
	if len(args) > 0 {
		name = args[0]
	}
	if cfg.Name == "" && cfg.Image == "" && name == "" {
		return fmt.Errorf("service name or image is not provided")
	}

	if cfg.Name == "" {
		cfg.Name = name
	}

	copyFile := internal.BoolFlag(&args, "--copy", "-c")
	force := internal.BoolFlag(&args, "--force", "-f")

	if force {
		if copyFile {
			if err := internal.CopyDockerCompose(cfg); err != nil {
				return fmt.Errorf("failed to copy docker compose file: %w", err)
			}
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

		cmdArgs := []string{"up", "-d", "--force-recreate"}
		if cfg.Name != "" {
			cmdArgs = append(cmdArgs, cfg.Name)
		}
		outPipe, errPipe, cancel, wait, err := internal.Execute(cfg, "docker compose", cmdArgs)
		if err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
		defer cancel()
		go io.Copy(os.Stdout, outPipe)
		go io.Copy(os.Stderr, errPipe)
		_ = wait()

		return nil
	}

	fmt.Printf("Starting temporary container for %s...\n", cfg.Name)

	cmdArgs := []string{"--project-name=shadow", "up", "-d", "--wait", cfg.Name}
	_, errPipe, cancel, wait, err := internal.Execute(cfg, "docker compose", cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	go io.Copy(os.Stderr, errPipe)
	if err := wait(); err != nil {
		return fmt.Errorf("failed to start temporary container: %w", err)
	}

	time.Sleep(1 * time.Second)

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
	if copyFile {
		if err := internal.CopyDockerCompose(cfg); err != nil {
			return fmt.Errorf("failed to copy docker compose file: %w", err)
		}
	}

	fmt.Printf("Temporary container started for %s, going to start new container...\n", cfg.Name)

	cmdArgs = []string{"up", "-d", "--force-recreate", "--wait", cfg.Name}
	_, errPipe, cancel, wait, err = internal.Execute(cfg, "docker compose", cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	go io.Copy(os.Stderr, errPipe)
	if err := wait(); err != nil {
		return fmt.Errorf("failed to start new container: %w", err)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("New container started, going to stop temporary container...")

	cmdArgs = []string{"--project-name=shadow", "stop", cfg.Name}
	_, errPipe, cancel, wait, err = internal.Execute(cfg, "docker compose", cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	go io.Copy(os.Stderr, errPipe)
	if err := wait(); err != nil {
		return fmt.Errorf("failed to stop temporary container: %w", err)
	}
	cmdArgs = []string{"--project-name=shadow", "rm", cfg.Name, "--force"}
	_, errPipe, cancel, wait, err = internal.Execute(cfg, "docker compose", cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	go io.Copy(os.Stderr, errPipe)
	if err := wait(); err != nil {
		return fmt.Errorf("failed to stop temporary container: %w", err)
	}

	fmt.Println("Graceful restart completed")

	return nil
}

func prompt() {
	fmt.Printf("-> Press any key to continue")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
