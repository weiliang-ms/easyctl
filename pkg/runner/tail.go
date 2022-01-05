/*
	MIT License

Copyright (c) 2020 xzx.weiliang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
package runner

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// TailFile all & real time
func (server ServerInternal) TailFile(path string, offset int64, whence int, stopCh <-chan struct{}) {
	// init sftp
	sftp, err := SftpConnect(server.Username, server.Password, server.Host, server.Port)
	if err != nil {
		log.Println(err)
	}

	f, err := sftp.OpenFile(path, os.O_RDONLY)
	if err != nil {
		panic(err)
	}

	f.Seek(offset, whence)

	for {
		select {
		case <-stopCh:
			_ = f.Close()
			return
		default:
		}

		buf := bufio.NewReader(f)

		for {
			select {
			case <-stopCh:
				_ = f.Close()
				return
			default:
			}
			line, err := buf.ReadString('\n')
			line = strings.TrimSpace(line)
			if line != "" {
				fmt.Printf("[%s] %s\n", server.Host, line)
			}
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Println("Read file error!", err)
					return
				}
			}
		}
	}
}
