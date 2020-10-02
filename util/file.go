package util

import (
	"os"
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
