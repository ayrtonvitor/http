package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal("Could not open file: %w", err)
	}
	defer file.Close()

	chunk := make([]byte, 8)
	for {
		n, err := file.Read(chunk)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Fatal("Error reading from file: %w", err)
			}
			break
		}

		fmt.Printf("read: %s\n", string(chunk[:n]))
	}
}
