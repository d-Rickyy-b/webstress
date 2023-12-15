package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"

	"github.com/d-Rickyy-b/webstress/internal/models"
	"github.com/d-Rickyy-b/webstress/internal/ui"
	"github.com/d-Rickyy-b/webstress/internal/webstress"
)

func main() {
	parser := argparse.NewParser("webstress", "Websocket stress tool developed in Go")
	parser.ExitOnHelp(true)

	remoteAddr := parser.String("a", "remote-addr", &argparse.Options{Required: true, Help: "remote address to connect to"})
	recoverError := parser.Flag("r", "recover", &argparse.Options{Required: false, Help: "recover from certain errors", Default: true})
	pingInterval := parser.Int("p", "ping-interval", &argparse.Options{Required: false, Help: "number of seconds between pings", Default: 30})
	workerCount := parser.Int("w", "worker-count", &argparse.Options{Required: false, Help: "number of workers to start", Default: 30})
	noUI := parser.Flag("", "noUI", &argparse.Options{Required: false, Help: "use to disable the UI", Default: false})
	_ = noUI

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

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
	stress.Init(*remoteAddr, *workerCount, *pingInterval, *recoverError)
	go stress.Start()

	cui.RegisterWebstress(stress)
	cui.Run()
}
