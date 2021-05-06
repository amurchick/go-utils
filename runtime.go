package utils

import (
	"os"
	"sync/atomic"

	"github.com/amurchick/go-utils/finalizer"
)

var AtExit = finalizer.New()
var AtAlarm = finalizer.New()

var exitInProgress uint32

func Exit(code ...interface{}) {
	if !atomic.CompareAndSwapUint32(&exitInProgress, 0, 1) {
		return
	}
	log.Info("stopping...")

	AtExit.Run()

	exitCode := 0
	needExit := true
	if len(code) == 1 {
		if _exitCode, ok := code[0].(int); ok {
			exitCode = _exitCode
		}
		if _needExit, ok := code[0].(bool); ok {
			needExit = _needExit
		}
	}
	log.Info("stopped!\n\n")
	if needExit {
		os.Exit(exitCode)
	}
}
