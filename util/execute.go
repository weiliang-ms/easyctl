package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ExecuteCmd(command string) (err error, result string) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", command)
	fmt.Printf("[shell] 执行语句：%s\n", command)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	shellErr := cmd.Run()
	if shellErr != nil {
		err = shellErr
	}
	return err, result
}
