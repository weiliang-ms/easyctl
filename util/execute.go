package util

import (
	"bufio"
	"bytes"
	"easyctl/constant"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func ExecuteCmd(command string) (err error, result string) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", command)
	log.Printf("%s 执行语句：%s\n", PrintCyan(constant.Shell), command)
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

func ExecuteCmdIgnoreErr(command string) (result string) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", command)
	fmt.Printf("%s 执行语句：%s\n", PrintCyan(constant.Shell), command)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	shellErr := cmd.Run()
	if shellErr != nil {
		fmt.Println(shellErr.Error())
	}
	return result
}

func ExecuteCmdAcceptResult(command string) (result string) {
	cmd := exec.Command("/bin/bash", "-c", command)
	PrintActionBanner([]string{constant.LoopbackAddress}, fmt.Sprintf("执行命令：%s", cmd))
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err.Error())
		return ""
	}
	// 正常日志
	logScan := bufio.NewScanner(stdout)
	go func() {
		for logScan.Scan() {
			//fmt.Println(logScan.Text())
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
		return ""
	}

	//fmt.Println("返回结果：",result)
	return result
}

func ExecuteCmdPrintStd(command string) {
	cmd := exec.Command("/bin/bash", "-c", command)
	PrintActionBanner([]string{constant.LoopbackAddress}, fmt.Sprintf("执行命令：%s", cmd))
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err.Error())
	}
	// 正常日志
	logScan := bufio.NewScanner(stdout)
	go func() {
		for logScan.Scan() {
			fmt.Println(logScan.Text())
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
		log.Fatal("执行失败...")
	}
}

func ExecuteCmdResult(command string) (result string, err error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return "", err
	}
	// 正常日志
	logScan := bufio.NewScanner(stdout)
	go func() {
		for logScan.Scan() {
			result = logScan.Text()
		}
	}()
	// 错误日志
	errBuf := bytes.NewBufferString("")
	scan := bufio.NewScanner(stderr)
	for scan.Scan() {
		s := scan.Text()
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

func ExecuteIgnoreStd(shell string) bool {
	cmd := exec.Command("/bin/bash", "-c", shell)
	if err := cmd.Start(); err != nil {
		return false
	}
	// 等待命令执行完
	cmd.Wait()
	if !cmd.ProcessState.Success() {
		// 执行失败，返回错误信息
		return false
	}

	return true
}

func Run(command string) int {
	cmd := exec.Command("/bin/bash", "-c", command)
	//fmt.Printf("执行命令：%s", cmd)

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal(err.Error())
	}
	// 正常日志
	logScan := bufio.NewScanner(stdout)
	go func() {
		for logScan.Scan() {
			fmt.Println(logScan.Text())
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
		log.Fatal("执行失败...")
	}
	return cmd.ProcessState.ExitCode()
}
