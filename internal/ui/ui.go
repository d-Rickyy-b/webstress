package ui

import (
	"github.com/rivo/tview"
	"time"
	"webstress/internal/models"
	"webstress/internal/webstress"
)

type UI struct {
	app          *tview.Application
	Logger       *models.Logger
	websocketBox *models.WebsocketView
	loggerView   *models.LoggerView
	graphView    *models.GraphView
	infoBox      *models.InfoBox
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
	ui.graphView = models.NewGraphView()
	ui.infoBox = models.NewInfoBox()

	// Layout the UI
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.websocketBox, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(ui.loggerView, 0, 1, false).
				AddItem(ui.graphView, 0, 1, false),
				0, 1, false).
			AddItem(ui.infoBox, 1, 0, false),
			0, 2, false)

	ui.app = app.SetRoot(flex, true).SetFocus(flex)
}

func (ui *UI) RegisterWebstress(ws *webstress.WebStress) {
	for i, worker := range ws.Workers {
		ui.websocketBox.Connections[i] = worker.WSData
	}
	ui.infoBox.Addr = ws.Addr
	ui.infoBox.Connections = ui.websocketBox.Connections
	ui.infoBox.MsgCounter = ws.MsgCounter
	ui.graphView.Counter = ws.MsgCounter
}

func (ui *UI) Run() {
	go func() {
		tick := time.NewTicker(200 * time.Millisecond)
		for {
			select {
			case <-tick.C:
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
