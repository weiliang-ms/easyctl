package file

import (
	"bufio"
	"fmt"
	"github.com/weiliang-ms/gotool/strings"
	"os"
)

func WriteWithIPS(slice strings.IPS, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	for _, v := range slice {
		_, err := fmt.Fprintln(w, v)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}
