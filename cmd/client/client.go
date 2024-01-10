package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var msgstack = make([]*tview.TextView, 0, 4)
var stack_idx int = 0

func handle_incoming_msg(conn *net.TCPConn, grid *tview.Grid, app *tview.Application) {
	reply := make([]byte, 255)

	for {
		n, err := conn.Read(reply)
		if err != nil {
			fmt.Println("Read from server failed:", err.Error())
			conn.Close()
			os.Exit(1)
		}
		tv := tview.NewTextView().SetDynamicColors(true)
    txt := fmt.Sprintf("From [red]%s[white] on [red]%s [green]=> [yellow]%s", conn.RemoteAddr().String(), time.Now().Local().Format(time.ANSIC), string(reply[:n]))
    tv.SetText(txt).SetBorder(false)
		tv.SetChangedFunc(func() { app.Draw() })
		msgstack = append(msgstack, tv)
		var rows = make([]int, len(msgstack))
		for i := range msgstack {
			rows[i] = msgstack[len(msgstack)-1-i].GetFieldHeight() + 1
			grid.AddItem(msgstack[len(msgstack)-1-i], i, 0, rows[i], 1, 0, 0, false)
		}
		grid.SetRows(rows...)
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
	grid := tview.NewGrid()
	grid.SetColumns(0).SetGap(0, 0).SetBorders(true).SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyUp:
			if stack_idx == 0 {
				return ev
			}

			stack_idx--

		case tcell.KeyDown:
			if stack_idx == len(msgstack)-1 {
				return ev
			}

			stack_idx++

		default:
			return ev
		}

		grid.SetOffset(stack_idx-1, 0)

		return nil
	})
	inputField.SetLabel("Enter Message: ").
		SetPlaceholder("Your message").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyCtrlC {
				app.Stop()
			} else if key == tcell.KeyEnter {
				_, err = conn.Write([]byte(inputField.GetText()))
        if err != nil {
          log.Fatal("Failed to write to server!\n")
        }
				inputField.SetText("")
			}
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(inputField, 1, 1, true).
		AddItem(grid, 0, 1, false)

	go handle_incoming_msg(conn, grid, app)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
