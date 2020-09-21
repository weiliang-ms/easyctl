package util

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
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

func ExecuteCmdAcceptResult(command string) (result string, err error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	fmt.Printf("[shell] 执行语句：%s\n", command)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Println("exec the cmd ", " failed")
		fmt.Println(err.Error())
		return "", err
	}
	// 正常日志
	logScan := bufio.NewScanner(stdout)
	go func() {
		for logScan.Scan() {
			fmt.Println(logScan.Text())
			result = logScan.Text()
		}
	}()
	// 错误日志
	errBuf := bytes.NewBufferString("")
	scan := bufio.NewScanner(stderr)
	for scan.Scan() {
		s := scan.Text()
		log.Println("build error: ", s)
		errBuf.WriteString(s)
		errBuf.WriteString("\n")
	}
	// 等待命令执行完
	cmd.Wait()
	if !cmd.ProcessState.Success() {
		// 执行失败，返回错误信息
		return "", errors.New(errBuf.String())
	}
	return result, nil
}
