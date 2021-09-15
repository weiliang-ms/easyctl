package ssh

import (
	"fmt"
	"testing"
)

func TestMakeSSHKeyPair(t *testing.T) {
	prv, pub, err := MakeSSHKeyPair()
	if err != nil {
		panic(err)
	}

	fmt.Println(prv)
	fmt.Println(pub)
}
