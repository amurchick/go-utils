package utils

import (
	"github.com/amurchick/go-utils/finalizer"
)

var AtExit = finalizer.New()
var AtAlarm = finalizer.New()
