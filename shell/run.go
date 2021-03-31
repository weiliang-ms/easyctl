package shell

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
)

type ExecResult struct {
	ExitCode int
	StdErr   string
	StdOut   string
}

func Run(command string) (re ExecResult) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", command)
	//log.Printf("%s 执行语句：%s\n", PrintCyan(constant.Shell), command)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err.Error())
	}

	// 标准输出
	logScan := bufio.NewScanner(stdout)
	go func() {
		for logScan.Scan() {
			log.Println(logScan.Text())
			re.StdOut = logScan.Text()
		}
	}()

	// 标准错误输出
	errBuf := bytes.NewBufferString("")
	scan := bufio.NewScanner(stderr)
	for scan.Scan() {
		s := scan.Text()
		log.Println("build error: ", s)
		errBuf.WriteString(s)
		errBuf.WriteString("\n")
		re.StdErr = logScan.Text()
	}

	// 等待命令执行完
	cmd.Wait()

	re.ExitCode = cmd.ProcessState.ExitCode()

	return re

}
