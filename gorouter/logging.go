package main

import (
	"fmt"
	"log"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	ERROR = "ERROR"
	WARN  = "WARN"
)

func logMessage(level string, msg string) {
	log.Printf("[%s] - %s", level, msg)
}

func logInfo(format string, v ...interface{}) {
	logMessage(INFO, fmt.Sprintf(format, v...))
}

func logWarn(format string, v ...interface{}) {
	logMessage(WARN, fmt.Sprintf(format, v...))
}

func logError(format string, e error, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if e != nil {
		msg = fmt.Sprintf("%s - %s", msg, e)
	}
	logMessage(ERROR, msg)
}

func logDebug(format string, v ...interface{}) {
	logMessage(DEBUG, fmt.Sprintf(format, v...))
}
