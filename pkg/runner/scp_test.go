package runner

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestScpErrorPathFile(t *testing.T) {
	item := ScpItem{SrcPath: "1.tt"}
	server := ServerInternal{}
	err := server.Scp(item)
	assert.EqualError(t, err, "CreateFile 1.tt: The system cannot find the file specified.")
}

func TestScpNilFile(t *testing.T) {
	f, _ := os.Create("1.txt")
	item := ScpItem{SrcPath: "1.txt"}
	server := ServerInternal{}
	err := server.Scp(item)
	assert.EqualError(t, err, "源文件: 1.txt 大小为0")
	_ = f.Close()
	_ = os.Remove("1.txt")
}

func TestConnectErr(t *testing.T) {
	f, _ := os.Create("1.txt")
	_, _ = f.Write([]byte("123"))
	item := ScpItem{SrcPath: "1.txt"}
	server := ServerInternal{}
	err := server.Scp(item)
	fmt.Println(err)

	assert.EqualError(t, err, "连接远程主机：失败 ->dial tcp :0: connectex: The requested address is not valid in its context.")
	_ = f.Close()
	_ = os.Remove("1.txt")
}
