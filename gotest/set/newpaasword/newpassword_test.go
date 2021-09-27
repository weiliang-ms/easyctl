package passwordless

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"os"
	"testing"
)

func TestParseNewPasswordConfig(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	c, err := set.ParseNewPasswordConfig(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", c.Password)
}

func TestNewPasswordScript(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	c, err := set.NewPasswordScript(b, set.NewPasswordTmpl)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", c)
}

func TestNewPassword(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	err := set.NewPassword(b, logger)
	if err != nil {
		panic(err)
	}
}
