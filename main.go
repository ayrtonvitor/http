package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal("Could not open file: %w", err)
	}
	defer file.Close()

	line := ""
	for {
		buf := make([]byte, 8)
		n, err := file.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Fatal("Error reading from file: %w", err)
			}
			if line != "" {
				fmt.Printf("read: %s%d\n", line, len(line))
			}
			break
		}
		parts := strings.Split(string(buf[:n]), "\n")
		for _, part := range parts[:len(parts)-1] {
			fmt.Printf("read: %s%s\n", line, part)
			line = ""
		}
		line += string(parts[len(parts)-1])
	}
}
