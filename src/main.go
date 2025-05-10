package main

import (
	"fmt"
	"net"
)

const ADDR string = "127.0.0.1"
const PORT string = ":8000"

func main() {
	// listen for incoming connections on port 8000
	ln, err := net.Listen("tcp", ADDR+PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("TCP listening on", ADDR+PORT)

	// accept incoming connections and handle them
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// handle the connections in a new goroutine
		go HandleConnection(conn)
	}
}
