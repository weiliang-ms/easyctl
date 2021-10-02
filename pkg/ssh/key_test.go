package ssh

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 生成密钥对测试用例
func TestMakeSSHKeyPair(t *testing.T) {
	prv, pub, err := MakeSSHKeyPair()
	assert.Nil(t, err)
	fmt.Println(prv)
	fmt.Println(pub)
}
