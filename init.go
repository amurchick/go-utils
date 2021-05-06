package utils

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/amurchick/go-utils/logger"
)

var log logger.Interface

func SetLogger(logger logger.Interface) {
	log = logger
}

func init() {
	log = logger.Log
	signalChannel := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP)
	signal.Notify(signalChannel,
		os.Interrupt,
		syscall.SIGALRM,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	go func() {
		for signal := range signalChannel {
			if signal == syscall.SIGALRM {
				log.Warn("Signal %q received", signal.String())
				go AtAlarm.Run()
			} else {
				log.Warn("signal %q received, exiting...", signal.String())
				Exit()
			}
		}
	}()
}
