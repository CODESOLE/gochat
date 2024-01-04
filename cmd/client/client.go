package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// var boxes []tview.Box

func handle_incoming_msg(conn *net.TCPConn, flex *tview.Flex, app *tview.Application) {
	reply := make([]byte, 255)

	for {
		n, err := conn.Read(reply)
		if err != nil {
			fmt.Println("Read from server failed:", err.Error())
			conn.Close()
			os.Exit(1)
		}
		tv := tview.NewTextView().SetChangedFunc(func() { app.Draw() })
    trimmed_str := strings.TrimSpace(string(reply[:n]))
		fmt.Fprintf(tv, "%s", trimmed_str)
		tv.SetBorder(true).SetTitle(time.Now().String())
		flex.AddItem(tv, 0, 1, false)
	}
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

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Couldn' t connect to server':", err.Error())
		os.Exit(1)
	}

	app := tview.NewApplication()
	inputField := tview.NewInputField()
	inputField.SetLabel("Enter Message: ").
		SetPlaceholder("Your message").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyCtrlC {
				app.Stop()
			} else if key == tcell.KeyEnter {
				conn.Write([]byte(inputField.GetText()))
				inputField.SetText("")
			}
		})
	inputField.SetBorder(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(inputField, 0, 1, true)

	go handle_incoming_msg(conn, flex, app)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
