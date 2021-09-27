package deny

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"strings"
)

func Selinux(config []byte, logger *logrus.Logger) error {
	return Item(config, logger, closeSELinuxShell)
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
