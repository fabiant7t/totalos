package command

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// MemoryModules returns memory module information
func MemoryModules(m remotecommand.Machine, cb ssh.HostKeyCallback) (modules []string, err error) {
	cmd := `dmidecode -t memory`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return modules, fmt.Errorf("Remote command MemoryModules failed: %w", err)
	}

	var inMemoryDevice bool
	var memoryDevices []Module
	var memoryDevice Module
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	for scanner.Scan() {
		trimmedLine := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(trimmedLine, "Memory Device") {
			// End scanning the previous memory device information and push
			// it to the list
			if inMemoryDevice {
				memoryDevices = append(memoryDevices, memoryDevice)
			}
			// Start scanning a new memory device
			inMemoryDevice = true
			memoryDevice = Module{}
			continue
		}
		if inMemoryDevice {
			splitted := strings.SplitAfterN(trimmedLine, ": ", 2)
			if len(splitted) >= 2 {
				key, _ := strings.CutSuffix(splitted[0], ": ")
				value := splitted[1]
				switch key {
				case "Type":
					memoryDevice.Type = value
				case "Size":
					re := regexp.MustCompile(`(?P<size>\d+) (?P<unit>\w+)`)
					if re.MatchString(value) {
						memoryDevice.Size = value
					} else { // Empty slot, skip until reaching next device information
						inMemoryDevice = false
						continue
					}
				case "Speed":
					memoryDevice.Speed = value
				}
			}
		}
	}
	// Push last
	if inMemoryDevice {
		memoryDevices = append(memoryDevices, memoryDevice)
	}
	// Any scanner errors?
	if err := scanner.Err(); err != nil {
		return modules, fmt.Errorf("Remote command MemoryModules failed: %w", err)
	}

	for _, md := range memoryDevices {
		modules = append(modules, md.String())
	}
	return modules, nil
}

type Module struct {
	Size  string
	Type  string
	Speed string
}

func (m *Module) String() string {
	return fmt.Sprintf("%s, %s, %s", m.Size, m.Type, m.Speed)
}
