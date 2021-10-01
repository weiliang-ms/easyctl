package runner

import (
	"bufio"
	"fmt"
	"github.com/pkg/sftp"
	"io"
	"log"
	"os"
	"strings"
)

// TailFile all & real time
func (server ServerInternal) TailFile(path string, offset int64, whence int, stopCh <-chan struct{}) {
	// init sftp
	sftp, err := sftpConnect(server.Username, server.Password, server.Host, server.Port)
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

func readAtTheBeginning(f *sftp.File, host string, stopCh <-chan struct{}) {
}

func readAtTheLatest(f *sftp.File, host string, stopCh <-chan struct{}) {
	for {
		//f.Seek(0, 2)
		select {
		case <-stopCh:
			_ = f.Close()
			return
		default:
		}

		f.Seek(0, 2)
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
				fmt.Printf("[%s] %s\n", host, line)
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
