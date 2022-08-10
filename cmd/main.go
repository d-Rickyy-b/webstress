package main

import (
	"flag"
	"strings"
	"webstress/internal/models"
	"webstress/internal/ui"
	"webstress/internal/webstress"
)

func main() {
	var workerCount = flag.Int("workerCount", 30, "number of workers to start")
	var pingInterval = flag.Int("pingInterval", 30, "number of seconds between pings")
	var remoteAddr = flag.String("remoteAddr", "ws://localhost:8080/", "remote address to connect to")
	var noUI = flag.Bool("noUI", false, "use to disable the UI")
	_ = noUI
	flag.Parse()

	if !strings.HasPrefix(*remoteAddr, "ws://") && !strings.HasPrefix(*remoteAddr, "wss://") {
		*remoteAddr = "wss://" + *remoteAddr
	}

	// TODO feature: flag to enable or disable CUI
	// if useCUI{
	// 	ui.StartUI()
	// }

	logger := models.NewLogger()
	cui := ui.NewUI(*workerCount, logger)

	stress := &webstress.WebStress{}
	stress.SetLogger(logger)
	stress.Init(*remoteAddr, *workerCount, *pingInterval)
	go stress.Start()

	cui.RegisterWebstress(stress)
	cui.Run()
}
