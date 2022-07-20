package models

import (
	"fmt"
	"log"
	"testing"
)

func TestLogger_Log(t *testing.T) {
	l := NewLogger()
	for i := 0; i < 100; i++ {
		l.Log(fmt.Sprintf("Test %d", i))
	}
}

func TestGetEntries(t *testing.T) {
	l := NewLogger()
	for i := 0; i < 35; i++ {
		l.Log(fmt.Sprintf("Test %d", i))
	}

	entries := l.GetLastEntries(10)
	log.Println(entries)

	entries = l.GetLastEntries(5)
	log.Println(entries)

	entries = l.GetLastEntries(100)
	log.Println(entries)
}

func TestCalculateRingbufferOffset(t *testing.T) {
	type args struct {
		bufferLen    int
		currentIndex int
		offset       int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{name: "Test negative in range", args: args{bufferLen: 6, currentIndex: 5, offset: -2}, want: 3},
		{name: "Test negative out of range", args: args{bufferLen: 6, currentIndex: 5, offset: -5}, want: 0},
		{name: "Test negative out of range 2", args: args{bufferLen: 6, currentIndex: 5, offset: -6}, want: 5},
		{name: "Test negative out of range 3", args: args{bufferLen: 10, currentIndex: 5, offset: -20}, want: 5},
		{name: "Test negative out of range 4", args: args{bufferLen: 10, currentIndex: 5, offset: -21}, want: 4},
		{name: "Test positive in range", args: args{bufferLen: 6, currentIndex: 4, offset: 1}, want: 5},
		{name: "Test positive in range 2", args: args{bufferLen: 100, currentIndex: 5, offset: 26}, want: 31},
		{name: "Test positive out of range", args: args{bufferLen: 6, currentIndex: 5, offset: 2}, want: 1},
		{name: "Test positive out of range 2", args: args{bufferLen: 10, currentIndex: 9, offset: 1}, want: 0},
		{name: "Test positive out of range 3", args: args{bufferLen: 30, currentIndex: 17, offset: 20}, want: 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateRingbufferOffset(tt.args.bufferLen, tt.args.currentIndex, tt.args.offset); got != tt.want {
				t.Errorf("calculateRingbufferOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}
