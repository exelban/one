package internal

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v3"
	"os"
)

func getPassphrase(key []byte) []byte {
	list := make(map[string][]byte)
	path := fmt.Sprintf("%s/passphrases", supportFolder())
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if err := yaml.Unmarshal(b, &list); err != nil {
			return nil
		}
		if bytes, ok := list[string(key)]; ok {
			return bytes
		}
	}

	fmt.Print("Passphrase: ")
	bytes, _ := terminal.ReadPassword(0)
	fmt.Println()

	list[string(key)] = bytes

	if b, err := yaml.Marshal(list); err == nil {
		if err := os.WriteFile(path, b, 0644); err != nil {
			fmt.Printf("error write passphrases: %v\n", err)
		}
	}

	return bytes
}
