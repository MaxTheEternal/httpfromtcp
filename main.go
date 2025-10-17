package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	line := ""
	strings := make(chan string)
	go func() {
		defer f.Close()
		defer close(strings)
		for {
			buf := make([]byte, 8)
			n, err := f.Read(buf)

			if err == io.EOF {
				break // End of file, break the loop
			}
			if err != nil {
				panic("upsi")
			}

			buf = buf[:n]
			if i := bytes.IndexByte(buf, '\n'); i != -1 {
				line += string(buf[:i])
				buf = buf[i+1:]
				strings <- line
				line = ""
			}

			line += string(buf)

		}
	}()
	return strings

}
func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		panic("upsi")
	}

	lines := getLinesChannel(file)
	for str := range lines {
		fmt.Printf("read: %s\n", str)
	}

}
