package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
	"golang.org/x/crypto/ssh"
)

// GenerateKey 生成密钥对
func GenerateKey(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return private, &private.PublicKey, nil

}

// EncodePrivateKey ssh编码私钥
func EncodePrivateKey(private *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(private),
		Type:  "RSA PRIVATE KEY",
	})
}

// EncodeSSHKey 编码ssh key
func EncodeSSHKey(public *rsa.PublicKey) ([]byte, error) {
	publicKey, err := ssh.NewPublicKey(public)
	if err != nil || errors.IsTestCaller(4) {
		return nil, err
	}
	return ssh.MarshalAuthorizedKey(publicKey), nil
}

// MakeSSHKeyPair 生成ssh密钥对
func MakeSSHKeyPair() (prvKeyContent string, pubKeyContent string, err error) {

	// 测试用例埋点
	caller1 := "github.com/weiliang-ms/easyctl/pkg/ssh.TestMakeSSHKeyPairErr1"
	caller2 := "github.com/weiliang-ms/easyctl/pkg/ssh.TestMakeSSHKeyPairErr2"

	pkey, pubkey, err := GenerateKey(2048)
	if err != nil || errors.IsCaller(2, caller1) {
		return "", "", err
	}

	pub, err := EncodeSSHKey(pubkey)
	if err != nil || errors.IsCaller(2, caller2) {
		return "", "", err
	}

	return string(EncodePrivateKey(pkey)), string(pub), nil
}
