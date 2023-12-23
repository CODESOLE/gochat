package main

import (
	"fmt"
	"log"
	"net"
  "github.com/CODESOLE/gochat/internal/core"
)

const Port = "1234"

func handle_client(conn net.Conn, msg_ch chan string) {
	buf := make(core.Payload, 255)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Client(%s) disconnected!\n", conn.RemoteAddr().String())
			conn.Close()
			return
		}
		msg_ch <- fmt.Sprintf("Client(%s) send: %s", conn.RemoteAddr().String(), string(buf[:n]))
	}
}

func server(msg_ch chan string) {
	for s := range msg_ch {
		fmt.Println(s)
	}
}

func main() {
	msg_ch := make(chan string)
	listener, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen to port %s: %s\n", Port, err.Error())
	}
	fmt.Printf("Listening to TCP connections on port %s ...\n", Port)
	go server(msg_ch)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("An error occured while accepting connection")
		}
		log.Printf("Client(%s) connected!\n", conn.RemoteAddr().String())
		go handle_client(conn, msg_ch)
	}
}
