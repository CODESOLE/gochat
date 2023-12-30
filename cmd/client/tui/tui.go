package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func RenderPromptUI(width, height int, s []byte) {
	p := widgets.NewParagraph()
	p.Text = string(s)
	p.Border = true
	p.SetRect(0, height-3, width, height)

	ui.Render(p)
}
