package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/exelban/one/internal"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

const done = "done"

func ContextCMD(cfg *internal.Config, args []string) error {
	if len(args) > 0 {
		if err := internal.ActivateContext(args[0]); err != nil {
			return err
		}
		fmt.Println(done)
		return nil
	}

	if ctx := internal.ActiveContext(); ctx != nil {
		fmt.Println(ctx.Name)
	} else {
		fmt.Println("local")
	}
	return nil
}
func AddContextCMD(cfg *internal.Config, args []string) error {
	name := internal.StringFlag(&args, "--name", "-n")
	host := internal.StringFlag(&args, "--host", "-h")
	username := internal.StringFlag(&args, "--username", "--user", "-u")
	password := internal.StringFlag(&args, "--password", "--pass", "-p")
	privateKey := internal.StringFlag(&args, "--private-key", "--key", "--pkey", "--privateKey")

	buildFile := internal.StringFlag(&args, "--build-file", "--file", "-f")
	buildPush := internal.BoolFlag(&args, "--build-push", "--push", "-p")
	buildPlatform := internal.StringFlag(&args, "--build-platforms", "--platforms")
	buildArgs := internal.StringFlag(&args, "--build-args", "--args")
	buildForceRecreate := internal.BoolFlag(&args, "--build-force-recreate", "--force", "-f")

	if name == "" {
		return fmt.Errorf("name is required")
	}

	ctx, err := internal.FindContext(name)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	}
	if ctx != nil {
		fmt.Printf("context %s already exists, do you want to update it? [y/N]: ", name)
		var confirm string
		_, _ = fmt.Scanln(&confirm)

		if strings.ToLower(confirm) != "y" {
			return nil
		}

		if err := ctx.Delete(); err != nil {
			return err
		}
	}

	context := &internal.Context{
		ID:   UUID(),
		Name: name,
		SSH: &internal.SSH{
			Host:       host,
			Username:   username,
			Password:   password,
			PrivateKey: privateKey,
		},
	}

	if buildFile != "" || buildPush || buildPlatform != "" || buildArgs != "" || buildForceRecreate {
		context.Build = &internal.Build{
			File:          buildFile,
			Push:          buildPush,
			Platforms:     buildPlatform,
			Args:          buildArgs,
			ForceRecreate: buildForceRecreate,
		}
	}

	if err := context.Save(); err != nil {
		return err
	}

	fmt.Println(done)
	return nil
}
func DeleteContextCMD(cfg *internal.Config, args []string) error {
	ctx, err := internal.FindContext(args[0])
	if err != nil {
		return err
	}
	if err := ctx.Delete(); err != nil {
		return err
	}
	fmt.Println(done)
	return nil
}
func ContextListCMD(cfg *internal.Config, args []string) error {
	list, err := internal.AllContexts()
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(list)
	if err != nil {
		return err
	}
	b = b[:len(b)-1]
	fmt.Println(string(b))

	return nil
}
func ActivateContextCMD(cfg *internal.Config, args []string) error {
	if err := internal.ActivateContext(args[0]); err != nil {
		return err
	}
	fmt.Println(done)
	return nil
}
func DeactivateContextCMD(cfg *internal.Config, args []string) error {
	ctx := internal.ActiveContext()
	if ctx == nil {
		return nil
	}
	ctx.Active = false
	if err := ctx.Save(); err != nil {
		return err
	}
	fmt.Println(done)
	return nil
}

func UUID() string {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		timestamp := time.Now().UnixNano()
		hexTimestamp := fmt.Sprintf("%x", timestamp)
		if len(hexTimestamp) > 12 {
			hexTimestamp = hexTimestamp[:12]
		} else {
			hexTimestamp = fmt.Sprintf("%012s", hexTimestamp)
		}
		return hexTimestamp
	}
	return hex.EncodeToString(buf)
}
