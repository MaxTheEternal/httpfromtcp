package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
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
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic("upsi")
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			panic("upsi")
		}

		fmt.Println("Connected")
		for str := range getLinesChannel(con) {
			fmt.Printf("read: %s\n", str)
		}
	}

}
