package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/CODESOLE/gochat/internal/core"
)

type Clients []net.Conn

var clients = make(Clients, 0, 100)
var mu sync.Mutex

func handle_client(conn net.Conn, msg_ch chan core.Payload) {
	pl := core.Payload{}
	pl.Msg = make([]byte, 255)
	pl.IpAddr = conn.RemoteAddr().String()
	pl.Conn = conn
	for {
		n, err := conn.Read(pl.Msg)
		pl.SendTime = time.Now().String()
		if err != nil {
			log.Printf("Client(%s) disconnected!\n", conn.RemoteAddr().String())
			conn.Close()
			mu.Lock()
			clients = slices.DeleteFunc(clients, func(c net.Conn) bool {
				return c.RemoteAddr().String() == conn.RemoteAddr().String()
			})
			mu.Unlock()
			return
		}
		log.Printf("Client(%s) send: %s", conn.RemoteAddr().String(), string(pl.Msg[:n]))
		msg_ch <- pl
	}
}

func server(msg_ch chan core.Payload) {
	for s := range msg_ch {
		for _, c := range clients {
			_, err := c.Write(s.Msg)
			if err != nil {
				s.Conn.Close()
				log.Printf("Cannot send message to %v\n", s.IpAddr)
				continue
			}
			log.Printf("Message sent from %v to %v", s.IpAddr, c.RemoteAddr().String())
		}
		clear(s.Msg)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Missing argument! You must specify the PORT to be listened in.")
	}
	port := os.Args[1]

	msg_ch := make(chan core.Payload)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Could not listen to port %s: %s\n", port, err.Error())
	}
	fmt.Printf("Listening to TCP connections on port %s ...\n", port)
	go server(msg_ch)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("An error occured while accepting connection")
			continue
		}
		clients = append(clients, conn)
		log.Printf("Client(%s) connected!\n", conn.RemoteAddr().String())
		go handle_client(conn, msg_ch)
	}
}
