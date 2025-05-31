package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ayrtonvitor/http/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Could not set up listener: %w", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		log.Printf("New connection accepted. Remote address: %s",
			conn.RemoteAddr().String())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("Request error: %s", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		log.Printf("Connection %s closed\n", conn.RemoteAddr().String())
	}
}
