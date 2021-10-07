package add

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestUser(t *testing.T) {
	// test error ParseNewUserConfig
	var b = `
new-user:
  - name: user01
  nologin: false
  password: "dddd"  #
  user-dir: ""      # default /home/username
`
	assert.Errorf(t, User(command.OperationItem{B: []byte(b), Logger: logrus.New()}),
		"yaml: unmarshal errors:\n  line 3: cannot unmarshal !!map into []runner.ServerExternal")

	// 测试参数合法性
	b = `
new-user:
  name: user01
  nologin: false
  password: "dddd"  #
  user-dir: ""      # default /home/username
`
	assert.Errorf(t, User(command.OperationItem{B: []byte(b), Logger: logrus.New()}), "密码长度不能小于6位")

	// 测试模板赋值传错误参数
	b = `
new-user:
  name: user01
  nologin: false
  password: "dd12dd"  #
  user-dir: "/ddd"      # default /home/username
`
	fmt.Println(User(command.OperationItem{B: []byte(b), Logger: logrus.New()}))
}
