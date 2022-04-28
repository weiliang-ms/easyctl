package runner

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/ssh/terminal"

	expect "github.com/Netflix/go-expect"
)

func getPassword(fd int) string {
	bytePassword, _ := terminal.ReadPassword(fd)

	return string(bytePassword)
}

func TestName(t *testing.T) {

	c, _ := expect.NewConsole()

	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.SendLine("hunter2")
	}()

	echoText := getPassword(int(c.Tty().Fd()))

	<-donec

	fmt.Printf("\nPassword from stdin: %s", echoText)
}
