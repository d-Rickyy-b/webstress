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
	lastCheckedMessages int64
	Connections         []*WebsocketData
	Addr                string
}

func (i *InfoBox) Draw(screen tcell.Screen) {
	i.Box.DrawForSubclass(screen, i)
	x, y, width, height := i.GetInnerRect()
	if width <= 0 || height <= 0 {
		return
	}

	totalMessages := i.CalcTotalMessages()
	rate := i.GetRate()

	format := "%s | Connections: [%d/%d] | Total Messages: %d | Messages/Second: %d/s"
	text := fmt.Sprintf(format, i.Addr, i.GetActiveConnections(), len(i.Connections), totalMessages, rate)

	tview.Print(screen, text, x+1, y, width-2, tview.AlignCenter, tcell.ColorWhite)

	i.lastChecked = time.Now()
	i.lastCheckedMessages = totalMessages
}

func (i InfoBox) CalcTotalMessages() (total int64) {
	for _, connection := range i.Connections {
		total += connection.MessageCount()
	}
	return total
}

func (i InfoBox) GetActiveConnections() (total int) {
	for _, connection := range i.Connections {
		if connection.Connected {
			total++
		}
	}
	return total
}

func (i InfoBox) GetRate() (rate int64) {
	return i.Connections[0].MessageRate()
}

func NewInfoBox() *InfoBox {
	return &InfoBox{
		Box: tview.NewBox(),
	}
}
