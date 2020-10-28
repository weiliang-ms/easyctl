package util

import (
	"fmt"
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
