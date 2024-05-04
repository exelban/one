package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Context struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	Active bool   `yaml:"active,omitempty"`
	SSH    *SSH   `yaml:"ssh,omitempty"`
}

func (c *Context) Save() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshall configuration: %w", err)
	}
	path := fmt.Sprintf("%s/%s.json", supportFolder(), c.ID)
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("error write configuration file: %w", err)
	}
	return nil
}
func (c *Context) Delete() error {
	return os.Remove(fmt.Sprintf("%s/%s.json", supportFolder(), c.ID))
}

func AllContexts() ([]*Context, error) {
	folder := supportFolder()
	if _, err := os.Stat(folder); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(folder, 0755); err != nil {
			return nil, fmt.Errorf("error create folder: %w", err)
		}
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("error read folder: %w", err)
	}

	contexts := []*Context{}
	for _, file := range files {
		if file.IsDir() || file.Name() == ".DS_Store" || !strings.Contains(file.Name(), ".json") {
			continue
		}

		context := &Context{}
		b, err := os.ReadFile(fmt.Sprintf("%s/%s", folder, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("error read file: %w", err)
		}
		if err := yaml.Unmarshal(b, context); err != nil {
			return nil, fmt.Errorf("error unmarshal file: %w", err)
		}
		contexts = append(contexts, context)
	}

	return contexts, nil
}
func ActiveContext() *Context {
	list, err := AllContexts()
	if err != nil {
		return nil
	}
	for _, context := range list {
		if context.Active {
			return context
		}
	}
	return nil
}
func FindContext(val string) (*Context, error) {
	list, err := AllContexts()
	if err != nil {
		return nil, err
	}

	for _, context := range list {
		if context.Name == val || context.ID == val {
			return context, nil
		}
	}

	return nil, fmt.Errorf("context %s not found", val)
}
func ActivateContext(val string) error {
	list, err := AllContexts()
	if err != nil {
		return err
	}
	for _, context := range list {
		if context.Active {
			context.Active = false
			if err := context.Save(); err != nil {
				return err
			}
		}
		if context.Name == val || context.ID == val {
			context.Active = true
			if err := context.Save(); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func supportFolder() string {
	return fmt.Sprintf("%s/Library/Application Support/one", os.Getenv("HOME"))
}
