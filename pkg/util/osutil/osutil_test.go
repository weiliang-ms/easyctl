package osutil

import (
	"fmt"
	"os"
	"testing"
)

func Test_MkDirPanicErr(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
		os.Remove("ddd")
	}()
	os.Mkdir("ddd", 0644)
	MkDirPanicErr("ddd", 0644)
}

func Test_WriteFilePanicErr_NilFileNameCase(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	WriteFilePanicErr("", []byte("ddd"))
}

func Test_WriteFilePanicErr(t *testing.T) {

	f, _ := os.Create("ddd.txt")

	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
		f.Close()
		os.Remove("ddd.txt")
	}()

	WriteFilePanicErr("ddd.txt", []byte("ddd"))
}
