package osutil

import (
	"io/fs"
	"os"
)

func MkDirPanicErr(path string, mode fs.FileMode) {
	err := os.Mkdir(path, mode)
	if err != nil {
		panic(err)
	}
}

func WriteFilePanicErr(filepath string, b []byte) int {
	f, err := os.Create(filepath)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	count, _ := f.Write(b)

	// todo: this case
	//if err != nil {
	//	panic(err)
	//}

	return count
}
