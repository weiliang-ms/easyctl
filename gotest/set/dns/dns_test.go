package passwordless

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"os"
	"testing"
)

func TestParseDnsConfig(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	c, err := set.ParseDnsConfig(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", c.DnsList)
}

func TestAddDnsScript(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	c, err := set.AddDnsScript(b, set.NewPasswordTmpl)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", c)
}

//
func TestSetDns(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	err := set.Dns(b, logger)
	if err != nil {
		panic(err)
	}
}
