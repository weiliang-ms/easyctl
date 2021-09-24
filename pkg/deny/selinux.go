package deny

import (
	"bufio"
	"strings"
)

func Selinux(config []byte, debug bool) error {
	return Item(config, debug, closeSELinuxShell)
}

// todo confirm
func confirm(reader *bufio.Reader) (string, error) {
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)

		if input != "" && (input == "yes" || input == "no") {
			return input, nil
		}
	}
}
