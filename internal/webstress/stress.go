package webstress

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"

	"github.com/d-Rickyy-b/webstress/internal/models"
)

var RecoverableError = errors.New("recoverable error")
var UnrecoverableError = errors.New("unrecoverable error")

type WebStress struct {
	Workers    []*Worker
	Addr       string
	logger     *models.Logger
	MsgCounter *models.MsgCounter
	wg         sync.WaitGroup
}

type Worker struct {
	addr         string
	pingInterval int
	WSData       *models.WebsocketData
	logger       *models.Logger
	wg           *sync.WaitGroup
	recover      bool
	rateLimit    int
}

func (webstress *WebStress) Init(url string, workerCount, pingInterval, rateLimit int, recover bool) {
	webstress.Addr = url
	webstress.logger.Log(fmt.Sprintf("Starting %d workers, pingInterval: %d\n", workerCount, pingInterval))
	webstress.MsgCounter = models.NewMsgCounter(5)

	for i := 0; i < workerCount; i++ {
		w := Worker{
			addr:         url,
			pingInterval: pingInterval,
			wg:           &webstress.wg,
			WSData:       models.NewWebsocketData(i, webstress.MsgCounter),
			logger:       webstress.logger,
			recover:      recover,
			rateLimit:    rateLimit,
		}
		webstress.Workers = append(webstress.Workers, &w)
	}
}

func (webstress *WebStress) Start() {
	// log.Printf("Starting %d workers, pingInterval: %d\n", workerCount, pingInterval)
	for _, worker := range webstress.Workers {
		go worker.runWithRecovery()
		time.Sleep(time.Millisecond * 10)
	}

	// log.Printf("Started %d workers, pingInterval: %d\n", workerCount, pingInterval)
	webstress.logger.Log(fmt.Sprintf("Started %d workers\n", len(webstress.Workers)))
	webstress.wg.Wait()
}

func (webstress *WebStress) SetLogger(l *models.Logger) {
	webstress.logger = l
}

func (w *Worker) runWithRecovery() {
	w.wg.Add(1)
	defer func() {
		w.WSData.Connected = false
		w.wg.Done()
	}()

	// Run worker until it crashes
	for {
		err := w.run()

		if !w.recover || (err != nil && errors.Is(err, UnrecoverableError)) {
			return
		}
		w.logger.Log(fmt.Sprintf("[Worker %d] Reconnecting...\n", w.WSData.ID+1))
	}
}

func (w *Worker) run() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	w.logger.Log(fmt.Sprintf("[Worker %d] connecting to %s\n", w.WSData.ID+1, w.addr))
	header := http.Header{"User-Agent": []string{"Webstress/1.0 (github.com/d-Rickyy-b/webstress)"}}
	c, _, err := websocket.DefaultDialer.Dial(w.addr, header)
	if err != nil {
		w.logger.Log(fmt.Sprintf("[Worker %d] Error while connecting to websocket: %v\n", w.WSData.ID+1, err))
		return UnrecoverableError
	}
	defer c.Close()

	c.SetCloseHandler(func(code int, text string) error {
		switch code {
		case websocket.CloseNoStatusReceived:
			w.logger.Log(fmt.Sprintf("[Worker %d] Stopping websocket - Server received no ping\n", w.WSData.ID+1))
		default:
			w.logger.Log(fmt.Sprintf("[Worker %d] Stopping websocket - %s\n", w.WSData.ID+1, text))
		}
		message := websocket.FormatCloseMessage(code, "Connection closed")
		return c.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
	})

	done := make(chan struct{})

	// Start background goroutine to read messages from the websocket
	go w.listen(done, c)

	w.WSData.Connected = true

	// Create ticker to send ping messages to the websocket on each interval
	pingTicker := time.NewTicker(time.Duration(w.pingInterval) * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-pingTicker.C:
			writeErr := c.WriteMessage(websocket.PingMessage, []byte{})
			if writeErr != nil {
				w.logger.Log(fmt.Sprintf("[Worker %d] Error while writing ping message: %v\n", w.WSData.ID+1, writeErr))
				return RecoverableError
			}
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			writeErr := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if writeErr != nil {
				w.logger.Log(fmt.Sprintf("[Worker %d] Error while writing close message: %v\n", w.WSData.ID+1, writeErr))
				return UnrecoverableError
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return UnrecoverableError
		}
	}
}

// listen listens for messages from the websocket and increments the message counter.
func (w *Worker) listen(done chan struct{}, c *websocket.Conn) {
	l := rate.NewLimiter(rate.Limit(w.rateLimit), 10)

	defer close(done)
	for {
		if w.rateLimit > 0 {
			waitErr := l.Wait(context.Background())
			if waitErr != nil {
				w.logger.Log(fmt.Sprintf("[Worker %d] Error while waiting for rate limit: %v\n", w.WSData.ID+1, waitErr))
				return
			}
		}

		_, _, readErr := c.ReadMessage()
		if readErr != nil {
			if websocket.IsUnexpectedCloseError(readErr, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				w.logger.Log(fmt.Sprintf("[Worker %d] Error while reading message: %v\n", w.WSData.ID+1, readErr))
				return
			}
			w.logger.Log(fmt.Sprintf("[Worker %d] Error while reading message: %v\n", w.WSData.ID+1, readErr))
			return
		}
		w.WSData.IncCounter()
	}
}
