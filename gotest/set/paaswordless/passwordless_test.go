package passwordless

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"os"
	"testing"
)

func TestGenRsaKey(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	err := set.PasswordLess(b, logger)
	if err != nil {
		panic(err)
	}
}

func TestMakeKeyPairScript(t *testing.T) {
	script, err := set.MakeKeyPairScript(set.PasswordLessTmpl)
	if err != nil {
		panic(err)
	}

	fmt.Println(script)
}
