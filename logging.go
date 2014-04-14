package cinotify

import (
	"log"
)

// Logger can be set to any logger instance to get output, otherwise all
// logging information will be thrown out
var Logger *log.Logger

// DoLog logs anything of interest.
func DoLog(args ...interface{}) {
	if Logger != nil {
		Logger.Print(args...)
	}
}

// DoLogf logs anything of interest with a format
func DoLogf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Printf(format, args...)
	}
}
