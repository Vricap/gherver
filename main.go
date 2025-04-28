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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// close the connection when we're done
	defer conn.Close()

	// read incoming data
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// send back the data
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the incoming data
	fmt.Printf("Received: %s", buf)
}
