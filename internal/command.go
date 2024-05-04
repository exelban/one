package internal

import (
	"fmt"
	"strings"
)

type Command struct {
	Command     []string
	Description string
	MaxArgs     int
	Handler     func(*Config, []string) error
	Subcommands []Command
}

func (c *Command) Is(v string) bool {
	for _, cmd := range c.Command {
		if cmd == v {
			return true
		}
	}
	return false
}

func (c *Command) Help(v []string) string {
	str := fmt.Sprintf("one %s command not found\n\nAvailable commands:", strings.Join(v, " "))
	for _, command := range c.Subcommands {
		str += fmt.Sprintf("\n  %s   		%s", command.Command[0], command.Description)
	}
	str += fmt.Sprintf(`

Usage:
  one %s [arguments] [flags]`, c.Command[0])
	return str
}
