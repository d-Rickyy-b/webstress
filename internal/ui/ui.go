package ui

import (
	"github.com/rivo/tview"
	"time"
	"webstress/internal/models"
	"webstress/internal/webstress"
)

type UI struct {
	app               *tview.Application
	Logger            *models.Logger
	websocketBox      *models.WebsocketView
	loggerView        *models.LoggerView
	messagesPerSecBox *tview.Box
	infoBox           *models.InfoBox
	addr              string
}

func NewUI(websocketCount int, logger *models.Logger) *UI {
	ui := &UI{}
	ui.Init(websocketCount)
	ui.SetLogger(logger)
	return ui
}

func (ui *UI) Init(websocketCount int) {
	app := tview.NewApplication()
	ui.loggerView = models.NewLoggerView()
	ui.websocketBox = models.NewWebsocketView(websocketCount)
	ui.messagesPerSecBox = tview.NewBox().SetBorder(true).SetTitle(" Messages / Second ")
	ui.infoBox = models.NewInfoBox()

	// Layout the UI
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.websocketBox, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(ui.loggerView, 0, 1, false).
				AddItem(ui.messagesPerSecBox, 0, 1, false),
				0, 1, false).
			AddItem(ui.infoBox, 1, 0, false),
			0, 2, false)

	ui.app = app.SetRoot(flex, true).SetFocus(flex)
}

func (ui *UI) RegisterWebstress(ws *webstress.WebStress) {
	ui.addr = ws.Addr
	ui.infoBox.Addr = ws.Addr
	for i, worker := range ws.Workers {
		ui.websocketBox.Connections[i] = worker.WSData
	}
	ui.infoBox.Connections = ui.websocketBox.Connections
}

func (ui *UI) Run() {
	go func() {
		tick := time.NewTicker(150 * time.Millisecond)
		for {
			select {
			case <-tick.C:
				//app.Draw()
				ui.app.QueueUpdateDraw(func() {})
			}
		}
	}()

	if err := ui.app.Run(); err != nil {
		panic(err)
	}
}

func (ui *UI) SetLogger(logger *models.Logger) {
	ui.loggerView.Logger = logger
	ui.Logger = logger
}
