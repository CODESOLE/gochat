package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"time"

	"github.com/CODESOLE/gochat/internal/core"
)

type Clients []net.Conn

var clients = make(Clients, 0, 100)

func handle_client(conn net.Conn, msg_ch chan core.Payload) {
	pl := core.Payload{}
	pl.Msg = make([]byte, 255)
	pl.IpAddr = conn.RemoteAddr().String()
	pl.Conn = conn
	pl.SendTime = time.Now().String()
	for {
		n, err := conn.Read(pl.Msg)
		if err != nil {
			log.Printf("Client(%s) disconnected!\n", conn.RemoteAddr().String())
			clients = slices.DeleteFunc(clients, func(c net.Conn) bool { return c.RemoteAddr().String() == conn.RemoteAddr().Network() })
			conn.Close()
			return
		}
		fmt.Printf("Client(%s) send: %s", conn.RemoteAddr().String(), string(pl.Msg[:n]))
		msg_ch <- pl
	}
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func server(msg_ch chan core.Payload) {
	for s := range msg_ch {
		others := filter[net.Conn](clients, func(c net.Conn) bool {
			return c.RemoteAddr().String() != s.IpAddr
		})
		for i := range others {
			fmt.Println(others[i].RemoteAddr().String())
			_, err := others[i].Write(s.Msg)
			if err != nil {
				s.Conn.Close()
				log.Printf("Cannot send message to %v\n", s.IpAddr)
				continue
			}
			log.Printf("Message sent from %v to %v", s.IpAddr, others[i].RemoteAddr().String())
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Absent argument! You must specify the PORT to be listened in.")
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
