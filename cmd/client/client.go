package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/CODESOLE/gochat/cmd/client/tui"
	"github.com/CODESOLE/gochat/internal/core"
	ui "github.com/gizak/termui/v3"
)

func handle_incoming_msg(conn *net.TCPConn) {
	reply := make([]byte, 255)

	for {
		n, err := conn.Read(reply)
    core.PushStack(reply[:n], conn.RemoteAddr().String(), time.Now().String())
		if err != nil {
			fmt.Println("Read from server failed:", err.Error())
			conn.Close()
			os.Exit(1)
		}
	}
}

func is_alphanum(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Missing argument! You must specify IP:PORT in this format: 'xxxx.yyyy.zzzz.wwww:pppp'")
	}
	servAddr := os.Args[1] // IP:PORT
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		log.Fatalln("ResolveTCPAddr failed: ", err.Error())
		os.Exit(1)
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Couldn' t connect to server':", err.Error())
		os.Exit(1)
	}
	go handle_incoming_msg(conn)

	buf := make([]byte, 255)
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Millisecond * 100).C
	var ui_width, ui_height int = ui.TerminalDimensions()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID { // event string/identifier
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				ui_width, ui_height = payload.Width, payload.Height
			}
			switch e.Type {
			case ui.KeyboardEvent: // handle all key presses
				eventID := e.ID // keypress string
				if eventID == "<Tab>" {
					for range [4]byte{} {
						buf = append(buf, ' ')
					}
				} else if eventID == "<Space>" {
					buf = append(buf, ' ')
				} else if eventID == "<Backspace>" || eventID == "<C-<Backspace>>" {
          buf = buf[:len(buf) - 1]
				}
				if len(eventID) == 1 && is_alphanum(eventID[0]) {
					buf = append(buf, eventID[0])
				} else if eventID == "<Enter>" {
					_, err = conn.Write(buf)
					if err != nil {
						conn.Close()
						fmt.Println("Write to server failed: ", err.Error())
						os.Exit(1)
					}
					buf = nil
				}
			}
		// use Go's built-in tickers for updating and drawing data
		case <-ticker:
			tui.RenderPromptUI(ui_width, ui_height, buf)
		}
	}
}
