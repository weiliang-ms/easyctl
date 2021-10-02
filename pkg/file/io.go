package file

import (
	"io/ioutil"
	"os"
)

// ReadAll 读取文件全部内容
func ReadAll(filePath string) string {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(b)
}
