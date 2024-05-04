package cmd

import (
	"fmt"
	"github.com/exelban/one/internal"
	"io"
	"os"
)

func LogsCMD(cfg *internal.Config, args []string) error {
	follow := internal.BoolFlag(&args, "-f", "--follow")

	container := cfg.Name
	if len(args) > 0 {
		container = args[0]
	}
	if container == "" {
		return fmt.Errorf("container id or name is required")
	}

	cmdName := "docker"
	cmdArgs := []string{"logs", container}
	if cfg.SSH != nil && cfg.SSH.SwarmMode {
		cmdArgs = []string{"service", "logs", container}
	}
	if follow {
		cmdArgs = append(cmdArgs, "-f")
	}

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
