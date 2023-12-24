package main

import (
	"bufio"
  "fmt"
	"log"
	"net"
	"os"
)

func handle_incoming_msg(conn *net.TCPConn, inmsgch chan string) {
	reply := make([]byte, 0, 255)

	for {
		_, err := conn.Read(reply)
		if err != nil {
			fmt.Println("Read from server failed:", err.Error())
			conn.Close()
			os.Exit(1)
		}
		r := string(reply)
		fmt.Println("reply from server=", r)
		inmsgch <- r
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Absent argument! You must specify IP:PORT in this format: 'xxxx.yyyy.zzzz.wwww:pppp'")
	}
	servAddr := os.Args[1] // IP:PORT
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		os.Exit(1)
	}
	inmsgch := make(chan string)
	go handle_incoming_msg(conn, inmsgch)

	for {
    reader := bufio.NewReader(os.Stdin)
    b, err := reader.ReadBytes('\n')
		if err != nil {
			os.Exit(1)
		}
		_, err = conn.Write(b)
		if err != nil {
			conn.Close()
			fmt.Println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		println("write to server = ", string(b))
	}

}
