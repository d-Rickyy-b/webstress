package models

import (
	"github.com/paulbellamy/ratecounter"
	"sync/atomic"
	"time"
)

type MsgCounter struct {
	rateCounter  *ratecounter.RateCounter
	messageCount uint64
	counterRate  int64
}

func NewMsgCounter(counterRate int64) *MsgCounter {
	interval := time.Duration(counterRate) * time.Second
	return &MsgCounter{rateCounter: ratecounter.NewRateCounter(interval), counterRate: counterRate}
}

func (m *MsgCounter) Incr() {
	m.rateCounter.Incr(1)
	atomic.AddUint64(&m.messageCount, 1)
}

func (m *MsgCounter) Count() uint64 {
	return atomic.LoadUint64(&m.messageCount)
}

func (m *MsgCounter) Rate() uint64 {
	return uint64(m.rateCounter.Rate() / m.counterRate)
}

func (m *MsgCounter) Set(value uint64) {
	atomic.StoreUint64(&m.messageCount, value)
}
