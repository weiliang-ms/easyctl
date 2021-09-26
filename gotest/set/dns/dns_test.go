package passwordless

import (
	"fmt"
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

	err := set.Dns(b, true)
	if err != nil {
		panic(err)
	}
}
