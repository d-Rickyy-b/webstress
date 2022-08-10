package models

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"sync/atomic"
	"webstress/internal/util"
)

type WebsocketData struct {
	ID           int
	messageCount uint64
	counter      *MsgCounter
	counterRate  uint64
	Connected    bool
}

func NewWebsocketData(id int, counter *MsgCounter) *WebsocketData {
	return &WebsocketData{ID: id, counter: counter, counterRate: 5}
}

func (w *WebsocketData) IncCounter() {
	w.counter.Incr()
	atomic.AddUint64(&w.messageCount, 1)
}

func (w *WebsocketData) MessageCount() uint64 {
	return atomic.LoadUint64(&w.messageCount)
}

// WebsocketView is the full view that contains all the individual websocket boxes
type WebsocketView struct {
	*tview.Box
	Connections    []*WebsocketData
	WebsocketCount int
}

func NewWebsocketView(websocketCount int) *WebsocketView {
	box := tview.NewBox()
	box.SetBorder(true).SetTitle(" Websockets ")

	return &WebsocketView{
		Box:            box,
		WebsocketCount: websocketCount,
		Connections:    make([]*WebsocketData, websocketCount),
	}
}

func (w WebsocketView) Draw(screen tcell.Screen) {
	w.Box.DrawForSubclass(screen, w)
	x, y, width, height := w.GetInnerRect()
	if width <= 0 || height <= 0 {
		return
	}
	numOfChars := util.NumOfChars(w.WebsocketCount)

	// Calclulate maximum length of all the websocket connections
	formatString := "%0" + strconv.Itoa(numOfChars) + "d:%d"

	var maxLength int
	for _, connection := range w.Connections {
		tmpStr := fmt.Sprintf(formatString, connection.ID, connection.MessageCount())
		tmpSize := len(tmpStr)
		maxLength = util.Max(maxLength, tmpSize)
	}

	// Calculate number of elements per row
	numPerRow := (width - 2) / (maxLength + 1)

	newStyle := tcell.Style{}
	newStyle = newStyle.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	for i, connection := range w.Connections {
		row := i / numPerRow
		col := ((i % numPerRow) * (maxLength + 1)) + 1

		tmpStr := fmt.Sprintf(formatString, connection.ID+1, connection.MessageCount())
		tview.Print(screen, tmpStr, x+col, y+row, maxLength+1, tview.AlignLeft, tcell.ColorWhite)

		if connection.Connected {
			newStyle = newStyle.Background(tcell.ColorGreen)
		} else {
			newStyle = newStyle.Background(tcell.ColorRed)
		}

		for i := 0; i < numOfChars; i++ {
			m, c, _, _ := screen.GetContent(x+col+i, y+row)
			screen.SetContent(x+col+i, y+row, m, c, newStyle)
		}
	}
}
