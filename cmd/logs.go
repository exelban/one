package cmd

import (
	"bufio"
	"encoding/json"
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

	id, err := nameToID(cfg, container)
	if err != nil {
		return fmt.Errorf("failed to get service id: %w", err)
	}

	cmdName := "docker"
	cmdArgs := []string{"logs", id}
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

func nameToID(cfg *internal.Config, name string) (string, error) {
	cmdArgs := []string{"ps", "--format='{{json .}}'"}

	outPipe, errPipe, cancel, wait, err := internal.Execute(cfg, "docker compose", cmdArgs)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	_ = wait()

	bytes := make([]byte, 4*1024)
	n, _ := errPipe.Read(bytes)
	if n != 0 {
		return "", fmt.Errorf("%s", string(bytes[:n]))
	}

	scanner := bufio.NewScanner(outPipe)
	for scanner.Scan() {
		b := scanner.Bytes()
		if b[0] == '\'' {
			b = b[1 : len(b)-1]
		}
		if b[len(b)-1] == '\'' {
			b = b[:len(b)-1]
		}

		type container struct {
			ID      string `json:"ID"`
			Service string `json:"Service"`
		}
		var c container
		if err := json.Unmarshal(b, &c); err != nil {
			return "", fmt.Errorf("failed to unmarshal container: %w", err)
		}
		if c.Service == name {
			return c.ID, nil
		}
	}

	return "", nil
}
