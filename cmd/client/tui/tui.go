package tui

import (
	"github.com/CODESOLE/gochat/internal/core"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func RenderPromptUI(width, height int, s []byte) {
	p := widgets.NewParagraph()
	p.Text = string(s)
	p.Border = true
	p.SetRect(0, height-3, width, height)

	ui.Render(p)

	for x := range core.MsgStack {
		core.MsgStack[x].View.SetRect(0, core.MsgStack[x].View.GetRect().Min.Y, width, core.MsgStack[x].View.GetRect().Max.Y)
		ui.Render(core.MsgStack[x].View)
	}
}
