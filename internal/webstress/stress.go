package webstress

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/paulbellamy/ratecounter"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
	"webstress/internal/models"
)

type WebStress struct {
	Workers []*Worker
	Addr    string
	logger  *models.Logger
	wg      sync.WaitGroup
}

type Worker struct {
	addr         string
	pingInterval int
	WSData       *models.WebsocketData
	logger       *models.Logger
	wg           *sync.WaitGroup
}

func (webstress *WebStress) Init(url string, workerCount int, pingInterval int) {
	webstress.Addr = url
	webstress.logger.Log(fmt.Sprintf("Starting %d workers, pingInterval: %d\n", workerCount, pingInterval))
	counter := ratecounter.NewRateCounter(5 * time.Second)

	for i := 0; i < workerCount; i++ {
		w := Worker{addr: url, pingInterval: pingInterval, wg: &webstress.wg, WSData: models.NewWebsocketData(i, counter), logger: webstress.logger}
		webstress.Workers = append(webstress.Workers, &w)
	}
}

func (webstress *WebStress) Start() {
	// log.Printf("Starting %d workers, pingInterval: %d\n", workerCount, pingInterval)
	for _, worker := range webstress.Workers {
		go worker.run()
		time.Sleep(time.Millisecond * 10)
	}

	// log.Printf("Started %d workers, pingInterval: %d\n", workerCount, pingInterval)
	webstress.logger.Log(fmt.Sprintf("Started %d workers\n", len(webstress.Workers)))
	webstress.wg.Wait()
}

func (webstress *WebStress) CountMessages() (count int64) {
	for _, w := range webstress.Workers {
		count += w.MessageCount()
	}
	return count
}

func (webstress *WebStress) GetConnections() (count int) {
	for _, worker := range webstress.Workers {
		if worker.WSData.Connected {
			count++
		}
	}
	return count
}

func (webstress *WebStress) SetLogger(l *models.Logger) {
	webstress.logger = l
}

func (w Worker) MessageCount() int64 {
	return w.WSData.MessageCount()
}

func (w *Worker) run() {
	w.wg.Add(1)
	defer func() {
		w.WSData.Connected = false
		w.wg.Done()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	w.logger.Log(fmt.Sprintf("[Worker %d] connecting to %s\n", w.WSData.ID+1, w.addr))
	header := http.Header{"User-Agent": []string{"Webstress/1.0 (github.com/d-Rickyy-b/webstress)"}}
	c, _, err := websocket.DefaultDialer.Dial(w.addr, header)
	if err != nil {
		w.logger.Log(fmt.Sprintf("[Worker %d] Error while connecting to websocket: %v\n", w.WSData.ID+1, err))
		return
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
	go func() {
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					w.logger.Log(fmt.Sprintf("[Worker %d] Error while reading message: %v\n", w.WSData.ID+1, err))
					return
				}
				w.logger.Log(fmt.Sprintf("[Worker %d] Error while reading message: %v\n", w.WSData.ID+1, err))
				return
			}
			w.WSData.IncCounter()
		}
	}()

	w.WSData.Connected = true

	// Create ticker to send ping messages to the websocket on each interval
	pingTicker := time.NewTicker(time.Duration(w.pingInterval) * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-done:
			return
		case <-pingTicker.C:
			err := c.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				w.logger.Log(fmt.Sprintf("[Worker %d] Error while writing ping message: %v\n", w.WSData.ID+1, err))
				return
			}
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				w.logger.Log(fmt.Sprintf("[Worker %d] Error while writing close message: %v\n", w.WSData.ID+1, err))
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
