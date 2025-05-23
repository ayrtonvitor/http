package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	add, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("Could not resolve udp address: %w", err)
	}
	conn, err := net.DialUDP("udp", nil, add)
	if err != nil {
		log.Fatal("Could not resolve udp connection: %w", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from std: %v", err)
			continue
		}
		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Printf("Error writing to udp conn: %v", err)
			continue
		}
	}
}
