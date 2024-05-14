package main

import (
	"fmt"
	"github.com/exelban/one/cmd"
	"github.com/exelban/one/internal"
	"os"
)

type cli struct {
	commands []internal.Command

	config  *internal.Config
	context *internal.Context
}

func main() {
	cli := &cli{
		config:  internal.LoadConfig(),
		context: internal.ActiveContext(),
		commands: []internal.Command{
			{
				Command:     []string{"build", "b"},
				Description: "Build docker image",
				Handler:     cmd.BuildCMD,
			},
			{
				Command:     []string{"start", "up", "install", "i"},
				Description: "Start service",
				Handler:     cmd.StartCMD,
			},
			{
				Command:     []string{"stop", "down", "delete", "d"},
				Description: "Stop service",
				Handler:     cmd.StopCMD,
			},
			{
				Command:     []string{"restart", "r", "upgrade", "u"},
				Description: "Restart service",
				Handler:     cmd.RestartCMD,
			},
			{
				Command:     []string{"logs", "l"},
				Description: "Service logs",
				Handler:     cmd.LogsCMD,
			},
			{
				Command:     []string{"list", "ps", "ls"},
				Description: "List services",
				Handler:     cmd.ListCMD,
			},
			{
				Command:     []string{"context", "ctx"},
				Description: "Manage contexts",
				MaxArgs:     1,
				Handler:     cmd.ContextCMD,
				Subcommands: []internal.Command{
					{
						Command:     []string{"add", "new", "create"},
						Description: "Add new context",
						Handler:     cmd.AddContextCMD,
					},
					{
						Command:     []string{"delete", "remove", "rm"},
						Description: "Delete context",
						Handler:     cmd.DeleteContextCMD,
					},
					{
						Command:     []string{"list", "ls", "ps"},
						Description: "List contexts",
						Handler:     cmd.ContextListCMD,
					},
					{
						Command:     []string{"activate", "use"},
						Description: "Activate context",
						Handler:     cmd.ActivateContextCMD,
					},
					{
						Command:     []string{"deactivate", "unuse"},
						Description: "Deactivate context",
						Handler:     cmd.DeactivateContextCMD,
					},
				},
			},
		},
	}

	if err := cli.run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var version = "unknown"

func (s *cli) run(args []string) error {
	if len(args) == 0 {
		fmt.Println(s.help(nil))
		return nil
	}

	if args[0] == "version" || args[0] == "--version" || args[0] == "v" || args[0] == "-v" {
		fmt.Println(version)
		return nil
	}
	ctx := internal.StringFlag(&args, "--ctx", "--context", "-c")
	if ctx != "" {
		s.config.Context = ctx
	}

	if s.config.SSH == nil {
		if s.config.Context != "" {
			ctx, err := internal.FindContext(s.config.Context)
			if err != nil {
				fmt.Printf("error find context: %v\n", err)
			} else if ctx != nil {
				s.config.SSH = ctx.SSH
				if ctx.Build != nil && s.config.Build == nil {
					s.config.Build = ctx.Build
				}
			}
		} else if activeCtx := internal.ActiveContext(); activeCtx != nil {
			s.config.SSH = activeCtx.SSH
		}
	}

	if c := internal.DockerCompose(s.config); c != nil {
		if s.config.Name == "" && c.Name != "" {
			s.config.Name = c.Name
		}
		if s.config.Build == nil && c.Build != nil {
			s.config.Build = c.Build
		}
	}

	for _, c := range s.commands {
		if c.Is(args[0]) {
			if len(args) > 1 && c.Subcommands != nil {
				for _, sc := range c.Subcommands {
					if sc.Is(args[1]) {
						return sc.Handler(s.config, args[2:])
					}
				}

				if c.MaxArgs > 0 && len(args[1:]) == c.MaxArgs {
					return c.Handler(s.config, args[1:])
				}

				fmt.Println(c.Help(args))
				return nil
			}
			return c.Handler(s.config, args[1:])
		}
	}

	fmt.Println(s.help(args))
	return nil
}

func (s *cli) help(args []string) string {
	str := "One ctl to rule them all. More info at https://github.com/exelban/one\n"
	if args != nil {
		str = fmt.Sprintf("one %s command not found\n", args[0])
	}
	str += "\nAvailable commands:"

	for _, command := range s.commands {
		str += fmt.Sprintf("\n  %s  	%s", command.Command[0], command.Description)
	}

	str += `

Usage:
  one [command] [arguments] [flags]`

	return str
}
