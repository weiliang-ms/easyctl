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

func TestGenerateKey(t *testing.T) {
	a, b, err := GenerateKey(0)
	assert.Nil(t, a)
	assert.Nil(t, b)
	assert.EqualError(t, err, "crypto/rsa: too few primes of given length to generate an RSA key")
}
