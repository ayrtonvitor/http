package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal("Could not open file: %w", err)
	}
	defer file.Close()

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

		lineChan := getLineChannel(conn)
		for line := range lineChan {
			fmt.Println(line)
		}

		log.Printf("Connection %s closed\n", conn.RemoteAddr().String())
	}
}

func getLineChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		line := ""
		for {
			buf := make([]byte, 8)
			n, err := f.Read(buf)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					log.Fatal("Error reading from file: %w", err)
				}
				if line != "" {
					lines <- line
				}
				return
			}
			parts := strings.Split(string(buf[:n]), "\n")
			for _, part := range parts[:len(parts)-1] {
				lines <- line + part
				line = ""
			}
			line += string(parts[len(parts)-1])
		}
	}()

	return lines
}
