package models

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Logger struct {
	buff  []string
	index int
	mu    sync.RWMutex
}

func (l *Logger) Log(text string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	newIndex := l.index % len(l.buff)
	l.buff[newIndex] = text
	l.index = newIndex + 1
}

func (l *Logger) GetSize() int {
	var size int
	for _, entry := range l.buff {
		if entry != "" {
			size++
		}
	}
	return size
}

func (l *Logger) GetLastEntries(num int) []string {
	num = min(num, l.GetSize())

	l.mu.RLock()
	defer l.mu.RUnlock()
	result := make([]string, num)

	newIndex := calculateRingbufferOffset(len(l.buff), l.index, -num)
	for i := 0; i < num; i++ {
		index := (newIndex + i) % len(l.buff)
		result[i] = l.buff[index]
	}
	return result
}

func calculateRingbufferOffset(bufferLen, currentIndex, offset int) int {
	newIndex := currentIndex + offset
	for newIndex < 0 {
		newIndex = bufferLen + newIndex
	}

	return newIndex % (bufferLen)
}

func (l *Logger) GetLastEntry() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.buff[l.index]
}

func (l *Logger) String() string {
	return fmt.Sprintf("%v", l.buff)
}

func NewLogger() *Logger {
	return &Logger{buff: make([]string, 100)}
}

type LoggerView struct {
	*tview.Box
	Logger *Logger
}

func NewLoggerView() *LoggerView {
	box := tview.NewBox()
	box.SetTitle(" Logger ").SetBorder(true)
	return &LoggerView{Logger: NewLogger(), Box: box}
}

func (l LoggerView) Draw(screen tcell.Screen) {
	l.Box.DrawForSubclass(screen, l)

	// Draw label.
	x, y, width, height := l.GetInnerRect()
	if width <= 0 || height <= 0 {
		return
	}

	entries := l.Logger.GetLastEntries(height)
	for i := 0; i < height; i++ {
		if i >= len(entries) {
			break
		}
		tview.Print(screen, entries[i], x+1, y+i, width-2, tview.AlignLeft, tcell.ColorWhite)
	}
}
