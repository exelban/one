package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/exelban/one/internal"
	"os"
	"strings"
)

func ListCMD(cfg *internal.Config, args []string) error {
	cmdName := "docker"
	cmdArgs := []string{"ps", "--format='{{json .}}'"}
	if cfg.SSH != nil && cfg.SSH.SwarmMode {
		cmdArgs = []string{"service", "ls", "--format='{{json .}}'"}
	}

	outPipe, errPipe, cancel, wait, err := internal.Execute(cfg, cmdName, cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	defer cancel()
	_ = wait()

	bytes := make([]byte, 1024)
	n, _ := errPipe.Read(bytes)
	if n != 0 {
		return fmt.Errorf("%s", string(bytes[:n]))
	}

	data := [][]string{}
	scanner := bufio.NewScanner(outPipe)
	for scanner.Scan() {
		b := scanner.Bytes()
		if b[0] == '\'' {
			b = b[1 : len(b)-1]
		}
		if b[len(b)-1] == '\'' {
			b = b[:len(b)-1]
		}
		d, err := parseContainer(b)
		if err != nil {
			return fmt.Errorf("failed to parse container: %w", err)
		}
		data = append(data, d)
	}

	if len(data) == 0 {
		return fmt.Errorf("no containers to show")
	}

	re := lipgloss.NewRenderer(os.Stdout)
	purple := lipgloss.Color("99")
	green := lipgloss.Color("10")
	red := lipgloss.Color("9")
	yellow := lipgloss.Color("11")

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
			}
			st := re.NewStyle().Padding(0, 1)
			if col == len(data[row-1])-1 {
				switch data[row-1][col] {
				case "running":
					return st.Foreground(green)
				case "exited", "dead":
					return st.Foreground(red)
				default:
					if strings.Contains(data[row-1][col], "/") {
						st = st.Width(7).Align(lipgloss.Center)
						arr := strings.Split(data[row-1][col], "/")
						if arr[0] == arr[1] {
							return st.Foreground(green)
						}
						return st.Foreground(red)
					}
					return st.Foreground(yellow)
				}
			}
			return st
		}).
		Headers("ID", "Name", "Image", "Tag", "Ports", "Status", "State").
		Rows(data...)

	fmt.Println(t)

	return nil
}

func parseContainer(b []byte) ([]string, error) {
	type service struct {
		ID       string `json:"ID"`
		Names    string `json:"Names"`
		Name     string `json:"Name"`
		Image    string `json:"Image"`
		Ports    string `json:"Ports"`
		Status   string `json:"Status"`
		State    string `json:"State"`
		Replicas string `json:"Replicas"`
		Mode     string `json:"Mode"`
	}

	var c service
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal container JSON: %w", err)
	}

	ports := []string{}
	if c.Ports != "" {
		unique := make(map[string]string)
		for _, p := range strings.Split(c.Ports, ",") {
			if strings.Contains(p, "->") {
				arr := strings.Split(p, "->")
				public := strings.Trim(arr[0], " ")
				public = strings.TrimPrefix(public, ":::")
				public = strings.TrimPrefix(public, "0.0.0.0:")
				public = strings.TrimPrefix(public, "*:")
				if len(arr) > 1 {
					private := strings.Trim(arr[1], " ")
					private = strings.TrimSuffix(private, "/tcp")
					unique[private] = public
				}
			} else {
				port := strings.Trim(p, " ")
				port = strings.TrimPrefix(port, ":::")
				port = strings.TrimSuffix(port, "/tcp")
				unique[port] = ""
			}
		}
		for k, v := range unique {
			if v == "" {
				ports = append(ports, k)
				continue
			}
			ports = append(ports, fmt.Sprintf("%s:%s", v, k))
		}
	}

	name := c.Names
	if c.Name != "" {
		name = c.Name
	}

	status := c.Status
	if c.Mode != "" {
		status = c.Mode
	}

	state := c.State
	if c.Replicas != "" {
		state = c.Replicas
	}

	imageArr := strings.Split(c.Image, ":")
	image := c.Image
	tag := "latest"
	if len(imageArr) > 0 {
		image = imageArr[0]
	}
	if len(imageArr) > 1 {
		tag = imageArr[1]
	}

	return []string{c.ID, name, image, tag, strings.Join(ports, ", "), status, state}, nil
}
