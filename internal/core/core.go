package core

import (
	"fmt"
	"net"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Payload struct {
	Msg      []byte
	Conn     net.Conn
	IpAddr   string
	SendTime string
}

type MsgBox struct {
	TimeStamp string
	RmtAdrr   string
	Msg       string
	View      *widgets.Paragraph
}

var MsgStack []MsgBox

func PushStack(msg []byte, ipaddr string, ts string) {
	msg_str := string(msg)
	p := widgets.NewParagraph()
	p.Text = msg_str
	p.Border = true
	p.TitleStyle.Bg = ui.ColorRed
	p.TitleStyle.Fg = ui.ColorBlack
	p.Title = fmt.Sprintf("From %s on %s", ipaddr, ts)

	w, h := ui.TerminalDimensions()

	msg_h := (len(msg_str) / (w - 2)) + 1
	height_shift_up_off := msg_h + 2 + 1 // +2 for borders

	for i := range MsgStack {
		MsgStack[i].View.SetRect(MsgStack[i].View.GetRect().Min.X, MsgStack[i].View.GetRect().Min.Y+height_shift_up_off, MsgStack[i].View.GetRect().Max.X, MsgStack[i].View.GetRect().Max.Y)
	}

	p.SetRect(0, h-3-(msg_h+2), w, (h-3-(msg_h+2))+msg_h+2) // -3 for bottom command prompt height
	MsgStack = append(MsgStack, MsgBox{ts, ipaddr, msg_str, p})
}
