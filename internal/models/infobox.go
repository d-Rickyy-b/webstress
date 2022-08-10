package models

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type InfoBox struct {
	*tview.Box
	lastChecked         time.Time
	lastCheckedMessages uint64
	Connections         []*WebsocketData
	MsgCounter          *MsgCounter
	Addr                string
}

func (i *InfoBox) Draw(screen tcell.Screen) {
	i.Box.DrawForSubclass(screen, i)
	x, y, width, height := i.GetInnerRect()
	if width <= 0 || height <= 0 {
		return
	}

	totalMessages := i.MsgCounter.Count()
	rate := i.MsgCounter.Rate()

	format := "%s | Connections: [%d/%d] | Total Messages: %d | Messages/Second: %d/s"
	text := fmt.Sprintf(format, i.Addr, i.GetActiveConnections(), len(i.Connections), totalMessages, rate)

	tview.Print(screen, text, x+1, y, width-2, tview.AlignCenter, tcell.ColorWhite)

	i.lastChecked = time.Now()
	i.lastCheckedMessages = totalMessages
}

func (i *InfoBox) GetActiveConnections() (total int) {
	for _, connection := range i.Connections {
		if connection.Connected {
			total++
		}
	}
	return total
}

func NewInfoBox() *InfoBox {
	return &InfoBox{
		Box: tview.NewBox(),
	}
}
