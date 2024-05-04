package cmd

import (
	"fmt"
	"github.com/exelban/one/internal"
	"io"
	"os"
)

func StartCMD(cfg *internal.Config, args []string) error {
	copyFile := internal.BoolFlag(&args, "--copy", "-c")
	if copyFile {
		if err := internal.CopyDockerCompose(cfg); err != nil {
			return fmt.Errorf("failed to copy docker-compose file: %w", err)
		}
	}

	cmdName := "docker-compose"
	cmdArgs := []string{"up", "--build", "-d"}

	name := cfg.Name
	if len(args) > 0 {
		name = args[0]
	}
	if name != "" {
		cmdArgs = append(cmdArgs, name)
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
