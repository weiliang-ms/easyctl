package util

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/constant"
	"os"
	"strings"
	"sync"
)

func CreateFile(filePath string, content string) {
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	_, writeErr := file.WriteString(content)
	if writeErr != nil {
		panic(writeErr)
	}
}

func OverwriteContent(filePath string, content string) {

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, writeErr := file.WriteString(content)
	if writeErr != nil {
		fmt.Println(err.Error())
	}
}

func FormatFileName(name string) string {
	array := strings.Split(name, "/")
	if len(array) > 0 {
		name = array[len(array)-1]
	}
	return name
}

func WriteFile(filePath string, b []byte, serverList []Server) {
	if len(serverList) == 0 {
		PrintActionBanner([]string{constant.LoopbackAddress}, fmt.Sprintf("写文件：%s", filePath))
		OverwriteContent(filePath, string(b))
	} else {
		wg := sync.WaitGroup{}
		wg.Add(len(serverList))
		for _, v := range serverList {
			PrintActionBanner([]string{v.Host}, fmt.Sprintf("写文件：%s:%s", v.Host, filePath))
			go RemoteWriteFileParallel(filePath, b, v, &wg)
		}
		wg.Wait()
	}
}
