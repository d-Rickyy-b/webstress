package models

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GraphView struct {
	*tview.Box
	Connections []*WebsocketData
	Counter     *MsgCounter
	buff        []uint64
	dataIndex   int
	mu          sync.RWMutex
}

func NewGraphView() *GraphView {
	box := tview.NewBox()
	box.SetBorder(true).SetTitle(" Messages / Second ")

	return &GraphView{
		Box:  box,
		buff: make([]uint64, 200),
	}
}

func (g *GraphView) Draw(screen tcell.Screen) {
	g.Box.DrawForSubclass(screen, g)
	x, y, width, height := g.GetInnerRect()
	if width <= 0 || height <= 0 {
		return
	}

	lastEntries := g.GetLastEntries(width)
	var maxValue, minValue uint64
	minValue = ^uint64(0) // Max uint64
	for i := 0; i < len(lastEntries); i++ {
		maxValue = max(maxValue, lastEntries[i])
		minValue = min(minValue, lastEntries[i])
	}

	if len(lastEntries) == 0 {
		minValue = 0
	}

	g.Box.SetTitle(fmt.Sprintf(" Messages / Second | Max: %d | Min: %d ", maxValue, minValue))

	stepSize := max((maxValue-minValue)/uint64(height-1), 1)
	barWidth := width / max(len(lastEntries), 1)

	g.mu.Lock()
	defer g.mu.Unlock()

	// Draw chart
	for i := 0; i < len(lastEntries); i++ {
		colHeight := max(int((lastEntries[i]-minValue)/stepSize), 1)
		colHeight = min(colHeight, height-1)

		for j := 0; j < colHeight; j++ {
			for k := 0; k < barWidth; k++ {
				// █
				tview.Print(screen, "⣿", x+(i*barWidth)+k, y+(height-1)-j, width+1, tview.AlignLeft, tcell.ColorWhite)
			}
		}
	}

	newIndex := g.dataIndex % len(g.buff)
	g.buff[newIndex] = g.Counter.Rate()
	g.dataIndex = newIndex + 1
}

func (g *GraphView) GetSize() int {
	var size int
	for _, entry := range g.buff {
		if entry != 0 {
			size++
		}
	}
	return size
}

func (g *GraphView) GetLastEntries(num int) []uint64 {
	num = min(num, g.GetSize())

	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]uint64, num)

	newIndex := calculateRingbufferOffset(len(g.buff), g.dataIndex, -num)
	for i := 0; i < num; i++ {
		index := (newIndex + i) % len(g.buff)
		result[i] = g.buff[index]
	}
	return result
}
