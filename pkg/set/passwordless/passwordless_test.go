package passwordless

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"testing"
)

func TestGenRsaKey(t *testing.T) {
	err := Config("config.yaml", util.Debug)
	if err != nil {
		panic(err)
	}
}

func TestMakeKeyPairScript(t *testing.T) {
	script, err := MakeKeyPairScript(passwordLessTmpl)
	if err != nil {
		panic(err)
	}

	fmt.Println(script)
}
