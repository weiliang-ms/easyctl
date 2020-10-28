package sys

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/util"
)

const defaultPassword = "user123"

func AddUser(username string, password string, login bool) {
	var cmd string

	// todo 校验username合法性
	if password == "" && login == true {
		password = defaultPassword
	} else if password != "" && login == false {
		password = ""
	}

	fmt.Printf("创建用户：%s 登录类型：%t 密码：%s\n", username, login, password)

	if login {
		cmd = fmt.Sprintf("useradd -m %s  && echo \"%s\" | passwd --stdin %s", username, password, username)
	} else {
		cmd = fmt.Sprintf("groupadd %s;useradd %s -g %s -s /sbin/nologin -M", username, username, username)
	}

	util.ExecuteCmd(cmd)
}

func AddSudoUser(username string) string {
	return "chmod +w /etc/sudoers;" +
		fmt.Sprintf("echo \"%sALL=(ALL) NOPASSWD: ALL\" >> /etc/sudoers;", username) +
		"chmod -w /etc/sudoers"
}
