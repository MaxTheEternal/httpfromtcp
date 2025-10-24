package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MaxTheEternal/httpfromtcp/internal/request"
)

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

		req, err := request.RequestFromReader(con)
		if err != nil {
			log.Fatal(err)
		}
		rl := req.RequestLine
		fmt.Printf("Request Line:\n- Method: %s\n- Target: %s\n- Version: %s\n", rl.Method, rl.RequestTarget, rl.HttpVersion)
		fmt.Printf("Headers:\n")

		for k, v := range req.Headers.Headers {
			fmt.Printf("- %s: %s", k, v)
		}
	}
}
