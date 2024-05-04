package internal

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Execute executes the command based on the configuration: locally or remotely (ssh).
func Execute(config *Config, cmd string, args []string) (io.Reader, io.Reader, func(), func() error, error) {
	if config.SSH != nil {
		return executeViaSSH(config.SSH, fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
	}
	return executeLocally(cmd, args)
}

// executeViaSSH executes the command via SSH.
func executeViaSSH(c *SSH, cmd string) (io.Reader, io.Reader, func(), func() error, error) {
	knownHostsCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read known_hosts file: %w", err)
	}

	conf := &ssh.ClientConfig{
		User:            c.Username,
		HostKeyCallback: knownHostsCallback,
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoSKECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
			ssh.KeyAlgoSKED25519,
			ssh.KeyAlgoRSASHA256,
			ssh.KeyAlgoRSASHA512,
		},
	}

	var conn *ssh.Client
	if c.PrivateKey != "" {
		b, err := os.ReadFile(c.PrivateKey)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to read private key: %w", err)
		}
		var signer ssh.Signer
		if c.Passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(b, []byte(c.Passphrase))
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("failed to parse private key with passphrase: %w", err)
			}
		} else {
			passphraseErr := &ssh.PassphraseMissingError{}
			if _, err := ssh.ParseRawPrivateKey(b); err != nil && passphraseErr.Error() == err.Error() {
				signer, err = ssh.ParsePrivateKeyWithPassphrase(b, getPassphrase(b))
				if err != nil {
					return nil, nil, nil, nil, fmt.Errorf("failed to parse id_rsa with passphrase: %w", err)
				}
			} else {
				signer, err = ssh.ParsePrivateKey(b)
				if err != nil {
					return nil, nil, nil, nil, fmt.Errorf("failed to parse id_rsa: %w", err)
				}
			}
		}

		conf.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
		conn, err = ssh.Dial("tcp", c.GetHost(), conf)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to dial: %w", err)
		}
	} else if c.Password != "" {
		conf.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		conf.Auth = []ssh.AuthMethod{
			ssh.Password(c.Password),
		}
		conn, err = ssh.Dial("tcp", c.GetHost(), conf)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to dial: %w", err)
		}
	} else {
		b, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to read id_rsa: %w", err)
		}

		var signer ssh.Signer
		passphraseErr := &ssh.PassphraseMissingError{}
		if _, err := ssh.ParseRawPrivateKey(b); err != nil && passphraseErr.Error() == err.Error() {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(b, getPassphrase(b))
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("failed to parse id_rsa with passphrase: %w", err)
			}
		} else {
			signer, err = ssh.ParsePrivateKey(b)
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("failed to parse id_rsa: %w", err)
			}
		}

		conf.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		conf.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
		conn, err = ssh.Dial("tcp", c.GetHost(), conf)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to dial: %w", err)
		}
	}

	sess, err := conn.NewSession()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	stdOut, err := sess.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get stdOut pipe: %w", err)
	}
	stdErr, err := sess.StderrPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get stdErr pipe: %w", err)
	}

	if err := sess.Start(cmd); err != nil {
		return nil, nil, nil, nil, err
	}

	return stdOut, stdErr, func() {
		_ = sess.Close()
		_ = conn.Close()
	}, sess.Wait, nil
}

// executeLocally executes the command locally.
func executeLocally(cmd string, args []string) (io.Reader, io.Reader, func(), func() error, error) {
	command := exec.Command(cmd, args...)
	stdOut, err := command.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	stdErr, err := command.StderrPipe()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := command.Start(); err != nil {
		return nil, nil, nil, nil, err
	}
	return stdOut, stdErr, func() {}, command.Wait, nil
}
