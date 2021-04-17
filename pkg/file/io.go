package file

import (
	"io/ioutil"
	"os"
)

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
